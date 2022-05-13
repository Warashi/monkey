package code_test

import (
	"testing"

	"github.com/Warashi/monkey/code"
	"github.com/stretchr/testify/assert"
)

func TestMake(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		op          code.Opcode
		operands    []int64
		want        code.Instructions
		assertError assert.ErrorAssertionFunc
	}{
		{"constant", code.OpConstant, []int64{0xFFFE}, code.Instructions{byte(code.OpConstant), 0xFF, 0xFE}, assert.NoError},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := code.Make(tt.op, tt.operands...)
			tt.assertError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
