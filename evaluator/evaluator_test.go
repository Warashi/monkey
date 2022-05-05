package evaluator_test

import (
	"testing"

	"github.com/Warashi/implement-interpreter-with-go/evaluator"
	"github.com/Warashi/implement-interpreter-with-go/lexer"
	"github.com/Warashi/implement-interpreter-with-go/object"
	"github.com/Warashi/implement-interpreter-with-go/parser"
	. "github.com/Warashi/implement-interpreter-with-go/testutil"
	"github.com/stretchr/testify/assert"
)

func Eval(input string) object.Object {
	return evaluator.Eval(parser.New(lexer.New(input)).Parse(), object.NewEnvironment())
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input string
		want  object.Object
	}{
		{input: "5", want: IntegerObject(5)},
		{input: "10", want: IntegerObject(10)},
		{input: "-5", want: IntegerObject(-5)},
		{input: "-10", want: IntegerObject(-10)},
		{input: "5 + 5 + 5 + 5 - 10", want: IntegerObject(10)},
		{input: "2 * 2 * 2 * 2 * 2", want: IntegerObject(32)},
		{input: "-50 + 100 + -50", want: IntegerObject(0)},
		{input: "5 * 2 + 10", want: IntegerObject(20)},
		{input: "5 + 2 * 10", want: IntegerObject(25)},
		{input: "20 + 2 * -10", want: IntegerObject(0)},
		{input: "50 / 2 * 2 + 10", want: IntegerObject(60)},
		{input: "2 * (5 + 10)", want: IntegerObject(30)},
		{input: "3 * 3 * 3 + 10", want: IntegerObject(37)},
		{input: "3 * (3 * 3) + 10", want: IntegerObject(37)},
		{input: "(5 + 10 * 2 + 15 / 3) * 2 + -10", want: IntegerObject(50)},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, Eval(tt.input))
		})
	}
}

func TestEvalStringExpression(t *testing.T) {
	tests := []struct {
		input string
		want  object.Object
	}{
		{input: `"Hello, world!"`, want: StringObject("Hello, world!")},
		{input: `"Hello," + " " + "world!"`, want: StringObject("Hello, world!")},
		{input: `"a" < "b"`, want: BooleanObject(true)},
		{input: `"a" > "b"`, want: BooleanObject(false)},
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
		{input: "true", want: BooleanObject(true)},
		{input: "false", want: BooleanObject(false)},
		{input: "true", want: BooleanObject(true)},
		{input: "false", want: BooleanObject(false)},
		{input: "1 < 2", want: BooleanObject(true)},
		{input: "1 > 2", want: BooleanObject(false)},
		{input: "1 < 1", want: BooleanObject(false)},
		{input: "1 > 1", want: BooleanObject(false)},
		{input: "1 == 1", want: BooleanObject(true)},
		{input: "1 != 1", want: BooleanObject(false)},
		{input: "1 == 2", want: BooleanObject(false)},
		{input: "1 != 2", want: BooleanObject(true)},
		{input: "true == true", want: BooleanObject(true)},
		{input: "false == false", want: BooleanObject(true)},
		{input: "true == false", want: BooleanObject(false)},
		{input: "true != false", want: BooleanObject(true)},
		{input: "false != true", want: BooleanObject(true)},
		{input: "(1 < 2) == true", want: BooleanObject(true)},
		{input: "(1 < 2) == false", want: BooleanObject(false)},
		{input: "(1 > 2) == true", want: BooleanObject(false)},
		{input: "(1 > 2) == false", want: BooleanObject(true)},
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
		{input: "!true", want: BooleanObject(false)},
		{input: "!false", want: BooleanObject(true)},
		{input: "!5", want: BooleanObject(false)},
		{input: "!!true", want: BooleanObject(true)},
		{input: "!!false", want: BooleanObject(false)},
		{input: "!!5", want: BooleanObject(true)},
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
		{"if (true) { 10 }", IntegerObject(10)},
		{"if (false) { 10 }", NullObject()},
		{"if (1) { 10 }", IntegerObject(10)},
		{"if (1 < 2) { 10 }", IntegerObject(10)},
		{"if (1 > 2) { 10 }", NullObject()},
		{"if (1 > 2) { 10 } else { 20 }", IntegerObject(20)},
		{"if (1 < 2) { 10 } else { 20 }", IntegerObject(10)},
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
		{input: "return 10;", want: IntegerObject(10)},
		{input: "return 10; 9;", want: IntegerObject(10)},
		{input: "return 2 * 5; 9;", want: IntegerObject(10)},
		{input: "9; return 2 * 5; 9;", want: IntegerObject(10)},
		{input: "if (10 > 1) { return 10; }", want: IntegerObject(10)},
		{input: "if (10 > 1) { if (10 > 1) { return 10; } return 1; }", want: IntegerObject(10)},
		{input: "let f = fn(x) { return x; x + 10; }; f(10);", want: IntegerObject(10)},
		{input: "let f = fn(x) { let result = x + 10; return result; return 10; }; f(10);", want: IntegerObject(20)},
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
		{input: "5 + true;", want: ErrorObject("type mismatch: Integer + Boolean")},
		{input: "5 + true; 5;", want: ErrorObject("type mismatch: Integer + Boolean")},
		{input: "-true", want: ErrorObject("unknown operator: -Boolean")},
		{input: "true + false;", want: ErrorObject("unknown operator: Boolean + Boolean")},
		{input: "true + false + true + false;", want: ErrorObject("unknown operator: Boolean + Boolean")},
		{input: "5; true + false; 5", want: ErrorObject("unknown operator: Boolean + Boolean")},
		{input: "if (10 > 1) { true + false; }", want: ErrorObject("unknown operator: Boolean + Boolean")},
		{input: "if (10 > 1) { if (10 > 1) { return true + false; } return 1; }", want: ErrorObject("unknown operator: Boolean + Boolean")},
		{input: "foobar", want: ErrorObject("identifier not found: foobar")},
		{input: `"Hello" - "world"`, want: ErrorObject("unknown operator: String - String")},
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
		{input: "let a = 5; a;", want: IntegerObject(5)},
		{input: "let a = 5 * 5; a;", want: IntegerObject(25)},
		{input: "let a = 5; let b = a; b;", want: IntegerObject(5)},
		{input: "let a = 5; let b = a; let c = a + b + 5; c;", want: IntegerObject(15)},
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
			want: FunctionObject(
				object.NewEnvironment(),
				BlockStatement(ExpressionStatement(InfixExpression(Plus, Identifier("x"), IntegerLiteral(2)))),
				Identifier("x"),
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
		{input: "let identity = fn(x) { x; }; identity(5);", want: IntegerObject(5)},
		{input: "let identity = fn(x) { return x; }; identity(5);", want: IntegerObject(5)},
		{input: "let double = fn(x) {  x * 2; }; double(5);", want: IntegerObject(10)},
		{input: "let add = fn(x, y) { x + y; }; add(5, 5);", want: IntegerObject(10)},
		{input: "let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", want: IntegerObject(20)},
		{input: "fn(x) { x; }(5);", want: IntegerObject(5)},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, Eval(tt.input))
		})
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input string
		want  object.Object
	}{
		{input: `len("")`, want: IntegerObject(0)},
		{input: `len("four")`, want: IntegerObject(4)},
		{input: `len("hello, world")`, want: IntegerObject(12)},
		{input: `len([1, 2])`, want: IntegerObject(2)},
		{input: `len(1)`, want: ErrorObject("argument to `len` not supported, got Integer")},
		{input: `len("one", "two")`, want: ErrorObject("wrong number of arguments. got=2, want=1")},
		{input: `len([])`, want: IntegerObject(0)},
		{input: `len([1])`, want: IntegerObject(1)},
		{input: `len([1, 2])`, want: IntegerObject(2)},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, Eval(tt.input))
		})
	}
}

