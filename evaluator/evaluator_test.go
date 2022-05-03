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
	return &object.Integer{Value: val}
}

func BooleanObject(t *testing.T, val bool) object.Object {
	t.Helper()
	return &object.Boolean{Value: val}
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
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, Eval(t, tt.input))
		})
	}
}
