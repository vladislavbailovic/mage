package tokenizing

import (
	"mage/shell"
	"testing"
)

func Test_TokenizeIncludes(t *testing.T) {
	lines, _ := shell.LoadFile("../../fixtures/includes.mg")
	tkn := NewTokenizer("../../fixtures/includes.mg", lines)
	rawTokens := tkn.Tokenize()
	// debug.Tokens(rawTokens)
	if len(rawTokens) <= 0 {
		t.Fatalf("should at least have some raw tokens")
	}
}
