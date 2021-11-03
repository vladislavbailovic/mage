package processing

import (
	"testing"

	"mage/shell"
	"mage/typedefs"
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

func Test_Tokenizer(t *testing.T) {
	lines, _ := shell.LoadFile("../fixtures/macro.mg")
	tkn := newTokenizer("macro.mg", lines)
	expected := 63
	tokens := tkn.tokenize()
	// debug.Tokens(tokens)
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

func Test_WordPositions(t *testing.T) {
	lines, _ := shell.LoadFile("../fixtures/macro.mg")
	tkn := newTokenizer("macro.mg", lines)
	tokens := tkn.tokenize()
	expecteds := [][]int{
		//[]int{1,1}, // "macro" gets nerfed
		[]int{1, 8},
		[]int{1, 13},
		[]int{1, 22},
		[]int{1, 27},
		[]int{1, 33},
		[]int{1, 38},
		[]int{1, 41},
		// macro OTHER $(M3)
		//[]int{2,1}, // "macro" gets nerfed
		[]int{2, 8},
		[]int{2, 14},
		// macro M3 $(M4)
		//[]int{3,1}, // "macro" gets nerfed
		[]int{3, 8},
		[]int{3, 11},
		// macro M4 $(NAME)
		//[]int{4,1}, // "macro" gets nerfed
		[]int{4, 8},
		[]int{4, 11},
		// macro M5 $(NAME)
		//[]int{5,1}, // "macro" gets nerfed
		[]int{5, 8},
		[]int{5, 11},
		// root: tmp/whatever.test
		[]int{7, 1},
		[]int{7, 6},
		// echo $(NAME)
		[]int{8, 2}, // has tab
		[]int{8, 7},
		// tmp/whatever.test:
		[]int{10, 1},
		// echo nay nya $(OTHER)
		[]int{11, 2}, // has tab
		[]int{11, 7},
		[]int{11, 11},
		[]int{11, 15},
		// sed -e 's/$(M5)/nana/g'
		[]int{12, 2}, // has tab
		[]int{12, 6},
		[]int{12, 9},
	}
	current := 0
	for idx, tk := range tokens {
		if typedefs.TOKEN_WORD != tk.Kind {
			continue
		}
		e := expecteds[current]
		if tk.Pos.Line != e[0] {
			t.Fatalf("expected %s to be on line %d, but got %d", tk.Value, e[0], tk.Pos.Line)
		}
		if tk.Pos.Char != e[1] {
			t.Fatalf("expected [%s](%d) to start on char %d, but got %d", tk.Value, idx, e[1], tk.Pos.Char)
		}
		current++
		if current > len(expecteds)-1 {
			break
		}
	}
}

func Test_TokenizerSetsProperPositions_MacroDfn(t *testing.T) {
	lines, _ := shell.LoadFile("../fixtures/macro.mg")
	tkn := newTokenizer("macro.mg", lines)
	tkn.tokenize()
	for _, tk := range tkn.filter(typedefs.TOKEN_MACRO_DFN_OPEN) {
		if tk.Pos.Char != 8 {
			t.Fatalf("macro dfn should start at 7 (because keyword gets dropped), got %d", tk.Pos.Char)
		}
	}
	prev := 0
	for _, tk := range tkn.filter(typedefs.TOKEN_MACRO_DFN_CLOSE) {
		if tk.Pos.Line <= prev {
			t.Fatalf("macro dfn should end after %d, got %d", prev, tk.Pos.Line)
		}
		prev += 1
	}
}

func Test_TokenizerSetsProperPositions_Command(t *testing.T) {
	lines, _ := shell.LoadFile("../fixtures/macro.mg")
	tkn := newTokenizer("macro.mg", lines)
	tkn.tokenize()
	for _, tk := range tkn.filter(typedefs.TOKEN_COMMAND_OPEN) {
		if tk.Pos.Char != 2 {
			t.Fatalf("command should start at 2 (after tab), got %d", tk.Pos.Char)
		}
	}
	prev := 0
	for _, tk := range tkn.filter(typedefs.TOKEN_COMMAND_CLOSE) {
		if tk.Pos.Line <= prev {
			t.Fatalf("command should end after %d, got %d", prev, tk.Pos.Line)
		}
		prev += 1
	}
}

func Test_TokenizerSetsProperPositions_MacroCall(t *testing.T) {
	lines, _ := shell.LoadFile("../fixtures/macro.mg")
	tkn := newTokenizer("macro.mg", lines)
	tkn.tokenize()

	co := tkn.filter(typedefs.TOKEN_MACRO_CALL_OPEN)
	if co[0].Pos.Char != 14 {
		t.Fatalf("first call should start at c13, got c%d", co[0].Pos.Char)
	}
	if co[0].Pos.Line != 2 {
		t.Fatalf("first call should start at l2, got l%d", co[0].Pos.Line)
	}
	if co[1].Pos.Char != 11 {
		t.Fatalf("2nd call should start at c10, got c%d", co[1].Pos.Char)
	}
	if co[1].Pos.Line != 3 {
		t.Fatalf("2nd call should start at l3, got l%d", co[1].Pos.Line)
	}
	if co[5].Pos.Line != 11 {
		t.Fatalf("5th call should start at l11, got l%d", co[5].Pos.Line)
	}
	if co[5].Pos.Char != 15 {
		t.Fatalf("5th call should start at c15, got c%d", co[5].Pos.Char)
	}
}
