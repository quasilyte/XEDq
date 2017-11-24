package xedq

import (
	"errors"
)

const (
	// Upper limit for instruction operands count.
	maxArgLimit = 6
)

var (
	errEncReqConvert = errors.New("encoder: request conversion failed")
	errEncArgc       = errors.New("encoder: invalid operands count")
	errNameLen       = errors.New("instruction name too long")
)

// InitTables prepares XED for encoding/decoding requests.
// TODO: specify what operations are safe/unsafe prior to this call.
func InitTables() {
	xedTablesInit()
}

// Encoder is x86 instructions assembler.
//
// Should be created with NewEncoder.
//
// Non thread-safe due to memory sharing with all spawned EncodeRequest's.
// Use Encoder.Copy to create a new instance of Encoder that have
// same settings as the original, but can be safely transfered to
// other goroutine.
type Encoder struct {
	tmpbuf buffer

	mode xedState
	err  error
}

// EncoderOption is a configuration function for NewEncoder.
type EncoderOption func(*Encoder)

// EncoderMode32 sets machine mode to 32bit.
// Predefined EncoderOption.
func EncoderMode32(enc *Encoder) { enc.mode = newXEDState32() }

// EncoderMode64 sets machine mode to 64bit.
// Predefined EncoderOption.
func EncoderMode64(enc *Encoder) { enc.mode = newXEDState64() }

// NewEncoder returns encoder that is configured by specified options.
//
// Default options:
//   + EncoderMode64
func NewEncoder(options ...EncoderOption) *Encoder {
	var enc Encoder

	// Set defaults.
	EncoderMode64(&enc)

	for i := range options {
		options[i](&enc)
	}

	return &enc
}

// LastError returns last executed encoding request error.
// Note that errors are not stacked up.
func (enc *Encoder) LastError() error {
	return enc.err
}

// Request creates new encoding request for instruction of specified name.
// See EncodingRequest.
func (enc *Encoder) Request(name string) *EncodeRequest {
	return &EncodeRequest{encoder: enc, name: name}
}

// encode assembles single encode request.
// On error, nil is returned and enc.err is set to associated error value.
// Returned slice is not shared.
func (enc *Encoder) encode(req *EncodeRequest) []byte {
	if len(req.name) > bufferCapacity {
		enc.err = errNameLen
		return nil
	}
	if req.argc < 0 || req.argc > maxArgLimit {
		enc.err = errEncArgc
		return nil
	}

	iclass := newXEDIclass(req.name, &enc.tmpbuf)
	inst := newXEDInst(&enc.mode, iclass, req)

	var code []byte
	code, enc.err = xedEncode(&inst, &enc.tmpbuf)
	return code
}
