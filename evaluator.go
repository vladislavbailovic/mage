package main

import (
	"fmt"
	"errors"
)

func getStack(startNode string, parser parser) ([]task, error) {
	parser.parse()
	stack, err := prepareEvaluationStack(startNode, parser, []task{})
	if err != nil {
		return nil, err
	}
	return stack, nil
}

func prepareEvaluationStack(taskName string, parser parser, stack []task) ([]task, error) {
	dfn, ok := parser.tasks[taskName + ":"]
	if !ok {
		return nil, errors.New("Unable to resolve task definition for: " + taskName)
	}
	var err error
	for _, dependency := range dfn.dependencies {
		stack, err = prepareEvaluationStack(dependency, parser, stack)
		if err != nil {
			errMsg := fmt.Errorf(
				"file %s, line %d (%s) %v",
				dfn.pos.file,
				dfn.pos.line,
				dfn.normalizedName,
				err,
			)
			return nil, errMsg
		}
	}
	item := newTask(dfn)
	stack = append(stack, item)

	return stack, nil
}
