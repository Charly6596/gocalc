package evaluator

import (
	"fmt"
	"gocalc/object"
	"gocalc/testing_utils"
	"testing"
)

func TestStdlib(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"typeof(5 / 0)",
			object.ERROR.String(),
		},
		{
			"typeof(5)",
			object.FLOAT.String(),
		},
		{
			"len([])",
			"0",
		},
		{
			"len([1, true])",
			"2",
		},
		{
			"head([1, true])",
			"1",
		},
		{
			"tail([1, true])",
			"True",
		},
		{
			"get([1, true], 1)",
			"True",
		},
		{
			"get([1, true], 0)",
			"1",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testingutils.Equals(t, tt.expectedMessage, evaluated.String(), "response")
	}
}

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
		ev := New()
		ev.Eval(tt.input)
		res := ev.Eval(ANS)
		testFloatObject(t, res, tt.expected)
	}

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

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"!true", false},
		{"!false", true},
		{"true == true", true},
		{"true != true", false},
		{"false == false", true},
		{"false != false", false},
		{"true == false", false},
		{"10 > 10", false},
		{"10 >= 10", true},
		{"10 <= 10", true},
		{"10 == 10", true},
		{"10 < 10", false},
		{"9 < 10", true},
		{"9 <= 10", true},
		{"9 >= 10", false},
		{"9 > 10", false},
		{"!(9 > 10)", true},
		{"!(9 <= 10)", false},
		{"!(true && true)", false},
		{"true && true", true},
		{"true || false", true},
		{"(3 > 5) || (5 < 10)", true},
		{"(3 < 5) && (5 < 10)", true},
		{"true && false || true", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) {
	result, ok := obj.(*object.Boolean)
	testingutils.Assert(t, ok, "obj is not %s, got %T (%+v)", object.BOOLEAN, obj, obj)
	testingutils.Equals(t, expected, result.Value, "result.Value")
}

func TestEvalListExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"[5]", []string{"5"}},
		{"[true, false, 9]", []string{"True", "False", "9"}},
		{"[5 + 5 + 5 + 5 - 10, -5]", []string{"10", "-5"}},
		{"[true, 2 * 2 * 2 * 2 * 2]", []string{"True", "32"}},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testListObject(t, evaluated, tt.expected)
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
	ev := New()
	res := ev.Eval(input)
	return res
}

func testFloatObject(t *testing.T, obj object.Object, expected float64) {
	result, ok := obj.(*object.Float)
	testingutils.Assert(t, ok, "obj is not %s, got %T (%+v)", object.FLOAT, obj, obj)
	testingutils.Equals(t, expected, result.Value, "result.Value")
}

func testListObject(t *testing.T, obj object.Object, expected []string) {
	result, ok := obj.(*object.List)
	testingutils.Assert(t, ok, "obj is not %s, got %T (%+v)", object.LIST, obj, obj)
	testingutils.Equals(t, len(expected), len(result.Values), "len(result.Values)")
	for i, o := range result.Values {
		testingutils.Equals(t, expected[i], o.String(), "result.Values["+fmt.Sprint(i)+"]")
	}
}
