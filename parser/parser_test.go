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

	expectedLen := 1

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

func TestIdentifierExpression(t *testing.T) {
	input := "myVar;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	assertNoParseErrors(t, p)
	testingutils.Equals(t, 1, len(program.Statements), "len(program.Statements)")

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	testingutils.Assert(t, ok, "program.Statements[0] not ast.ExpressionStatement. got=%T", program.Statements[0])

	ident, ok := stmt.Expression.(*ast.Identifier)
	exp := "myVar"
	testingutils.Assert(t, ok, "smt not *ast.Identifier. got=%T", stmt.Expression)
	testingutils.Equals(t, exp, ident.Value, "ident.Value")
	testingutils.Equals(t, exp, ident.TokenLiteral(), "ident.TokenLiteral()")
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "10;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	assertNoParseErrors(t, p)
	testingutils.Equals(t, 1, len(program.Statements), "len(program.Statements)")

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	testingutils.Assert(t, ok, "program.Statements[0] not ast.ExpressionStatement. got=%T", program.Statements[0])

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	exp := int64(10)
	testingutils.Assert(t, ok, "stmt not *ast.IntegerLiteral. got=%T", stmt.Expression)
	testingutils.Equals(t, exp, literal.Value, "literal.Value")

	exp2 := "10"
	testingutils.Equals(t, exp2, literal.TokenLiteral(), "literal.TokenLiteral()")
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
