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
		Token:      Token(e),
		Expression: e,
	}
}