func TestArrayLiteral(t *testing.T) {
	tests := []struct {
		input string
		want  object.Object
	}{
		{input: `[1, 2 * 3, 4 + 5]`, want: ArrayObject(IntegerObject(1), IntegerObject(6), IntegerObject(9))},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, Eval(tt.input))
		})
	}
}

func TestArrayIndexing(t *testing.T) {
	tests := []struct {
		input string
		want  object.Object
	}{
		{input: `[1][0]`, want: IntegerObject(1)},
		{input: `[0, 1][1]`, want: IntegerObject(1)},
		{input: `[0, 1, 2][1]`, want: IntegerObject(1)},
		{input: `let x = [1]; x[0]`, want: IntegerObject(1)},
		{input: `let x = [0, 1]; x[1]`, want: IntegerObject(1)},
		{input: `let x = [0, 1, 2]; x[1]`, want: IntegerObject(1)},
		{input: `let x = [0, 1, 2]; x[0] + x[1] + x[2]`, want: IntegerObject(3)},
		{input: `let i = 0; [1][i]`, want: IntegerObject(1)},
		{input: `[1, 2, 3][1 + 1]`, want: IntegerObject(3)},
		{input: `let x = [1, 2, 3]; let i = x[1]; x[i]`, want: IntegerObject(3)},
		{input: `let x = [1, 2, 3]; x[x[1]]`, want: IntegerObject(3)},
		{input: `[1, 2, 3][3]`, want: ErrorObject("index out of range. index=3, len=3")},
		{input: `[1, 2, 3][-1]`, want: ErrorObject("index out of range. index=-1, len=3")},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, Eval(tt.input))
		})
	}
}

func TestHashLiteral(t *testing.T) {
	tests := []struct {
		input string
		want  object.Object
	}{
		{
			input: `{"one": 1, "two": 2, 3: "three", 3+1:("fo" + "ur")}`,
			want: HashObject(map[object.Hashable]object.Object{
				StringObject("one"): IntegerObject(1),
				StringObject("two"): IntegerObject(2),
				IntegerObject(3):    StringObject("three"),
				IntegerObject(4):    StringObject("four"),
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, Eval(tt.input))
		})
	}
}

func TestHashIndexing(t *testing.T) {
	tests := []struct {
		input string
		want  object.Object
	}{
		{input: `{"one": 1, "two": 2, 3: "three", 3+1:("fo" + "ur")}["one"]`, want: IntegerObject(1)},
		{input: `{"one": 1, "two": 2, 3: "three", 3+1:("fo" + "ur")}["t" + "wo"]`, want: IntegerObject(2)},
		{input: `{"one": 1, "two": 2, 3: "three", 3+1:("fo" + "ur")}[1+2]`, want: StringObject("three")},
		{input: `let m = {"one": 1, "two": 2, 3: "three", 3+1:("fo" + "ur")}; m[1+2]`, want: StringObject("three")},
		{input: `let m = {"one": 1, 1: "ONE", "two": 2, 3: "three", 3+1:("fo" + "ur")};let k = m["one"] m[k]`, want: StringObject("ONE")},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, Eval(tt.input))
		})
	}
}
