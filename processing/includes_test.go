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

func Test_PreprocessorAppliesIncludes(t *testing.T) {
	filepath := "../fixtures/includes.mg"
	lines, _ := shell.LoadFile(filepath)
	tkn := newTokenizer(filepath, lines)
	proc := NewPreprocessor(tkn.tokenize())
	err := proc.doIncludes()
	if err != nil {
		t.Log(err)
		t.Fatalf("includes preprocessing should have been a success")
	}

	err = proc.doMacros()
	if err != nil {
		t.Log(err)
		t.Fatalf("macros expansion failed")
	}

	// debug.Tokens(proc.tokens)
	if 102 != len(proc.tokens) {
		t.Fatalf("expected exactly 102 tokens with includes, got %d", len(proc.tokens))
	}
}

func Test_RecursiveInclusionShouldErrorOut(t *testing.T) {
	filepath := "../fixtures/recursive-inclusion.mg"
	lines, _ := shell.LoadFile(filepath)
	tkn := newTokenizer(filepath, lines)
	proc := NewPreprocessor(tkn.tokenize())
	err := proc.doIncludes()
	if err == nil {
		t.Log(err)
		t.Fatalf("expected to fail after too many inclusions")
	}
}

func Test_MultilevelInclusionShouldWork(t *testing.T) {
	filepath := "../fixtures/double-include-parent.mg"
	lines, _ := shell.LoadFile(filepath)
	tkn := newTokenizer(filepath, lines)
	proc := NewPreprocessor(tkn.tokenize())
	err := proc.doIncludes()
	if err != nil {
		t.Log(err)
		t.Fatalf("multi-level include should work")
	}
	if len(proc.tokens) != 87 {
		t.Fatalf("expected exactly 87 tokens in multi-level include, got %d", len(proc.tokens))
	}
}
