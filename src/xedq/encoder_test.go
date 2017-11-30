package xedq

import (
	"testing"
)

func init() {
	InitTables()
}

// TODO:
// - test fixed disp width effect.
// - test vmem (vsib) instructions.

func TestEncoderMode32Eosz8(t *testing.T) {
	encoder := NewEncoder(EncoderMode32)

	req := func(name string) *EncodeRequest {
		return encoder.Request(name).SetEosz8()
	}

	runEncoderTests(t, map[string][]*EncodeRequest{
		"0433":   {req("ADD").Reg("AL").Uint8(0x33)},
		"88c8":   {req("MOV").Reg("AL").Reg("CL")},
		"88c1":   {req("MOV").Reg("CL").Reg("AL")},
		"30d3":   {req("XOR").Reg("BL").Reg("DL")},
		"30da":   {req("XOR").Reg("DL").Reg("BL")},
		"82f10f": {req("XOR").Reg("CL").Uint8(0x0f)},
		"340f":   {req("XOR").Reg("AL").Uint8(0x0f)},
		"30c0":   {req("XOR").Reg("AL").Reg("AL")},
	})
}

func TestEncoderMode32Eosz16(t *testing.T) {
	encoder := NewEncoder(EncoderMode32)

	req := func(name string) *EncodeRequest {
		return encoder.Request(name).SetEosz16()
	}

	runEncoderTests(t, map[string][]*EncodeRequest{
		"6683c033": {req("ADD").Reg("AX").Uint8(0x33)},
		"6689c8":   {req("MOV").Reg("AX").Reg("CX")},
		"6689c1":   {req("MOV").Reg("CX").Reg("AX")},
		"6631d3":   {req("XOR").Reg("BX").Reg("DX")},
		"6631da":   {req("XOR").Reg("DX").Reg("BX")},
		"6683f10f": {req("XOR").Reg("CX").Uint8(0x0f)},
		"6683f00f": {req("XOR").Reg("AX").Uint8(0x0f)},
		"6631c0":   {req("XOR").Reg("AX").Reg("AX")},
	})
}

func TestEncoderMode32Eosz32(t *testing.T) {
	encoder := NewEncoder(EncoderMode32)

	req := func(name string) *EncodeRequest {
		return encoder.Request(name).SetEosz32()
	}

	runEncoderTests(t, map[string][]*EncodeRequest{
		"83c077":       {req("ADD").Reg("EAX").Uint8(0x77)},
		"0511223344":   {req("ADD").Reg("EAX").Uint32(0x44332211)},
		"89c8":         {req("MOV").Reg("EAX").Reg("ECX")},
		"89c1":         {req("MOV").Reg("ECX").Reg("EAX")},
		"31d3":         {req("XOR").Reg("EBX").Reg("EDX")},
		"31da":         {req("XOR").Reg("EDX").Reg("EBX")},
		"83f10f":       {req("XOR").Reg("ECX").Uint8(0x0f)},
		"81f1f0f00000": {req("XOR").Reg("ECX").Uint32(0xf0f0)},
		"83f00f":       {req("XOR").Reg("EAX").Uint8(0x0f)},
		"35f0f00000":   {req("XOR").Reg("EAX").Uint32(0xf0f0)},
		"31c0":         {req("XOR").Reg("EAX").Reg("EAX")},
	})
}

func TestEncoderMode64Eosz8(t *testing.T) {
	encoder := NewEncoder(EncoderMode64)

	req := func(name string) *EncodeRequest {
		return encoder.Request(name).SetEosz8()
	}

	runEncoderTests(t, map[string][]*EncodeRequest{
		"0433":   {req("ADD").Reg("AL").Uint8(0x33)},
		"88c8":   {req("MOV").Reg("AL").Reg("CL")},
		"88c1":   {req("MOV").Reg("CL").Reg("AL")},
		"30d3":   {req("XOR").Reg("BL").Reg("DL")},
		"30da":   {req("XOR").Reg("DL").Reg("BL")},
		"80f10f": {req("XOR").Reg("CL").Uint8(0x0f)},
		"340f":   {req("XOR").Reg("AL").Uint8(0x0f)},
		"30c0":   {req("XOR").Reg("AL").Reg("AL")},
	})
}

