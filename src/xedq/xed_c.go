package xedq

/*
#cgo LDFLAGS: -lxed
#include <xed/xed-interface.h>
*/
import "C"

import (
	"unsafe"
)

// This is the only file that imports "C".
//
// Most definitions are prefixed with "xed" to make it clear
// that those functions are low-level and can be unsafe.

// xedRegister holds XED register ID.
type xedRegister uint16

// Constants that define register indexes from XED xed_reg_enum_t.
const (
	regGPR8  = xedRegister(C.XED_REG_AL)
	regGPR16 = xedRegister(C.XED_REG_AX)
	regGPR32 = xedRegister(C.XED_REG_EAX)
	regGPR64 = xedRegister(C.XED_REG_RAX)
	regXMM   = xedRegister(C.XED_REG_XMM0)
	regYMM   = xedRegister(C.XED_REG_YMM0)
	regZMM   = xedRegister(C.XED_REG_ZMM0)
	regK     = xedRegister(C.XED_REG_K0)
)

const (
	xedRegInvalid    = xedRegister(C.XED_REG_INVALID)
	xedIclassInvalid = xedIclass(C.XED_ICLASS_INVALID)
)

var (
	errEmpty = xedError(C.XED_ERROR_NONE)
)

// String returns reg name.
func (reg xedRegister) String() string {
	return C.GoString(C.xed_reg_enum_t2str(C.xed_reg_enum_t(reg)))
}

func xedTablesInit() { C.xed_tables_init() }

type (
	xedState  C.xed_state_t
	xedInst   C.xed_encoder_instruction_t
	xedIclass C.xed_iclass_enum_t
	xedError  C.xed_error_enum_t
)

func (inst *xedInst) CPtr() *C.xed_encoder_instruction_t {
	return (*C.xed_encoder_instruction_t)(inst)
}

func newXEDState32() xedState {
	var state C.xed_state_t
	C.xed_state_zero(&state)
	state.stack_addr_width = C.XED_ADDRESS_WIDTH_32b
	state.mmode = C.XED_MACHINE_MODE_LEGACY_32
	return xedState(state)
}

func (err xedError) CValue() C.xed_error_enum_t {
	return C.xed_error_enum_t(err)
}

func (err xedError) Empty() bool {
	return err == C.XED_ERROR_NONE
}

func (err xedError) Error() string {
	// TODO: check if XED does bound check for error codes.
	return "XED error: " + C.GoString(C.xed_error_enum_t2str(err.CValue()))
}

func (state xedState) CValue() C.xed_state_t {
	return C.xed_state_t(state)
}

func newXEDState64() xedState {
	var state C.xed_state_t
	C.xed_state_zero(&state)
	state.stack_addr_width = C.XED_ADDRESS_WIDTH_64b
	state.mmode = C.XED_MACHINE_MODE_LONG_64
	return xedState(state)
}

func (iclass xedIclass) String() string {
	return C.GoString(C.xed_iclass_enum_t2str(iclass.CValue()))
}

func (iclass xedIclass) CValue() C.xed_iclass_enum_t {
	return C.xed_iclass_enum_t(iclass)
}

func newXEDIclass(name string, tmpbuf *buffer) xedIclass {
	tmpbuf.SetCString(name)
	iclass := C.str2xed_iclass_enum_t(tmpbuf.CString())
	return xedIclass(iclass)
}

func newXEDInst(state *xedState, req *EncodeRequest) xedInst {
	var inst C.xed_encoder_instruction_t

	iclass := req.iclass

	var eosz C.xed_uint_t
	switch req.eosz {
	case eosz8:
		eosz = 8
	case eosz16:
		eosz = 16
	case eosz32:
		eosz = 32
	case eosz64:
		eosz = 64
	default:
		//TODO: try to set with respect to DF64 and DF32.
		eosz = 32
	}

	// It is possible to initialize inst operands directly,
	// but that is more likely to break than xed_instN API,
	// which is explicitly public.
	switch req.argc {
	default:
		panic("unexpected args count")
	case 0:
		C.xed_inst0(&inst, state.CValue(), iclass.CValue(), eosz)
	case 1:
		C.xed_inst1(&inst, state.CValue(), iclass.CValue(), eosz,
			xedOperand(req, 0))
	case 2:
		C.xed_inst2(&inst, state.CValue(), iclass.CValue(), eosz,
			xedOperand(req, 0),
			xedOperand(req, 1))
	case 3:
		C.xed_inst3(&inst, state.CValue(), iclass.CValue(), eosz,
			xedOperand(req, 0),
			xedOperand(req, 1),
			xedOperand(req, 2))
	case 4:
		C.xed_inst4(&inst, state.CValue(), iclass.CValue(), eosz,
			xedOperand(req, 0),
			xedOperand(req, 1),
			xedOperand(req, 2),
			xedOperand(req, 3))
	case 5:
		C.xed_inst5(&inst, state.CValue(), iclass.CValue(), eosz,
			xedOperand(req, 0),
			xedOperand(req, 1),
			xedOperand(req, 2),
			xedOperand(req, 3),
			xedOperand(req, 4))
	}

	return xedInst(inst)
}

