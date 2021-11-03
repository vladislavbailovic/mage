package processor

import (
	"testing"
	"strings"
)

func Test_Preprocessor(t *testing.T) {
	lines, _ := loadFile("../fixtures/macro.mg")
	tkn := newTokenizer("macro.mg", strings.Join(lines, "\n"))
	tokens, err := preprocess(tkn.tokenize())
	if err != nil {
		t.Fatalf("preprocessing error: %s", err)
	}
	for _,tk := range tokens {
		t.Logf("> %s (%s)\n", tk.value, toktype(tk.kind))
	}
}

