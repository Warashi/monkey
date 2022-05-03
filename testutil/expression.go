package testutil

import (
	"strconv"

	"github.com/Warashi/implement-interpreter-with-go/ast"
	"github.com/Warashi/implement-interpreter-with-go/token"
)

func InfixExpression(op token.Token, left, right ast.Expression) *ast.InfixExpression {
	return &ast.InfixExpression{
		Token:    op,
		Operator: op.Literal,
		Left:     left,
		Right:    right,
	}
}

func Identifier(name string) *ast.Identifier {
	return &ast.Identifier{
		Token: token.Token{Type: token.IDENT, Literal: name},
		Value: name,
	}
}

func IntegerLiteral(val int64) *ast.IntegerLiteral {
	return &ast.IntegerLiteral{
		Token: token.Token{Type: token.INT, Literal: strconv.FormatInt(val, 10)},
		Value: val,
	}
}
