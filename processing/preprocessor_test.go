package processing

import (
	"strings"
	"testing"
)

func Test_PreprocessorMacros(t *testing.T) {
	lines, _ := loadFile("../fixtures/macro.mg")
	tkn := newTokenizer("macro.mg", strings.Join(lines, "\n"))
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
