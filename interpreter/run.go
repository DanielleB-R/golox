package interpreter

import (
	"bufio"
	"fmt"
	"os"
)

func RunFile(path string) {
	script, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file", path)
		os.Exit(1)
	}
	err = run(string(script))

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(65)
	}
}

// NOTE: This doesn't work! Solve that in a bit
func RunPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		err := run(scanner.Text())
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func run(source string) error {
	scanner := NewSourceScanner(source)
	tokens, err := scanner.ScanTokens()
	if err != nil {
		return err
	}

	parser := NewParser(tokens)
	statements, err := parser.Parse()
	if err != nil {
		return err
	}

	// NOTE: We'll need to have a global interpreter to make the repl work right
	interpreter := NewInterpreter()
	interpreter.Interpret(statements)

	return nil
}
