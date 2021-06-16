package lexer

import (
	"gocalc/token"
)

type Lexer struct {
	input        string
	position     int
	nextPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.nextPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.nextPosition]
	}

	l.position = l.nextPosition
	l.nextPosition++
}

func (l *Lexer) NextToken() token.Token {

	l.eatWhitespaces()

	if isDigit(l.ch) {
		return token.Token{Type: token.INT, Literal: l.readWhile(isDigit)}
	}

	if isLetter(l.ch) {
		return token.Token{Type: token.IDENTIFIER, Literal: l.readWhile(isLetter)}
	}

	defer l.readChar()

	switch l.ch {
	case '=':
		return token.New(token.ASSIGN, l.ch)
	case ';':
		return token.New(token.SEMICOLON, l.ch)
	case '(':
		return token.New(token.LPAREN, l.ch)
	case ')':
		return token.New(token.RPAREN, l.ch)
	case ',':
		return token.New(token.COMMA, l.ch)
	case '+':
		return token.New(token.PLUS, l.ch)
	case '-':
		return token.New(token.MINUS, l.ch)
	case '*':
		return token.New(token.ASTERISK, l.ch)
	case '/':
		return token.New(token.SLASH, l.ch)
	case '^':
		return token.New(token.CARET, l.ch)
	case 0:
		return token.Token{Type: token.EOF, Literal: ""}
	}

	return token.New(token.ILLEGAL, l.ch)
}

type bytePredicate func(ch byte) bool

func (l *Lexer) readWhile(pred bytePredicate) string {
	start := l.position
	for pred(l.ch) {
		l.readChar()
	}

	return l.input[start:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || isDigit(ch)
}

func (l *Lexer) eatWhitespaces() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\r' || l.ch == '\t' {
		l.readChar()
	}
}
