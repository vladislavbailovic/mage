package main

import (
	"flag"
	"fmt"
	"mage/evaluation"
	"mage/processing"
	"os"
)

const (
	FIXTURE      string = "fixtures/macro.mg"
	RECORD_STORE string = "tmp/test.json"
	ROOT_TASK    string = "root"
)

func main() {
	file := flag.String("f", FIXTURE, "File to process")
	flag.Parse()

	proc := processing.NewProcessor(*file)
	dfns, errD := proc.GetTasks()
	if errD != nil {
		fmt.Println(errD)
		os.Exit(1)
	}

	var rootTask string
	if len(flag.Args()) == 0 {
		var errF error
		rootTask, errF = proc.GetFirstTaskName()
		if errF != nil {
			fmt.Println(errF)
			os.Exit(1)
		}
	} else {
		rootTask = flag.Args()[0]
	}

	tasks, errT := evaluation.GetEvaluationStack(rootTask, dfns)
	if errT != nil {
		fmt.Println(errT)
		os.Exit(1)
	}

	fmt.Println(evaluation.Compile(tasks))
}
