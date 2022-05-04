package parser_test

import (
	"strconv"
	"testing"

	"github.com/Warashi/implement-interpreter-with-go/ast"
	"github.com/Warashi/implement-interpreter-with-go/lexer"
	"github.com/Warashi/implement-interpreter-with-go/parser"
	"github.com/Warashi/implement-interpreter-with-go/parser/testdata"
	. "github.com/Warashi/implement-interpreter-with-go/testutil"
	"github.com/Warashi/implement-interpreter-with-go/token"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLetStatement(t *testing.T) {
	p := parser.New(lexer.New(testdata.Let))
	program := p.Parse()
	require.Empty(t, p.Errors())
	require.NotNil(t, program)

	wants := []ast.Statement{
		LetStatement(Identifier("x"), IntegerLiteral(5)),
		LetStatement(Identifier("y"), IntegerLiteral(10)),
		LetStatement(Identifier("foobar"), IntegerLiteral(838383)),
	}
	assert.Equal(t, wants, program.Statements)
}

func TestReturnStatement(t *testing.T) {
	p := parser.New(lexer.New(testdata.Return))
	program := p.Parse()
	require.Empty(t, p.Errors())
	require.NotNil(t, program)

	wants := []ast.Statement{
		ReturnStatement(IntegerLiteral(5)),
		ReturnStatement(IntegerLiteral(10)),
		ReturnStatement(IntegerLiteral(993322)),
	}
	assert.Equal(t, wants, program.Statements)
}

func TestIdentifierExpression(t *testing.T) {
	p := parser.New(lexer.New(testdata.IdentifierExpression))
	program := p.Parse()
	require.Empty(t, p.Errors())
	require.NotNil(t, program)

	wants := []ast.Statement{ExpressionStatement(Identifier("foobar"))}
	assert.Equal(t, wants, program.Statements)
}

func TestIntegerLiteralExpression(t *testing.T) {
	p := parser.New(lexer.New(testdata.IntegerLiteralExpression))
	program := p.Parse()
	require.Empty(t, p.Errors())
	require.NotNil(t, program)

	wants := []ast.Statement{ExpressionStatement(IntegerLiteral(5))}
	assert.Equal(t, wants, program.Statements)
}

func TestStringLiteralExpression(t *testing.T) {
	p := parser.New(lexer.New(testdata.StringLiteralExpression))
	program := p.Parse()
	require.Empty(t, p.Errors())
	require.NotNil(t, program)

	wants := []ast.Statement{ExpressionStatement(StringLiteral("hello, world"))}
	assert.Equal(t, wants, program.Statements)
}

func TestBooleanLiteralExpression(t *testing.T) {
	p := parser.New(lexer.New(testdata.BooleanLiteralExpression))
	program := p.Parse()
	require.Empty(t, p.Errors())
	require.NotNil(t, program)

	true := ExpressionStatement(True)
	false := ExpressionStatement(False)
	wants := []ast.Statement{true, false, true, false}
	assert.Equal(t, wants, program.Statements)
}

func TestFunctionLiteralExpression(t *testing.T) {
	p := parser.New(lexer.New(testdata.FunctionLiteralExpression))
	program := p.Parse()
	require.Empty(t, p.Errors())
	require.NotNil(t, program)

	x := Identifier("x")
	y := Identifier("y")
	fn := FunctionLiteral(
		BlockStatement(ExpressionStatement(InfixExpression(Plus, x, y))),
		x, y)
	wants := []ast.Statement{ExpressionStatement(fn)}
	assert.Equal(t, wants, program.Statements)
}

func TestFunctionParameterParsing(t *testing.T) {
	buildWant := func(params ...string) []ast.Statement {
		var ids []*ast.Identifier
		for _, p := range params {
			ids = append(ids, Identifier(p))
		}
		return []ast.Statement{ExpressionStatement(FunctionLiteral(BlockStatement(), ids...))}
	}

	tests := []struct {
		input string
		want  []ast.Statement
	}{
		{input: "fn() {};", want: buildWant()},
		{input: "fn(x) {};", want: buildWant("x")},
		{input: "fn(x, y) {};", want: buildWant("x", "y")},
		{input: "fn(x, y, z) {};", want: buildWant("x", "y", "z")},
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

func TestArrayLiteralParsing(t *testing.T) {
	tests := []struct {
		input string
		want  []ast.Statement
	}{
		{
			input: "[1, 2 * 3, 4 + 5]",
			want: []ast.Statement{ExpressionStatement(
				ArrayLiteral(
					IntegerLiteral(1),
					InfixExpression(Asterisk, IntegerLiteral(2), IntegerLiteral(3)),
					InfixExpression(Plus, IntegerLiteral(4), IntegerLiteral(5)),
				),
			)},
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

func TestHashLiteralParsing(t *testing.T) {
	tests := []struct {
		input string
		want  []ast.Statement
	}{
		{
			input: `{"one": 1, "two": 2, "three": 3}`,
			want: []ast.Statement{ExpressionStatement(
				HashLiteral(map[ast.Expression]ast.Expression{
					StringLiteral("one"):   IntegerLiteral(1),
					StringLiteral("two"):   IntegerLiteral(2),
					StringLiteral("three"): IntegerLiteral(3),
				}),
			)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := parser.New(lexer.New(tt.input))
			program := p.Parse()
			require.Empty(t, p.Errors())
			require.NotNil(t, program)
			if !cmp.Equal(tt.want, program.Statements) {
				t.Errorf("got=%v, want=%v, diff=%v", program.Statements, tt.want, cmp.Diff(tt.want, program.Statements))
			}
		})
	}
}

func TestPrefixExpression(t *testing.T) {
	pre := func(op token.Token, v int64) ast.Statement {
		return ExpressionStatement(PrefixExpression(op, IntegerLiteral(v)))
	}
	tests := []struct {
		input string
		want  []ast.Statement
	}{
		{input: "!5", want: []ast.Statement{pre(Bang, 5)}},
		{input: "-15", want: []ast.Statement{pre(Minus, 15)}},
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
	in := func(op token.Token, left, right int64) ast.Statement {
		return ExpressionStatement(InfixExpression(op, IntegerLiteral(left), IntegerLiteral(right)))
	}
	tests := []struct {
		input string
		want  []ast.Statement
	}{
		{input: "5 + 6", want: []ast.Statement{in(Plus, 5, 6)}},
		{input: "5 - 6", want: []ast.Statement{in(Minus, 5, 6)}},
		{input: "5 * 6", want: []ast.Statement{in(Asterisk, 5, 6)}},
		{input: "5 / 6", want: []ast.Statement{in(Slash, 5, 6)}},
		{input: "5 > 6", want: []ast.Statement{in(GT, 5, 6)}},
		{input: "5 < 6", want: []ast.Statement{in(LT, 5, 6)}},
		{input: "5 == 6", want: []ast.Statement{in(Equal, 5, 6)}},
		{input: "5 != 6", want: []ast.Statement{in(NotEqual, 5, 6)}},
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
		{input: "a * [1, 2, 3, 4][b * c] * d", want: "((a * ([1, 2, 3, 4][(b * c)])) * d)"},
		{input: "add(a * b[2], b[1], 2 * [1, 2][1])", want: "add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))"},
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
			input: "if (x < y) { x }",
			want: []ast.Statement{
				ExpressionStatement(
					IfExpression(
						InfixExpression(LT, Identifier("x"), Identifier("y")),
						BlockStatement(ExpressionStatement(Identifier("x"))),
						nil,
					),
				),
			},
		},
		{
			input: "if (x < y) { x } else { y }",
			want: []ast.Statement{
				ExpressionStatement(
					IfExpression(
						InfixExpression(LT, Identifier("x"), Identifier("y")),
						BlockStatement(ExpressionStatement(Identifier("x"))),
						BlockStatement(ExpressionStatement(Identifier("y"))),
					),
				),
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

func TestCallExpression(t *testing.T) {
	buildWant := func(fn string, args ...ast.Expression) ast.Statement {
		return ExpressionStatement(CallExpression(Identifier(fn), args...))
	}

	tests := []struct {
		input string
		want  []ast.Statement
	}{
		{
			input: "add(1, 2 * 3, 4 + 5)",
			want: []ast.Statement{
				buildWant("add",
					IntegerLiteral(1),
					InfixExpression(Asterisk, IntegerLiteral(2), IntegerLiteral(3)),
					InfixExpression(Plus, IntegerLiteral(4), IntegerLiteral(5))),
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

func TestIndexExpressionParsing(t *testing.T) {
	tests := []struct {
		input string
		want  []ast.Statement
	}{
		{
			input: "myArray[1 + 2]",
			want: []ast.Statement{ExpressionStatement(
				IndexExpression(Identifier("myArray"), InfixExpression(Plus, IntegerLiteral(1), IntegerLiteral(2))),
			)},
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
