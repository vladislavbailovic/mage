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

	tasks, errT := evaluation.GetEvaluationStack("root", dfns)
	if errT != nil {
		fmt.Println(errT)
		os.Exit(1)
	}

	fmt.Println(evaluation.Compile(tasks))
}
