package evaluator

import (
	"github.com/Warashi/implement-interpreter-with-go/ast"
	"github.com/Warashi/implement-interpreter-with-go/object"
)

func Eval(n ast.Node) object.Object {
	switch n := n.(type) {
	case *ast.Program:
		return evalStatements(n.Statements)
	case *ast.ExpressionStatement:
		return Eval(n.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}
	case *ast.BooleanLiteral:
		return &object.Boolean{Value: n.Value}
	default:
		return nil
	}
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt)
	}
	return result
}
