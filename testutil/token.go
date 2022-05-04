package testutil

import (
	"fmt"

	"github.com/Warashi/implement-interpreter-with-go/ast"
	"github.com/Warashi/implement-interpreter-with-go/token"
)

func Token(typ token.Type, literal string) token.Token {
	return token.Token{Type: typ, Literal: literal}
}

func TokenOf(n ast.Node) token.Token {
	switch n := n.(type) {
	case *ast.CallExpression:
		return TokenOf(n.Function)
	case *ast.InfixExpression:
		return TokenOf(n.Left)
	case *ast.Identifier:
		return n.Token
	case *ast.IntegerLiteral:
		return n.Token
	case *ast.BooleanLiteral:
		return n.Token
	case *ast.FunctionLiteral:
		return n.Token
	case *ast.PrefixExpression:
		return n.Token
	case *ast.IfExpression:
		return n.Token
	default:
		panic(fmt.Sprintf("unknown type: %T", n))
	}
}
