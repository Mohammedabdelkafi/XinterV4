package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func run() {
	a := true
	calcMode := false
	debug := false
	variables := make(map[string]interface{})
	commands := []string{}
	reader := bufio.NewReader(os.Stdin)
	for a {
		fmt.Print("Xinter ==> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		switch text {
		case "calc":
			calcMode = true
			fmt.Println("Calc mode activated")
		case "decalc":
			calcMode = false
			fmt.Println("Calc mode deactivated")
		case "dev":
			debug = true
			fmt.Println("Developer mode activated")
		case "undev":
			debug = false
			fmt.Println("Developer mode deactivated")
		case "run":
			fmt.Println("Running all commands...")
			for _, cmd := range commands {
				lexer := NewLexer(cmd, debug)
				NewParser(lexer.tokens, calcMode, debug, variables)
			}
		case "exit":
			fmt.Println("Exiting")
			a = false
		default:
			commands = append(commands, text)
			lexer := NewLexer(text, debug)
			NewParser(lexer.tokens, calcMode, debug, variables)
		}
	}
}

func main() {
	run()
}
