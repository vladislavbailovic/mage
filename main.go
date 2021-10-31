package main

import "fmt"

func main() {
	parser := newParser("fixtures/simple.mg")
	parser.parse()
	myAge := int64(0)

	stack := prepareEvaluationStack("root", parser, []task{})
	for _, t := range stack {
		if t.getAge() > myAge {
			fmt.Println("... skip newer", t.getName())
			continue
		}
		fmt.Println(">", t.getName())
		for idx, command := range t.getCommands() {
			fmt.Println("\t -", idx, ":", command)
		}
	}
}

func prepareEvaluationStack(taskName string, parser parser, stack []task) []task {
	dfn, ok := parser.tasks[taskName + ":"]
	if !ok {
		panic("Unable to resolve task definition for: " + taskName)
	}
	for _, dependency := range dfn.dependencies {
		stack = prepareEvaluationStack(dependency, parser, stack)
	}
	item := newTask(dfn)
	stack = append(stack, item)

	return stack
}
