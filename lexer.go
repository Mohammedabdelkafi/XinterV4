package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Lexer struct {
	src      []rune
	tokens   []Token
	pos      int
	currChar rune
	digits   string
	letters  string
	debug    bool
}

func NewLexer(text string, debug bool) *Lexer {
	l := &Lexer{
		src:     []rune(text),
		tokens:  []Token{},
		pos:     -1,
		digits:  "1234567890",
		letters: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_",
		debug:   debug,
	}
	l.advance()
	l.tokenize()
	return l
}

func (l *Lexer) advance() {
	l.pos++
	if l.pos < len(l.src) {
		l.currChar = l.src[l.pos]
	} else {
		l.currChar = 0 // Null character
	}
	if l.debug {
		fmt.Printf("Lexer advance: pos=%d, curr_char=%c\n", l.pos, l.currChar)
	}
}

func (l *Lexer) tokenize() {
	tokenMap := map[rune]string{
		'+': "PLUS", '-': "MINUS", '*': "MULTIPLY", '/': "DIVIDE",
		'=': "EQUALS", '&': "AND", '|': "OR", '!': "NOT",
		'(': "LPAREN", ')': "RPAREN",
	}
	for l.currChar != 0 {
		if tokenType, found := tokenMap[l.currChar]; found {
			l.tokens = append(l.tokens, Token{Type: tokenType, Value: l.currChar})
			l.advance()
		} else if l.currChar == '"' {
			l.parseString()
		} else if l.currChar == ' ' {
			l.advance() // Skip whitespace
		} else if strings.ContainsRune(l.digits, l.currChar) {
			l.parseNum()
		} else if strings.ContainsRune(l.letters, l.currChar) {
			l.parseVar()
		} else {
			panic(fmt.Sprintf("Unexpected character: %c", l.currChar))
		}
	}
	if l.debug {
		fmt.Printf("Lexer tokens: %v\n", l.tokens)
	}
}

func (l *Lexer) parseNum() {
	numStr := ""
	for l.currChar != 0 && strings.ContainsRune(l.digits, l.currChar) {
		numStr += string(l.currChar)
		l.advance()
	}
	num, _ := strconv.Atoi(numStr)
	l.tokens = append(l.tokens, Token{Type: "NUMBER", Value: num})
}

func (l *Lexer) parseString() {
	strVal := ""
	l.advance() // Skip the opening quote
	for l.currChar != 0 && l.currChar != '"' {
		strVal += string(l.currChar)
		l.advance()
	}
	if l.currChar == '"' {
		l.advance() // Skip the closing quote
		l.tokens = append(l.tokens, Token{Type: "STRING", Value: strVal})
	} else {
		panic("Unterminated string literal")
	}
}

func (l *Lexer) parseVar() {
	varStr := ""
	for l.currChar != 0 && strings.ContainsRune(l.letters+l.digits, l.currChar) {
		varStr += string(l.currChar)
		l.advance()
	}
	if varStr == "true" || varStr == "false" {
		l.tokens = append(l.tokens, Token{Type: "BOOLEAN", Value: varStr == "true"})
	} else {
		l.tokens = append(l.tokens, Token{Type: "IDENTIFIER", Value: varStr})
	}
}
