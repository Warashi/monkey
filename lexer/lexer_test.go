package lexer_test

import (
	"strconv"
	"testing"

	"github.com/Warashi/implement-interpreter-with-go/lexer"
	"github.com/Warashi/implement-interpreter-with-go/token"
	"github.com/stretchr/testify/assert"
)

func TestNextToken(t *testing.T) {
	tok := func(t token.Type, l string) token.Token { return token.Token{Type: t, Literal: l} }

	tests := []struct {
		name  string
		input string
		wants []token.Token
	}{
		{
			name:  "symbols",
			input: `=+(){},;`,
			wants: []token.Token{
				tok(token.ASSIGN, "="),
				tok(token.PLUS, "+"),
				tok(token.LPAREN, "("),
				tok(token.RPAREN, ")"),
				tok(token.LBRACE, "{"),
				tok(token.RBRACE, "}"),
				tok(token.COMMA, ","),
				tok(token.SEMICOLON, ";"),
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			for waste, want := range tt.wants {
				waste, want := waste, want
				t.Run(strconv.Itoa(waste), func(t *testing.T) {
					l := lexer.New(tt.input)
					for i := 0; i < waste; i++ {
						l.NextToken()
					}
					got := l.NextToken()
					assert.Equal(t, got, want)
				})
			}
		})
	}
}
