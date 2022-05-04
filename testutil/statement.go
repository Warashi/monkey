package testutil

import (
	"github.com/Warashi/implement-interpreter-with-go/ast"
	"github.com/Warashi/implement-interpreter-with-go/token"
)

func BlockStatement(statements ...ast.Statement) *ast.BlockStatement {
	return &ast.BlockStatement{
		Token:      token.Token{Type: token.LBRACE, Literal: "{"},
		Statements: statements,
	}
}

func ExpressionStatement(e ast.Expression) *ast.ExpressionStatement {
	return &ast.ExpressionStatement{
		Token:      TokenOf(e),
		Expression: e,
	}
}

func LetStatement(name *ast.Identifier, value ast.Expression) *ast.LetStatement {
	return &ast.LetStatement{
		Token: token.Token{Type: token.LET, Literal: "let"},
		Name:  name,
		Value: value,
	}
}

func ReturnStatement(value ast.Expression) *ast.ReturnStatement {
	return &ast.ReturnStatement{
		Token: token.Token{Type: token.RETURN, Literal: "return"},
		Value: value,
	}
}
