package tokenizing

import (
	"mage/shell"
	"testing"
)

func Test_PreprocessorShellcallMacros(t *testing.T) {
	lines, _ := shell.LoadFile("../../fixtures/shell-macros.mg")
	tkn := NewTokenizer("../../fixtures/shell-macros.mg", lines)
	rawTokens := tkn.Tokenize()
	// debug.Tokens(tokens)
	if len(rawTokens) <= 0 {
		t.Fatalf("should at least have some raw tokens")
	}
}
