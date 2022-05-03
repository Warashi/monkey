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
	pre := func(t token.Type, op string, v int64) ast.Statement {
		return &ast.ExpressionStatement{
			Token: token.Token{Type: t, Literal: op},
			Expression: &ast.PrefixExpression{
				Token:    token.Token{Type: t, Literal: op},
				Operator: op,
				Right: &ast.IntegerLiteral{
					Token: token.Token{Type: token.INT, Literal: strconv.FormatInt(v, 10)},
					Value: v,
				},
			},
		}
	}
	tests := []struct {
		input string
		want  []ast.Statement
	}{
		{input: "!5", want: []ast.Statement{pre(token.BANG, "!", 5)}},
		{input: "-15", want: []ast.Statement{pre(token.MINUS, "-", 15)}},
	}
	for i, tt := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			p := parser.New(lexer.New(tt.input))
			program := p.Parse()
			require.Empty(t, p.Errors())
			require.NotNil(t, program)
			assert.Equal(t, tt.want, program.Statements)
		})
	}
}
