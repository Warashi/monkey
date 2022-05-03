package testutil

import (
	"strconv"

	"github.com/Warashi/implement-interpreter-with-go/ast"
	"github.com/Warashi/implement-interpreter-with-go/object"
	"github.com/Warashi/implement-interpreter-with-go/token"
)

func IntegerObject(val int64) object.Object {
	return object.Integer{Value: val}
}

func BooleanObject(val bool) object.Object {
	return object.Boolean{Value: val}
}

func NullObject() object.Object {
	return object.Null{}
}

func ErrorObject(message string) object.Object {
	return object.Error{Message: message}
}

func FunctionObject(env object.Environment, body *ast.BlockStatement, params ...*ast.Identifier) object.Object {
	return object.Function{
		Parameters: params,
		Body:       body,
		Env:        env,
	}
}

func BlockStatement(statements ...ast.Statement) *ast.BlockStatement {
	return &ast.BlockStatement{
		Token:      token.Token{Type: token.LBRACE, Literal: "{"},
		Statements: statements,
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

func ExpressionStatement(e ast.Expression) *ast.ExpressionStatement {
	return &ast.ExpressionStatement{
		Token:      Token(e),
		Expression: e,
	}
}

func IntegerLiteral(val int64) *ast.IntegerLiteral {
	return &ast.IntegerLiteral{
		Token: token.Token{Type: token.INT, Literal: strconv.FormatInt(val, 10)},
		Value: val,
	}
}

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
