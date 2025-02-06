package main

import (
	"fmt"
	"strings"
)

func main() {
	tokenize()
}

func tokenize() {
	str := "Hello World"
	res := strings.Split(str, " ")
	fmt.Printf("%v\n", res)
}
