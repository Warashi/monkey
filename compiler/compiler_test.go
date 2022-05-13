package compiler_test

import (
	"bytes"
	"testing"

	"github.com/Warashi/monkey/code"
	"github.com/Warashi/monkey/compiler"
	"github.com/Warashi/monkey/lexer"
	"github.com/Warashi/monkey/object"
	"github.com/Warashi/monkey/parser"
	. "github.com/Warashi/monkey/testutil"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testcase struct {
	name  string
	input string
	want  compiler.Bytecode
}

func TestIntegerArithmetric(t *testing.T) {
	t.Parallel()
	type (
		c = compiler.Bytecode
		i = code.Instructions
		o = []object.Object
	)
	var (
		cat   = ConcatInstructions
		instr = MakeInstructions
		int   = IntegerObject
	)

	tests := []testcase{
		{"1+2", "1 + 2", c{cat(instr(t, code.OpConstant, 0), instr(t, code.OpConstant, 1), instr(t, code.OpAdd)), o{int(1), int(2)}}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			program := parser.New(lexer.New(tt.input)).Parse()

			compiler := compiler.New()
			require.NoError(t, compiler.Compile(program))
			if want, got := tt.want, compiler.Bytecode(); !cmp.Equal(want, got) {
				t.Error(cmp.Diff(want, got))
			}
		})
	}
}

func TestInstructionsString(t *testing.T) {
	want := `0000 OpAdd
0001 OpConstant 2
0004 OpConstant 65535
`

	got := ConcatInstructions(
		MakeInstructions(t, code.OpAdd),
		MakeInstructions(t, code.OpConstant, 2),
		MakeInstructions(t, code.OpConstant, 65535),
	)

	assert.Equal(t, want, got.String())
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		name      string
		op        code.Opcode
		operands  []int64
		bytesRead int
	}{}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			def, err := code.Lookup(tt.op)
			require.NoError(t, err)

			instr := MakeInstructions(t, tt.op, tt.operands...)
			read, n, err := code.ReadOperands(def, bytes.NewReader(instr[1:]))
			assert.NoError(t, err)
			assert.Equal(t, tt.bytesRead, n)
			assert.Equal(t, tt.operands, read)
		})
	}
}
