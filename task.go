package main

import (
	"os"
)

type executionItem struct {
	pos sourcePosition
	name string
	spec []string
}

type task interface {
	getAge() int
	getName() string
	getCommands() []string
}

type ruleTask struct { executionItem }
type fileTask struct { executionItem }

func (r executionItem)getName() string {
	return r.name
}
func (r executionItem)getCommands() []string {
	return r.spec
}

func (t ruleTask)getAge() int {
	return 0
}
func (t fileTask)getAge() int {
	return 1
}

func newTask(dfn taskDefinition) task {
	_, err := os.Stat(dfn.name[:len(dfn.name)-1])
	if err != nil {
		return ruleTask{executionItem{dfn.pos, dfn.name, dfn.commands}}
	}
	return fileTask{executionItem{dfn.pos, dfn.name, dfn.commands}}
}
