package compiler

import (
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

func (c *Compiler) Compile(program ast.Node) error {
	return nil
}

func (c *Compiler) Bytecode() Bytecode {
	return Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}
