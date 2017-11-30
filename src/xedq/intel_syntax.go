package xedq

import (
	"errors"
	"strconv"
	"strings"
)

var memScaleMap = map[string]uint8{
	"1": 1,
	"2": 2,
	"4": 4,
	"8": 8,
}

// IntelMemExprParse parses Intel-like syntax for memory operands.
// Implements MemExprParseFunc signature.
//
// expr can be in these forms:
//   "BASE"
//   "BASE±DISP"
//   "BASE+INDEX"
//   "BASE+INDEX±DISP"
//   "BASE+INDEX*SCALE"
//   "BASE+INDEX*SCALE±DISP"
//   "INDEX*SCALE"
//   "INDEX*SCALE±DISP"
// BASE and INDEX are register names.
// SCALE can be 1, 2, 4 or 8.
// DISP is integer in decimal or hex format.
// For hex, use "0x" prefix. Only lower case a-f letters are accepted.
// No whitespace is allowed.
func IntelMemExprParse(expr string) (Ptr, error) {
	var ptr Ptr

	signPos := -1
dispLoop:
	for i := len(expr) - 1; i >= 0; i-- {
		switch expr[i] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			continue
		case 'a', 'b', 'c', 'd', 'e', 'f', 'x':
			continue
		case '+', '-':
			signPos = i
			break dispLoop
		default:
			break dispLoop
		}
	}
	if signPos != -1 {
		dispBase := 10
		dispExpr := expr[signPos+1:]
		if len(dispExpr) >= 3 && dispExpr[0] == '0' && dispExpr[1] == 'x' {
			dispExpr = dispExpr[len("0x"):]
			dispBase = 16
		}
		disp, err := strconv.ParseInt(dispExpr, dispBase, 64)
		if err != nil {
			return ptr, errors.New("disp parse error: " + err.Error())
		}
		ptr.Disp = int32(disp)
		if expr[signPos] == '-' {
			ptr.Disp = -ptr.Disp
		}

		expr = expr[:signPos]
	}

	indexPos := strings.IndexByte(expr, '+')
	scalePos := strings.IndexByte(expr, '*')

	switch {
	case indexPos == -1 && scalePos == -1:
		// [base].
		ptr.Base = expr
	case indexPos == -1 && scalePos != -1:
		// [index*scale].
		ptr.Index = expr[:scalePos]
		ptr.Scale = memScaleMap[expr[scalePos+1:]]
	case indexPos != -1 && scalePos == -1:
		// [base+index].
		ptr.Base = expr[:indexPos]
		ptr.Index = expr[indexPos+1:]
	default:
		// [base+index*scale].
		ptr.Base = expr[:indexPos]
		ptr.Index = expr[indexPos+1 : scalePos]
		ptr.Scale = memScaleMap[expr[scalePos+1:]]
	}

	return ptr, nil
}
