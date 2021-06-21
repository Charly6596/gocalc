package parser

import (
	"fmt"
	"gocalc/ast"
	"gocalc/lexer"
	"gocalc/testing_utils"
	"testing"
)

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	assertNoParseErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	testingutils.Assert(t, ok, "program.Statements[0] not ast.ExpressionStatement. got=%T", program.Statements[0])

	exp, ok := stmt.Expression.(*ast.CallExpression)
	testingutils.Assert(t, ok, "stmt.Expression not *ast.PrefixExpression. got=%T", stmt.Expression)

	testIdentifier(t, exp.Function, "add")
	testingutils.Equals(t, 3, len(exp.Arguments), "exp.Arguments")
	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}
func TestErrorMessages(t *testing.T) {
	input := `
x = 15;
y =!= 20;
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

	exp := "myVar"
	testIdentifier(t, stmt.Expression, exp)
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

func TestBooleanLiteralExpression(t *testing.T) {
	input := "true;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	assertNoParseErrors(t, p)
	testingutils.Equals(t, 1, len(program.Statements), "len(program.Statements)")

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	testingutils.Assert(t, ok, "program.Statements[0] not ast.ExpressionStatement. got=%T", program.Statements[0])

	literal, ok := stmt.Expression.(*ast.BooleanLiteral)
	exp := true
	testingutils.Assert(t, ok, "stmt not *ast.BooleanLiteral. got=%T", stmt.Expression)
	testingutils.Equals(t, exp, literal.Value, "literal.Value")

	exp2 := "true"
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
		{"15 && 15;", "&&", 15, 15},
		{"15 || 15;", "||", 15, 15},
		{"15 >= 15;", ">=", 15, 15},
		{"15 > 15;", ">", 15, 15},
		{"15 <= 15;", "<=", 15, 15},
		{"15 < 15;", "<", 15, 15},
		{"15 == 15;", "==", 15, 15},
		{"15 != 15;", "!=", 15, 15},
	}
	for _, it := range infixTests {
		l := lexer.New(it.input)
		p := New(l)
		program := p.ParseProgram()
		assertNoParseErrors(t, p)
		testingutils.Equals(t, 1, len(program.Statements), "incorrect program.Statements lenght")

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		testingutils.Assert(t, ok, "s not *ast.ExpressionStatement. got=%T", program.Statements[0])

		testInfixExpression(t, stmt.Expression, it.leftValue, it.operator, it.rightValue)
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
			"!true == false",
			"((!true) == false)",
		},
		{
			"3 >= 5",
			"(3 >= 5)",
		},
		{
			"3 > 5",
			"(3 > 5)",
		},
		{
			"3 < 5",
			"(3 < 5)",
		},
		{
			"3 <= 5",
			"(3 <= 5)",
		},
		{
			"3 == 5",
			"(3 == 5)",
		},
		{
			"3 != 5",
			"(3 != 5)",
		},
		{
			"!(3 != 5)",
			"(!(3 != 5))",
		},
		{
			"3.05 + 4; -5.3 * 5.33",
			"(3.05 + 4)((-5.3) * 5.33)",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
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

func testIdentifier(t *testing.T, expression ast.Expression, expected string) {
	ident, ok := expression.(*ast.Identifier)
	testingutils.Assert(t, ok, "smt not *ast.Identifier. got=%T", expression)
	testingutils.Equals(t, expected, ident.Value, "ident.Value")
	testingutils.Equals(t, expected, ident.TokenLiteral(), "ident.TokenLiteral()")
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) {
	switch v := expected.(type) {
	case int:
	case int64:
	case float32:
		testFloatLiteral(t, exp, float64(v))
	case float64:
		testFloatLiteral(t, exp, v)
	case string:
		testIdentifier(t, exp, v)
	default:
		t.Errorf("type of exp not handled. got=%T", expected)
	}
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) {

	res, ok := exp.(*ast.InfixExpression)
	testingutils.Assert(t, ok, "stmt.Expression not *ast.PrefixExpression. got=%T", exp)
	testingutils.Equals(t, operator, res.Operator, "exp.Operator")
	testLiteralExpression(t, res.Right, right)
	testLiteralExpression(t, res.Left, left)
}
