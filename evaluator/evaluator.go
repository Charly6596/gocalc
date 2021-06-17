package evaluator

import (
	"fmt"
	"gocalc/ast"
	"gocalc/object"
	"math"
)

var (
	NULL = &object.Null{}
	ANS  = "ans"
)

type Environment struct {
	store map[string]object.Object
}

func NewEnvironment() *Environment {
	s := make(map[string]object.Object)
	return &Environment{store: s}
}

func (e *Environment) Get(name string) (object.Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}
func (e *Environment) Set(name string, val object.Object) object.Object {
	e.store[name] = val
	return val
}

func (e *Environment) String() string {
	return fmt.Sprint(e.store)
}

func (e *Environment) Eval(program ast.Node) object.Object {
	res := eval(program, e)
	if !isError(res) && res != nil {
		e.Set(ANS, res)
	}

	return res
}

func eval(node ast.Node, env *Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.AssignmentStatement:
		val := eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return nil
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.ExpressionStatement:
		return eval(node.Expression, env)
	case *ast.PrefixExpression:
		r := eval(node.Right, env)
		if isError(r) {
			return r
		}

		return evalPrefixExpression(node.Operator, r)
	case *ast.InfixExpression:
		r := eval(node.Right, env)
		if isError(r) {
			return r
		}

		l := eval(node.Left, env)
		if isError(l) {
			return l
		}

		return evalInfixExpression(node.Operator, l, r)
	case *ast.FloatLiteral:
		return &object.Float{Value: node.Value}
	}
	return NULL
}

func evalIdentifier(node *ast.Identifier, env *Environment) object.Object {
	val, ok := env.Get(node.Value)

	if !ok {
		return newError(object.IDENTIFIER_NOT_FOUND_ERROR, node.Value)
	}

	return val
}

func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ERROR
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case isFloat(left) && isFloat(right):
		return evalInfixExpressionFloat(operator, left, right)
	default:
		return newError(object.UNKNOWN_INFIX_OPERATOR_ERROR, left.Type(), operator, right.Type())
	}
}

func isFloat(obj object.Object) bool {
	return obj.Type() == object.FLOAT
}

func evalInfixExpressionFloat(operator string, left, right object.Object) object.Object {
	x1, x2 := left.(*object.Float).Value, right.(*object.Float).Value

	switch operator {
	case "+":
		return &object.Float{Value: x1 + x2}
	case "-":
		return &object.Float{Value: x1 - x2}
	case "*":
		return &object.Float{Value: x1 * x2}
	case "^":
		return &object.Float{Value: math.Pow(x1, x2)}
	case "/":
		if x2 == 0 {
			return object.DivideByZeroError(left, right)
		}

		return &object.Float{Value: x1 / x2}
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
	if num, ok := object.ToFloat(right); ok {
		return &object.Float{Value: -num.Value}
	}

	return newError(object.UNKNOWN_PREFIX_OPERATOR_ERROR, "-", right.Type())
}

func evalProgram(program *ast.Program, env *Environment) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = eval(statement, env)
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
