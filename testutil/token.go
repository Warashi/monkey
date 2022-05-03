package testutil

import (
	"github.com/Warashi/implement-interpreter-with-go/ast"
	"github.com/Warashi/implement-interpreter-with-go/token"
)

func Token(n ast.Node) token.Token {
	switch n := n.(type) {
	case *ast.InfixExpression:
		return Token(n.Left)
	case *ast.Identifier:
		return n.Token
	default:
		return token.Token{}
	}
}
