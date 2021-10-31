package main

import (
	"fmt"
	"os"
)

const (
	FIXTURE string = "fixtures/simple.mg"
	RECORD_STORE string = "tmp/test.json"
	ROOT_TASK string = "root"
)

func main() {
	parser, err := newParser(FIXTURE)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	store := newRecordStore(RECORD_STORE)
	myAge := int64(store.getTime(ROOT_TASK))
	stack, err := getStack(ROOT_TASK, parser)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	evaluateStack(stack, myAge)

	store.recordTime("root")
	store.save()
}

