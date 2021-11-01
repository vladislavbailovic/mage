package main

import (
	"fmt"
	"os"

	"mage/processor"
	"mage/ruleset"
)

const (
	FIXTURE string = "fixtures/macro.mg"
	RECORD_STORE string = "tmp/test.json"
	ROOT_TASK string = "root"
)

func main() {
	// evaluateRules()
}

func evaluateRules() {
	parser, err := processor.NewParser(FIXTURE)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	store := ruleset.NewRecordStore(RECORD_STORE)
	myAge := int64(store.GetTime(ROOT_TASK))
	stack, err := ruleset.GetStack(ROOT_TASK, parser)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ruleset.EvaluateStack(stack, myAge, store)

	store.RecordTime("root")
	store.Save()
}

