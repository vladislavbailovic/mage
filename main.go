package main

import (
	"fmt"
	"os"
)

func main() {
	parser, err := newParser("fixtures/simple.mg")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	myAge := int64(0)
	stack, err := getStack("root", parser)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, t := range stack {
		if t.getAge() > myAge {
			fmt.Println("... skip newer", t.getName())
			continue
		}
		fmt.Println(">", t.getName())
		for idx, command := range t.getCommands() {
			fmt.Println("\t -", idx, ":", command)
		}
	}
}

