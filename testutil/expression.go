package testutil

import (
	"strconv"

	"github.com/Warashi/implement-interpreter-with-go/ast"
	"github.com/Warashi/implement-interpreter-with-go/token"
)

var (
	True = &ast.BooleanLiteral{
		Token: token.Token{Type: token.TRUE, Literal: "true"},
		Value: true,
	}
	False = &ast.BooleanLiteral{
		Token: token.Token{Type: token.FALSE, Literal: "false"},
		Value: false,
	}
	Asterisk = Token(token.ASTERISK, "*")
	Slash    = Token(token.SLASH, "/")
	Plus     = Token(token.PLUS, "+")
	Minus    = Token(token.MINUS, "-")
	GT       = Token(token.GT, ">")
	LT       = Token(token.LT, "<")
	Equal    = Token(token.EQ, "==")
	NotEqual = Token(token.NOT_EQ, "!=")
	Bang     = Token(token.BANG, "!")
)

func PrefixExpression(op token.Token, right ast.Expression) *ast.PrefixExpression {
	return &ast.PrefixExpression{
		Token:    op,
		Operator: op.Literal,
		Right:    right,
	}
}

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

func StringLiteral(val string) *ast.StringLiteral {
	return &ast.StringLiteral{
		Token: token.Token{Type: token.STRING, Literal: val},
		Value: val,
	}
}

func FunctionLiteral(body *ast.BlockStatement, params ...*ast.Identifier) *ast.FunctionLiteral {
	return &ast.FunctionLiteral{
		Token:      token.Token{Type: token.FUNCTION, Literal: "fn"},
		Parameters: params,
		Body:       body,
	}
}

func ArrayLiteral(elements ...ast.Expression) *ast.ArrayLiteral {
	return &ast.ArrayLiteral{
		Token:    token.Token{Type: token.LBLACKET, Literal: "["},
		Elements: elements,
	}
}

func IfExpression(cond ast.Expression, cons, alt *ast.BlockStatement) *ast.IfExpression {
	return &ast.IfExpression{
		Token:       token.Token{Type: token.IF, Literal: "if"},
		Condition:   cond,
		Consequence: cons,
		Alternative: alt,
	}
}

func CallExpression(fn *ast.Identifier, args ...ast.Expression) *ast.CallExpression {
	return &ast.CallExpression{
		Token:     Token(token.LPAREN, "("),
		Function:  fn,
		Arguments: args,
	}
}
