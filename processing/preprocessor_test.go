package processing

import (
	"mage/shell"
	"testing"
)

func Test_PreprocessorMacros(t *testing.T) {
	lines, _ := shell.LoadFile("../fixtures/macro.mg")
	tkn := newTokenizer("macro.mg", lines)
	rawTokens := tkn.tokenize()
	tokens, err := preprocess(rawTokens)
	if err != nil {
		t.Fatalf("preprocessing error: %s", err)
	}
	// debugTokens(tokens)
	if len(tokens) >= len(rawTokens) {
		t.Fatalf("expected macros to be removed from tokens, got %d", len(tokens))
	}
}
