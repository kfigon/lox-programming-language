package main

import (
	"bufio"
	"fmt"
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
	t, err := lex(string(b))
	if err != nil {
		fmt.Println("Got error:", err)
		return
	}
	fmt.Println(t)
}

func interpreterMode() {
	fmt.Println("Welcome to lox interpreter")
	fmt.Println("type 'quit' to exit")
	for true {
		fmt.Print("> ")
		in := bufio.NewReader(os.Stdin)
		line, _ := in.ReadString('\n')
		line = strings.TrimSuffix(line, "\r\n")
		line = strings.TrimSuffix(line, "\n")

		if line == "quit" || line == "exit" {
			fmt.Println("Bye")
			return
		} else if line != "" {
			t, err := lex(line)
			if err != nil {
				fmt.Println("got lexer error: ", err)
				continue
			}
			exp, errs := NewParser(t).Parse()
			if len(errs) > 0 {
				fmt.Println("got parser errors: ", errs)
				continue
			}			
			got, err := interpret(exp)
			if err != nil {
				fmt.Println("got interpreter error:",err)
				continue
			}
			for _,v := range got {
				fmt.Println(*(v.v))
			}
		}
	}
}
