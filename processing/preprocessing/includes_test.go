package preprocessing

import (
	"mage/debug"
	"mage/processing/tokenizing"
	"mage/shell"
	"testing"
)

func Test_PreprocessorDoesIncludes(t *testing.T) {
	filepath := "../../fixtures/includes.mg"
	lines, _ := shell.LoadFile(filepath)
	tkn := tokenizing.NewTokenizer(filepath, lines)
	proc := newPreprocessor(tkn.Tokenize())
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
	filepath := "../../fixtures/includes.mg"
	lines, _ := shell.LoadFile(filepath)
	tkn := tokenizing.NewTokenizer(filepath, lines)
	proc := newPreprocessor(tkn.Tokenize())
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

	debug.Tokens(proc.tokens)
	if 102 != len(proc.tokens) {
		t.Fatalf("expected exactly 102 tokens with includes, got %d", len(proc.tokens))
	}
}

func Test_RecursiveInclusionShouldErrorOut(t *testing.T) {
	filepath := "../../fixtures/recursive-inclusion.mg"
	lines, _ := shell.LoadFile(filepath)
	tkn := tokenizing.NewTokenizer(filepath, lines)
	proc := newPreprocessor(tkn.Tokenize())
	err := proc.doIncludes()
	if err == nil {
		t.Log(err)
		t.Fatalf("expected to fail after too many inclusions")
	}
}

func Test_MultilevelInclusionShouldWork(t *testing.T) {
	filepath := "../../fixtures/double-include-parent.mg"
	lines, _ := shell.LoadFile(filepath)
	tkn := tokenizing.NewTokenizer(filepath, lines)
	proc := newPreprocessor(tkn.Tokenize())
	err := proc.doIncludes()
	if err != nil {
		t.Log(err)
		t.Fatalf("multi-level include should work")
	}
	if len(proc.tokens) != 87 {
		t.Fatalf("expected exactly 87 tokens in multi-level include, got %d", len(proc.tokens))
	}
}
