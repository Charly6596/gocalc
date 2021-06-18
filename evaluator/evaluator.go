package evaluator

import (
	"fmt"
	"gocalc/ast"
	"gocalc/environment"
	"gocalc/lexer"
	"gocalc/object"
	"gocalc/parser"
	"math"
	"strings"
)

var (
	NULL = &object.Null{}
	ANS  = "ans"
)

type Evaluator struct {
	global *environment.Environment
	lexer  *lexer.Lexer
	parser *parser.Parser
}

// TODO: Libraries
var (
	n_typeof  = &object.NativeFunction{Function: nativeTypeof, Name: "typeof"}
	n_typeofS = &object.NativeFunction{Function: nativeTypeofS, Name: "typeofS"}
)

func New() *Evaluator {
	ev := &Evaluator{}
	ev.global = environment.New()

	// TODO: Libraries
	ev.global.Set("typeof", n_typeof)
	ev.global.Set("typeofS", n_typeofS)
	return ev
}

// TODO: Libraries, return multiple values
func nativeTypeofS(objs ...object.Object) object.Object {
	if len(objs) == 0 {
		return &object.String{Value: object.NATIVE_FUNCTION.Stringf("typeofS")}
	}
	obj := objs[0]
	return &object.String{Value: obj.TypeS()}
}

// TODO: Libraries, return multiple values
func nativeTypeof(objs ...object.Object) object.Object {
	if len(objs) == 0 {
		return &object.Type{Value: object.NATIVE_FUNCTION}
	}
	obj := objs[0]
	return &object.Type{Value: obj.Type()}
}

func (ev *Evaluator) Eval(input string) object.Object {
	lexer := lexer.New(input)
	parser := parser.New(lexer)
	program := parser.ParseProgram()

	if parser.HasErrors() {
		errs := strings.Join(parser.Errors(), "\n\t\t")
		return newError(object.SYNTAX_ERROR, errs)
	}

	res := ev.Program(program)
	if !isError(res) {
		ev.global.Set(ANS, res)
	}
	return res
}

func (ev *Evaluator) Program(program *ast.Program) object.Object {
	var result object.Object
	for _, statement := range program.Statements {
		result = ev.evaluate(statement)
		switch result := result.(type) {
		case *object.Error:
			return result
		}
	}
	return result
}

func (ev *Evaluator) Identifier(id *ast.Identifier) object.Object {
	val, ok := ev.global.Get(id.Value)

	if !ok {
		return newError(object.IDENTIFIER_NOT_FOUND_ERROR, id.Value)
	}

	return val
}

func (ev *Evaluator) FloatLiteral(fl *ast.FloatLiteral) object.Object {
	return &object.Float{Value: fl.Value}
}

func (ev *Evaluator) AssignmentStatement(as *ast.AssignmentStatement) object.Object {
	val := ev.evaluate(as.Value)
	if isError(val) {
		return val
	}

	ev.global.Set(as.Name.Value, val)

	return NULL
}

func (ev *Evaluator) ExpressionStatement(es *ast.ExpressionStatement) object.Object {
	return ev.evaluate(es.Expression)
}

func (ev *Evaluator) PrefixExpression(pe *ast.PrefixExpression) object.Object {
	r := ev.evaluate(pe.Right)

	if isError(r) {
		return r
	}

	return evalPrefixExpression(pe.Operator, r)

}

func (ev *Evaluator) InfixExpression(ie *ast.InfixExpression) object.Object {
	r := ev.evaluate(ie.Right)

	if isError(r) {
		return r
	}

	l := ev.evaluate(ie.Left)

	if isError(l) {
		return l
	}

	return evalInfixExpression(ie.Operator, l, r)
}

func (ev *Evaluator) CallExpression(ce *ast.CallExpression) object.Object {
	val, ok := ev.global.Get(ce.Function.TokenLiteral())

	if !ok {
		return newError(object.IDENTIFIER_NOT_FOUND_ERROR, ce.Function.TokenLiteral())
	}

	switch fn := val.(type) {
	case *object.NativeFunction:
		args := ev.evalExpressions(ce.Arguments)
		return fn.Function(args...)
	}

	return val
}

func getError(objs []object.Object) (err object.Object, ok bool) {
	ok = true
	if len(objs) == 1 && objs[0].Type() == object.ERROR {
		ok = false
		err = objs[0]
	}

	return
}

func (ev *Evaluator) evalExpressions(expressions []ast.Expression) (result []object.Object) {
	for _, e := range expressions {
		exp := ev.evaluate(e)
		if isError(exp) {
			return []object.Object{exp}
		}
		result = append(result, exp)
	}

	return
}

func (ev *Evaluator) evaluate(node ast.Node) object.Object {
	return node.Accept(ev)
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

func newError(msg string, v ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(msg, v...)}
}
