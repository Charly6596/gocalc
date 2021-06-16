package evaluator

import (
	"fmt"
	"gocalc/ast"
	"gocalc/object"
)

var (
	NULL = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.PrefixExpression:
		r := Eval(node.Right)
		if isError(r) {
			return r
		}

		return evalPrefixExpression(node.Operator, r)
	case *ast.InfixExpression:
		r := Eval(node.Right)
		if isError(r) {
			return r
		}

		l := Eval(node.Left)
		if isError(l) {
			return l
		}

		return evalInfixExpression(node.Operator, l, r)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}

	return NULL
}

func isError(obj object.Object) bool {
	return obj.Type() == object.ERROR
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case isInteger(left) && isInteger(right):
		return evalInfixExpressionInteger(operator, left, right)
	default:
		return newError(object.UNKNOWN_INFIX_OPERATOR_ERROR, left.Type(), operator, right.Type())
	}
}

func isInteger(obj object.Object) bool {
	return obj.Type() == object.INTEGER
}

func evalInfixExpressionInteger(operator string, left, right object.Object) object.Object {
	x1, x2 := left.(*object.Integer).Value, right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: x1 + x2}
	case "-":
		return &object.Integer{Value: x1 - x2}
	case "*":
		return &object.Integer{Value: x1 * x2}
	case "/":
		if x2 == 0 {
			return object.DivideByZeroError(x1, x2)
		}

		return &object.Integer{Value: x1 / x2}
	}

	return newError(object.UNKNOWN_INFIX_OPERATOR_ERROR, left.Type(), operator, right.Type())

}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	if operator == "-" {
		return evalMinusOperator(right)
	}

	return newError(object.UNKNOWN_PREFIX_OPERATOR_ERROR, operator, right.Type())
}

func evalMinusOperator(right object.Object) object.Object {
	if num, ok := object.ToInteger(right); ok {
		return &object.Integer{Value: -num.Value}
	}

	return newError(object.UNKNOWN_PREFIX_OPERATOR_ERROR, "-", right.Type())
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = Eval(statement)
		switch result := result.(type) {
		case *object.Error:
			return result
		}
	}
	return result
}

func newError(msg string, v ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(msg, v...)}
}
