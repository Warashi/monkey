package vm_test

import (
	"testing"

	"github.com/Warashi/monkey/compiler"
	"github.com/Warashi/monkey/lexer"
	"github.com/Warashi/monkey/object"
	"github.com/Warashi/monkey/parser"
	. "github.com/Warashi/monkey/testutil"
	"github.com/Warashi/monkey/vm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testcase struct {
	name  string
	input string
	want  object.Object
}

func TestIntegerArithmetric(t *testing.T) {
	t.Parallel()
	tests := []testcase{
		{"constant/1", "1", IntegerObject(1)},
		{"constant/2", "2", IntegerObject(2)},
		{"plus", "1 + 2", IntegerObject(3)},
		{"minus", "1 - 2", IntegerObject(-1)},
		{"asterisk", "1 * 2", IntegerObject(2)},
		{"slash", "6 / 2", IntegerObject(3)},
		{"multi-calculation/1", "50 / 2 * 2 + 10 - 5", IntegerObject(55)},
		{"multi-calculation/2", "5 + 5 + 5 + 5 - 10", IntegerObject(10)},
		{"multi-calculation/3", "2 * 2 * 2 * 2 * 2", IntegerObject(32)},
		{"multi-calculation/4", "5 * 2 + 10", IntegerObject(20)},
		{"multi-calculation/5", "5 + 2 * 10", IntegerObject(25)},
		{"multi-calculation/6", "5 * (2 + 10)", IntegerObject(60)},
		{"prefix-op/minus/-5", "-5", IntegerObject(-5)},
		{"prefix-op/minus/-10", "-10", IntegerObject(-10)},
		{"multi-calculation/7", "-50 + 100 + -50", IntegerObject(0)},
		{"multi-calculation/8", "(5 + 10 * 2 + 15 / 3) * 2 + -10", IntegerObject(50)},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			compiler := compiler.New()
			require.NoError(t, compiler.Compile(parser.New(lexer.New(tt.input)).Parse()))

			vm := vm.New(compiler.Bytecode())
			require.NoError(t, vm.Run())

			assert.Equal(t, tt.want, vm.LastPopedStackElem())
		})
	}
}

func TestBooleanExpressions(t *testing.T) {
	t.Parallel()
	tests := []testcase{
		{"true", "true", BooleanObject(true)},
		{"false", "false", BooleanObject(false)},
		{"lt/int/true", "1 < 2", BooleanObject(true)},
		{"lt/int/false", "1 < 1", BooleanObject(false)},
		{"gt/int/true", "2 > 1", BooleanObject(true)},
		{"gt/int/false", "1 > 1", BooleanObject(false)},
		{"eq/int/true", "1 == 1", BooleanObject(true)},
		{"eq/int/false", "1 == 2", BooleanObject(false)},
		{"neq/int/true", "1 != 2", BooleanObject(true)},
		{"neq/int/false", "1 != 1", BooleanObject(false)},
		{"eq/bool/true", "true == true", BooleanObject(true)},
		{"eq/bool/false", "true == false", BooleanObject(false)},
		{"neq/bool/true", "true != false", BooleanObject(true)},
		{"neq/bool/false", "true != true", BooleanObject(false)},
		{"composite/true", "(1 < 2) == true", BooleanObject(true)},
		{"composite/false", "(1 < 2) == false", BooleanObject(false)},
		{"composite/true", "(1 > 2) == false", BooleanObject(true)},
		{"composite/false", "(1 > 2) == true", BooleanObject(false)},
		{"prefix/bang/true", "!false", BooleanObject(true)},
		{"prefix/bang/false", "!true", BooleanObject(false)},
		{"prefix/bang/true", "!!true", BooleanObject(true)},
		{"prefix/bang/false", "!!false", BooleanObject(false)},
		{"prefix/bang/false", "!5", BooleanObject(false)},
		{"prefix/bang/true", "!!5", BooleanObject(true)},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			compiler := compiler.New()
			require.NoError(t, compiler.Compile(parser.New(lexer.New(tt.input)).Parse()))

			vm := vm.New(compiler.Bytecode())
			require.NoError(t, vm.Run())

			assert.Equal(t, tt.want, vm.LastPopedStackElem())
		})
	}
}
