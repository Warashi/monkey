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
		return evalStatements(n.Statements)
	case *ast.ExpressionStatement:
		return Eval(n.Expression)
	case *ast.IntegerLiteral:
		return object.Integer{Value: n.Value}
	case *ast.BooleanLiteral:
		return booleanObject(n.Value)
	case *ast.PrefixExpression:
		return evalPrefixExpression(n.Operator, Eval(n.Right))
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

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt)
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
	if right.Type() != object.INTEGER_OBJ {
		return NULL
	}
	return object.Integer{Value: -right.(object.Integer).Value}
}
