package xedq

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// TODO: may want to reduce coupling between Encoder and EncodeRequest.
// Benchmarks are needed to make decision making easier.
// Also need to track allocations.

// argTag represents instruction argument (operand) class.
type argTag uint32

// All valid argument classes.
const (
	argEmpty argTag = iota
	argReg
	argMem8  // BYTE PTR
	argMem16 // WORD PTR
	argMem32 // DWORD PTR
	argMem64 // QWORD PTR
	argUint8
	argUint32
	argUint64
	argInt8
	argInt32
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
	name    string   // Name of the instruction

	// Immediate operand payload.
	// TODO: add second immediate field?
	imm uint64

	// Memory argument fields.
	// There is always at most one memory argument.

	memDisp      uint64   // Displacement amount
	memBase      register // SIB - B
	memIndex     register // SIB - I
	memScale     uint8    // SIB - S
	memDispWidth uint8    // 8/16/32/64

	// Holds each argument type.
	tags [maxArgLimit]argTag

	// Actual number of arguments set.
	argc uint8

	eosz effectiveOperandSize

	// Register arguments.
	regs [maxArgLimit]register
}

// Reg pushes register with name regName to arguments list.
func (req *EncodeRequest) Reg(regName string) *EncodeRequest {
	req.pushReg(registerByName(regName, &req.encoder.tmpbuf))
	return req
}

// Mem8 pushes 8bit memory indirect to arguments list.
//
// sibExpr format is akin to Intel addressing syntax:
//   - "BASE"
//   - "BASE+INDEX"
//   - "INDEX*SCALE"
//   - "BASE+INDEX*SCALE"
// BASE and INDEX are any valid registers.
// SCALE can be 1/2/4/8.
//
// Examples for sibExpr:
//   "RAX"
//   "RAX+RCX"
//   "RDX*2"
//   "RAX+RCX*4"
//   "RAX+XMM0*4" // VSIB
//
// To set base/index/scale individually, one can call Mem8 with empty string
// argument, and then use SetMemBase/SetMemIndex/SetMemScale to set
// them ony-by-one.
func (req *EncodeRequest) Mem8(sibExpr string) *EncodeRequest {
	req.memBase, req.memIndex, req.memScale = parseSIBExpr(sibExpr, &req.encoder.tmpbuf)
	req.pushTag(argMem8)
	return req
}

// Mem16 pushes 16bit memory indirect to arguments list.
// sibExpr has same format as in Mem8.
func (req *EncodeRequest) Mem16(sibExpr string) *EncodeRequest {
	req.memBase, req.memIndex, req.memScale = parseSIBExpr(sibExpr, &req.encoder.tmpbuf)
	req.pushTag(argMem16)
	return req
}

// Mem32 pushes 32bit memory indirect to arguments list.
// sibExpr has same format as in Mem8.
func (req *EncodeRequest) Mem32(sibExpr string) *EncodeRequest {
	req.memBase, req.memIndex, req.memScale = parseSIBExpr(sibExpr, &req.encoder.tmpbuf)
	req.pushTag(argMem32)
	return req
}

// Mem64 pushes 64bit memory indirect to arguments list.
// sibExpr has same format as in Mem8.
func (req *EncodeRequest) Mem64(sibExpr string) *EncodeRequest {
	req.memBase, req.memIndex, req.memScale = parseSIBExpr(sibExpr, &req.encoder.tmpbuf)
	req.pushTag(argMem64)
	return req
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

// Int32 pushes 32bit signed immediate to argument list.
// Notice: current implementation is limited to single immediate, so
// instructions like ENTER are not encodable yet.
func (req *EncodeRequest) Int32(v int32) *EncodeRequest {
	req.imm = uint64(v)
	req.pushTag(argInt32)
	return req
}

// SetMemBase assigns regName as a base addressing register.
func (req *EncodeRequest) SetMemBase(regName string) *EncodeRequest {
	req.memBase = registerByName(regName, &req.encoder.tmpbuf)
	return req
}

// SetMemIndex assigns regName as scaled indexing register.
func (req *EncodeRequest) SetMemIndex(regName string) *EncodeRequest {
	req.memIndex = registerByName(regName, &req.encoder.tmpbuf)
	return req
}

// SetMemScale sets indexing register scaling to specified value.
func (req *EncodeRequest) SetMemScale(scale int) *EncodeRequest {
	req.memScale = uint8(scale)
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

// SetDisp8 sets memory operand displacement to 8bit value disp.
func (req *EncodeRequest) SetDisp8(disp uint8) *EncodeRequest {
	req.memDisp = uint64(disp)
	req.memDispWidth = 8
	return req
}

// SetDisp32 sets memory operand displacement to 32bit value disp.
func (req *EncodeRequest) SetDisp32(disp uint32) *EncodeRequest {
	req.memDisp = uint64(disp)
	req.memDispWidth = 32
	return req
}

// Encode executes encode request and returns result "as it".
func (req *EncodeRequest) Encode() []byte {
	n := req.encoder.encode(req)
	code := make([]byte, n)
	copy(code, req.encoder.tmpbuf.data[:])
	return code
}

// EncodeTo is like Encode, but instead of allocating new byte slice,
// it writes output to w.
// Returns w.Write result.
func (req *EncodeRequest) EncodeTo(w io.Writer) (int, error) {
	n := req.encoder.encode(req)
	return w.Write(req.encoder.tmpbuf.data[:n])
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
	if req.argc == 0 {
		return req.name
	}

	args := make([]string, req.argc)

	for i := 0; i < int(req.argc); i++ {
		switch req.tags[i] {
		default:
			args[i] = "<?>"
		case argReg:
			args[i] = req.regs[i].String()
		case argMem32:
			args[i] = "mem32" + sibString(req.memBase, req.memIndex, req.memScale)
		case argMem64:
			args[i] = "mem64" + sibString(req.memBase, req.memIndex, req.memScale)
		case argUint8:
			args[i] = fmt.Sprintf("uint8(%#x)", req.imm)
		case argUint32:
			args[i] = fmt.Sprintf("uint32(%#x)", req.imm)
		case argUint64:
			args[i] = fmt.Sprintf("uint64(%#x)", req.imm)
		case argInt8:
			args[i] = fmt.Sprintf("int8(%#x)", int32(req.imm))
		case argInt32:
			args[i] = fmt.Sprintf("int32(%#x)", int32(req.imm))
		}
	}

	return req.name + " " + strings.Join(args, ", ")
}

func (req *EncodeRequest) pushTag(tag argTag) {
	req.tags[req.argc] = tag
	req.argc++
}

func (req *EncodeRequest) pushReg(reg register) {
	req.regs[req.argc] = reg
	req.pushTag(argReg)
}
