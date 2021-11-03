package evaluation

// import (
// 	"fmt"
// 	"errors"

// 	"mage/processor"
// )

// func GetStack(startNode string, parser processor.Parser) ([]task, error) {
// 	parser.Parse()
// 	stack, err := prepareEvaluationStack(startNode, parser, []task{})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return stack, nil
// }

// func prepareEvaluationStack(taskName string, parser processor.Parser, stack []task) ([]task, error) {
// 	dfn, ok := parser.Tasks[taskName + ":"]
// 	if !ok {
// 		return nil, errors.New("Unable to resolve task definition for: " + taskName)
// 	}
// 	var err error
// 	for _, dependency := range dfn.Dependencies {
// 		stack, err = prepareEvaluationStack(dependency, parser, stack)
// 		if err != nil {
// 			errMsg := fmt.Errorf(
// 				"file %s, line %d (%s) %v",
// 				dfn.Pos.File,
// 				dfn.Pos.Line,
// 				dfn.NormalizedName,
// 				err,
// 			)
// 			return nil, errMsg
// 		}
// 	}
// 	item := newTask(dfn)
// 	stack = append(stack, item)

// 	return stack, nil
// }

// func EvaluateStack(stack []task, epoch int64, store *recordStore) {
// 	for _, t := range stack {
// 		age := t.getAge()
// 		if age == 0 {
// 			age = int64(store.GetTime(t.getName()))
// 		}

// 		if age <= epoch {
// 			fmt.Println("... skip older", t.getName())
// 			continue
// 		}
// 		fmt.Println(">", t.getName())
// 		evaluateTask(t, store)
// 	}
// 	store.Save()
// }

// func evaluateTask(t task, store *recordStore) {
// 	for idx, command := range t.getCommands() {
// 		fmt.Println("\t -", idx, ":", command)
// 	}
// 	store.RecordTime(t.getName())
// }
