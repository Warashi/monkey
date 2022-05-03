package evaluator_test

import (
	"testing"

	"github.com/Warashi/implement-interpreter-with-go/evaluator"
	"github.com/Warashi/implement-interpreter-with-go/lexer"
	"github.com/Warashi/implement-interpreter-with-go/object"
	"github.com/Warashi/implement-interpreter-with-go/parser"
	"github.com/stretchr/testify/assert"
)

func Eval(t *testing.T, input string) object.Object {
	t.Helper()
	return evaluator.Eval(parser.New(lexer.New(input)).Parse())
}

func IntegerObject(t *testing.T, val int64) object.Object {
	t.Helper()
	return object.Integer{Value: val}
}

func BooleanObject(t *testing.T, val bool) object.Object {
	t.Helper()
	return object.Boolean{Value: val}
}

func NullObject(t *testing.T) object.Object {
	return object.Null{}
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input string
		want  object.Object
	}{
		{input: "5", want: IntegerObject(t, 5)},
		{input: "10", want: IntegerObject(t, 10)},
		{input: "-5", want: IntegerObject(t, -5)},
		{input: "-10", want: IntegerObject(t, -10)},
		{input: "5 + 5 + 5 + 5 - 10", want: IntegerObject(t, 10)},
		{input: "2 * 2 * 2 * 2 * 2", want: IntegerObject(t, 32)},
		{input: "-50 + 100 + -50", want: IntegerObject(t, 0)},
		{input: "5 * 2 + 10", want: IntegerObject(t, 20)},
		{input: "5 + 2 * 10", want: IntegerObject(t, 25)},
		{input: "20 + 2 * -10", want: IntegerObject(t, 0)},
		{input: "50 / 2 * 2 + 10", want: IntegerObject(t, 60)},
		{input: "2 * (5 + 10)", want: IntegerObject(t, 30)},
		{input: "3 * 3 * 3 + 10", want: IntegerObject(t, 37)},
		{input: "3 * (3 * 3) + 10", want: IntegerObject(t, 37)},
		{input: "(5 + 10 * 2 + 15 / 3) * 2 + -10", want: IntegerObject(t, 50)},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, Eval(t, tt.input))
		})
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input string
		want  object.Object
	}{
		{input: "true", want: BooleanObject(t, true)},
		{input: "false", want: BooleanObject(t, false)},
		{input: "true", want: BooleanObject(t, true)},
		{input: "false", want: BooleanObject(t, false)},
		{input: "1 < 2", want: BooleanObject(t, true)},
		{input: "1 > 2", want: BooleanObject(t, false)},
		{input: "1 < 1", want: BooleanObject(t, false)},
		{input: "1 > 1", want: BooleanObject(t, false)},
		{input: "1 == 1", want: BooleanObject(t, true)},
		{input: "1 != 1", want: BooleanObject(t, false)},
		{input: "1 == 2", want: BooleanObject(t, false)},
		{input: "1 != 2", want: BooleanObject(t, true)},
		{input: "true == true", want: BooleanObject(t, true)},
		{input: "false == false", want: BooleanObject(t, true)},
		{input: "true == false", want: BooleanObject(t, false)},
		{input: "true != false", want: BooleanObject(t, true)},
		{input: "false != true", want: BooleanObject(t, true)},
		{input: "(1 < 2) == true", want: BooleanObject(t, true)},
		{input: "(1 < 2) == false", want: BooleanObject(t, false)},
		{input: "(1 > 2) == true", want: BooleanObject(t, false)},
		{input: "(1 > 2) == false", want: BooleanObject(t, true)},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, Eval(t, tt.input))
		})
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input string
		want  object.Object
	}{
		{input: "!true", want: BooleanObject(t, false)},
		{input: "!false", want: BooleanObject(t, true)},
		{input: "!5", want: BooleanObject(t, false)},
		{input: "!!true", want: BooleanObject(t, true)},
		{input: "!!false", want: BooleanObject(t, false)},
		{input: "!!5", want: BooleanObject(t, true)},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, Eval(t, tt.input))
		})
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input string
		want  object.Object
	}{
		{"if (true) { 10 }", IntegerObject(t, 10)},
		{"if (false) { 10 }", NullObject(t)},
		{"if (1) { 10 }", IntegerObject(t, 10)},
		{"if (1 < 2) { 10 }", IntegerObject(t, 10)},
		{"if (1 > 2) { 10 }", NullObject(t)},
		{"if (1 > 2) { 10 } else { 20 }", IntegerObject(t, 20)},
		{"if (1 < 2) { 10 } else { 20 }", IntegerObject(t, 10)},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, Eval(t, tt.input))
		})
	}
}
