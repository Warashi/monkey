package evaluator

import (
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
		return evalPrefixExpression(n.Operator, Eval(n.Right))
	case *ast.InfixExpression:
		return evalInfixExpression(n.Operator, Eval(n.Left), Eval(n.Right))
	case *ast.BlockStatement:
		return evalBlockStatement(n)
	case *ast.IfExpression:
		return evalIfExpression(n)
	case *ast.ReturnStatement:
		return object.Return{Value: Eval(n.Value)}
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

func evalProgram(p *ast.Program) object.Object {
	var result object.Object
	for _, stmt := range p.Statements {
		result = Eval(stmt)
		if v, ok := result.(object.Return); ok {
			return v.Value
		}
	}
	return result
}

func evalBlockStatement(s *ast.BlockStatement) object.Object {
	var result object.Object
	for _, stmt := range s.Statements {
		result = Eval(stmt)
		if result != nil && result.Type() == object.TypeReturn {
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
		return NULL
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
		return NULL
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
	default:
		return NULL
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
		return NULL
	}
}

func evalIfExpression(n *ast.IfExpression) object.Object {
	if isTruthy(Eval(n.Condition)) {
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
