package token

type TokenType byte

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL TokenType = iota
	EOF

	literal_beg
	IDENT  // x, x2, y
	FLOAT  // 10.14
	IMAG   // 10.14i
	CHAR   // 'a'
	STRING // "abc"
	literal_end

	operator_beg
	SEMICOLON
	COMMA
	PERIOD
	LPAREN
	RPAREN
	ASSIGN
	PLUS
	MINUS
	ASTERISK
	SLASH
	CARET
	operator_end

	keyword_beg
	IMPORT
	TYPE
	keyword_end
)

var tokenNames = []string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	// Literals
	IDENT:  "IDENT",
	FLOAT:  "FLOAT",
	IMAG:   "IMAG",
	CHAR:   "CHAR",
	STRING: "STRING",

	// Delimiters
	SEMICOLON: ";",
	COMMA:     ",",
	PERIOD:    ".",
	LPAREN:    "(",
	RPAREN:    ")",

	// Operators
	ASSIGN:   "=",
	PLUS:     "+",
	MINUS:    "-",
	ASTERISK: "*",
	SLASH:    "/",
	CARET:    "^",

	// Keywords
	IMPORT: "import",
	TYPE:   "type",
}

func (tt TokenType) String() string { return tokenNames[tt] }

func New(tokenType TokenType, ch byte) Token { return Token{Type: tokenType, Literal: string(ch)} }

func (t *Token) IsIllegal() bool { return t.Type == ILLEGAL }
