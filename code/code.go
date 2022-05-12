package code

import "fmt"

//go:generate go run golang.org/x/tools/cmd/stringer -type Opcode -trimprefix Op
type (
	Instructions []byte
	Opcode       byte
)

const (
	_ Opcode = iota
	OpConstant
)

type Definition struct {
	Name          string
	OperandWitdth []int
}

var definitions = map[Opcode]Definition{
	OpConstant: {"OpConstant", []int{2}},
}

func Lookup(op Opcode) (Definition, error) {
	def, ok := definitions[op]
	if !ok {
		return Definition{}, fmt.Errorf("%s not found.", op)
	}
	return def, nil
}
