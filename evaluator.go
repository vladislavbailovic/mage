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

func evaluateStack(stack []task, epoch int64) {
	store := newRecordStore(RECORD_STORE)
	for _, t := range stack {
		age := t.getAge()
		if age == 0 {
			age = int64(store.getTime(t.getName()))
		}

		if age <= epoch {
			fmt.Println("... skip older", t.getName())
			continue
		}
		fmt.Println(">", t.getName())
		evaluateTask(t, store)
	}
	store.save()
}

func evaluateTask(t task, store *recordStore) {
	for idx, command := range t.getCommands() {
		fmt.Println("\t -", idx, ":", command)
	}
	store.recordTime(t.getName())
}
