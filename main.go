package main

import (
	"fmt"
	"os"

	"github.com/DanielleB-R/golox/interpreter"
)

func main() {
	if len(os.Args) > 2 {
		fmt.Fprintln(os.Stderr, "Usage: golox [script]")
		os.Exit(64)
	} else if len(os.Args) == 2 {
		interpreter.RunFile(os.Args[1])
	} else {
		interpreter.RunPrompt()
	}

}
