package main

import (
	"fmt"
)

type Parser struct {
	tokens    []Token
	idx       int
	currTok   Token
	calcMode  bool
	debug     bool
	variables map[string]interface{}
}

func NewParser(tokens []Token, calcMode, debug bool, variables map[string]interface{}) *Parser {
	p := &Parser{
		tokens:    tokens,
		idx:       -1,
		calcMode:  calcMode,
		debug:     debug,
		variables: variables,
	}
	p.advance()
	p.parse()
	return p
}

func (p *Parser) advance() {
	p.idx++
	if p.idx < len(p.tokens) {
		p.currTok = p.tokens[p.idx]
	} else {
		p.currTok = Token{}
	}
	if p.debug {
		fmt.Printf("Parser advance: idx=%d, curr_tok=%v\n", p.idx, p.currTok)
	}
}

func (p *Parser) parse() {
	for p.currTok.Type != "" {
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
		fmt.Println("Result:", result)
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
		if val, found := p.variables[p.currTok.Value.(string)]; found {
			result = val
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
