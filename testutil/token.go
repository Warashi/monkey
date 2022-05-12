package testutil

import (
	"reflect"

	"github.com/Warashi/monkey/ast"
	"github.com/Warashi/monkey/token"
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
	case *ast.IndexExpression:
		return TokenOf(n.Left)
	default:
		return reflect.ValueOf(n).Elem().FieldByName("Token").Interface().(token.Token)
	}
}
