package typedefs

import (
	"mage/shell"
)

type executionItem struct {
	pos  SourcePosition
	name string
	spec []string
}

type Task interface {
	GetAge() int64
	GetName() string
	GetCommands() []string
}

type ruleTask struct{ executionItem }
type fileTask struct{ executionItem }

func (r executionItem) GetName() string {
	return r.name
}
func (r executionItem) GetCommands() []string {
	return r.spec
}

func (t ruleTask) GetAge() int64 {
	return 0
}
func (t fileTask) GetAge() int64 {
	return shell.GetFileMtime(t.name)
}

func NewTask(dfn TaskDefinition) Task {
	if shell.FileExists(dfn.Name) {
		return fileTask{executionItem{dfn.Pos, dfn.Name, dfn.Commands}}
	}
	return ruleTask{executionItem{dfn.Pos, dfn.Name, dfn.Commands}}
}
