package evaluator

import (
	"fmt"
	"gocalc/ast"
	"gocalc/lexer"
	"gocalc/object"
	"gocalc/parser"
	"gocalc/testing_utils"
	"testing"
)

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 / 0",
			fmt.Sprintf(object.DIVIDE_BY_ZERO, 5, 0),
		},
		{
			"a",
			fmt.Sprintf(object.IDENTIFIER_NOT_FOUND_ERROR, "a"),
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		testingutils.Assert(t, ok, "no error object returned, got %T", evaluated)
		testingutils.Equals(t, tt.expectedMessage, errObj.Message, "Error message")
	}
}

func TestAnsExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"a = 5; a;", 5},
		{"a = 5 * 5; a;", 25},
		{"a = 5; b = a; b;", 5},
		{"a = 5; b = a; c = a + b + 5; c;", 15},
	}
	for _, tt := range tests {
		env := NewEnvironment()
		program := parseInput(tt.input)
		env.Eval(program)
		program = parseInput(ANS)
		res := env.Eval(program)

		testFloatObject(t, res, tt.expected)
	}

}

func parseInput(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func TestAssignmentStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"a = 5; a;", 5},
		{"a = 5 * 5; a;", 25},
		{"a = 5.5; b = a; b;", 5.5},
		{"a = 5.5; b = a; c = a + b + 5; c;", 16},
	}
	for _, tt := range tests {
		testFloatObject(t, testEval(tt.input), tt.expected)
	}
}

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5", 5},
		{"10.5", 10.5},
		{"999", 999},
		{"-999", -999},
		{"-10.5", -10.5},
		{"-5", -5},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10.5", 37.5},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := NewEnvironment()
	return env.Eval(program)
}

func testFloatObject(t *testing.T, obj object.Object, expected float64) {
	result, ok := obj.(*object.Float)
	testingutils.Assert(t, ok, "obj is not %s, got %T (%+v)", object.FLOAT, obj, obj)
	testingutils.Equals(t, expected, result.Value, "result.Value")
}
