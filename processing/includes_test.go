package processing

import (
	"mage/shell"
	"testing"
)

func Test_TokenizeIncludes(t *testing.T) {
	lines, _ := shell.LoadFile("../fixtures/includes.mg")
	tkn := newTokenizer("../fixtures/includes.mg", lines)
	rawTokens := tkn.tokenize()
	// debug.Tokens(rawTokens)
	if len(rawTokens) <= 0 {
		t.Fatalf("should at least have some raw tokens")
	}
}

func Test_PreprocessIncludes(t *testing.T) {
	filepath := "../fixtures/includes.mg"
	lines, _ := shell.LoadFile(filepath)
	tkn := newTokenizer(filepath, lines)
	rawTokens := tkn.tokenize()
	tokens, err := preprocessIncludes(rawTokens)
	if err != nil {
		t.Fatalf("preprocessing includes error: %s", err)
	}
	// debug.Tokens(tokens)
	if 128 != len(tokens) {
		t.Fatalf("expected exactly 128 tokens with includes, got %d", len(tokens))
	}
}

func Test_PreprocessorDoesIncludes(t *testing.T) {
	filepath := "../fixtures/includes.mg"
	lines, _ := shell.LoadFile(filepath)
	tkn := newTokenizer(filepath, lines)
	proc := NewPreprocessor(tkn.tokenize())
	err := proc.doIncludes()
	if err != nil {
		t.Fatalf("preprocessing includes error: %s", err)
	}
	// debug.Tokens(tokens)
	if 128 != len(proc.tokens) {
		t.Fatalf("expected exactly 128 tokens with includes, got %d", len(proc.tokens))
	}
}

func Test_ApplyIncludes(t *testing.T) {
	filepath := "../fixtures/includes.mg"
	lines, _ := shell.LoadFile(filepath)
	tkn := newTokenizer(filepath, lines)
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

	// debug.Tokens(tokens)
	if 102 != len(tokens) {
		t.Fatalf("expected exactly 102 tokens with includes, got %d", len(tokens))
	}
}

func Test_RecursiveInclusionShouldErrorOut(t *testing.T) {
	filepath := "../fixtures/recursive-inclusion.mg"
	lines, _ := shell.LoadFile(filepath)
	tkn := newTokenizer(filepath, lines)
	_, err := preprocessIncludes(tkn.tokenize())
	if err == nil {
		t.Log(err)
		t.Fatalf("expected to fail after too many inclusions")
	}
}

func Test_MultilevelInclusionShouldWork(t *testing.T) {
	filepath := "../fixtures/double-include-parent.mg"
	lines, _ := shell.LoadFile(filepath)
	tkn := newTokenizer(filepath, lines)
	tokens, err := preprocessIncludes(tkn.tokenize())
	if err != nil {
		t.Log(err)
		t.Fatalf("multi-level include should work")
	}
	if len(tokens) != 87 {
		t.Fatalf("expected exactly 87 tokens in multi-level include, got %d", len(tokens))
	}
}
