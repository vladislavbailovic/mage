package evaluation

import (
	"fmt"
	"mage/typedefs"
	"strings"
)

// Accepts a set of task definitions and start point, and
// orders them into an execution stack

func getEvaluationStackFrom(start string, dfns map[string]typedefs.TaskDefinition) ([]task, error) {
	stack := []task{}
	return getEvaluationSubstackFrom(start, dfns, stack)
}

func getEvaluationSubstackFrom(start string, dfns map[string]typedefs.TaskDefinition, stack []task) ([]task, error) {
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

	item := newTask(root)
	stack = append(stack, item)

	return stack, nil
}
