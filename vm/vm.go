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

var (
	True  = object.Boolean{Value: true}
	False = object.Boolean{Value: false}
)

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
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			if err := vm.executeBinaryOperation(op); err != nil {
				return fmt.Errorf("vm.executeBinaryOperation: %w", err)
			}
		case code.OpTrue:
			if err := vm.push(True); err != nil {
				return fmt.Errorf("vm.push: %w", err)
			}
		case code.OpFalse:
			if err := vm.push(False); err != nil {
				return fmt.Errorf("vm.push: %w", err)
			}
		case code.OpPop:
			if _, err := vm.pop(); err != nil {
				return fmt.Errorf("vm.pop: %w", err)
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

func (vm *VM) LastPopedStackElem() object.Object {
	return vm.stack[vm.sp]
}

func (vm *VM) executeBinaryOperation(op code.Opcode) error {
	right, err := vm.pop()
	if err != nil {
		return fmt.Errorf("vm.pop: %w", err)
	}
	left, err := vm.pop()
	if err != nil {
		return fmt.Errorf("vm.pop: %w", err)
	}
	switch {
	case left.Type() == object.TypeInteger && right.Type() == object.TypeInteger:
		left, right := left.(object.Integer), right.(object.Integer)
		if err := vm.executeBinaryIntegerOperation(op, left, right); err != nil {
			return fmt.Errorf("vm.executeBinaryIntegerOperation: %w", err)
		}
	default:
		return fmt.Errorf("unsupported types: op=%s, left: %s, right: %s", op.String(), left.Type().String(), right.Type().String())
	}
	return nil
}

func (vm *VM) executeBinaryIntegerOperation(op code.Opcode, left, right object.Integer) error {
	var result int64
	switch op {
	case code.OpAdd:
		result = left.Value + right.Value
	case code.OpSub:
		result = left.Value - right.Value
	case code.OpMul:
		result = left.Value * right.Value
	case code.OpDiv:
		result = left.Value / right.Value
	default:
		return fmt.Errorf("uknown operator: %s", op.String())
	}

	if err := vm.push(object.Integer{Value: result}); err != nil {
		return fmt.Errorf("vm.push: %w", err)
	}
	return nil
}
