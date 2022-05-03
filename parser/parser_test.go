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

func TestBooleanLiteralExpression(t *testing.T) {
	p := parser.New(lexer.New(testdata.BooleanLiteralExpression))
	program := p.Parse()
	require.Empty(t, p.Errors())
	require.NotNil(t, program)

	trueExp := &ast.ExpressionStatement{
		Token: token.Token{Type: token.TRUE, Literal: "true"},
		Expression: &ast.BooleanLiteral{
			Token: token.Token{Type: token.TRUE, Literal: "true"},
			Value: true,
		},
	}
	falseExp := &ast.ExpressionStatement{
		Token: token.Token{Type: token.FALSE, Literal: "false"},
		Expression: &ast.BooleanLiteral{
			Token: token.Token{Type: token.FALSE, Literal: "false"},
			Value: false,
		},
	}
	wants := []ast.Statement{trueExp, falseExp, trueExp, falseExp}
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

func TestInfixExpression(t *testing.T) {
	in := func(t token.Type, op string, left, right int64) ast.Statement {
		l := &ast.IntegerLiteral{
			Token: token.Token{Type: token.INT, Literal: strconv.FormatInt(left, 10)},
			Value: left,
		}
		r := &ast.IntegerLiteral{
			Token: token.Token{Type: token.INT, Literal: strconv.FormatInt(right, 10)},
			Value: right,
		}
		return &ast.ExpressionStatement{
			Token: l.Token,
			Expression: &ast.InfixExpression{
				Token:    token.Token{Type: t, Literal: op},
				Operator: op,
				Left:     l,
				Right:    r,
			},
		}
	}
	tests := []struct {
		input string
		want  []ast.Statement
	}{
		{input: "5 + 6", want: []ast.Statement{in(token.PLUS, "+", 5, 6)}},
		{input: "5 - 6", want: []ast.Statement{in(token.MINUS, "-", 5, 6)}},
		{input: "5 * 6", want: []ast.Statement{in(token.ASTERISK, "*", 5, 6)}},
		{input: "5 / 6", want: []ast.Statement{in(token.SLASH, "/", 5, 6)}},
		{input: "5 > 6", want: []ast.Statement{in(token.GT, ">", 5, 6)}},
		{input: "5 < 6", want: []ast.Statement{in(token.LT, "<", 5, 6)}},
		{input: "5 == 6", want: []ast.Statement{in(token.EQ, "==", 5, 6)}},
		{input: "5 != 6", want: []ast.Statement{in(token.NOT_EQ, "!=", 5, 6)}},
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

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{input: "-a * b", want: "((-a) * b)"},
		{input: "!-a", want: "(!(-a))"},
		{input: "a + b + c", want: "((a + b) + c)"},
		{input: "a + b - c", want: "((a + b) - c)"},
		{input: "a * b * c", want: "((a * b) * c)"},
		{input: "a * b / c", want: "((a * b) / c)"},
		{input: "a + b / c", want: "(a + (b / c))"},
		{input: "a + b * c + d / e - f", want: "(((a + (b * c)) + (d / e)) - f)"},
		{input: "3 + 4; -5 * 5", want: "(3 + 4)((-5) * 5)"},
		{input: "5 > 4 == 3 < 4", want: "((5 > 4) == (3 < 4))"},
		{input: "5 < 4 != 3 > 4", want: "((5 < 4) != (3 > 4))"},
		{input: "3 + 4 * 5 == 3 * 1 + 4 * 5", want: "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{input: "true", want: "true"},
		{input: "false", want: "false"},
		{input: "3 > 5 == false", want: "((3 > 5) == false)"},
		{input: "3 < 5 == true", want: "((3 < 5) == true)"},
		{input: "1 + (2 + 3) + 4", want: "((1 + (2 + 3)) + 4)"},
		{input: "(5 + 5) * 2", want: "((5 + 5) * 2)"},
		{input: "2 / (5 + 5)", want: "(2 / (5 + 5))"},
		{input: "(5 + 5) * 2 * (5 + 5)", want: "(((5 + 5) * 2) * (5 + 5))"},
		{input: "-(5 + 5)", want: "(-(5 + 5))"},
		{input: "!(true == true)", want: "(!(true == true))"},
		{input: "a + add(b * c) + d", want: "((a + add((b * c))) + d)"},
		{input: "add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", want: "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{input: "add(a + b + c * d / f + g)", want: "add((((a + b) + ((c * d) / f)) + g))"},
		{input: "", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := parser.New(lexer.New(tt.input))
			program := p.Parse()
			require.Empty(t, p.Errors())
			require.NotNil(t, program)
			assert.Equal(t, tt.want, program.String())
		})
	}
}

func TestIfExpression(t *testing.T) {
	tests := []struct {
		input string
		want  []ast.Statement
	}{
		{
			/*
			   +   Consequence: (*ast.BlockStatement)({
			   +    Token: (token.Token) {
			   +     Type: (token.Type) 18,
			   +     Literal: (string) (len=1) "{"
			   +    },
			   +    Statements: ([]ast.Statement) (len=1) {
			   +     (*ast.ExpressionStatement)({
			   +      Token: (token.Token) {
			   +       Type: (token.Type) 2,
			   +       Literal: (string) (len=1) "x"
			   +      },
			   +      Expression: (*ast.Identifier)({
			   +       Token: (token.Token) {
			   +        Type: (token.Type) 2,
			   +        Literal: (string) (len=1) "x"
			   +       },
			   +       Value: (string) (len=1) "x"
			   +      })
			   +     })
			   +    }
			   +   }),
			*/
			input: "if (x < y) { x }",
			want: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.IF, Literal: "if"},
					Expression: &ast.IfExpression{
						Token: token.Token{Type: token.IF, Literal: "if"},
						Condition: &ast.InfixExpression{
							Token:    token.Token{Type: token.LT, Literal: "<"},
							Operator: "<",
							Left: &ast.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "x"},
								Value: "x",
							},
							Right: &ast.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "y"},
								Value: "y",
							},
						},
						Consequence: &ast.BlockStatement{
							Token: token.Token{Type: token.LBRACE, Literal: "{"},
							Statements: []ast.Statement{
								&ast.ExpressionStatement{
									Token: token.Token{
										Type:    token.IDENT,
										Literal: "x",
									},
									Expression: &ast.Identifier{
										Token: token.Token{
											Type:    token.IDENT,
											Literal: "x",
										},
										Value: "x",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			input: "if (x < y) { x } else { y }",
			want: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.IF, Literal: "if"},
					Expression: &ast.IfExpression{
						Token: token.Token{Type: token.IF, Literal: "if"},
						Condition: &ast.InfixExpression{
							Token:    token.Token{Type: token.LT, Literal: "<"},
							Operator: "<",
							Left: &ast.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "x"},
								Value: "x",
							},
							Right: &ast.Identifier{
								Token: token.Token{Type: token.IDENT, Literal: "y"},
								Value: "y",
							},
						},
						Consequence: &ast.BlockStatement{
							Token: token.Token{Type: token.LBRACE, Literal: "{"},
							Statements: []ast.Statement{
								&ast.ExpressionStatement{
									Token: token.Token{
										Type:    token.IDENT,
										Literal: "x",
									},
									Expression: &ast.Identifier{
										Token: token.Token{
											Type:    token.IDENT,
											Literal: "x",
										},
										Value: "x",
									},
								},
							},
						},
						Alternative: &ast.BlockStatement{
							Token: token.Token{Type: token.LBRACE, Literal: "{"},
							Statements: []ast.Statement{
								&ast.ExpressionStatement{
									Token: token.Token{
										Type:    token.IDENT,
										Literal: "y",
									},
									Expression: &ast.Identifier{
										Token: token.Token{
											Type:    token.IDENT,
											Literal: "y",
										},
										Value: "y",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := parser.New(lexer.New(tt.input))
			program := p.Parse()
			require.Empty(t, p.Errors())
			require.NotNil(t, program)
			assert.Equal(t, tt.want, program.Statements)
		})
	}
}
