package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	INT     = "INT"
	MUL     = "MUL"
	DO      = "DO"
	DONT    = "DON'T"
	LPAREN  = "("
	RPAREN  = ")"
	COMMA   = ","
	IDENT   = "IDENT"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"mul":     MUL,
	"do()":    DO,
	"don't()": DONT,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
