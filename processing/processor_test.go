package processing

import (
	"mage/debug"
	"mage/processing/preprocessing"
	"mage/processing/tokenizing"
	"mage/shell"
	"testing"
)

func Test_Process(t *testing.T) {
	lines, _ := shell.LoadFile("../fixtures/macro.mg")
	tkn := tokenizing.NewTokenizer("../fixtures/macro.mg", lines)
	rawTokens := tkn.Tokenize()
	tokens, _ := preprocessing.Preprocess(rawTokens)
	dfns, err := process(tokens)
	if err != nil {
		debug.Tokens(tokens)
		t.Log(err)
		t.Fatalf("expected processing to succeed")
	}
	//debug.TaskDefinitions(dfns)
	if len(dfns) != 2 {
		t.Fatalf("expected 2 task definitions, but got %d", len(dfns))
	}
}

func Test_Process_RedefiningRulesCausesError(t *testing.T) {
	lines, _ := shell.LoadFile("../fixtures/rule-conflict.mg")
	tkn := tokenizing.NewTokenizer("../fixtures/rule-conflict.mg", lines)
	rawTokens := tkn.Tokenize()
	tokens, _ := preprocessing.Preprocess(rawTokens)
	_, err := process(tokens)
	if err == nil {
		t.Log(err)
		t.Fatalf("expected rule name conflict")
	}
}
