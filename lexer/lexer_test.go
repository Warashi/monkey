package lexer_test

import (
	"testing"

	"github.com/Warashi/implement-interpreter-with-go/lexer"
	"github.com/Warashi/implement-interpreter-with-go/token"
	"github.com/stretchr/testify/assert"
)

func TestNextToken(t *testing.T) {
	const input = `=+(){},;`

	tests := []struct {
		name  string
		input string
		waste int
		want  token.Token
	}{
		{name: token.ASSIGN.String(), waste: 0, want: token.Token{token.ASSIGN, "="}},
		{name: token.PLUS.String(), waste: 1, want: token.Token{token.PLUS, "+"}},
		{name: token.LPAREN.String(), waste: 2, want: token.Token{token.LPAREN, "("}},
		{name: token.RPAREN.String(), waste: 3, want: token.Token{token.RPAREN, ")"}},
		{name: token.LBRACE.String(), waste: 4, want: token.Token{token.LBRACE, "{"}},
		{name: token.RBRACE.String(), waste: 5, want: token.Token{token.RBRACE, "}"}},
		{name: token.COMMA.String(), waste: 6, want: token.Token{token.COMMA, ","}},
		{name: token.SEMICOLON.String(), waste: 7, want: token.Token{token.SEMICOLON, ";"}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			l := lexer.New(input)
			for i := 0; i < tt.waste; i++ {
				l.NextToken()
			}
			got := l.NextToken()
			assert.Equal(t, got, tt.want)
		})

	}
}
