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
	// native functions
	nf_typeof  = newNativeFunction(nativeTypeof, "typeof")
	nf_typeofS = newNativeFunction(nativeTypeofS, "typeofS")
	nf_inspect = newNativeFunction(nativeInspect, "inspect")
	arr_len    = newNativeFunction(arrLen, "len")

	mathfn_sin   = newNativeFunction(math2NativeFn(math.Sin), "sin")
	mathfn_cos   = newNativeFunction(math2NativeFn(math.Cos), "cos")
	mathfn_ln    = newNativeFunction(math2NativeFn(math.Log), "ln")
	mathfn_log2  = newNativeFunction(math2NativeFn(math.Log2), "log2")
	mathfn_log10 = newNativeFunction(math2NativeFn(math.Log10), "log10")
	mathfn_sqrt  = newNativeFunction(math2NativeFn(math.Sqrt), "sqrt")
	mathc_e      = newFloat(math.E)
	mathc_pi     = newFloat(math.Pi)
	mathc_phi    = newFloat(math.Phi)
)

func newNativeFunction(fn NativeFn, name string) *NativeFunction {
	return &NativeFunction{Function: fn, Name: name}
}

func newFloat(val float64) *object.Float {
	return &object.Float{Value: val}
}

func New() *Evaluator {
	ev := &Evaluator{}
	ev.global = environment.New()

	// TODO: Libraries
	ev.global.Set("typeof", nf_typeof)
	ev.global.Set("typeofS", nf_typeofS)
	ev.global.Set("inspect", nf_inspect)
	ev.global.Set("len", arr_len)
	ev.global.Set("ln", mathfn_ln)
	ev.global.Set("log2", mathfn_log2)
	ev.global.Set("log10", mathfn_log10)
	ev.global.Set("log", mathfn_log10)
	ev.global.Set("sqrt", mathfn_sqrt)
	ev.global.Set("sin", mathfn_sin)
	ev.global.Set("cos", mathfn_cos)
	ev.global.Set("e", mathc_e)
	ev.global.Set("pi", mathc_pi)
	ev.global.Set("phi", mathc_phi)
	return ev
}

// TODO: Libraries, return multiple values
func mathSin(ev *Evaluator, objs ...object.Object) object.Object {
	num, ok := objs[0].(*object.Float)
	if !ok {
		// TODO: handle error
		return NULL
	}

	return &object.Float{Value: math.Sin(num.Value)}
}

type mathFn func(float64) float64

func math2NativeFn(fn mathFn) NativeFn {
	return func(ev *Evaluator, objs ...object.Object) object.Object {
		num, ok := objs[0].(*object.Float)
		if !ok {
			// TODO: handle error
			return NULL
		}
		return &object.Float{Value: fn(num.Value)}
	}
}

func mathLog(ev *Evaluator, objs ...object.Object) object.Object {
	num, ok := objs[0].(*object.Float)
	if !ok {
		// TODO: handle error
		return NULL
	}
	return &object.Float{Value: math.Log(num.Value)}
}

func nativeInspect(ev *Evaluator, objs ...object.Object) object.Object {
	if len(objs) == 0 {
		return object.NewString(fmt.Sprint(ev.global))
	}

	// TODO: take string argument and inspect that function
	return NULL
}

func nativeTypeofS(ev *Evaluator, objs ...object.Object) object.Object {
	if len(objs) == 0 {
		return &object.String{Value: object.NATIVE_FUNCTION.Stringf("typeofS")}
	}
	obj := objs[0]
	return &object.String{Value: obj.TypeS()}
}

func nativeTypeof(ev *Evaluator, objs ...object.Object) object.Object {
	if len(objs) == 0 {
		return &object.Type{Value: object.NATIVE_FUNCTION}
	}
	obj := objs[0]
	return &object.Type{Value: obj.Type()}
}

