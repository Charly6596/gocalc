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
				Token: token.Token{Type: token.ASSIGNMENT, Literal: "x1"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "x1"},
					Value: "x1",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENTIFIER, Literal: "x2"},
					Value: "x2",
				},
			},
		},
	}

	expected := "x1 = x2;"
	actual := program.String()
	testingutils.Equals(t, expected, actual, "program.String()")
}
