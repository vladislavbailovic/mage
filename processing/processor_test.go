package processing

import (
	"mage/shell"
	"testing"
)

func Test_Process(t *testing.T) {
	lines, _ := shell.LoadFile("../fixtures/macro.mg")
	tkn := newTokenizer("../fixtures/macro.mg", lines)
	rawTokens := tkn.tokenize()
	tokens, _ := preprocess(rawTokens)
	dfns, err := process(tokens)
	if err != nil {
		t.Fatalf("expected processing to succeed")
	}
	//debug.TaskDefinitions(dfns)
	if len(dfns) != 2 {
		t.Fatalf("expected 2 task definitions, but got %d", len(dfns))
	}
}

func Test_Process_RedefiningRulesCausesError(t *testing.T) {
	lines, _ := shell.LoadFile("../fixtures/rule-conflict.mg")
	tkn := newTokenizer("../fixtures/rule-conflict.mg", lines)
	rawTokens := tkn.tokenize()
	tokens, _ := preprocess(rawTokens)
	_, err := process(tokens)
	if err == nil {
		t.Log(err)
		t.Fatalf("expected rule name conflict")
	}
}
