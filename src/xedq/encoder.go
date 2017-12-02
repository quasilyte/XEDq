package xedq

import (
	"errors"
	"io"
)

const (
	// Upper limit for instruction operands count.
	maxArgLimit = 6
)

var (
	errEncReqConvert = errors.New("encoder: request conversion failed")
)

// MemExprParseFunc is a type of function that is used by Encoder
// to handle MemExpr arguments.
//
// Expr argument is some kind of effective address expression
// which syntax can vary.
type MemExprParseFunc func(expr string) (Ptr, error)

// Encoder is x86 instructions assembler.
//
// Should be created with NewEncoder.
// Should be copied with Encoder.Copy.
//
// Not thread-safe. Copy encoder, or share with mutex.
type Encoder struct {
	dispWidth int

	tmpbuf buffer

	mode xedState
	err  error

	MemExprParser MemExprParseFunc
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
	enc.mode = newXEDState64()
	enc.MemExprParser = IntelMemExprParse

	for i := range options {
		options[i](&enc)
	}

	return &enc
}

// Copy returns Encoder deep copy.
func (enc Encoder) Copy() Encoder {
	// Implementation may change in future.
	// Clients should not expect simple bitewise copy
	// to be a valid deep copy forever.
	return enc
}

// Err returns the last executed encoding request error.
func (enc *Encoder) Err() error {
	return enc.err
}

// SetDispWidth changes displacement encoding strategy.
//
// width values:
//   0  - use smallest displacement size
//   8  - 8bit displacement
//   32 - 32bit displacement
// All other values are treated as 0.
func (enc Encoder) SetDispWidth(width uint8) {
	switch width {
	case 8, 32:
		enc.dispWidth = int(width)
	default:
		enc.dispWidth = 0
	}
}

// Request creates new encoding request for instruction of specified name.
// See EncodingRequest.
func (enc *Encoder) Request(name string) *EncodeRequest {
	if len(name) > bufferCapacity {
		return &EncodeRequest{encoder: enc, iclass: xedIclassInvalid}
	}
	return &EncodeRequest{
		encoder: enc,
		iclass:  newXEDIclass(name, &enc.tmpbuf),
	}
}

// encode assembles req and returns result in freshly allocated slice of bytes.
func (enc *Encoder) encode(req *EncodeRequest) []byte {
	var n int
	inst := newXEDInst(&enc.mode, req)
	n, enc.err = xedEncode(&inst, &enc.tmpbuf)
	code := make([]byte, n)
	copy(code, enc.tmpbuf.data[:])
	return code
}

// encodeTo assembles req and writes result to w.
func (enc *Encoder) encodeTo(w io.Writer, req *EncodeRequest) (int, error) {
	var n int
	inst := newXEDInst(&enc.mode, req)
	n, enc.err = xedEncode(&inst, &enc.tmpbuf)
	return w.Write(enc.tmpbuf.data[:n])
}
