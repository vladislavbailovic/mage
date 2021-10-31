package main

import "fmt"

type task struct {
	pos sourcePosition
	name string
	spec []string
}

func main() {
	parser := newParser("main.go")
	parser.parse()

	stack := prepareEvaluationStack("root", parser, []task{})
	fmt.Println(stack)
}

func prepareEvaluationStack(taskName string, parser parser, stack []task) []task {
	dfn, ok := parser.tasks[taskName + ":"]
	if !ok {
		panic("Unable to resolve task definition for: " + taskName)
	}
	for _, dependency := range dfn.dependencies {
		stack = prepareEvaluationStack(dependency, parser, stack)
	}
	item := task{
		dfn.pos,
		dfn.name,
		dfn.commands,
	}
	stack = append(stack, item)

	return stack
}
