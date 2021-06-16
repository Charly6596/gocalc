package parser

import (
	"fmt"
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

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input      string
		operator   string
		rightValue int64
	}{
		{"-15;", "-", 15},
	}
	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		assertNoParseErrors(t, p)
		testingutils.Equals(t, 1, len(program.Statements), "incorrect program.Statements lenght")

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		testingutils.Assert(t, ok, "s not *ast.ExpressionStatement. got=%T", program.Statements[0])

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		testingutils.Assert(t, ok, "stmt.Expression not *ast.PrefixExpression. got=%T", stmt.Expression)
		testingutils.Equals(t, tt.operator, exp.Operator, "exp.Operator")

		testIntegerLiteral(t, exp.Right, tt.rightValue)
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		operator   string
		leftValue  int64
		rightValue int64
	}{
		{"15 - 15;", "-", 15, 15},
		{"15 + 15;", "+", 15, 15},
		{"15 * 15;", "*", 15, 15},
		{"15 / 15;", "/", 15, 15},
		{"15 ^ 15;", "^", 15, 15},
	}
	for _, it := range infixTests {
		l := lexer.New(it.input)
		p := New(l)
		program := p.ParseProgram()
		assertNoParseErrors(t, p)
		testingutils.Equals(t, 1, len(program.Statements), "incorrect program.Statements lenght")

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		testingutils.Assert(t, ok, "s not *ast.ExpressionStatement. got=%T", program.Statements[0])

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		testingutils.Assert(t, ok, "stmt.Expression not *ast.PrefixExpression. got=%T", stmt.Expression)
		testingutils.Equals(t, it.operator, exp.Operator, "exp.Operator")
		testIntegerLiteral(t, exp.Right, it.rightValue)
		testIntegerLiteral(t, exp.Left, it.leftValue)
	}
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, value int64) {
	integ, ok := exp.(*ast.IntegerLiteral)
	testingutils.Assert(t, ok, "exp not *ast.IntegerLiteral. got=%T", exp)
	testingutils.Equals(t, value, integ.Value, "integ.Value")
	testingutils.Equals(t, fmt.Sprintf("%d", value), integ.TokenLiteral(), "integ.TokenLiteral()")
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		assertNoParseErrors(t, p)
		actual := program.String()
		testingutils.Equals(t, tt.expected, actual, "program.String()")
	}
}
