package processing

import (
	"path"
	"strings"
	"testing"
)

func Test_TokenizeIncludes(t *testing.T) {
	lines, _ := loadFile("../fixtures/includes.mg")
	tkn := newTokenizer("macro.mg", strings.Join(lines, "\n"))
	rawTokens := tkn.tokenize()
	// debugTokens(rawTokens)
	if len(rawTokens) <= 0 {
		t.Fatalf("should at least have some raw tokens")
	}
}

func Test_PreprocessIncludes(t *testing.T) {
	filepath := "../fixtures/includes.mg"
	lines, _ := loadFile(filepath)
	tkn := newTokenizer(path.Base(filepath), strings.Join(lines, "\n"))
	rawTokens := tkn.tokenize()
	tokens, err := preprocessIncludes(rawTokens)
	if err != nil {
		t.Fatalf("preprocessing includes error: %s", err)
	}
	// debugTokens(tokens)
	if 128 != len(tokens) {
		t.Fatalf("expected exactly 128 tokens with includes, got %d", len(tokens))
	}
}

func Test_ApplyIncludes(t *testing.T) {
	filepath := "../fixtures/includes.mg"
	lines, _ := loadFile(filepath)
	tkn := newTokenizer(path.Base(filepath), strings.Join(lines, "\n"))
	rawTokens, err := preprocessIncludes(tkn.tokenize())
	if err != nil {
		t.Log(err)
		t.Fatalf("includes preprocessing should have been a success")
	}

	tokens, err := preprocessMacros(rawTokens)
	if err != nil {
		t.Log(err)
		t.Fatalf("macros expansion failed")
	}

	// debugTokens(tokens)
	if 102 != len(tokens) {
		t.Fatalf("expected exactly 102 tokens with includes, got %d", len(tokens))
	}
}
