package parser

import (
	"gocalc/ast"
	"gocalc/lexer"
	"gocalc/testing_utils"
	"testing"
)

func TestErrorMessages(t *testing.T) {
	input := `
x = 15;
y == 20;
`

	l := lexer.New(input)
	p := New(l)
	p.ParseProgram()

	expectedLen := 2

	testingutils.Equals(t, expectedLen, len(p.errors), "parser.errors")
}

func TestAssignmentStatements(t *testing.T) {
	input := `
x = 15;
y = 20;
x123 = 99999;
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	testingutils.Assert(t, program != nil, "ParseProgram() returned nil")
	assertNoParseErrors(t, p)
	testingutils.Assert(t, len(program.Statements) == 3, "program.Statements does not contain 3 statement. got=%d", len(program.Statements))

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"x123"},
	}

	for i, tt := range tests {
		s := program.Statements[i]
		name := tt.expectedIdentifier

		testingutils.Equals(t, s.TokenLiteral(), name, "s.TokenLiteral()")
		assignment, ok := s.(*ast.AssignmentStatement)
		testingutils.Assert(t, ok, "s not *ast.AssignmentExpression. got=%T", s)
		testingutils.Equals(t, name, assignment.Name.Value, "s.Name.Value")
		testingutils.Equals(t, name, assignment.Name.TokenLiteral(), "s.Name")
	}
}

func assertNoParseErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