func xedEncode(inst *xedInst, dstbuf *buffer) (int, error) {
	var req C.xed_encoder_request_t
	C.xed_encoder_request_zero_set_mode(&req, &inst.mode)
	ok := C.xed_convert_to_encoder_request(&req, inst.CPtr())
	if ok == 0 {
		return 0, errEncReqConvert
	}

	codeLen := C.uint(0)
	err := xedError(C.xed_encode(
		&req,
		dstbuf.CBytes(),
		C.uint(dstbuf.Cap()),
		&codeLen,
	))
	if !err.Empty() {
		return 0, err
	}

	return int(codeLen), nil
}

func xedMemOperand(req *EncodeRequest, bitSize int) C.xed_encoder_operand_t {
	var disp C.xed_enc_displacement_t
	disp.displacement = C.xed_uint64_t(req.ptr.Disp)
	switch req.dispWidth {
	case 8:
		disp.displacement_bits = 8
	case 32:
		disp.displacement_bits = 32
	default:
		if req.ptr.Disp == 0 {
			disp.displacement_bits = 0
		} else if req.ptr.Disp >= -128 && req.ptr.Disp <= 127 {
			disp.displacement_bits = 8
		} else {
			disp.displacement_bits = 32
		}
	}
	base := registerByName[req.ptr.Base]
	index := registerByName[req.ptr.Index]
	return C.xed_mem_bisd(
		C.xed_reg_enum_t(base),
		C.xed_reg_enum_t(index),
		C.xed_uint_t(req.ptr.Scale),
		disp,
		C.xed_uint_t(bitSize))
}

func xedOperand(req *EncodeRequest, index int) C.xed_encoder_operand_t {
	switch req.tags[index] {
	case argUint8:
		return C.xed_imm0(C.xed_uint64_t(req.imm), 8)
	case argUint32:
		return C.xed_imm0(C.xed_uint64_t(req.imm), 32)
	case argInt8:
		return C.xed_simm0(C.xed_int32_t(req.imm), 8)
	case argInt16:
		return C.xed_simm0(C.xed_int32_t(req.imm), 16)
	case argInt32:
		return C.xed_simm0(C.xed_int32_t(req.imm), 32)
	case argMem:
		return xedMemOperand(req, int(req.memWidth))
	case argRel8:
		return C.xed_relbr(C.xed_int32_t(req.rel), 8)
	case argRel16:
		return C.xed_relbr(C.xed_int32_t(req.rel), 16)
	case argRel32:
		return C.xed_relbr(C.xed_int32_t(req.rel), 32)

	default:
		return C.xed_reg(C.xed_reg_enum_t(req.regs[index]))
	}
}

const (
	// Should be big enough to hold:
	// - longest ICLASS.
	// - longest encoding string (XED_MAX_INSTRUCTION_BYTES).
	bufferCapacity = 48
)

type buffer struct {
	data [bufferCapacity]byte
}

func (b *buffer) Cap() int { return bufferCapacity }

func (b *buffer) SetCString(s string) {
	copy(b.data[:], s)
	b.data[len(s)] = 0 // '\0' terminator
}

func (b *buffer) CString() *C.char {
	return (*C.char)(unsafe.Pointer(&b.data[0]))
}

func (b *buffer) CBytes() *C.uint8_t {
	return (*C.uint8_t)(unsafe.Pointer(&b.data[0]))
}

func (b *buffer) GoBytes(length int) []byte {
	goBytes := make([]byte, length)
	copy(goBytes, b.data[:])
	return goBytes
}
