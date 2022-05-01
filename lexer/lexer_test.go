package lexer_test

import (
	"strconv"
	"testing"

	"github.com/Warashi/implement-interpreter-with-go/lexer"
	"github.com/Warashi/implement-interpreter-with-go/token"
	"github.com/stretchr/testify/assert"
)

func TestNextToken(t *testing.T) {
	const input = `=+(){},;`
	tok := func(t token.Type, l string) token.Token { return token.Token{Type: t, Literal: l} }

	wants := []token.Token{
		tok(token.ASSIGN, "="),
		tok(token.PLUS, "+"),
		tok(token.LPAREN, "("),
		tok(token.RPAREN, ")"),
		tok(token.LBRACE, "{"),
		tok(token.RBRACE, "}"),
		tok(token.COMMA, ","),
		tok(token.SEMICOLON, ";"),
	}

	for waste, want := range wants {
		waste, want := waste, want
		t.Run(strconv.Itoa(waste), func(t *testing.T) {
			l := lexer.New(input)
			for i := 0; i < waste; i++ {
				l.NextToken()
			}
			got := l.NextToken()
			assert.Equal(t, got, want)
		})

	}
}
