package evaluation

import (
	"mage/processing"
	"testing"
)

func Test_Ordering_StackGettingFailsForInvalidRootTask(t *testing.T) {
	proc := processing.NewProcessor("../fixtures/simple.mg")
	dfns, _ := proc.GetTasks()
	_, err := getEvaluationStackFrom("non-existent-task", dfns)
	if err == nil {
		t.Fatalf("stack getting should fail for invalid first task")
	}
}

func Test_Ordering(t *testing.T) {
	// dfns, err := processing.ProcessFile("../fixtures/simple.mg")
	proc := processing.NewProcessor("../fixtures/simple.mg")
	dfns, err := proc.GetTasks()
	if err != nil {
		t.Log(err)
		t.Fatalf("expected processing to be a success")
	}

	tasks, err := getEvaluationStackFrom("root", dfns)
	expected := []string{
		"parser.go",
		"dependency1",
		"tmp",
		"tmp/not-created-yet.go",
		"dep-dependency1",
		"dependency2",
		"root",
	}

	if len(expected) != len(tasks) {
		t.Fatalf("expected %d tasks, but got %d", len(expected), len(tasks))
	}

	for idx, taskName := range expected {
		actual := tasks[idx]
		if taskName != actual.GetName() {
			t.Fatalf("expected %s at position %d - got %s instead", taskName, idx, actual.GetName())
		}
	}
}
