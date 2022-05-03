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

func Eval(n ast.Node) object.Object {
	switch n := n.(type) {
	case *ast.Program:
		return evalProgram(n)
	case *ast.ExpressionStatement:
		return Eval(n.Expression)
	case *ast.IntegerLiteral:
		return object.Integer{Value: n.Value}
	case *ast.BooleanLiteral:
		return booleanObject(n.Value)
	case *ast.PrefixExpression:
		right := Eval(n.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(n.Operator, right)
	case *ast.InfixExpression:
		left := Eval(n.Left)
		if isError(left) {
			return left
		}
		right := Eval(n.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(n.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(n)
	case *ast.IfExpression:
		return evalIfExpression(n)
	case *ast.ReturnStatement:
		result := Eval(n.Value)
		if isError(result) {
			return result
		}
		return object.Return{Value: result}
	default:
		return nil
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

func evalProgram(p *ast.Program) object.Object {
	var result object.Object
	for _, stmt := range p.Statements {
		result = Eval(stmt)
		switch result := result.(type) {
		case object.Return:
			return result.Value
		case object.Error:
			return result
		}
	}
	return result
}

func evalBlockStatement(s *ast.BlockStatement) object.Object {
	var result object.Object
	for _, stmt := range s.Statements {
		result = Eval(stmt)
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

func evalIfExpression(n *ast.IfExpression) object.Object {
	cond := Eval(n.Condition)
	if isError(cond) {
		return cond
	}
	if isTruthy(cond) {
		return Eval(n.Consequence)
	}
	if n.Alternative != nil {
		return Eval(n.Alternative)
	}
	return NULL
}

func isTruthy(o object.Object) bool {
	switch o {
	case NULL, FALSE:
		return false
	default:
		return true
	}
}
