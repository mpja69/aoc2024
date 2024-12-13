package main

import "strings"

type Parser struct {
	input               string
	pos, readPos, start int
	ch                  byte
	tokens              chan Token
}

type Token struct {
	Type tokenType
	val  string
}

const (
	EOF = iota
	TEXT
	MUL
	DO
	DONT
)

type tokenType int

type stateFn func(p *Parser) stateFn

func ParserNew(str string) (*Parser, chan Token) {
	p := &Parser{
		input:  str,
		tokens: make(chan Token),
	}
	p.run()
	return p, p.tokens
}

func (p *Parser) emit(t tokenType) {
	p.tokens <- Token{t, p.input[p.start:p.pos]}
	p.start = p.pos
}

func (p *Parser) next() byte {
	if p.pos >= len(p.input) {
		return EOF
	}
	ch := p.input[p.pos]
	p.pos++
	return ch
}

func (p *Parser) peek() byte {
	ch := p.next()
	p.backup()
	return ch
}

func (p *Parser) backup() {
	p.pos--
}

func (p *Parser) ignore() {
	p.start = p.pos
}

func (p *Parser) accept(valid string) bool {
	if strings.IndexByte(valid, p.next()) >= 0 {
		return true
	}
	p.backup()
	return false
}

func (p *Parser) acceptRun(valid string) {
	for strings.IndexByte(valid, p.next()) >= 0 {
	}
	p.backup()
}

func (p *Parser) run() {
	for state := startState; state != nil; {
		state = state(p)
	}
	close(p.tokens)
}

func startState(p *Parser) stateFn {
	for {
		if strings.HasPrefix(p.input[p.start:], "mul") {
			if p.pos > p.start {
				p.emit(TEXT)
			}
			return mulState
		}
		if strings.HasPrefix(p.input[p.start:], "do()") {
			if p.pos > p.start {
				p.emit(TEXT)
			}
			return doState
		}
		if strings.HasPrefix(p.input[p.start:], "don't()") {
			if p.pos > p.start {
				p.emit(TEXT)
			}
			return dontState
		}
		if p.next() == EOF {
			break
		}
	}
	if p.pos > p.start {
		p.emit(TEXT)
	}
	p.emit(EOF)
	return nil
}

func mulState(p *Parser) stateFn {
	p.pos += len("mul")
	if p.next() != '(' {
		return startState
	}
	ch := p.next()
	if ch < '0' || ch > '9' {
		return startState
	}
	// parseNumber(p)???
	return startState
}
func doState(p *Parser) stateFn {
	return startState
}
func dontState(p *Parser) stateFn {
	return startState
}
