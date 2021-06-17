package ast

import (
	testingutils "gocalc/testing_utils"
	"gocalc/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&AssignmentStatement{
				Token: token.Token{Type: token.ASSIGN, Literal: "x1"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "x1"},
					Value: "x1",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "x2"},
					Value: "x2",
				},
			},
			&ExpressionStatement{
				Token: token.Token{Type: token.FLOAT, Literal: "0.5"},
				Expression: &FloatLiteral{
					Token: token.Token{Type: token.FLOAT, Literal: "0.5"},
					Value: 0.5,
				},
			},
		},
	}

	expected := "x1 = x2;0.5"
	actual := program.String()
	testingutils.Equals(t, expected, actual, "program.String()")
}
