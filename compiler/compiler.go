package compiler

import (
	"fmt"

	"github.com/Warashi/monkey/ast"
	"github.com/Warashi/monkey/code"
	"github.com/Warashi/monkey/object"
)

type (
	Bytecode struct {
		Instructions code.Instructions
		Constants    []object.Object
	}
	Compiler struct {
		instructions code.Instructions
		constants    []object.Object
	}
)

func New() *Compiler {
	return &Compiler{}
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, stmt := range node.Statements {
			if err := c.Compile(stmt); err != nil {
				return fmt.Errorf("c.Compile(%T): %w", node, err)
			}
		}
	case *ast.ExpressionStatement:
		if err := c.Compile(node.Expression); err != nil {
			return fmt.Errorf("c.Compile(%T): %w", node, err)
		}
	case *ast.InfixExpression:
		if err := c.Compile(node.Left); err != nil {
			return fmt.Errorf("c.Compile(%T): %w", node, err)
		}
		if err := c.Compile(node.Right); err != nil {
			return fmt.Errorf("c.Compile(%T): %w", node, err)
		}
	case *ast.IntegerLiteral:
		if _, err := c.emit(code.OpConstant, c.addConstant(object.Integer{Value: node.Value})); err != nil {
			return fmt.Errorf("c.emit: %w", err)
		}
	default:
		return fmt.Errorf("unknown type: %T", node)
	}
	return nil
}

func (c *Compiler) Bytecode() Bytecode {
	return Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

func (c *Compiler) addConstant(obj object.Object) int64 {
	c.constants = append(c.constants, obj)
	return int64(len(c.constants) - 1)
}

func (c *Compiler) addInstruction(ins code.Instructions) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}

func (c *Compiler) emit(op code.Opcode, operands ...int64) (int, error) {
	ins, err := code.Make(op, operands...)
	if err != nil {
		return 0, fmt.Errorf("code.Make: %w", err)
	}
	return c.addInstruction(ins), nil
}
