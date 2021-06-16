package evaluator

import (
	"fmt"
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
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		testingutils.Assert(t, ok, "no error object returned, got %T", evaluated)
		testingutils.Equals(t, tt.expectedMessage, errObj.Message, "Error message")
	}
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"999", 999},
		{"-999", -999},
		{"-10", -10},
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
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) {
	result, ok := obj.(*object.Integer)
	testingutils.Assert(t, ok, "obj is not %s, got %T (%+v)", object.INTEGER, obj, obj)
	testingutils.Equals(t, expected, result.Value, "result.Value")
}
