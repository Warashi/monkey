package parser_test

import (
	"strconv"
	"testing"

	"github.com/Warashi/implement-interpreter-with-go/ast"
	"github.com/Warashi/implement-interpreter-with-go/lexer"
	"github.com/Warashi/implement-interpreter-with-go/parser"
	"github.com/Warashi/implement-interpreter-with-go/parser/testdata"
	"github.com/Warashi/implement-interpreter-with-go/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLetStatement(t *testing.T) {
	p := parser.New(lexer.New(testdata.Let))
	program := p.Parse()
	require.Empty(t, p.Errors())
	require.NotNil(t, program)
	assert.Equal(t, testdata.Let, program.String())

	let := func(name string) ast.Statement {
		return &ast.LetStatement{
			Token: token.Token{Type: token.LET, Literal: "let"},
			Name:  &ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: name}, Value: name}}
	}
	wants := []ast.Statement{let("x"), let("y"), let("foobar")}
	assert.Equal(t, wants, program.Statements)
}

func TestReturnStatement(t *testing.T) {
	p := parser.New(lexer.New(testdata.Return))
	program := p.Parse()
	require.Empty(t, p.Errors())
	require.NotNil(t, program)
	assert.Equal(t, testdata.Return, program.String())

	ret := func(name string) ast.Statement {
		return &ast.ReturnStatement{
			Token: token.Token{Type: token.RETURN, Literal: "return"},
		}
	}
	wants := []ast.Statement{ret("5"), ret("10"), ret("993322")}
	assert.Equal(t, wants, program.Statements)
}

func TestIdentifierExpression(t *testing.T) {
	p := parser.New(lexer.New(testdata.IdentifierExpression))
	program := p.Parse()
	require.Empty(t, p.Errors())
	require.NotNil(t, program)
	assert.Equal(t, testdata.IdentifierExpression, program.String())

	id := func(name string) ast.Statement {
		return &ast.ExpressionStatement{
			Token: token.Token{Type: token.IDENT, Literal: name},
			Expression: &ast.Identifier{
				Token: token.Token{Type: token.IDENT, Literal: name},
				Value: name,
			},
		}
	}
	wants := []ast.Statement{id("foobar")}
	assert.Equal(t, wants, program.Statements)
}

func TestIntegerLiteralExpression(t *testing.T) {
	p := parser.New(lexer.New(testdata.IntegerLiteralExpression))
	program := p.Parse()
	require.Empty(t, p.Errors())
	require.NotNil(t, program)
	assert.Equal(t, testdata.IntegerLiteralExpression, program.String())

	integer := func(value int64) ast.Statement {
		return &ast.ExpressionStatement{
			Token: token.Token{Type: token.INT, Literal: strconv.FormatInt(value, 10)},
			Expression: &ast.IntegerLiteral{
				Token: token.Token{Type: token.INT, Literal: strconv.FormatInt(value, 10)},
				Value: value,
			},
		}
	}
	wants := []ast.Statement{integer(5)}
	assert.Equal(t, wants, program.Statements)
}
