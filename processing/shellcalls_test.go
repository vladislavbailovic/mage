package processing

import (
	"mage/shell"
	"testing"
)

func Test_PreprocessorShellcallMacros(t *testing.T) {
	lines, _ := shell.LoadFile("../fixtures/shell-macros.mg")
	tkn := newTokenizer("../fixtures/shell-macros.mg", lines)
	rawTokens := tkn.tokenize()
	tokens, err := preprocess(rawTokens)
	if err != nil {
		t.Fatalf("preprocessing error: %s", err)
	}
	t.Log(err, tokens)
	// debug.Tokens(tokens)
	// if len(tokens) >= len(rawTokens) {
	// 	t.Fatalf("expected macros to be removed from tokens, got %d", len(tokens))
	// }
}
