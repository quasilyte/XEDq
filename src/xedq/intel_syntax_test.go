package xedq

import (
	"testing"
)

func TestIntelMemExprParse(t *testing.T) {
	tests := map[string]Ptr{
		// Base only.
		"RCX":        {"RCX", "", 0, 0},
		"R8":         {"R8", "", 0, 0},
		"RCX+0":      {"RCX", "", 0, 0},
		"RCX+1750":   {"RCX", "", 0, 1750},
		"RCX+0x0":    {"RCX", "", 0, 0},
		"RCX+0xf0":   {"RCX", "", 0, 0xf0},
		"RCX-0x01fa": {"RCX", "", 0, -0x1fa},
		// Index*Scale.
		"RAX*2":        {"", "RAX", 2, 0},
		"R9*8":         {"", "R9", 8, 0},
		"RAX*2+0":      {"", "RAX", 2, 0},
		"RAX*2+1750":   {"", "RAX", 2, 1750},
		"RAX*2+0x0":    {"", "RAX", 2, 0},
		"RAX*2+0xf0":   {"", "RAX", 2, 0xf0},
		"RAX*2-0x01fa": {"", "RAX", 2, -0x1fa},
		// Base+Index.
		"RAX+RCX":        {"RAX", "RCX", 0, 0},
		"R9+R8":          {"R9", "R8", 0, 0},
		"RAX+RCX+0":      {"RAX", "RCX", 0, 0},
		"RAX+RCX+1750":   {"RAX", "RCX", 0, 1750},
		"RAX+RCX+0x0":    {"RAX", "RCX", 0, 0},
		"RAX+RCX+0xf0":   {"RAX", "RCX", 0, 0xf0},
		"RAX+RCX-0x01fa": {"RAX", "RCX", 0, -0x1fa},
		// Base+Index*Scale.
		"RAX+RCX*4":        {"RAX", "RCX", 4, 0},
		"R9+R8*1":          {"R9", "R8", 1, 0},
		"RAX+RCX*4+0":      {"RAX", "RCX", 4, 0},
		"RAX+RCX*4+1750":   {"RAX", "RCX", 4, 1750},
		"RAX+RCX*4+0x0":    {"RAX", "RCX", 4, 0},
		"RAX+RCX*4+0xf0":   {"RAX", "RCX", 4, 0xf0},
		"RAX+RCX*4-0x01fa": {"RAX", "RCX", 4, -0x1fa},
	}

	for expr, want := range tests {
		have, err := IntelMemExprParse(expr)
		if err != nil {
			t.Errorf("Parse(%q): error:\n%v", expr, err)
		}
		if have != want {
			t.Errorf("Parse(%q): output mismatch:\nhave: %#v\nwant: %#v",
				expr, have, want)
		}
	}
}
