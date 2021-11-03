package evaluation

import (
	"mage/processing"
	"testing"
)

func Test_ExecuteStack_HappyPath(t *testing.T) {
	dfns, _ := processing.ProcessFile("../fixtures/includes.mg")
	tasks, _ := getEvaluationStackFrom("root", dfns)

	outputs, err := Execute(tasks)
	if err != nil {
		t.Log(err)
		t.Fatalf("execution should succeed")
	}

	if len(outputs) != 7 {
		t.Fatalf("expected exactly 7 outputs from commands, got %d", len(outputs))
	}
}

func Test_ExecuteStack_HappyPathWithCommends(t *testing.T) {
	dfns, _ := processing.ProcessFile("../fixtures/run-with-comments.mg")
	tasks, _ := getEvaluationStackFrom("root", dfns)

	outputs, err := Execute(tasks)
	if err != nil {
		t.Log(err)
		t.Fatalf("execution should succeed")
	}

	if len(outputs) != 2 {
		t.Fatalf("expected exactly 7 outputs from commands, got %d", len(outputs))
	}
	if outputs[0] != "" {
		t.Fatalf("expected empty output for comment, got [%v]", outputs[0])
	}
	if outputs[1] != "Yay\n" {
		t.Fatalf("expected echoed output for echo, got [%v]", outputs[1])
	}
}

func Test_ExecuteStack_InvalidCommand(t *testing.T) {
	dfns, _ := processing.ProcessFile("../fixtures/invalid-commands.mg")
	tasks, _ := getEvaluationStackFrom("root", dfns)

	_, err := Execute(tasks)
	if err == nil {
		t.Fatalf("execution should have failed")
	}
}
