package processor

import (
	"testing"
	"strings"

	"mage/typedefs"
)

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

func Test_WordPositions(t *testing.T) {
	lines, _ := loadFile("../fixtures/macro.mg")
	tkn := newTokenizer("macro.mg", strings.Join(lines, "\n"))
	tokens := tkn.tokenize()
	expecteds := [][]int{
		//[]int{1,1}, // "macro" gets nerfed
		[]int{1,7},
		[]int{1,12},
		[]int{1,21},
		[]int{1,26},
		[]int{1,32},
		[]int{1,37},
		[]int{1,40},
	// macro OTHER $(M3)
		//[]int{2,1}, // "macro" gets nerfed
		[]int{2,7},
		[]int{2,13},
	// macro M3 $(M4)
		//[]int{3,1}, // "macro" gets nerfed
		[]int{3,7},
		[]int{3,10},
	// macro M4 $(NAME)
		//[]int{4,1}, // "macro" gets nerfed
		[]int{4,7},
		[]int{4,10},
	// macro M5 $(NAME)
		//[]int{5,1}, // "macro" gets nerfed
		[]int{5,7},
		[]int{5,10},
	// root: tmp/whatever.test
		[]int{7,1},
		[]int{7,6},
	// echo $(NAME)
		[]int{8,2}, // has tab
		[]int{8,7},
	// tmp/whatever.test:
		[]int{10,1},
	// echo nay nya $(OTHER)
		[]int{11,2}, // has tab
		[]int{11,7},
		[]int{11,11},
		[]int{11,15},
	// sed -e 's/$(M5)/nana/g'
		[]int{12,2}, // has tab
		[]int{12,6},
		[]int{12,9},
	}
	current := 0
	for idx,tk := range tokens {
		if typedefs.TOKEN_WORD != tk.kind {
			continue
		}
		e := expecteds[current]
		if tk.pos.Line != e[0] {
			t.Fatalf("expected %s to be on line %d, but got %d", tk.value, e[0], tk.pos.Line)
		}
		if tk.pos.Char != e[1] {
			t.Fatalf("expected [%s](%d) to start on char %d, but got %d", tk.value, idx, e[1], tk.pos.Char)
		}
		current++
		if current > len(expecteds) - 1 {
			break
		}
	}
}

func Test_TokenizerSetsProperPositions_MacroDfn(t *testing.T) {
	lines, _ := loadFile("../fixtures/macro.mg")
	tkn := newTokenizer("macro.mg", strings.Join(lines, "\n"))
	tkn.tokenize()
	for _,tk := range tkn.filter(typedefs.TOKEN_MACRO_DFN_OPEN) {
		if tk.pos.Char != 7 {
			t.Fatalf("macro dfn should start at 7 (because keyword gets dropped), got %d", tk.pos.Char)
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
		if tk.pos.Char != 2 {
			t.Fatalf("command should start at 2 (after tab), got %d", tk.pos.Char)
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

func Test_TokenizerSetsProperPositions_MacroCall(t *testing.T) {
	lines, _ := loadFile("../fixtures/macro.mg")
	tkn := newTokenizer("macro.mg", strings.Join(lines, "\n"))
	tkn.tokenize()

	co := tkn.filter(typedefs.TOKEN_MACRO_CALL_OPEN)
	if co[0].pos.Char != 13 {
		t.Fatalf("first call should start at c13, got c%d", co[0].pos.Char)
	}
	if co[0].pos.Line != 2 {
		t.Fatalf("first call should start at l2, got l%d", co[0].pos.Line)
	}
	if co[1].pos.Char != 10 {
		t.Fatalf("2nd call should start at c10, got c%d", co[1].pos.Char)
	}
	if co[1].pos.Line != 3 {
		t.Fatalf("2nd call should start at l3, got l%d", co[1].pos.Line)
	}
	if co[5].pos.Line != 11 {
		t.Fatalf("5th call should start at l11, got l%d", co[5].pos.Line)
	}
	if co[5].pos.Char != 15 {
		t.Fatalf("5th call should start at c15, got c%d", co[5].pos.Char)
	}
}
