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
x= 0.5;
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()

	testingutils.Assert(t, program != nil, "ParseProgram() returned nil")
	assertNoParseErrors(t, p)
	testingutils.Assert(t, len(program.Statements) == 4, "program.Statements does not contain 3 statement. got=%d", len(program.Statements))

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"x123"},
		{"x"},
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

func TestFloatLiteralExpression(t *testing.T) {
	input := "0.1;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	assertNoParseErrors(t, p)
	testingutils.Equals(t, 1, len(program.Statements), "len(program.Statements)")

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	testingutils.Assert(t, ok, "program.Statements[0] not ast.ExpressionStatement. got=%T", program.Statements[0])

	literal, ok := stmt.Expression.(*ast.FloatLiteral)
	exp := float64(0.1)
	testingutils.Assert(t, ok, "stmt not *ast.FloatLiteral. got=%T", stmt.Expression)
	testingutils.Equals(t, exp, literal.Value, "literal.Value")

	exp2 := "0.1"
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
		rightValue float64
	}{
		{"-15;", "-", 15},
		{"-.5;", "-", 0.5},
		{"-0.5;", "-", 0.5},
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

		testFloatLiteral(t, exp.Right, tt.rightValue)
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		operator   string
		leftValue  float64
		rightValue float64
	}{
		{"0.2 - 15;", "-", 0.2, 15},
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
		testFloatLiteral(t, exp.Right, it.rightValue)
		testFloatLiteral(t, exp.Left, it.leftValue)
	}
}

func testFloatLiteral(t *testing.T, exp ast.Expression, value float64) {
	num, ok := exp.(*ast.FloatLiteral)
	testingutils.Assert(t, ok, "exp not *ast.FloatLiteral. got=%T", exp)
	testingutils.Equals(t, value, num.Value, "num.Value")
	testingutils.Equals(t, fmt.Sprint(value), num.TokenLiteral(), "num.TokenLiteral()")
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
		{
			"3.05 + 4; -5.3 * 5.33",
			"(3.05 + 4)((-5.3) * 5.33)",
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
