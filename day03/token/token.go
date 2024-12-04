package token

const (
	EOF    = "EOF"
	INT    = "INT"
	MUL    = "MUL"
	LPAREN = "("
	RPAREN = ")"
	COMMA  = ","
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}
