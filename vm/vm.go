package vm

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/Warashi/monkey/code"
	"github.com/Warashi/monkey/compiler"
	"github.com/Warashi/monkey/object"
)

const StackSize = 1 << 11

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int // Always points to the next value. Top of stack is stack[sp-1]
}

func New(code compiler.Bytecode) *VM {
	return &VM{
		constants:    code.Constants,
		instructions: code.Instructions,

		stack: make([]object.Object, StackSize),
		sp:    0,
	}
}

func (vm *VM) Run() error {
	r := bytes.NewReader(vm.instructions)
	for {
		op, err := code.ReadOpcode(r)
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("code.ReadOpcode: %w", err)
		}

		switch op {
		case code.OpConstant:
			idx, err := code.ReadUint16(r)
			if err != nil {
				return fmt.Errorf("code.ReadUint16: %w", err)
			}
			if err := vm.push(vm.constants[idx]); err != nil {
				return fmt.Errorf("vm.push: %w", err)
			}
		case code.OpAdd:
			left, err := vm.pop()
			if err != nil {
				return fmt.Errorf("vm.pop: %w", err)
			}
			right, err := vm.pop()
			if err != nil {
				return fmt.Errorf("vm.pop: %w", err)
			}
			result, err := add(left, right)
			if err != nil {
				return fmt.Errorf("add: %w", err)
			}
			if err := vm.push(result); err != nil {
				return fmt.Errorf("vm.push: %w", err)
			}
		default:
			return fmt.Errorf("unknown opcode: %s", op.String())
		}
	}
}

func (vm *VM) push(obj object.Object) error {
	if vm.sp >= StackSize {
		return errors.New("stack overflow")
	}
	vm.stack[vm.sp] = obj
	vm.sp++
	return nil
}

func (vm *VM) pop() (object.Object, error) {
	if vm.sp == 0 {
		return nil, fmt.Errorf("stack underflow")
	}
	obj := vm.stack[vm.sp-1]
	vm.sp--
	return obj, nil
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func add(left, right object.Object) (object.Object, error) {
	return object.Integer{Value: left.(object.Integer).Value + right.(object.Integer).Value}, nil
}
