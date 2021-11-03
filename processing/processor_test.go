package processing

import (
	"strings"
	"testing"
)

func Test_Process(t *testing.T) {
	lines, _ := loadFile("../fixtures/macro.mg")
	tkn := newTokenizer("macro.mg", strings.Join(lines, "\n"))
	rawTokens := tkn.tokenize()
	tokens, _ := preprocess(rawTokens)
	dfns, err := process(tokens)
	if err != nil {
		t.Fatalf("expected processing to succeed")
	}
	//debugTaskDefinitions(dfns)
	if len(dfns) != 2 {
		t.Fatalf("expected 2 task definitions, but got %d", len(dfns))
	}
}

func Test_Process_RedefiningRulesCausesError(t *testing.T) {
	lines, _ := loadFile("../fixtures/rule-conflict.mg")
	tkn := newTokenizer("rule-conflict.mg", strings.Join(lines, "\n"))
	rawTokens := tkn.tokenize()
	tokens, _ := preprocess(rawTokens)
	_, err := process(tokens)
	if err == nil {
		t.Log(err)
		t.Fatalf("expected rule name conflict")
	}
}
