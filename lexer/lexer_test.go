package lexer_test

import (
	_ "embed"
	"strconv"
	"testing"

	"github.com/Warashi/implement-interpreter-with-go/lexer"
	"github.com/Warashi/implement-interpreter-with-go/lexer/testdata"
	. "github.com/Warashi/implement-interpreter-with-go/testutil"
	"github.com/Warashi/implement-interpreter-with-go/token"
	"github.com/stretchr/testify/assert"
)

func TestNextToken(t *testing.T) {
	tests := []struct {
		name  string
		input string
		wants []token.Token
	}{
		{
			name:  "symbols",
			input: testdata.Symbols,
			wants: []token.Token{
				Token(token.ASSIGN, "="),
				Token(token.PLUS, "+"),
				Token(token.LPAREN, "("),
				Token(token.RPAREN, ")"),
				Token(token.LBRACE, "{"),
				Token(token.RBRACE, "}"),
				Token(token.COMMA, ","),
				Token(token.SEMICOLON, ";"),
			},
		},
		{
			name:  "first",
			input: testdata.First,
			wants: []token.Token{
				Token(token.LET, "let"),
				Token(token.IDENT, "five"),
				Token(token.ASSIGN, "="),
				Token(token.INT, "5"),
				Token(token.SEMICOLON, ";"),
				Token(token.LET, "let"),
				Token(token.IDENT, "ten"),
				Token(token.ASSIGN, "="),
				Token(token.INT, "10"),
				Token(token.SEMICOLON, ";"),
				Token(token.LET, "let"),
				Token(token.IDENT, "add"),
				Token(token.ASSIGN, "="),
				Token(token.FUNCTION, "fn"),
				Token(token.LPAREN, "("),
				Token(token.IDENT, "x"),
				Token(token.COMMA, ","),
				Token(token.IDENT, "y"),
				Token(token.RPAREN, ")"),
				Token(token.LBRACE, "{"),
				Token(token.IDENT, "x"),
				Token(token.PLUS, "+"),
				Token(token.IDENT, "y"),
				Token(token.SEMICOLON, ";"),
				Token(token.RBRACE, "}"),
				Token(token.SEMICOLON, ";"),
				Token(token.LET, "let"),
				Token(token.IDENT, "result"),
				Token(token.ASSIGN, "="),
				Token(token.IDENT, "add"),
				Token(token.LPAREN, "("),
				Token(token.IDENT, "five"),
				Token(token.COMMA, ","),
				Token(token.IDENT, "ten"),
				Token(token.RPAREN, ")"),
				Token(token.SEMICOLON, ";"),
				Token(token.EOF, ""),
			},
		},
		{
			name:  "second",
			input: testdata.Second,
			wants: []token.Token{
				Token(token.LET, "let"),
				Token(token.IDENT, "five"),
				Token(token.ASSIGN, "="),
				Token(token.INT, "5"),
				Token(token.SEMICOLON, ";"),
				Token(token.LET, "let"),
				Token(token.IDENT, "ten"),
				Token(token.ASSIGN, "="),
				Token(token.INT, "10"),
				Token(token.SEMICOLON, ";"),
				Token(token.LET, "let"),
				Token(token.IDENT, "add"),
				Token(token.ASSIGN, "="),
				Token(token.FUNCTION, "fn"),
				Token(token.LPAREN, "("),
				Token(token.IDENT, "x"),
				Token(token.COMMA, ","),
				Token(token.IDENT, "y"),
				Token(token.RPAREN, ")"),
				Token(token.LBRACE, "{"),
				Token(token.IDENT, "x"),
				Token(token.PLUS, "+"),
				Token(token.IDENT, "y"),
				Token(token.SEMICOLON, ";"),
				Token(token.RBRACE, "}"),
				Token(token.SEMICOLON, ";"),
				Token(token.LET, "let"),
				Token(token.IDENT, "result"),
				Token(token.ASSIGN, "="),
				Token(token.IDENT, "add"),
				Token(token.LPAREN, "("),
				Token(token.IDENT, "five"),
				Token(token.COMMA, ","),
				Token(token.IDENT, "ten"),
				Token(token.RPAREN, ")"),
				Token(token.SEMICOLON, ";"),
				Token(token.BANG, "!"),
				Token(token.MINUS, "-"),
				Token(token.SLASH, "/"),
				Token(token.ASTERISK, "*"),
				Token(token.INT, "5"),
				Token(token.SEMICOLON, ";"),
				Token(token.INT, "5"),
				Token(token.LT, "<"),
				Token(token.INT, "10"),
				Token(token.GT, ">"),
				Token(token.INT, "5"),
				Token(token.SEMICOLON, ";"),
				Token(token.IF, "if"),
				Token(token.LPAREN, "("),
				Token(token.INT, "5"),
				Token(token.LT, "<"),
				Token(token.INT, "10"),
				Token(token.RPAREN, ")"),
				Token(token.LBRACE, "{"),
				Token(token.RETURN, "return"),
				Token(token.TRUE, "true"),
				Token(token.SEMICOLON, ";"),
				Token(token.RBRACE, "}"),
				Token(token.ELSE, "else"),
				Token(token.LBRACE, "{"),
				Token(token.RETURN, "return"),
				Token(token.FALSE, "false"),
				Token(token.SEMICOLON, ";"),
				Token(token.RBRACE, "}"),
				Token(token.INT, "10"),
				Token(token.EQ, "=="),
				Token(token.INT, "10"),
				Token(token.SEMICOLON, ";"),
				Token(token.INT, "10"),
				Token(token.NOT_EQ, "!="),
				Token(token.INT, "9"),
				Token(token.SEMICOLON, ";"),
				Token(token.EOF, ""),
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
