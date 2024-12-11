package lexer

import (
	"aoc2024/day03/token"
	"strings"
)

type Lexer struct {
	input                  string
	Position, readPosition int
	ch                     byte
}

func New(line string) *Lexer {
	l := &Lexer{input: line}
	l.readChar()
	return l
}

// Lexing the charstream to find actual tokens.
// (The emitted tokens are parsed later according to a grammar.)
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	// NOTE: Not needed here...since all other chars are skipable
	// l.skipWhitespace()

	switch l.ch {
	case 0:
		tok = newToken(token.EOF, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	default:
		if tok, ok := l.checkDo(); ok {
			return tok
		}
		if tok, ok := l.checkDont(); ok {
			return tok
		}
		if tok, ok := l.checkMul(); ok {
			return tok
		}
		if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		}
		// NOTE: Not needed here...since we only look for explicit tokens
		// if isLetter(l.ch) {
		// 	tok.Literal = l.readIdentifier()
		// 	tok.Type = token.LookupIdent(tok.Literal)
		// 	return tok
		// }
	}

	l.readChar()
	return tok
}

func (l *Lexer) checkMul() (token.Token, bool) {
	if strings.HasPrefix(l.input[l.Position:], "mul") {
		tok := token.Token{Type: token.MUL, Literal: "mul"}
		l.eatNchars(len(tok.Literal))
		return tok, true
	}
	return token.Token{}, false
}

func (l *Lexer) checkDo() (token.Token, bool) {
	if strings.HasPrefix(l.input[l.Position:], "do()") {
		tok := token.Token{Type: token.DO, Literal: "do()"}
		l.eatNchars(len(tok.Literal))
		return tok, true
	}
	return token.Token{}, false
}
func (l *Lexer) checkDont() (token.Token, bool) {
	if strings.HasPrefix(l.input[l.Position:], "don't()") {
		tok := token.Token{Type: token.DONT, Literal: "don't()"}
		l.eatNchars(len(tok.Literal))
		return tok, true
	}
	return token.Token{}, false
}

// helper methods

func (l *Lexer) readIdentifier() string {
	position := l.Position
	for isLetter(l.ch) { //&& l.ch != 'm' && l.ch != 'd' {
		l.readChar()
	}
	return l.input[position:l.Position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
	position := l.Position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.Position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// Like "eat", "consume", "skip"
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.Position = l.readPosition
	l.readPosition++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// INFO: New. Like inverse of readChar
func (l *Lexer) backup() {
	l.readPosition = l.Position
	l.Position--
	l.ch = l.input[l.readPosition]
}

// INFO: New. Like multiple readChar:s
func (l *Lexer) eatNchars(n int) {
	for i := 0; i < n; i++ {
		l.readChar()
	}
}
