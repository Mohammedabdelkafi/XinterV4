package main

import (
	"fmt"
)

func main() {
	x := "let a  = 7"
	slice(x)
}
func advance(text string) {
	print("i'm bored")
}
func slice(text string) {
	str := text
	var res []string
	for _, ch := range str {
		res = append(res, fmt.Sprintf("%q", ch))
	}
	fmt.Println()
}
func Tokenize() {}
