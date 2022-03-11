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
	run(string(script))
}

// NOTE: This doesn't work! Solve that in a bit
func RunPrompt() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		run(scanner.Text())
	}
}

func run(source string) {
	scanner := NewSourceScanner(source)
	tokens := scanner.ScanTokens()

	for token := range tokens {
		fmt.Println(token)
	}
}
