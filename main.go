package main

import (
	"bufio"
	"fmt"
	"lox/interpreter"
	"lox/lexer"
	"lox/parser"
	"os"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		interpreterMode()
	} else if len(os.Args) == 2 {
		fileMode(os.Args[1])
	} else {
		fmt.Println("Invalid number of arguments")
	}
}

func fileMode(fileName string) {
	b, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Cant open file %v: %v\n", fileName, err)
		return
	}
	stmts, err := parse(string(b))
	if err != nil {
		fmt.Println(err)
		return
	}

	in := interpreter.NewInterpreter()
	for _, s := range stmts {
		if err := s.AcceptStatement(in); err != nil {
			fmt.Println("got error:", err)
			return
		}
	}
}

func interpreterMode() {
	fmt.Println("Welcome to lox interpreter")
	fmt.Println("type 'quit' to exit")
	in := interpreter.NewInterpreter()

	for true {
		fmt.Print("> ")
		line := getLine()

		if line == "quit" || line == "exit" {
			fmt.Println("Bye")
			return
		} else if line != "" {
			
			stmts, err := parse(line)
			if err != nil {
				fmt.Println(err)
				return
			}

			for _, s := range stmts {
				s.AcceptStatement(in)
			}

		}
	}
}

func getLine() string {
	in := bufio.NewReader(os.Stdin)
	line, _ := in.ReadString('\n')
	line = strings.TrimSuffix(line, "\r\n")
	line = strings.TrimSuffix(line, "\n")
	return line
}

func parse(input string) ([]parser.Statement, error) {
	toks, err := lexer.Lex(input)
	if err != nil {
		return nil, fmt.Errorf("lexer error: %w", err)
	}
	got, errs := parser.NewParser(toks).Parse()
	if len(errs) != 0 {
		v := []string{}
		for _, e := range errs {
			v = append(v, e.Error())
		}
		return nil, fmt.Errorf("parser errors: %s", strings.Join(v, ","))
	}
	
	return got, nil
}