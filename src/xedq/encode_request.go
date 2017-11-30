package xedq

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// argTag represents instruction argument (operand) class.
type argTag uint32

// Supported argument classes.
const (
	argEmpty argTag = iota
	argReg
	argMem
	argUint8
	argUint32
	argUint64
	argInt8
	argInt16
	argInt32
	argRel8
	argRel16
	argRel32
)

// effectiveOperandSize is XED's EOSZ attribute.
// Specifies instruction data size.
type effectiveOperandSize uint8

// All valid effective operand sizes.
const (
	eoszDefault effectiveOperandSize = iota
	eosz8
	eosz16
	eosz32
	eosz64
)

func (eosz effectiveOperandSize) String() string {
	switch eosz {
	case eosz8:
		return "8"
	case eosz16:
		return "16"
	case eosz32:
		return "32"
	case eosz64:
		return "64"
	default:
		return "??"
	}
}

// EncodeRequest is an instruction builder.
//
// Methods that have no "Set" prefix push argument
// into instruction being constructed.
// Setters only set specified field and are safe
// to be called repeatedly.
//
// When all options and arguments are set,
// instruction can be created by one of the Encode methods.
type EncodeRequest struct {
	encoder *Encoder // Encoder that spawned this EncodeRequest
	iclass  xedIclass

	// Immediate operand payload.
	// TODO: add second immediate field?
	imm uint64

	rel int32

	ptr      Ptr
	memWidth uint16

	// Holds each argument type.
	tags [maxArgLimit]argTag

	// Actual number of arguments set.
	argc uint8

	eosz effectiveOperandSize

	// Register arguments.
	regs [maxArgLimit]xedRegister
}

// Reg pushes register with name regName to arguments list.
func (req *EncodeRequest) Reg(regName string) *EncodeRequest {
	req.pushReg(registerByName[regName])
	return req
}

// Mem pushes memory indirect to arguments list.
//
// Width is a pointer size in bits.
// Common values are:
//   8   | BYTE PTR
//   16  | WORD PTR
//   32  | DWORD PTR
//   64  | QWORD PTR
//   80  | TBYTE PTR (x87)
//   128 | XMMWORD PTR
//   256 | YMMWORD PTR
//   512 | ZMMWORD PTR
func (req *EncodeRequest) Mem(width uint16, ptr Ptr) *EncodeRequest {
	req.pushTag(argMem)
	req.ptr = ptr
	req.memWidth = width
	return req
}

// MemExpr is like Mem, but uses mem expr string to specify effective address.
// expr format/syntax depends on the Encoder.MemExprParser.
//
// Panics if expr string is malformed.
func (req *EncodeRequest) MemExpr(width uint16, expr string) *EncodeRequest {
	ptr, err := req.encoder.MemExprParser(expr)
	if err != nil {
		panic(err)
	}
	return req.Mem(width, ptr)
}

// Uint8 pushes 8bit unsigned immediate to argument list.
// Notice: current implementation is limited to single immediate, so
// instructions like ENTER are not encodable yet.
func (req *EncodeRequest) Uint8(v uint8) *EncodeRequest {
	req.imm = uint64(v)
	req.pushTag(argUint8)
	return req
}

// Uint32 pushes 32bit unsigned immediate to argument list.
// Notice: current implementation is limited to single immediate, so
// instructions like ENTER are not encodable yet.
func (req *EncodeRequest) Uint32(v uint32) *EncodeRequest {
	req.imm = uint64(v)
	req.pushTag(argUint32)
	return req
}

// Int8 pushes 8bit signed immediate to argument list.
// Notice: current implementation is limited to single immediate, so
// instructions like ENTER are not encodable yet.
func (req *EncodeRequest) Int8(v int8) *EncodeRequest {
	req.imm = uint64(v)
	req.pushTag(argInt8)
	return req
}

// Int16 pushes 16bit signed immediate to argument list.
// Notice: current implementation is limited to single immediate, so
// instructions like ENTER are not encodable yet.
func (req *EncodeRequest) Int16(v int16) *EncodeRequest {
	req.imm = uint64(v)
	req.pushTag(argInt16)
	return req
}

