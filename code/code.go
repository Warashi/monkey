package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

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

func Make(op Opcode, operands ...int64) ([]byte, error) {
	def, err := Lookup(op)
	if err != nil {
		return nil, fmt.Errorf("Lookup: %w", err)
	}
	instLen := 1
	for _, w := range def.OperandWitdth {
		instLen += w
	}
	buf := new(bytes.Buffer)
	buf.Grow(instLen)
	binary.Write(buf, binary.BigEndian, op)

	for i, o := range operands {
		width := def.OperandWitdth[i]
		switch width {
		case 2:
			binary.Write(buf, binary.BigEndian, uint16(o))
		}
	}

	return buf.Bytes(), nil
}