func arrLen(ev *Evaluator, objs ...object.Object) object.Object {
	if len(objs) == 0 {
		return newFloat(0)
	}
	obj, ok := objs[0].(*object.List)

	if !ok {
		return newError("Len can only be applied to lists. Got %s", obj.Type())
	}

	return newFloat(float64(len(obj.Values)))

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

func (ev *Evaluator) BooleanLiteral(fl *ast.BooleanLiteral) object.Object {
	return &object.Boolean{Value: fl.Value}
}

func (ev *Evaluator) AssignmentStatement(as *ast.AssignmentStatement) object.Object {
	val := ev.evaluate(as.Value)
	if isError(val) {
		return val
	}

	ev.global.Set(as.Name.Value, val)

	return nil
}

func (ev *Evaluator) ListLiteral(ll *ast.ListLiteral) object.Object {
	return &object.List{Values: ev.evalExpressions(ll.Values)}
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
	case *NativeFunction:
		args := ev.evalExpressions(ce.Arguments)
		return fn.Function(ev, args...)
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
	case isBoolean(left) && isBoolean(right):
		return evalInfixExpressionBoolean(operator, left, right)
	default:
		return newError(object.UNKNOWN_INFIX_OPERATOR_ERROR, left.Type(), operator, right.Type())
	}
}

func isFloat(obj object.Object) bool {
	return obj.Type() == object.FLOAT
}

func isBoolean(obj object.Object) bool {
	return obj.Type() == object.BOOLEAN
}

func evalInfixExpressionBoolean(operator string, left, right object.Object) object.Object {
	x1, x2 := left.(*object.Boolean).Value, right.(*object.Boolean).Value

	switch operator {
	case "==":
		return newBool(x1 == x2)
	case "!=":
		return newBool(x1 != x2)
	case "&&":
		return newBool(x1 && x2)
	case "||":
		return newBool(x1 || x2)
	}

	return newError(object.UNKNOWN_INFIX_OPERATOR_ERROR, left.Type(), operator, right.Type())
}

func newBool(val bool) *object.Boolean {
	return &object.Boolean{Value: val}
}

func evalInfixExpressionFloat(operator string, left, right object.Object) object.Object {
	x1, x2 := left.(*object.Float).Value, right.(*object.Float).Value

	switch operator {
	case "+":
		return newFloat(x1 + x2)
	case "-":
		return newFloat(x1 - x2)
	case "*":
		return newFloat(x1 * x2)
	case "^":
		return newFloat(math.Pow(x1, x2))
	case "/":
		if x2 == 0 {
			return object.DivideByZeroError(left, right)
		}
		return newFloat(x1 / x2)
	case ">=":
		return newBool(x1 >= x2)
	case ">":
		return newBool(x1 > x2)
	case "<":
		return newBool(x1 < x2)
	case "<=":
		return newBool(x1 <= x2)
	case "==":
		return newBool(x1 == x2)
	case "!=":
		return newBool(x1 != x2)
	}

	return newError(object.UNKNOWN_INFIX_OPERATOR_ERROR, left.Type(), operator, right.Type())

}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch {
	case isFloat(right):
		return evalPrefixExpressionFloat(operator, right)
	case isBoolean(right):
		return evalPrefixExpressionBoolean(operator, right)
	default:
		return newError(object.UNKNOWN_PREFIX_OPERATOR_ERROR, operator, right.Type())
	}
}

func evalPrefixExpressionBoolean(operator string, right object.Object) object.Object {
	x1 := right.(*object.Boolean).Value
	switch operator {
	case "!":
		return newBool(!x1)
	default:
		return newError(object.UNKNOWN_PREFIX_OPERATOR_ERROR, operator, right.Type())
	}
}

func evalPrefixExpressionFloat(operator string, right object.Object) object.Object {
	x1 := right.(*object.Float).Value
	switch operator {
	case "-":
		return newFloat(-x1)
	}

	return newError(object.UNKNOWN_PREFIX_OPERATOR_ERROR, operator, right.Type())
}

func newError(msg string, v ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(msg, v...)}
}
