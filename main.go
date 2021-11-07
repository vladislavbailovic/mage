package main

import (
	"flag"
	"fmt"
	"mage/evaluation"
	"mage/processing"
	"mage/typedefs"
	"os"
)

const (
	FIXTURE      string = "fixtures/macro.mg"
	RECORD_STORE string = ".Magefile"
	ROOT_TASK    string = "root"
)

func main() {
	var err error

	file := flag.String("f", FIXTURE, "File to process")
	flag.Parse()

	var dfns map[string]typedefs.TaskDefinition
	proc := processing.NewProcessor(*file)
	dfns, err = proc.GetTasks()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var rootTask string
	if len(flag.Args()) == 0 {
		rootTask, err = proc.GetFirstTaskName()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		rootTask = flag.Args()[0]
	}

	stack := evaluation.NewStack(rootTask, dfns)
	records := evaluation.NewRecordStore(RECORD_STORE)
	stack.SetRecords(records)

	var tasks []typedefs.Task
	tasks, err = stack.Evaluate()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(evaluation.Compile(tasks))
	stack.Record()
}
