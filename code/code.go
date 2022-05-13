package code

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"strings"
)

//go:generate go run golang.org/x/tools/cmd/stringer -type Opcode -trimprefix Op
type (
	Instructions []byte
	Opcode       byte
)

const (
	_ Opcode = iota
	OpConstant
	OpAdd
)

type Definition struct {
	Name          string
	OperandWitdth []int
}

var definitions = map[Opcode]Definition{
	OpConstant: {"OpConstant", []int{2}},
	OpAdd:      {"OpAdd", nil},
}

func Lookup(op Opcode) (Definition, error) {
	def, ok := definitions[op]
	if !ok {
		return Definition{}, fmt.Errorf("%s not found.", op)
	}
	return def, nil
}

func Make(op Opcode, operands ...int64) (Instructions, error) {
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

func (ins Instructions) String() string {
	r := bytes.NewReader(ins)
	w := new(strings.Builder)
	read := 0
	for {
		op, err := ReadOpcode(r)
		if errors.Is(err, io.EOF) {
			return w.String()
		}
		if err != nil {
			fmt.Fprintf(w, "ERROR: %s\n", err.Error())
			return w.String()
		}
		def, err := Lookup(op)
		if err != nil {
			fmt.Fprintf(w, "ERROR: %s\n", err.Error())
			return w.String()
		}

		operands, n, err := ReadOperands(def, r)
		if err != nil {
			fmt.Fprintf(w, "ERROR: %s\n", err.Error())
		}

		fmt.Fprintf(w, "%04d %s\n", read, ins.fmtInstruction(def, operands))

		read += 1 + n
	}
}

func (ins Instructions) fmtInstruction(def Definition, operands []int64) string {
	operandCount := len(def.OperandWitdth)
	if operandCount != len(operands) {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n", len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	default:
		return fmt.Sprintf("ERROR: unhandled operandCount for %s", def.Name)
	}
}

func ReadOpcode(r io.Reader) (Opcode, error) {
	op := make([]byte, 1)
	_, err := r.Read(op)
	if err != nil {
		return 0, fmt.Errorf("r.Read: %w", err)
	}
	return Opcode(op[0]), nil
}

func ReadOperands(def Definition, r io.Reader) ([]int64, int, error) {
	operands := make([]int64, len(def.OperandWitdth))
	read := 0
	for i, width := range def.OperandWitdth {
		switch width {
		case 2:
			var err error
			operands[i], err = ReadUint16(r)
			if err != nil {
				return nil, 0, fmt.Errorf("ReadUint16: %w", err)
			}
		}
		read += width
	}

	return operands, read, nil
}

func ReadUint16(r io.Reader) (int64, error) {
	var read uint16
	if err := binary.Read(r, binary.BigEndian, &read); err != nil {
		return 0, fmt.Errorf("binary.Read: %w", err)
	}
	return int64(read), nil
}
