package main

import (
	"fmt"
	"os"

	"github.com/DanielleB-R/golox/interpreter"
)

func main() {
	if len(os.Args) > 1 {
		fmt.Println("Usage: golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 1 {
		interpreter.RunFile(os.Args[0])
	} else {
		interpreter.RunPrompt()
	}

}
