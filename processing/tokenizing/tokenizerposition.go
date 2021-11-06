package tokenizing

import (
	"mage/typedefs"
)

type tokenizerPosition struct {
	source      string
	cursor      int
	currentLine int
	currentChar int
}

func (tp *tokenizerPosition) advanceCursor(chr int) {
	tp.cursor += chr
}
func (tp *tokenizerPosition) advanceChar(chr int) {
	tp.currentChar += chr
}
func (tp *tokenizerPosition) advanceLine(chr int) {
	tp.currentLine += chr
	tp.currentChar = 1
}
func (tp *tokenizerPosition) advance(chr int) {
	tp.advanceCursor(chr)
	tp.advanceChar(chr)
}
func (tp tokenizerPosition) getPosition() typedefs.SourcePosition {
	return typedefs.SourcePosition{
		tp.source,
		tp.currentLine,
		tp.currentChar,
	}
}