func TestEncoderMode64Eosz16(t *testing.T) {
	encoder := NewEncoder(EncoderMode64)

	req := func(name string) *EncodeRequest {
		return encoder.Request(name).SetEosz16()
	}

	runEncoderTests(t, map[string][]*EncodeRequest{
		"6683c033":     {req("ADD").Reg("AX").Uint8(0x33)},
		"6689c8":       {req("MOV").Reg("AX").Reg("CX")},
		"6689c1":       {req("MOV").Reg("CX").Reg("AX")},
		"6631d3":       {req("XOR").Reg("BX").Reg("DX")},
		"6631da":       {req("XOR").Reg("DX").Reg("BX")},
		"6683f10f":     {req("XOR").Reg("CX").Uint8(0x0f)},
		"6683f00f":     {req("XOR").Reg("AX").Uint8(0x0f)},
		"6631c0":       {req("XOR").Reg("AX").Reg("AX")},
		"666781003333": {req("ADD").MemExpr(16, "EAX").Int16(0x3333)},
		"66678d00":     {req("LEA").Reg("AX").MemExpr(64, "EAX")},
		"66678d400f":   {req("LEA").Reg("AX").MemExpr(64, "EAX+0x0f")},
		"66678d40f1":   {req("LEA").Reg("AX").MemExpr(64, "EAX-0x0f")},
	})
}

func TestEncoderMode64Eosz32(t *testing.T) {
	encoder := NewEncoder(EncoderMode64)

	req := func(name string) *EncodeRequest {
		return encoder.Request(name).SetEosz32()
	}

	runEncoderTests(t, map[string][]*EncodeRequest{
		"83c077":         {req("ADD").Reg("EAX").Uint8(0x77)},
		"0511223344":     {req("ADD").Reg("EAX").Uint32(0x44332211)},
		"89c8":           {req("MOV").Reg("EAX").Reg("ECX")},
		"89c1":           {req("MOV").Reg("ECX").Reg("EAX")},
		"4531d0":         {req("XOR").Reg("R8D").Reg("R10D")},
		"4531c2":         {req("XOR").Reg("R10D").Reg("R8D")},
		"4183f70f":       {req("XOR").Reg("R15D").Uint8(0x0f)},
		"4181f7f0f00000": {req("XOR").Reg("R15D").Uint32(0xf0f0)},
		"678b0491":       {req("MOV").Reg("EAX").MemExpr(32, "ECX+EDX*4")},
		"678b449144":     {req("MOV").Reg("EAX").MemExpr(32, "ECX+EDX*4+0x44")},
		"678d0488":       {req("LEA").Reg("EAX").MemExpr(32, "EAX+ECX*4")},
	})
}

func TestEncoderMode64Eosz64(t *testing.T) {
	encoder := NewEncoder(EncoderMode64)

	req := func(name string) *EncodeRequest {
		return encoder.Request(name).SetEosz64()
	}

	runEncoderTests(t, map[string][]*EncodeRequest{
		"4883c077":       {req("ADD").Reg("RAX").Uint8(0x77)},
		"480511223344":   {req("ADD").Reg("RAX").Uint32(0x44332211)},
		"4889c8":         {req("MOV").Reg("RAX").Reg("RCX")},
		"4889c1":         {req("MOV").Reg("RCX").Reg("RAX")},
		"4d31d0":         {req("XOR").Reg("R8").Reg("R10")},
		"4d31c2":         {req("XOR").Reg("R10").Reg("R8")},
		"4983f70f":       {req("XOR").Reg("R15").Uint8(0x0f)},
		"4981f7f0f00000": {req("XOR").Reg("R15").Uint32(0xf0f0)},
		"488b04c8":       {req("MOV").Reg("RAX").MemExpr(64, "RAX+RCX*8")},
		"488b44c844":     {req("MOV").Reg("RAX").MemExpr(64, "RAX+RCX*8+0x44")},
		"488d0409":       {req("LEA").Reg("RAX").MemExpr(64, "RCX+RCX")},
	})
}

func runEncoderTests(t *testing.T, tests map[string][]*EncodeRequest) {
	for encoding, requests := range tests {
		for _, req := range requests {
			have := req.EncodeHexString()
			err := req.encoder.LastError()
			if err != nil {
				t.Errorf("%q encoding error:\n%s\n%s",
					encoding, req, err.Error())
				continue
			}

			want := encoding
			if have != want {
				t.Errorf("encoding mismatch:\n%s\nhave: %q\nwant: %q",
					req, have, want)
			}
		}
	}
}
