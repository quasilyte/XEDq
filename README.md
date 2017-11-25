# XEDq

**XEDq** brings [Intel XED](https://github.com/intelxed/xed) powers into [Go](https://golang.org/) space.

## Warning: not even alpha

This project is under active development and should not be used anywhere.  
When the time is right, this warning will go away.

## What it is

Convenient library built on top of **Intel XED** for encoding/decoding of 
[X86](https://ru.wikipedia.org/wiki/X86) instructions.

Focused around two use cases:
1. Manually-coded expressions: succinct and readable API, easy to use.
2. Programmatically built expressions: composable API, easy to extend.

## Examples

Simple usage:

```go
encoder := xedq.NewEncoder()

add := encoder.Request("ADD").Reg("EAX").Mem32("EDX+ECX*4")
fmt.Println(add.EncodeHexString()) // => "6703048a"
fmt.Println(add.Encode())          // => [103 3 4 138]
fmt.Println(add.String())          // => "ADD EAX, mem32[EDX+ECX*4]"

// AVX512 instruction.
vaddpd := encoder.Request("VADDPD").Reg("XMM0").Reg("K4").Reg("XMM10").Reg("XMM20")
fmt.Println(vaddpd.EncodeHexString()) // => "62b1ad0c58c4"
fmt.Println(vaddpd.Encode())          // => [98 177 173 12 88 196]
fmt.Println(vaddpd.String())          // => "VADDPD XMM0, K4, XMM10, XMM20"
```

For more examples, see [encoder tests](src/xedq/encoder_test.go).
