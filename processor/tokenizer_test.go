package processor

import (
	"testing"
	"strings"

	"mage/typedefs"
)

func Test_Tokenizer(t *testing.T) {
	lines, _ := loadFile("../fixtures/macro.mg")
	tkn := newTokenizer("macro.mg", strings.Join(lines, "\n"))
	expected := 63
	tokens := tkn.tokenize()
	for _, tk := range tokens {
		t.Log(tk)
	}
	if expected != len(tokens) {
		t.Fatalf("expected %d tokens, but got %d", expected, len(tokens))
	}
	tokMacros := tkn.filter(typedefs.TOKEN_MACRO_DFN_OPEN)
	if 5 != len(tokMacros) {
		t.Fatalf("there should be 5 macros, not %d", len(tokMacros))
	}
	tokRules := tkn.filter(typedefs.TOKEN_RULE_OPEN)
	if 2 != len(tokRules) {
		t.Fatalf("there should be 2 rule dfns, not %d", len(tokRules))
	}
}

func Test_TokenizerSetsProperPositions_MacroDfn(t *testing.T) {
	lines, _ := loadFile("../fixtures/macro.mg")
	tkn := newTokenizer("macro.mg", strings.Join(lines, "\n"))
	tkn.tokenize()
	for _,tk := range tkn.filter(typedefs.TOKEN_MACRO_DFN_OPEN) {
		if tk.pos.Char != 1 {
			t.Fatalf("macro dfn should start at 1, got %d", tk.pos.Char)
		}
	}
	prev := 0
	for _,tk := range tkn.filter(typedefs.TOKEN_MACRO_DFN_CLOSE) {
		if tk.pos.Line <= prev {
			t.Fatalf("macro dfn should end after %d, got %d", prev, tk.pos.Line)
		}
		prev += 1
	}
}

func Test_TokenizerSetsProperPositions_Command(t *testing.T) {
	lines, _ := loadFile("../fixtures/macro.mg")
	tkn := newTokenizer("macro.mg", strings.Join(lines, "\n"))
	tkn.tokenize()
	for _, tk := range tkn.filter(typedefs.TOKEN_COMMAND_OPEN) {
		if tk.pos.Char != 1 {
			t.Fatalf("command should start at 1, got %d", tk.pos.Char)
		}
	}
	prev := 0
	for _,tk := range tkn.filter(typedefs.TOKEN_COMMAND_CLOSE) {
		if tk.pos.Line <= prev {
			t.Fatalf("command should end after %d, got %d", prev, tk.pos.Line)
		}
		prev += 1
	}
}

func Test_TokenizerPosition(t *testing.T) {
	pos := tokenizerPosition{"test", 0, 0, 0 }
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

// func Test_Transform(t *testing.T) {
// 	lines, _ := loadFile("../fixtures/macro.mg")
// 	tokens := tokenize("test", strings.Join(lines, "\n"))
// 	transformedTokens := transform(tokens)
// 	for _, tk := range transformedTokens {
// 		t.Log(tk)
// 	}
// }
