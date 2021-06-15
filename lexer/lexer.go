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
	case 0:
		return token.Token{Type: token.EOF, Literal: ""}
	}

	return token.New(token.ILLEGAL, l.ch)
}
