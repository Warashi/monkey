package testutil

import (
	"testing"

	"github.com/Warashi/monkey/code"
)

func ConcatInstructions(is ...code.Instructions) code.Instructions {
	var cat code.Instructions
	for _, i := range is {
		cat = append(cat, i...)
	}
	return cat
}

func MakeInstructions(t *testing.T, op code.Opcode, operands ...int64) code.Instructions {
	t.Helper()
	return NoError(t, func() (code.Instructions, error) { return code.Make(op, operands...) })
}
