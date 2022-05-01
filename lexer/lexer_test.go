package lexer_test

import (
	_ "embed"
	"strconv"
	"testing"

	"github.com/Warashi/implement-interpreter-with-go/lexer"
	"github.com/Warashi/implement-interpreter-with-go/lexer/testdata"
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
			input: testdata.Symbols,
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
		{
			name:  "first",
			input: testdata.First,
			wants: []token.Token{
				tok(token.LET, "let"),
				tok(token.IDENT, "five"),
				tok(token.ASSIGN, "="),
				tok(token.INT, "5"),
				tok(token.SEMICOLON, ";"),
				tok(token.LET, "let"),
				tok(token.IDENT, "ten"),
				tok(token.ASSIGN, "="),
				tok(token.INT, "10"),
				tok(token.SEMICOLON, ";"),
				tok(token.LET, "let"),
				tok(token.IDENT, "add"),
				tok(token.ASSIGN, "="),
				tok(token.FUNCTION, "fn"),
				tok(token.LPAREN, "("),
				tok(token.IDENT, "x"),
				tok(token.COMMA, ","),
				tok(token.IDENT, "y"),
				tok(token.RPAREN, ")"),
				tok(token.LBRACE, "{"),
				tok(token.IDENT, "x"),
				tok(token.PLUS, "+"),
				tok(token.IDENT, "y"),
				tok(token.SEMICOLON, ";"),
				tok(token.RBRACE, "}"),
				tok(token.SEMICOLON, ";"),
				tok(token.LET, "let"),
				tok(token.IDENT, "result"),
				tok(token.ASSIGN, "="),
				tok(token.IDENT, "add"),
				tok(token.LPAREN, "("),
				tok(token.IDENT, "five"),
				tok(token.COMMA, ","),
				tok(token.IDENT, "ten"),
				tok(token.RPAREN, ")"),
				tok(token.SEMICOLON, ";"),
				tok(token.EOF, ""),
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
					assert.Equal(t, want, got)
				})
			}
		})
	}
}
