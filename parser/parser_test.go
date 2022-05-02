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
	p := parser.New(lexer.New(testdata.Let)).Parse()
	require.NotNil(t, p)
	require.Len(t, p.Statements, 3)

	let := func(name string) ast.Statement {
		return &ast.LetStatement{Name: &ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: name}, Value: name}}
	}
	wants := []ast.Statement{let("x"), let("y"), let("foobar")}
	for i, want := range wants {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			assert.Equal(t, want, p.Statements[i])
		})
	}
}
