package lexer

import (
	"gocalc/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
x1 = 3 + 5^3; (3 * 8 ); x1; 0.5; .5 > @; 
    abc != true   ; true && false || false;
    [true, false]
    `
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "x1"},
		{token.ASSIGN, "="},
		{token.FLOAT, "3"},
		{token.PLUS, "+"},
		{token.FLOAT, "5"},
		{token.CARET, "^"},
		{token.FLOAT, "3"},
		{token.SEMICOLON, ";"},
		{token.LPAREN, "("},
		{token.FLOAT, "3"},
		{token.ASTERISK, "*"},
		{token.FLOAT, "8"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "x1"},
		{token.SEMICOLON, ";"},
		{token.FLOAT, "0.5"},
		{token.SEMICOLON, ";"},
		{token.FLOAT, "0.5"},
		{token.GT, ">"},
		{token.ILLEGAL, "@"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "abc"},
		{token.NOT_EQ, "!="},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.TRUE, "true"},
		{token.AND, "&&"},
		{token.FALSE, "false"},
		{token.OR, "||"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.LBRACK, "["},
		{token.TRUE, "true"},
		{token.COMMA, ","},
		{token.FALSE, "false"},
		{token.RBRACK, "]"},
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
