package processor

import (
	"testing"
	"strings"
)

func Test_Preprocessor(t *testing.T) {
	lines, _ := loadFile("../fixtures/macro.mg")
	tkn := newTokenizer("macro.mg", strings.Join(lines, "\n"))
	rawTokens := tkn.tokenize()
	tokens, err := preprocess(rawTokens)
	if err != nil {
		t.Fatalf("preprocessing error: %s", err)
	}
	for _,tk := range tokens {
		t.Logf("> %s (%s)\n", tk.Value, toktype(tk.Kind))
	}
	if len(tokens) >= len(rawTokens) {
		t.Fatalf("expected macros to be removed from tokens, got %d", len(tokens))
	}
}
