package preprocessing

import (
	"mage/processing/tokenizing"
	"mage/shell"
	"testing"
)

func Test_PreprocessorShellcallMacros(t *testing.T) {
	lines, _ := shell.LoadFile("../../fixtures/shell-macros.mg")
	tkn := tokenizing.NewTokenizer("../../fixtures/shell-macros.mg", lines)
	rawTokens := tkn.Tokenize()
	tokens, err := Preprocess(rawTokens)
	if err != nil {
		t.Fatalf("preprocessing error: %s", err)
	}
	t.Log(err, tokens)
	// debug.Tokens(tokens)
	// if len(tokens) >= len(rawTokens) {
	// 	t.Fatalf("expected macros to be removed from tokens, got %d", len(tokens))
	// }
}
