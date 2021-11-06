package preprocessing

import (
	"mage/processing/tokenizing"
	"mage/shell"
	"testing"
)

func Test_PreprocessorMacros(t *testing.T) {
	lines, _ := shell.LoadFile("../../fixtures/macro.mg")
	tkn := tokenizing.NewTokenizer("../../fixtures/macro.mg", lines)
	rawTokens := tkn.Tokenize()
	tokens, err := Preprocess(rawTokens)
	if err != nil {
		t.Fatalf("preprocessing error: %s", err)
	}
	// debug.Tokens(tokens)
	if len(tokens) >= len(rawTokens) {
		t.Fatalf("expected macros to be removed from tokens, got %d", len(tokens))
	}
}
