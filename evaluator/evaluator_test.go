package evaluator_test

import (
	"testing"

	"github.com/Warashi/implement-interpreter-with-go/evaluator"
	"github.com/Warashi/implement-interpreter-with-go/lexer"
	"github.com/Warashi/implement-interpreter-with-go/object"
	"github.com/Warashi/implement-interpreter-with-go/parser"
	"github.com/Warashi/implement-interpreter-with-go/testutil"
	"github.com/Warashi/implement-interpreter-with-go/token"
	"github.com/stretchr/testify/assert"
)

var (
	Plus = token.Token{Type: token.PLUS, Literal: "+"}
)

func Eval(input string) object.Object {
	return evaluator.Eval(parser.New(lexer.New(input)).Parse(), object.NewEnvironment())
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input string
		want  object.Object
	}{
		{input: "5", want: testutil.IntegerObject(5)},
		{input: "10", want: testutil.IntegerObject(10)},
		{input: "-5", want: testutil.IntegerObject(-5)},
		{input: "-10", want: testutil.IntegerObject(-10)},
		{input: "5 + 5 + 5 + 5 - 10", want: testutil.IntegerObject(10)},
		{input: "2 * 2 * 2 * 2 * 2", want: testutil.IntegerObject(32)},
		{input: "-50 + 100 + -50", want: testutil.IntegerObject(0)},
		{input: "5 * 2 + 10", want: testutil.IntegerObject(20)},
		{input: "5 + 2 * 10", want: testutil.IntegerObject(25)},
		{input: "20 + 2 * -10", want: testutil.IntegerObject(0)},
		{input: "50 / 2 * 2 + 10", want: testutil.IntegerObject(60)},
		{input: "2 * (5 + 10)", want: testutil.IntegerObject(30)},
		{input: "3 * 3 * 3 + 10", want: testutil.IntegerObject(37)},
		{input: "3 * (3 * 3) + 10", want: testutil.IntegerObject(37)},
		{input: "(5 + 10 * 2 + 15 / 3) * 2 + -10", want: testutil.IntegerObject(50)},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, Eval(tt.input))
		})
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input string
		want  object.Object
	}{
		{input: "true", want: testutil.BooleanObject(true)},
		{input: "false", want: testutil.BooleanObject(false)},
		{input: "true", want: testutil.BooleanObject(true)},
		{input: "false", want: testutil.BooleanObject(false)},
		{input: "1 < 2", want: testutil.BooleanObject(true)},
		{input: "1 > 2", want: testutil.BooleanObject(false)},
		{input: "1 < 1", want: testutil.BooleanObject(false)},
		{input: "1 > 1", want: testutil.BooleanObject(false)},
		{input: "1 == 1", want: testutil.BooleanObject(true)},
		{input: "1 != 1", want: testutil.BooleanObject(false)},
		{input: "1 == 2", want: testutil.BooleanObject(false)},
		{input: "1 != 2", want: testutil.BooleanObject(true)},
		{input: "true == true", want: testutil.BooleanObject(true)},
		{input: "false == false", want: testutil.BooleanObject(true)},
		{input: "true == false", want: testutil.BooleanObject(false)},
		{input: "true != false", want: testutil.BooleanObject(true)},
		{input: "false != true", want: testutil.BooleanObject(true)},
		{input: "(1 < 2) == true", want: testutil.BooleanObject(true)},
		{input: "(1 < 2) == false", want: testutil.BooleanObject(false)},
		{input: "(1 > 2) == true", want: testutil.BooleanObject(false)},
		{input: "(1 > 2) == false", want: testutil.BooleanObject(true)},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, Eval(tt.input))
		})
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input string
		want  object.Object
	}{
		{input: "!true", want: testutil.BooleanObject(false)},
		{input: "!false", want: testutil.BooleanObject(true)},
		{input: "!5", want: testutil.BooleanObject(false)},
		{input: "!!true", want: testutil.BooleanObject(true)},
		{input: "!!false", want: testutil.BooleanObject(false)},
		{input: "!!5", want: testutil.BooleanObject(true)},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, Eval(tt.input))
		})
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input string
		want  object.Object
	}{
		{"if (true) { 10 }", testutil.IntegerObject(10)},
		{"if (false) { 10 }", testutil.NullObject()},
		{"if (1) { 10 }", testutil.IntegerObject(10)},
		{"if (1 < 2) { 10 }", testutil.IntegerObject(10)},
		{"if (1 > 2) { 10 }", testutil.NullObject()},
		{"if (1 > 2) { 10 } else { 20 }", testutil.IntegerObject(20)},
		{"if (1 < 2) { 10 } else { 20 }", testutil.IntegerObject(10)},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, Eval(tt.input))
		})
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input string
		want  object.Object
	}{
		{input: "return 10;", want: testutil.IntegerObject(10)},
		{input: "return 10; 9;", want: testutil.IntegerObject(10)},
		{input: "return 2 * 5; 9;", want: testutil.IntegerObject(10)},
		{input: "9; return 2 * 5; 9;", want: testutil.IntegerObject(10)},
		{input: "if (10 > 1) { return 10; }", want: testutil.IntegerObject(10)},
		{input: "if (10 > 1) { if (10 > 1) { return 10; } return 1; }", want: testutil.IntegerObject(10)},
		{input: "let f = fn(x) { return x; x + 10; }; f(10);", want: testutil.IntegerObject(10)},
		{input: "let f = fn(x) { let result = x + 10; return result; return 10; }; f(10);", want: testutil.IntegerObject(20)},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, Eval(tt.input))
		})
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input string
		want  object.Object
	}{
		{input: "5 + true;", want: testutil.ErrorObject("type mismatch: Integer + Boolean")},
		{input: "5 + true; 5;", want: testutil.ErrorObject("type mismatch: Integer + Boolean")},
		{input: "-true", want: testutil.ErrorObject("unknown operator: -Boolean")},
		{input: "true + false;", want: testutil.ErrorObject("unknown operator: Boolean + Boolean")},
		{input: "true + false + true + false;", want: testutil.ErrorObject("unknown operator: Boolean + Boolean")},
		{input: "5; true + false; 5", want: testutil.ErrorObject("unknown operator: Boolean + Boolean")},
		{input: "if (10 > 1) { true + false; }", want: testutil.ErrorObject("unknown operator: Boolean + Boolean")},
		{input: "if (10 > 1) { if (10 > 1) { return true + false; } return 1; }", want: testutil.ErrorObject("unknown operator: Boolean + Boolean")},
		{input: "foobar", want: testutil.ErrorObject("identifier not found: foobar")},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, Eval(tt.input))
		})
	}
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input string
		want  object.Object
	}{
		{input: "let a = 5; a;", want: testutil.IntegerObject(5)},
		{input: "let a = 5 * 5; a;", want: testutil.IntegerObject(25)},
		{input: "let a = 5; let b = a; b;", want: testutil.IntegerObject(5)},
		{input: "let a = 5; let b = a; let c = a + b + 5; c;", want: testutil.IntegerObject(15)},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, Eval(tt.input))
		})
	}
}

func TestFunctionObject(t *testing.T) {
	tests := []struct {
		input string
		want  object.Object
	}{
		{
			input: "fn(x) { x + 2; };",
			want: testutil.FunctionObject(
				object.NewEnvironment(),
				testutil.BlockStatement(testutil.ExpressionStatement(testutil.InfixExpression(Plus, testutil.Identifier("x"), testutil.IntegerLiteral(2)))),
				testutil.Identifier("x"),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, Eval(tt.input))
		})
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input string
		want  object.Object
	}{
		{input: "let identity = fn(x) { x; }; identity(5);", want: testutil.IntegerObject(5)},
		{input: "let identity = fn(x) { return x; }; identity(5);", want: testutil.IntegerObject(5)},
		{input: "let double = fn(x) {  x * 2; }; double(5);", want: testutil.IntegerObject(10)},
		{input: "let add = fn(x, y) { x + y; }; add(5, 5);", want: testutil.IntegerObject(10)},
		{input: "let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", want: testutil.IntegerObject(20)},
		{input: "fn(x) { x; }(5);", want: testutil.IntegerObject(5)},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, Eval(tt.input))
		})
	}
}
