package evaluation

import (
	"mage/processing"
	"strings"
	"testing"
)

func Test_GetCompiledStatements(t *testing.T) {
	proc := processing.NewProcessor("../fixtures/includes.mg")
	dfns, _ := proc.GetTasks()
	tasks, _ := getEvaluationStackFrom("root", dfns)

	outputs := GetCompiledStatements(tasks)

	if len(outputs) != 17 {
		t.Log(outputs)
		t.Fatalf("expected exactly 17 outputs from commands, got %d", len(outputs))
	}
}

func Test_Compile(t *testing.T) {
	proc := processing.NewProcessor("../fixtures/includes.mg")
	dfns, _ := proc.GetTasks()
	tasks, _ := getEvaluationStackFrom("root", dfns)

	statements := GetCompiledStatements(tasks)
	out := Compile(tasks)

	if strings.Join(statements, "\n") != out {
		t.Log(statements)
		t.Fatalf("expected compiled output to match newline-separated statements")
	}

	if len(out) != 236 {
		t.Log(out)
		t.Fatalf("expected output to be exactly 236 chars, got %d", len(out))
	}
}
