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
	getAge() int64
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

func (t ruleTask)getAge() int64 {
	return 0
}
func (t fileTask)getAge() int64 {
	fpath := t.name[:len(t.name)-1]
	f, err := os.Stat(fpath)
	if err != nil {
		return 0
	}
	return f.ModTime().Unix()
}

func newTask(dfn taskDefinition) task {
	_, err := os.Stat(dfn.name[:len(dfn.name)-1])
	if err != nil {
		return ruleTask{executionItem{dfn.pos, dfn.name, dfn.commands}}
	}
	return fileTask{executionItem{dfn.pos, dfn.name, dfn.commands}}
}
