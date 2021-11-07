package evaluation

import (
	"fmt"
	"mage/typedefs"
	"strings"
)

// Accepts a set of task definitions and start point, and
// orders them into an execution stack

func GetEvaluationStack(start string, dfns map[string]typedefs.TaskDefinition) ([]typedefs.Task, error) {
	stack := NewStack(start, dfns)
	return stack.Evaluate()
}

type Stack struct {
	dfns    map[string]typedefs.TaskDefinition
	root    string
	time    typedefs.Epoch
	records *recordStore
}

func NewStack(start string, dfns map[string]typedefs.TaskDefinition) *Stack {
	records := NewRecordStore("")
	return &Stack{dfns, start, typedefs.Epoch(0), records}
}

func (s *Stack) SetEpoch(t typedefs.Epoch) {
	s.time = t
}

func (s *Stack) SetRoot(r string) {
	s.root = r
}

func (s *Stack) SetRecords(rs *recordStore) {
	s.records = rs
}

func (s Stack) Evaluate() ([]typedefs.Task, error) {
	stack := []typedefs.Task{}
	return s.evaluateSubstack(s.root, stack)
}

func (s Stack) Record() {
	s.records.save()
}

func (s Stack) evaluateSubstack(start string, stack []typedefs.Task) ([]typedefs.Task, error) {
	root, ok := s.dfns[start]
	if !ok {
		return nil, fmt.Errorf("unable to resolve descending root task: [%s]", start)
	}

	item := newTask(root, s.records)
	if !s.withinEpoch(item) {
		return stack, nil
	}

	var err error
	for i := len(root.Dependencies) - 1; i >= 0; i-- {
		dependency := strings.TrimSpace(root.Dependencies[i]) // @TODO: lexer issue, fix
		stack, err = s.evaluateSubstack(dependency, stack)
		if err != nil {
			return nil, err
		}
	}

	stack = append(stack, item)

	return stack, nil
}

func (s Stack) withinEpoch(what typedefs.Milestone) bool {
	return what.GetMilestone() >= s.time
}
