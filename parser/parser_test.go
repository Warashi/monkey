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

func TestPrefixExpression(t *testing.T) {
	tests := []struct {
		i   string
		tok token.Token
		op  string
		v   int64
	}{
		{i: "!5", tok: token.Token{Type: token.BANG, Literal: "!"}, op: "!", v: 5},
		{i: "-15", tok: token.Token{Type: token.MINUS, Literal: "-"}, op: "-", v: 15},
	}
	for _, tt := range tests {
		p := parser.New(lexer.New(tt.i))
		program := p.Parse()
		require.Empty(t, p.Errors())
		require.NotNil(t, program)
		want := []ast.Statement{&ast.ExpressionStatement{
			Token: tt.tok,
			Expression: &ast.PrefixExpression{
				Token:    tt.tok,
				Operator: tt.op,
				Right:    &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: strconv.FormatInt(tt.v, 10)}, Value: tt.v},
			},
		}}
		assert.Equal(t, want, program.Statements)
	}
}
