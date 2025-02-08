package main

import (
	"fmt"
)

func main() {
	x := "let a  = 7"
	slice(x)
}
func advance(text string, position int32) {
	pos := position
	pos+=1
	curr_char :=
}
func slice(text string) {
	str := text
	var res []string
	for _, ch := range str {
		res = append(res, fmt.Sprintf("%q", ch))
	}
	fmt.Println()
}
