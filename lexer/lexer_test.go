package lexer

import (
	"gocalc/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
x = 3 + 5^3; (3 * 8 ); x;
    `
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENTIFIER, "x"},
		{token.ASSIGN, "="},
		{token.INT, "3"},
		{token.PLUS, "+"},
		{token.INT, "5"},
		{token.CARET, "^"},
		{token.INT, "3"},
		{token.SEMICOLON, ";"},
		{token.LPAREN, "("},
		{token.INT, "3"},
		{token.ASTERISK, "*"},
		{token.INT, "8"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.IDENTIFIER, "x"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
