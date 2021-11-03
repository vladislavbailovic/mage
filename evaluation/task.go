package evaluation

import (
	"mage/shell"
	"mage/typedefs"
)

type executionItem struct {
	pos  typedefs.SourcePosition
	name string
	spec []string
}

type task interface {
	getAge() int64
	getName() string
	getCommands() []string
}

type ruleTask struct{ executionItem }
type fileTask struct{ executionItem }

func (r executionItem) getName() string {
	return r.name
}
func (r executionItem) getCommands() []string {
	return r.spec
}

func (t ruleTask) getAge() int64 {
	return 0
}
func (t fileTask) getAge() int64 {
	return shell.GetFileMtime(t.name)
}

func newTask(dfn typedefs.TaskDefinition) task {
	if shell.FileExists(dfn.Name) {
		return fileTask{executionItem{dfn.Pos, dfn.Name, dfn.Commands}}
	}
	return ruleTask{executionItem{dfn.Pos, dfn.Name, dfn.Commands}}
}
