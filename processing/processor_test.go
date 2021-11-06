package processing

import (
	"testing"
)

func Test_Process(t *testing.T) {
	proc := NewProcessor("../fixtures/macro.mg")
	dfns, err := proc.GetTasks()
	if err != nil {
		t.Log(err)
		t.Fatalf("expected processing to succeed")
	}
	//debug.TaskDefinitions(dfns)
	if len(dfns) != 2 {
		t.Fatalf("expected 2 task definitions, but got %d", len(dfns))
	}
}

func Test_Process_RedefiningRulesCausesError(t *testing.T) {
	proc := NewProcessor("../fixtures/rule-conflict.mg")
	_, err := proc.GetTasks()
	if err == nil {
		t.Log(err)
		t.Fatalf("expected rule name conflict")
	}
}

func Test_Process_GetFirstRule_ErrorsOnEmptyStack(t *testing.T) {
	proc := NewProcessor("../fixtures/includes.mg")
	_, err := proc.GetFirstTaskName()
	if err == nil {
		t.Fatalf("expected to error out if asked for first task before processing everything")
	}
}

func Test_Process_GetFirstRule_ResolvesFirstRule(t *testing.T) {
	proc := NewProcessor("../fixtures/includes.mg")
	proc.GetTasks()
	rule, err := proc.GetFirstTaskName()
	if err != nil {
		t.Fatalf("expected rule to resolve")
	}
	if rule != "conflict-rule" { // because inclusion comes first
		t.Fatalf("expected first task to be 'conflict-rule', got %s", rule)
	}
}
