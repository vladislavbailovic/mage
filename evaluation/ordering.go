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
	evaluated bool
	dfns      map[string]typedefs.TaskDefinition
	root      string
	time      typedefs.Epoch
	records   *recordStore
	tasks     []typedefs.Task
}

func NewStack(start string, dfns map[string]typedefs.TaskDefinition) *Stack {
	records := NewRecordStore("")
	return &Stack{false, dfns, start, typedefs.Epoch(-1), records, []typedefs.Task{}}
}

func (s *Stack) SetEpoch(t typedefs.Epoch) {
	s.time = t
}

func (s *Stack) SetRoot(r string) {
	s.root = r
}

func (s *Stack) SetRecords(rs *recordStore) {
	s.records = rs
	ts := s.records.getTime(s.root)
	if ts > 0 {
		s.time = ts
	}
}

func (s *Stack) Evaluate() ([]typedefs.Task, error) {
	if s.evaluated {
		return s.tasks, nil
	}

	stack := []typedefs.Task{}
	s.evaluated = true
	tasks, err := s.evaluateSubstack(s.root, stack)
	if err != nil {
		return nil, err
	}

	s.tasks = tasks
	return tasks, nil
}

func (s Stack) Record() error {
	if !s.evaluated {
		return fmt.Errorf("unable to record times before evaluating tasks stack")
	}
	for _, task := range s.tasks {
		s.records.recordTime(task.GetName())
	}
	err := s.records.save()
	return err
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
	return what.GetMilestone() > s.time
}
