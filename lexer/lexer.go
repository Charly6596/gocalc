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
	l.ch = l.peekChar()
	l.advanceChar()
}

func (l *Lexer) advanceChar() {
	l.position = l.nextPosition
	l.nextPosition++
}

func (l *Lexer) peekChar() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.nextPosition]
	}
}

func (l *Lexer) NextToken() token.Token {

	l.eatWhitespaces()

	if isDigit(l.ch) {
		res := l.readWhile(isDigit)
		if res[0] == '.' {
			res = "0" + res
		}
		return token.NewExt(token.FLOAT, res)
	}

	if isLetter(l.ch) {
		t := token.Token{Literal: l.readWhile(isLetter)}
		if kw, ok := token.TryGetKeyword(t.Literal); ok {
			t.Type = kw
		} else {
			t.Type = token.IDENT
		}

		return t
	}

	defer l.readChar()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.advanceChar()
			return token.NewExt(token.EQ, "==")
		}
		return token.New(token.ASSIGN, l.ch)
	case '!':
		if l.peekChar() == '=' {
			l.advanceChar()
			return token.NewExt(token.NOT_EQ, "!=")
		}
		return token.New(token.BANG, l.ch)
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
	case '>':
		if l.peekChar() == '=' {
			l.advanceChar()
			return token.NewExt(token.GT_EQ, ">=")
		}
		return token.New(token.GT, l.ch)
	case '<':
		if l.peekChar() == '=' {
			l.advanceChar()
			return token.NewExt(token.LT_EQ, "<=")
		}
		return token.New(token.LT, l.ch)

	case 0:
		return token.NewExt(token.EOF, "")
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
	return '0' <= ch && ch <= '9' || ch == '.'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || isDigit(ch)
}

func (l *Lexer) eatWhitespaces() {
	for l.ch == ' ' || l.ch == '\n' || l.ch == '\r' || l.ch == '\t' {
		l.readChar()
	}
}
