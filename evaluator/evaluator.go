package evaluator

import (
	"github.com/Warashi/implement-interpreter-with-go/ast"
	"github.com/Warashi/implement-interpreter-with-go/object"
)

var (
	TRUE  = object.Boolean{Value: true}
	FALSE = object.Boolean{Value: false}
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
