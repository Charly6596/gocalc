package token

type TokenType byte

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL TokenType = iota
	EOF

	IDENTIFIER
	ASSIGNMENT

	INT

	SEMICOLON
	COMMA
	LPAREN
	RPAREN

	ASSIGN
	PLUS
	MINUS
	ASTERISK
	SLASH
	CARET
)

var tokenNames = []string{
	"ILLEGAL", //	ILLEGAL
	"EOF",     //	EOF

	"IDENTIFIER", // IDENTIFIER
	"ASSIGNMENT", // ASSIGNMENT

	"INT", // INT

	/* DELIMITERS */
	";", // SEMICOLON
	",", //	COMMA
	"(", //	LPAREN
	")", //	RPAREN

	/* OPERATORS */
	"=", //	ASSIGN
	"+", //	PLUS
	"-", //	MINUS
	"*", //	ASTERISK
	"/", //	SLASH
	"^", //	CARET
}

func (tt TokenType) String() string {
	return tokenNames[tt]
}

func New(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}
