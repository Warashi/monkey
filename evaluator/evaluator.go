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
	case left.Type() != right.Type():
		return newErrorf("type mismatch: %s %s %s", left.Type(), op, right.Type())
	default:
		return newErrorf("unknown operator: %s %s %s", left.Type(), op, right.Type())
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
	return newErrorf("identifier not found: %s", n.Value)
}

func isTruthy(o object.Object) bool {
	switch o {
	case NULL, FALSE:
		return false
	default:
		return true
	}
}