// Int32 pushes 32bit signed immediate to argument list.
// Notice: current implementation is limited to single immediate, so
// instructions like ENTER are not encodable yet.
func (req *EncodeRequest) Int32(v int32) *EncodeRequest {
	req.imm = uint64(v)
	req.pushTag(argInt32)
	return req
}

// Rel8 pushes 8bit branch displacement.
func (req *EncodeRequest) Rel8(v int8) *EncodeRequest {
	req.rel = int32(v)
	req.pushTag(argRel8)
	return req
}

// Rel16 pushes 16bit branch displacement.
func (req *EncodeRequest) Rel16(v int16) *EncodeRequest {
	req.rel = int32(v)
	req.pushTag(argRel16)
	return req
}

// Rel32 pushes 32bit branch displacement.
func (req *EncodeRequest) Rel32(v int32) *EncodeRequest {
	req.rel = v
	req.pushTag(argRel32)
	return req
}

// SetEosz8 sets instruction effective operand size to 8bit.
func (req *EncodeRequest) SetEosz8() *EncodeRequest {
	req.eosz = eosz8
	return req
}

// SetEosz16 sets instruction effective operand size to 16bit.
func (req *EncodeRequest) SetEosz16() *EncodeRequest {
	req.eosz = eosz16
	return req
}

// SetEosz32 sets instruction effective operand size to 32bit.
func (req *EncodeRequest) SetEosz32() *EncodeRequest {
	req.eosz = eosz32
	return req
}

// SetEosz64 sets instruction effective operand size to 64bit.
func (req *EncodeRequest) SetEosz64() *EncodeRequest {
	req.eosz = eosz64
	return req
}

// Encode executes encode request and returns result "as it".
func (req *EncodeRequest) Encode() []byte {
	return req.encoder.encode(req)
}

// EncodeTo is like Encode, but instead of allocating new byte slice,
// it writes output to w.
// Returns w.Write() result.
func (req *EncodeRequest) EncodeTo(w io.Writer) (int, error) {
	return req.encoder.encodeTo(w, req)
}

// EncodeHexString executes encode request and formats result as a hex string.
func (req *EncodeRequest) EncodeHexString() string {
	code := req.Encode()
	if code == nil {
		return ""
	}
	var buf bytes.Buffer
	for i := range code {
		fmt.Fprintf(&buf, "%02x", code[i])
	}
	return buf.String()
}

// String returns assembly-like instruction representation.
// Intended for debugging and pretty-printing (useful in tests).
func (req *EncodeRequest) String() string {
	var name string
	eosz := req.eosz.String()
	if eosz != "" {
		name = req.iclass.String() + "/" + eosz
	} else {
		name = req.iclass.String()
	}
	if req.argc == 0 {
		return name
	}

	args := make([]string, req.argc)

	for i := 0; i < int(req.argc); i++ {
		switch req.tags[i] {
		default:
			args[i] = "??"
		case argReg:
			args[i] = req.regs[i].String()
		case argMem:
			args[i] = memExprString(req)
		case argUint8:
			args[i] = fmt.Sprintf("uint8(%#x)", req.imm)
		case argUint32:
			args[i] = fmt.Sprintf("uint32(%#x)", req.imm)
		case argUint64:
			args[i] = fmt.Sprintf("uint64(%#x)", req.imm)
		case argInt8:
			args[i] = fmt.Sprintf("int8(%#x)", int32(req.imm))
		case argInt16:
			args[i] = fmt.Sprintf("int16(%#x)", int32(req.imm))
		case argInt32:
			args[i] = fmt.Sprintf("int32(%#x)", int32(req.imm))
		case argRel8:
			args[i] = fmt.Sprintf("rel8(%#x)", req.rel)
		case argRel16:
			args[i] = fmt.Sprintf("rel16(%#x)", req.rel)
		case argRel32:
			args[i] = fmt.Sprintf("rel32(%#x)", req.rel)
		}
	}

	return name + " " + strings.Join(args, ", ")
}

func (req *EncodeRequest) pushTag(tag argTag) {
	req.tags[req.argc] = tag
	req.argc++
}

func (req *EncodeRequest) pushReg(reg xedRegister) {
	req.regs[req.argc] = reg
	req.pushTag(argReg)
}
