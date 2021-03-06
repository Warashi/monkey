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
	type ()
	var (
		cat   = ConcatInstructions
		instr = MakeInstructions
		int   = IntegerObject
	)

	tests := []testcase{
		{
			name:  "plus",
			input: "1 + 2",
			want: compiler.Bytecode{
				Instructions: cat(
					instr(t, code.OpConstant, 0),
					instr(t, code.OpConstant, 1),
					instr(t, code.OpAdd),
					instr(t, code.OpPop),
				),
				Constants: []object.Object{
					int(1),
					int(2),
				},
			},
		},
		{
			name:  "minus",
			input: "1 - 2",
			want: compiler.Bytecode{
				Instructions: cat(
					instr(t, code.OpConstant, 0),
					instr(t, code.OpConstant, 1),
					instr(t, code.OpSub),
					instr(t, code.OpPop),
				),
				Constants: []object.Object{
					int(1),
					int(2),
				},
			},
		},
		{
			name:  "asterisk",
			input: "1 * 2",
			want: compiler.Bytecode{
				Instructions: cat(
					instr(t, code.OpConstant, 0),
					instr(t, code.OpConstant, 1),
					instr(t, code.OpMul),
					instr(t, code.OpPop),
				),
				Constants: []object.Object{
					int(1),
					int(2),
				},
			},
		},
		{
			name:  "slash",
			input: "1 / 2",
			want: compiler.Bytecode{
				Instructions: cat(
					instr(t, code.OpConstant, 0),
					instr(t, code.OpConstant, 1),
					instr(t, code.OpDiv),
					instr(t, code.OpPop),
				),
				Constants: []object.Object{
					int(1),
					int(2),
				},
			},
		},
		{
			name:  "semicolon",
			input: "1; 2",
			want: compiler.Bytecode{
				Instructions: cat(
					instr(t, code.OpConstant, 0),
					instr(t, code.OpPop),
					instr(t, code.OpConstant, 1),
					instr(t, code.OpPop),
				),
				Constants: []object.Object{
					int(1),
					int(2),
				},
			},
		},
		{
			name:  "minus",
			input: "-1",
			want: compiler.Bytecode{
				Instructions: cat(
					instr(t, code.OpConstant, 0),
					instr(t, code.OpMinus),
					instr(t, code.OpPop),
				),
				Constants: []object.Object{
					int(1),
				},
			},
		},
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

func TestBooleanExpressions(t *testing.T) {
	t.Parallel()
	var (
		cat   = ConcatInstructions
		instr = MakeInstructions
		int   = IntegerObject
	)

	tests := []testcase{
		{"true", "true", compiler.Bytecode{cat(instr(t, code.OpTrue), instr(t, code.OpPop)), nil}},
		{"false", "false", compiler.Bytecode{cat(instr(t, code.OpFalse), instr(t, code.OpPop)), nil}},
		{"gt", "1 > 2", compiler.Bytecode{
			Instructions: cat(
				instr(t, code.OpConstant, 0),
				instr(t, code.OpConstant, 1),
				instr(t, code.OpGreaterThan),
				instr(t, code.OpPop),
			),
			Constants: []object.Object{
				int(1),
				int(2),
			},
		}},
		{"lt", "1 < 2", compiler.Bytecode{
			Instructions: cat(
				instr(t, code.OpConstant, 0),
				instr(t, code.OpConstant, 1),
				instr(t, code.OpGreaterThan),
				instr(t, code.OpPop),
			),
			Constants: []object.Object{
				int(2),
				int(1),
			},
		}},
		{"eq", "1 == 2", compiler.Bytecode{
			Instructions: cat(
				instr(t, code.OpConstant, 0),
				instr(t, code.OpConstant, 1),
				instr(t, code.OpEqual),
				instr(t, code.OpPop),
			),
			Constants: []object.Object{
				int(1),
				int(2),
			},
		}},
		{"neq", "1 != 2", compiler.Bytecode{
			Instructions: cat(
				instr(t, code.OpConstant, 0),
				instr(t, code.OpConstant, 1),
				instr(t, code.OpNotEqual),
				instr(t, code.OpPop),
			),
			Constants: []object.Object{
				int(1),
				int(2),
			},
		}},
		{"eq", "true == false", compiler.Bytecode{
			Instructions: cat(
				instr(t, code.OpTrue),
				instr(t, code.OpFalse),
				instr(t, code.OpEqual),
				instr(t, code.OpPop),
			),
			Constants: nil,
		}},
		{"neq", "true != false", compiler.Bytecode{
			Instructions: cat(
				instr(t, code.OpTrue),
				instr(t, code.OpFalse),
				instr(t, code.OpNotEqual),
				instr(t, code.OpPop),
			),
			Constants: nil,
		}},
		{"bang", "!true", compiler.Bytecode{
			Instructions: cat(
				instr(t, code.OpTrue),
				instr(t, code.OpBang),
				instr(t, code.OpPop),
			),
			Constants: nil,
		}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			program := parser.New(lexer.New(tt.input)).Parse()

			compiler := compiler.New()
			require.NoError(t, compiler.Compile(program))
			if want, got := tt.want, compiler.Bytecode(); !cmp.Equal(want, got) {
				if !cmp.Equal(want.Instructions, got.Instructions) {
					t.Log(cmp.Diff(want.Instructions.String(), got.Instructions.String()))
				}
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

func TestConditionals(t *testing.T) {
	var (
		cat   = ConcatInstructions
		instr = MakeInstructions
		int   = IntegerObject
	)
	tests := []testcase{
		{
			name:  "if-true",
			input: "if (true) { 10 }; 3333;",
			want: compiler.Bytecode{
				Constants: []object.Object{int(10), int(3333)},
				Instructions: cat(
					// 0000
					instr(t, code.OpTrue),
					// 0001
					instr(t, code.OpJumpNotTruthy, 10),
					// 0004
					instr(t, code.OpConstant, 0),
					// 0007
					instr(t, code.OpJump, 11),
					// 0010
					instr(t, code.OpNull),
					// 0011
					instr(t, code.OpPop),
					// 0012
					instr(t, code.OpConstant, 1),
					// 0015
					instr(t, code.OpPop),
				),
			},
		},
		{
			name:  "if-true-else",
			input: "if (true) { 10 } else { 20 }; 3333;",
			want: compiler.Bytecode{
				Constants: []object.Object{int(10), int(20), int(3333)},
				Instructions: cat(
					instr(t, code.OpTrue),
					instr(t, code.OpJumpNotTruthy, 10),
					instr(t, code.OpConstant, 0),
					instr(t, code.OpJump, 13),
					instr(t, code.OpConstant, 1),
					instr(t, code.OpPop),
					instr(t, code.OpConstant, 2),
					instr(t, code.OpPop),
				),
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			program := parser.New(lexer.New(tt.input)).Parse()

			compiler := compiler.New()
			require.NoError(t, compiler.Compile(program))
			if want, got := tt.want, compiler.Bytecode(); !cmp.Equal(want, got) {
				if !cmp.Equal(want.Instructions, got.Instructions) {
					t.Log(cmp.Diff(want.Instructions.String(), got.Instructions.String()))
				}
				t.Error(cmp.Diff(want, got))
			}
		})
	}
}
