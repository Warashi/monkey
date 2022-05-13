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
		{"plus", "1 + 2", IntegerObject(2)}, // FIXME
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			compiler := compiler.New()
			require.NoError(t, compiler.Compile(parser.New(lexer.New(tt.input)).Parse()))

			vm := vm.New(compiler.Bytecode())
			require.NoError(t, vm.Run())

			assert.Equal(t, tt.want, vm.StackTop())
		})
	}
}
