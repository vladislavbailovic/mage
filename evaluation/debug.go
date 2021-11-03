package evaluation

import "fmt"

func debugTasks(tasks []task) {
	for i, tsk := range tasks {
		for j, cmd := range tsk.getCommands() {
			fmt.Printf("> task %d (%s), command %d: [%s]\n", i, tsk.getName(), j, cmd)
		}
	}
}
