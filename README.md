# XEDq

**XEDq** brings [Intel XED](https://github.com/intelxed/xed) powers into [Go](https://golang.org/) space.

## Warning: not even alpha

This project is under active development and should not be used anywhere.  
When the time is right, this warning will go away.

## What it is

Convenient library built on top of **Intel XED** for encoding/decoding of 
[X86](https://ru.wikipedia.org/wiki/X86) instructions.

**XEDq** is a good choice for building command line utility that
works with X86 assembly or machine code.

Example use cases:
- External assembler/disassembler validation:
  - Tests data/input/output generation.
  - End-to-end analysis.
- Framework for X86 assembler/disassembler written in Go.

Main focus is around usability and convenience, rather than performance:
- Fluent interface.
- Integrated memory expression parser.
- Automatical resolving of some input parameters (can be disabled).

## Examples

Simple usage:

```go
encoder := xedq.NewEncoder()

add := encoder.Request("ADD").Reg("EAX").MemExpr("DWORD PTR [EDX+ECX*4]")
fmt.Println(add.EncodeHexString()) // => "6703048a"
fmt.Println(add.Encode())          // => [103 3 4 138]
fmt.Println(add.String())          // => "ADD EAX, DWORD PTR [EDX+ECX*4]"

// AVX512 instruction.
vaddpd := encoder.Request("VADDPD").Reg("XMM0").Reg("K4").Reg("XMM10").Reg("XMM20")
fmt.Println(vaddpd.EncodeHexString()) // => "62b1ad0c58c4"
fmt.Println(vaddpd.Encode())          // => [98 177 173 12 88 196]
fmt.Println(vaddpd.String())          // => "VADDPD XMM0, K4, XMM10, XMM20"
```

For more examples, see [encoder tests](src/xedq/encoder_test.go).
