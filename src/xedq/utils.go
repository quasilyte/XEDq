package xedq

import (
	"bytes"
	"fmt"
)

// Functions that are hard to categorize or associate with any
// other source file.

func memExprString(req *EncodeRequest) string {
	var buf bytes.Buffer

	base := req.ptr.Base
	index := req.ptr.Index
	scale := req.ptr.Scale
	disp := req.ptr.Disp

	fmt.Fprintf(&buf, "mem%d", req.memWidth)

	buf.WriteByte('[')
	switch {
	case index == "" && scale == 0:
		fmt.Fprintf(&buf, "%s", base)
	case base == "" && index != "" && scale != 0:
		fmt.Fprintf(&buf, "%s*%d", index, scale)
	case index != "" && scale == 0:
		fmt.Fprintf(&buf, "%s+%s", base, index)
	default:
		fmt.Fprintf(&buf, "%s+%s*%d", base, index, scale)
	}
	switch {
	case disp > 0:
		fmt.Fprintf(&buf, "+%#x", disp)
	case disp < 0:
		fmt.Fprintf(&buf, "%#x", disp)
	}
	buf.WriteByte(']')

	return buf.String()
}
