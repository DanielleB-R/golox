package interpreter

import (
	"bufio"
	"fmt"
	"os"
)

func RunFile(path string) {
	script, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading file", path)
		os.Exit(1)
	}
	err = run(string(script), NewInterpreter())

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(65)
	}
}

func RunPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	interpreter := NewInterpreter()
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		err := run(scanner.Text(), interpreter)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func run(source string, interpreter *Interpreter) error {
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

	resolver := NewResolver(interpreter)
	err = resolver.Resolve(statements)
	if err != nil {
		return err
	}

	return interpreter.Interpret(statements)
}
