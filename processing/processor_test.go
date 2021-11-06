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
