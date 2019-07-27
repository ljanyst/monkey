package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ljanyst/edu-interp/pkg/evaluator"
)

const PROMPT = ">> "

func startRepl() {
	fmt.Print("This is a monkey evaluator\n")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			break
		}

		line := scanner.Text()
		err := evaluator.EvalString(line)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}
	}

	fmt.Print("Bye!\n")
}

func run(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
	}
	defer file.Close()

	err = evaluator.EvalReader(file)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) == 1 {
		startRepl()
	} else if len(os.Args) == 2 {
		run(os.Args[1])
	} else {
		fmt.Print("Usage:\n")
		fmt.Printf("    %s - take commands from stdin\n", os.Args[0])
		fmt.Printf("    %s filename.monkey - evaluate filename.monkey\n", os.Args[0])
	}
}
