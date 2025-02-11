package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type Token struct {
	Type  string
	Value interface{}
}

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
	lexer := &Lexer{
		src:     []rune(text),
		tokens:  []Token{},
		pos:     -1,
		digits:  "1234567890",
		letters: "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_",
		debug:   debug,
	}
	lexer.advance()
	lexer.tokenize()
	return lexer
}

func (l *Lexer) advance() {
	l.pos++
	if l.pos < len(l.src) {
		l.currChar = l.src[l.pos]
	} else {
		l.currChar = 0
	}
	if l.debug {
		fmt.Printf("Lexer advance: pos=%d, curr_char=%q\n", l.pos, l.currChar)
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
			l.tokens = append(l.tokens, Token{Type: tokenType, Value: string(l.currChar)})
			l.advance()
		} else if l.currChar == '"' {
			l.parseString()
		} else if unicode.IsSpace(l.currChar) {
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
		fmt.Printf("Lexer tokens: %+v\n", l.tokens)
	}
}

func (l *Lexer) parseNum() {
	numStr := ""
	for l.currChar != 0 && strings.ContainsRune(l.digits, l.currChar) {
		numStr += string(l.currChar)
		l.advance()
	}
	value, _ := strconv.Atoi(numStr)
	l.tokens = append(l.tokens, Token{Type: "NUMBER", Value: value})
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
	for l.currChar != 0 && (strings.ContainsRune(l.letters, l.currChar) || strings.ContainsRune(l.digits, l.currChar)) {
		varStr += string(l.currChar)
		l.advance()
	}
	if varStr == "true" || varStr == "false" {
		value := varStr == "true"
		l.tokens = append(l.tokens, Token{Type: "BOOLEAN", Value: value})
	} else {
		l.tokens = append(l.tokens, Token{Type: "IDENTIFIER", Value: varStr})
	}
}

type Parser struct {
	tokens    []Token
	idx       int
	currTok   Token
	calcMode  bool
	debug     bool
	variables map[string]interface{}
}

func NewParser(tokens []Token, calcMode bool, debug bool, variables map[string]interface{}) *Parser {
	parser := &Parser{
		tokens:    tokens,
		idx:       -1,
		calcMode:  calcMode,
		debug:     debug,
		variables: variables,
	}
	parser.advance()
	parser.parse()
	return parser
}

func (p *Parser) advance() {
	p.idx++
	if p.idx < len(p.tokens) {
		p.currTok = p.tokens[p.idx]
	} else {
		p.currTok = Token{}
	}
	if p.debug {
		fmt.Printf("Parser advance: idx=%d, curr_tok=%+v\n", p.idx, p.currTok)
	}
}

func (p *Parser) parse() {
	for p.currTok != (Token{}) {
		if p.currTok.Type == "IDENTIFIER" {
			if p.idx+1 < len(p.tokens) && p.tokens[p.idx+1].Type == "EQUALS" {
				p.assign()
			} else if p.currTok.Value == "log" {
				p.advance()
				p.print()
			} else {
				p.expr()
			}
		} else {
			p.expr()
		}
	}
}

func (p *Parser) assign() {
	varName := p.currTok.Value.(string)
	p.advance() // skip var name
	p.advance() // skip EQUALS
	value := p.expr()
	p.variables[varName] = value
	if p.debug {
		fmt.Printf("Assigned: %s = %v\n", varName, value)
	}
}

func (p *Parser) expr() interface{} {
	result := p.term()
	for p.currTok.Type == "PLUS" || p.currTok.Type == "MINUS" || p.currTok.Type == "AND" || p.currTok.Type == "OR" {
		op := p.currTok.Type
		p.advance()
		switch op {
		case "PLUS":
			result = result.(int) + p.term().(int)
		case "MINUS":
			result = result.(int) - p.term().(int)
		case "AND":
			result = result.(bool) && p.term().(bool)
		case "OR":
			result = result.(bool) || p.term().(bool)
		}
	}
	if p.calcMode {
		fmt.Println("Result: ", result)
	}
	return result
}

func (p *Parser) term() interface{} {
	result := p.factor()
	for p.currTok.Type == "MULTIPLY" || p.currTok.Type == "DIVIDE" {
		op := p.currTok.Type
		p.advance()
		switch op {
		case "MULTIPLY":
			result = result.(int) * p.factor().(int)
		case "DIVIDE":
			result = result.(int) / p.factor().(int)
		}
	}
	return result
}

func (p *Parser) factor() interface{} {
	var result interface{}
	switch p.currTok.Type {
	case "NUMBER":
		result = p.currTok.Value
		p.advance()
	case "STRING":
		result = p.currTok.Value
		p.advance()
	case "BOOLEAN":
		result = p.currTok.Value
		p.advance()
	case "IDENTIFIER":
		if value, found := p.variables[p.currTok.Value.(string)]; found {
			result = value
			p.advance()
		} else {
			panic(fmt.Sprintf("Undefined variable: %s", p.currTok.Value))
		}
	case "NOT":
		p.advance()
		result = !p.factor().(bool)
	case "LPAREN":
		p.advance()
		result = p.expr()
		if p.currTok.Type == "RPAREN" {
			p.advance()
		} else {
			panic("Expected ')'")
		}
	default:
		panic(fmt.Sprintf("Unexpected token: %v", p.currTok))
	}
	return result
}

func (p *Parser) print() {
	value := p.expr()
	fmt.Println(value)
}

func run() {
	calcMode := false
	debug := false
	variables := make(map[string]interface{})
	commands := []string{}

	reader := bufio.NewReader(os.Stdin)
	for {
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
			for _, cmd := range commands {
				lexer := NewLexer(cmd, debug)
				NewParser(lexer.tokens, calcMode, debug, variables)
			}
		case "exit":
			fmt.Println("exiting")
			return
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
