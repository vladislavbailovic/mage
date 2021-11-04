package evaluation

import (
	"fmt"
	"mage/typedefs"
	"strings"
)

// Accepts a set of task definitions and start point, and
// orders them into an execution stack

// @TODO pick entry point task
func GetEvaluationStack(start string, dfns map[string]typedefs.TaskDefinition) ([]typedefs.Task, error) {
	return getEvaluationStackFrom(start, dfns)
}

func getEvaluationStackFrom(start string, dfns map[string]typedefs.TaskDefinition) ([]typedefs.Task, error) {
	stack := []typedefs.Task{}
	return getEvaluationSubstackFrom(start, dfns, stack)
}

func getEvaluationSubstackFrom(start string, dfns map[string]typedefs.TaskDefinition, stack []typedefs.Task) ([]typedefs.Task, error) {
	root, ok := dfns[start]
	if !ok {
		return nil, fmt.Errorf("unable to resolve descending root task: [%s]", start)
	}

	var err error
	for i := len(root.Dependencies) - 1; i >= 0; i-- {
		dependency := strings.TrimSpace(root.Dependencies[i]) // @TODO: lexer issue, fix
		stack, err = getEvaluationSubstackFrom(dependency, dfns, stack)
		if err != nil {
			return nil, err
		}
	}

	item := typedefs.NewTask(root)
	stack = append(stack, item)

	return stack, nil
}
