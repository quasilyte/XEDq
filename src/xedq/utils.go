package xedq

import (
	"fmt"
	"strings"
)

// Functions that are hard to categorize or associate with any
// other source file.

// registerCache holds results of registerByName invocations
// to speedup register name resolution.
// Only valid register names are stored.
var registerCache = make(map[string]register, 256)

// registerByName maps regName to associated register object.
// If regName is not recognized, xedRegInvalid is returned.
func registerByName(regName string, tmpbuf *buffer) register {
	reg, ok := registerCache[regName]
	if ok {
		return reg
	}
	reg = newXEDRegister(regName, tmpbuf)
	if reg != xedRegInvalid {
		registerCache[regName] = reg
	}
	return reg
}

var memScaleMap = map[string]uint8{
	"1": 1,
	"2": 2,
	"4": 4,
	"8": 8,
}

func sibString(base, index register, scale uint8) string {
	switch {
	case index == 0 && scale == 0:
		return fmt.Sprintf("[%s]", base)
	case base == 0 && index != 0 && scale != 0:
		return fmt.Sprintf("[%s*%d]", index, scale)
	case index != 0 && scale == 0:
		return fmt.Sprintf("[%s+%s]", base, index)
	default:
		return fmt.Sprintf("[%s+%s*%d]", base, index, scale)
	}
}

func parseSIBExpr(sibExpr string, tmpbuf *buffer) (base, index register, scale uint8) {
	indexPos := strings.IndexByte(sibExpr, '+')
	scalePos := strings.IndexByte(sibExpr, '*')

	switch {
	case indexPos == -1 && scalePos == -1:
		// [base].
		base = registerByName(sibExpr, tmpbuf)
	case indexPos == -1 && scalePos != -1:
		// [index*scale].
		index = registerByName(sibExpr[:scalePos], tmpbuf)
		scale = memScaleMap[sibExpr[scalePos+1:]]
	case indexPos != -1 && scalePos == -1:
		// [base+index].
		base = registerByName(sibExpr[:indexPos], tmpbuf)
		index = registerByName(sibExpr[indexPos+1:], tmpbuf)
	default:
		// [base+index*scale].
		base = registerByName(sibExpr[:indexPos], tmpbuf)
		index = registerByName(sibExpr[indexPos+1:scalePos], tmpbuf)
		scale = memScaleMap[sibExpr[scalePos+1:]]
	}

	return base, index, scale

}
