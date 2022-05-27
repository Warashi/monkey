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
	Null  = object.Null{}
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
		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			if err := vm.executeComparison(op); err != nil {
				return fmt.Errorf("vm.executeComparison: %w", err)
			}
		case code.OpBang:
			if err := vm.executeBangOperator(); err != nil {
				return fmt.Errorf("vm.executeBangOperator: %w", err)
			}
		case code.OpMinus:
			if err := vm.executeMinusOperator(); err != nil {
				return fmt.Errorf("vm.executeMinusOperator: %w", err)
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
		case code.OpJump:
			pos, err := code.ReadUint16(r)
			if err != nil {
				return fmt.Errorf("code.ReadUint16: %w", err)
			}
			r.Seek(pos, 0)
		case code.OpJumpNotTruthy:
			pos, err := code.ReadUint16(r)

			condition, err := vm.pop()
			if err != nil {
				return fmt.Errorf("vm.pop: %w", err)
			}

			if !isTruthy(condition) {
				r.Seek(pos, 0)
			}
		case code.OpNull:
			if err := vm.push(Null); err != nil {
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

func (vm *VM) executeComparison(op code.Opcode) error {
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
		if err := vm.executeIntegerComparison(op, left, right); err != nil {
			return fmt.Errorf("vm.executeIntegerComparison: %w", err)
		}
	case left.Type() == object.TypeBoolean && right.Type() == object.TypeBoolean:
		left, right := left.(object.Boolean), right.(object.Boolean)
		if err := vm.executeBooleanComparison(op, left, right); err != nil {
			return fmt.Errorf("vm.executeBooleanComparison: %w", err)
		}
	default:
		return fmt.Errorf("unsupported types: op=%s, left: %s, right: %s", op.String(), left.Type().String(), right.Type().String())
	}
	return nil
}

func (vm *VM) executeIntegerComparison(op code.Opcode, left, right object.Integer) error {
	var result bool
	switch op {
	case code.OpEqual:
		result = left.Value == right.Value
	case code.OpNotEqual:
		result = left.Value != right.Value
	case code.OpGreaterThan:
		result = left.Value > right.Value
	default:
		return fmt.Errorf("uknown operator: %s", op.String())
	}
	if err := vm.push(booleanObject(result)); err != nil {
		return fmt.Errorf("vm.push: %w", err)
	}
	return nil
}

func (vm *VM) executeBooleanComparison(op code.Opcode, left, right object.Boolean) error {
	var result bool
	switch op {
	case code.OpEqual:
		result = left.Value == right.Value
	case code.OpNotEqual:
		result = left.Value != right.Value
	default:
		return fmt.Errorf("uknown operator: %s", op.String())
	}
	if err := vm.push(booleanObject(result)); err != nil {
		return fmt.Errorf("vm.push: %w", err)
	}
	return nil
}

func (vm *VM) executeBangOperator() error {
	operand, err := vm.pop()
	if err != nil {
		return fmt.Errorf("vm.pop: %w", err)
	}
	if isTruthy(operand) {
		if err := vm.push(False); err != nil {
			return fmt.Errorf("vm.push: %w", err)
		}
	} else {
		if err := vm.push(True); err != nil {
			return fmt.Errorf("vm.push: %w", err)
		}
	}
	return nil
}

func (vm *VM) executeMinusOperator() error {
	operand, err := vm.pop()
	if err != nil {
		return fmt.Errorf("vm.pop: %w", err)
	}
	if operand.Type() != object.TypeInteger {
		return fmt.Errorf("unsupported type for negation: %s", operand.Type())
	}
	val := operand.(object.Integer).Value
	if err := vm.push(object.Integer{Value: -val}); err != nil {
		return fmt.Errorf("vm.push: %w", err)
	}
	return nil
}

func booleanObject(value bool) object.Boolean {
	switch value {
	case true:
		return True
	case false:
		return False
	}
	panic("unreachable")
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case object.Boolean:
		return obj.Value
	case object.Null:
		return false
	default:
		return true
	}
}
