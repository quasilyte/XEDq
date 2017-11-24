package xedq

import (
	"testing"
)

func TestParseSIB(t *testing.T) {
	var tmpbuf buffer

	expressions := []string{
		"RAX",
		"RCX",
		"R8D",
		"RDX+RAX",
		"RDX+RCX",
		"RAX*2",
		"R10D*8",
		"RDX+RAX*2",
		"RDX+RCX*4",
		"RAX+XMM0*8",
		"RBP+XMM9*8",
		"R8+YMM20",
		"R10+ZMM10*2",
	}

	for _, expr := range expressions {
		have := sibString(parseSIBExpr(expr, &tmpbuf))
		want := "[" + expr + "]"
		if have != want {
			t.Errorf("parse(%s):\nhave: %v\nwant: %v",
				expr, have, want)
		}
	}
}
