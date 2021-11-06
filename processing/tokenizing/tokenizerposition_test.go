package tokenizing

import (
	"testing"
)

func Test_TokenizerPosition(t *testing.T) {
	pos := tokenizerPosition{"test", 0, 0, 0}
	pos.advance(161)
	if pos.currentChar != 161 {
		t.Fatalf("advance should move char to 161, got %d", pos.currentChar)
	}
	if pos.cursor != 161 {
		t.Fatalf("advance should move cursor to 161, got %d", pos.cursor)
	}
	if pos.currentLine != 0 {
		t.Fatalf("advance should NOT move line at all, got %d", pos.currentLine)
	}

	pos.advanceLine(1)
	if pos.currentLine != 1 {
		t.Fatalf("advance should move line to 1, got %d", pos.currentLine)
	}
	if pos.currentChar != 1 {
		t.Fatalf("advance line should move char to 1, got %d", pos.currentChar)
	}

	sp := pos.getPosition()
	if sp.File != pos.source {
		t.Fatalf("source mismatch, expected %s, got %s", pos.source, sp.File)
	}
	if sp.Line != pos.currentLine {
		t.Fatalf("line mismatch, expected %d, got %d", pos.currentLine, sp.Line)
	}
	if sp.Char != pos.currentChar {
		t.Fatalf("char mismatch, expected %d, got %d", pos.currentChar, sp.Char)
	}
}
