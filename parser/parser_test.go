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
	require.Len(t, program.Statements, 3)
	assert.Equal(t, testdata.Let, program.String())

	let := func(name string) ast.Statement {
		return &ast.LetStatement{
			Token: token.Token{Type: token.LET, Literal: "let"},
			Name:  &ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: name}, Value: name}}
	}
	wants := []ast.Statement{let("x"), let("y"), let("foobar")}
	for i, want := range wants {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, want, program.Statements[i])
		})
	}
}

func TestReturnStatement(t *testing.T) {
	p := parser.New(lexer.New(testdata.Return))
	program := p.Parse()
	require.Empty(t, p.Errors())
	require.NotNil(t, program)
	require.Len(t, program.Statements, 3)
	assert.Equal(t, testdata.Return, program.String())

	ret := func(name string) ast.Statement {
		return &ast.ReturnStatement{
			Token: token.Token{Type: token.RETURN, Literal: "return"},
		}
	}
	wants := []ast.Statement{ret("5"), ret("10"), ret("993322")}
	for i, want := range wants {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, want, program.Statements[i])
		})
	}
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
