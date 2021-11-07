package debug

import (
	"fmt"
	"mage/typedefs"
)

func Tasks(tasks []typedefs.Task) {
	for i, tsk := range tasks {
		for j, cmd := range tsk.GetCommands() {
			fmt.Printf("> task %d (%s), command %d: [%s]\n", i, tsk.GetName(), j, cmd)
		}
	}
}

func Records(records map[string]typedefs.Record) {
	for _, rec := range records {
		fmt.Println("[%s]: %d\n", rec.Name, rec.Timestamp)
	}
}
