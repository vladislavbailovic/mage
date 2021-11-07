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

type ruleTask struct {
	executionItem
	records recordStore
}
type fileTask struct{ executionItem }

func (r executionItem) GetName() string {
	return r.name
}
func (r executionItem) GetCommands() []string {
	return r.spec
}

func (t ruleTask) GetMilestone() typedefs.Epoch {
	return t.records.getTime(t.GetName())
}
func (t fileTask) GetMilestone() typedefs.Epoch {
	return typedefs.Epoch(shell.GetFileMtime(t.name))
}

func newTask(dfn typedefs.TaskDefinition, records recordStore) typedefs.Task {
	if shell.FileExists(dfn.Name) {
		return fileTask{executionItem{dfn.Pos, dfn.Name, dfn.Commands}}
	}
	return ruleTask{executionItem{dfn.Pos, dfn.Name, dfn.Commands}, records}
}
