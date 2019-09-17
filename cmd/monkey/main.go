package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ljanyst/monkey/pkg/evaluator"
)

const PROMPT = ">> "

func startRepl() {
	fmt.Print("This is a monkey evaluator\n")
	scanner := bufio.NewScanner(os.Stdin)
	c := evaluator.NewContext()
	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			break
		}

		line := scanner.Text()
		obj, err := evaluator.EvalString(line, c, "stdin")
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
		} else {
			fmt.Printf("%s\n", obj.Inspect())
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

	c := evaluator.NewContext()
	_, err = evaluator.EvalReader(file, c, filename)
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
