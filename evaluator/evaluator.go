package evaluator

import (
	"fmt"

	"github.com/Warashi/implement-interpreter-with-go/ast"
	"github.com/Warashi/implement-interpreter-with-go/object"
)

var (
	TRUE  = object.Boolean{Value: true}
	FALSE = object.Boolean{Value: false}
	NULL  = object.Null{}
)

func Eval(n ast.Node, env object.Environment) object.Object {
	switch n := n.(type) {
	case *ast.Program:
		return evalProgram(n, env)
	case *ast.ExpressionStatement:
		return Eval(n.Expression, env)
	case *ast.IntegerLiteral:
		return object.Integer{Value: n.Value}
	case *ast.StringLiteral:
		return object.String{Value: n.Value}
	case *ast.BooleanLiteral:
		return booleanObject(n.Value)
	case *ast.PrefixExpression:
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(n.Operator, right)
	case *ast.InfixExpression:
		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(n.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(n, env)
	case *ast.IfExpression:
		return evalIfExpression(n, env)
	case *ast.ReturnStatement:
		result := Eval(n.Value, env)
		if isError(result) {
			return result
		}
		return object.Return{Value: result}
	case *ast.LetStatement:
		result := Eval(n.Value, env)
		if isError(result) {
			return result
		}
		return env.Set(n.Name.Value, result)
	case *ast.Identifier:
		return evalIdentifier(n, env)
	case *ast.FunctionLiteral:
		return object.Function{Parameters: n.Parameters, Body: n.Body, Env: env}
	case *ast.CallExpression:
		fn := Eval(n.Function, env)
		if isError(fn) {
			return fn
		}
		args := evalExpresssions(n.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunciton(fn, args)
	case *ast.ArrayLiteral:
		elements := evalExpresssions(n.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return object.Array{Elements: elements}
	case *ast.HashLiteral:
		pairs := make(map[object.Hashable]object.Object, len(n.Pairs))
		for k, v := range n.Pairs {
			key := Eval(k, env)
			if isError(key) {
				return key
			}
			keyHashable, ok := key.(object.Hashable)
			if !ok {
				return newErrorf("%s cannot used as hash key", key.Type())
			}
			value := Eval(v, env)
			if isError(value) {
				return value
			}
			pairs[keyHashable] = value
		}
		return object.Hash{Pairs: pairs}
	case *ast.IndexExpression:
		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalIndexExpression(left, right)
	default:
		return newErrorf("unknown node: %T", n)
	}
}

func booleanObject(val bool) object.Boolean {
	if val {
		return TRUE
	}
	return FALSE
}

func newErrorf(format string, a ...any) object.Object {
	return object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(o object.Object) bool {
	return o.Type() == object.TypeError
}

func evalProgram(p *ast.Program, env object.Environment) object.Object {
	var result object.Object
	for _, stmt := range p.Statements {
		result = Eval(stmt, env)
		switch result := result.(type) {
		case object.Return:
			return result.Value
		case object.Error:
			return result
		}
	}
	return result
}

func evalBlockStatement(s *ast.BlockStatement, env object.Environment) object.Object {
	var result object.Object
	for _, stmt := range s.Statements {
		result = Eval(stmt, env)
		if t := result.Type(); t == object.TypeReturn || t == object.TypeError {
			return result
		}
	}
	return result
}

func evalPrefixExpression(op string, right object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return newErrorf("unknown operator %s", op)
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.TypeInteger {
		return newErrorf("unknown operator: -%s", right.Type())
	}
	return object.Integer{Value: -right.(object.Integer).Value}
}

func evalInfixExpression(op string, left, right object.Object) object.Object {
	switch {
	case op == "==":
		return booleanObject(left == right)
	case op == "!=":
		return booleanObject(left != right)
	case left.Type() == object.TypeInteger && right.Type() == object.TypeInteger:
		return evalIntegerInfixExpression(op, left.(object.Integer), right.(object.Integer))
	case left.Type() == object.TypeString && right.Type() == object.TypeString:
		return evalStringInfixExpression(op, left.(object.String), right.(object.String))
	case left.Type() != right.Type():
		return newErrorf("type mismatch: %s %s %s", left.Type(), op, right.Type())
	default:
		return newErrorf("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalIndexExpression(left, right object.Object) object.Object {
	switch {
	case left.Type() == object.TypeArray && right.Type() == object.TypeInteger:
		left, right := left.(object.Array), right.(object.Integer)
		if right.Value < 0 || int64(len(left.Elements)) <= right.Value {
			return newErrorf("index out of range. index=%d, len=%d", right.Value, len(left.Elements))
		}
		return left.Elements[right.Value]
	case left.Type() == object.TypeHash:
		left := left.(object.Hash)
		rightHashable, ok := right.(object.Hashable)
		if !ok {
			return newErrorf("%s cannot used as hash key", right.Type())
		}
		val, ok := left.Pairs[rightHashable]
		if !ok {
			return newErrorf("key not found. key=%s", right.Inspect())
		}
		return val
	default:
		return newErrorf("type mismatch: %s[%s]", left.Type(), right.Type())
	}
}

func evalIntegerInfixExpression(op string, left, right object.Integer) object.Object {
	switch op {
	case "+":
		return object.Integer{Value: left.Value + right.Value}
	case "-":
		return object.Integer{Value: left.Value - right.Value}
	case "*":
		return object.Integer{Value: left.Value * right.Value}
	case "/":
		return object.Integer{Value: left.Value / right.Value}
	case "<":
		return booleanObject(left.Value < right.Value)
	case ">":
		return booleanObject(left.Value > right.Value)
	default:
		return newErrorf("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalStringInfixExpression(op string, left, right object.String) object.Object {
	switch op {
	case "+":
		return object.String{Value: left.Value + right.Value}
	case "<":
		return booleanObject(left.Value < right.Value)
	case ">":
		return booleanObject(left.Value > right.Value)
	default:
		return newErrorf("unknown operator: %s %s %s", left.Type(), op, right.Type())
	}
}

func evalIfExpression(n *ast.IfExpression, env object.Environment) object.Object {
	cond := Eval(n.Condition, env)
	if isError(cond) {
		return cond
	}
	if isTruthy(cond) {
		return Eval(n.Consequence, env)
	}
	if n.Alternative != nil {
		return Eval(n.Alternative, env)
	}
	return NULL
}

func evalIdentifier(n *ast.Identifier, env object.Environment) object.Object {
	if val, ok := env.Get(n.Value); ok {
		return val
	}
	if val, ok := builtins[n.Value]; ok {
		return val
	}
	return newErrorf("identifier not found: %s", n.Value)
}

func evalExpresssions(e []ast.Expression, env object.Environment) []object.Object {
	result := make([]object.Object, 0, len(e))
	for _, e := range e {
		r := Eval(e, env)
		if isError(r) {
			return []object.Object{r}
		}
		result = append(result, r)
	}
	return result
}

func applyFunciton(fn object.Object, args []object.Object) object.Object {
	switch fn.Type() {
	case object.TypeFunction:
		f := fn.(object.Function)
		return unwrapReturnValue(Eval(f.Body, extendFunctionEnv(f, args)))
	case object.TypeBuiltin:
		f := fn.(object.Builtin)
		return f.Fn(args...)
	default:
		return newErrorf("not a function: %s", fn.Type())
	}
}

func unwrapReturnValue(o object.Object) object.Object {
	if o.Type() == object.TypeReturn {
		return o.(object.Return).Value
	}
	return o
}

func extendFunctionEnv(fn object.Function, args []object.Object) object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for i, param := range fn.Parameters {
		env.Set(param.Value, args[i])
	}
	return env
}

func isTruthy(o object.Object) bool {
	switch o {
	case NULL, FALSE:
		return false
	default:
		return true
	}
}
