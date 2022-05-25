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
	EmittedInstruction struct {
		Opcode   code.Opcode
		Position int
	}
	Compiler struct {
		instructions        code.Instructions
		constants           []object.Object
		lastInstruction     EmittedInstruction
		previousInstruction EmittedInstruction
	}
)

func New() *Compiler {
	return new(Compiler)
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
		if _, err := c.emit(code.OpPop); err != nil {
			return fmt.Errorf("c.emit: %w", err)
		}
	case *ast.BlockStatement:
		for _, s := range node.Statements {
			if err := c.Compile(s); err != nil {
				return fmt.Errorf("c.Compile(%T): %w", node, err)
			}
		}
	case *ast.IfExpression:
		if err := c.Compile(node.Condition); err != nil {
			return fmt.Errorf("c.Compile(%T): %w", node, err)
		}
		// TODO
		jumpNotTruthyPos, err := c.emit(code.OpJumpNotTruthy, 9999)
		if err != nil {
			return fmt.Errorf("c.emit: %w", err)
		}

		if err := c.Compile(node.Consequence); err != nil {
			return fmt.Errorf("c.Compile(%T): %w", node, err)
		}

		if c.lastInstructionIsPop() {
			c.removeLastPop()
		}

		if node.Alternative == nil {
			afterConsequencePos := len(c.instructions)
			c.changeOperand(jumpNotTruthyPos, int64(afterConsequencePos))
		} else {
			jumpPos, err := c.emit(code.OpJump, 9999)
			if err != nil {
				return fmt.Errorf("c.emit: %w", err)
			}
			afterConsequencePos := len(c.instructions)
			c.changeOperand(jumpNotTruthyPos, int64(afterConsequencePos))

			if err := c.Compile(node.Alternative); err != nil {
				return fmt.Errorf("c.Compile(%T): %w", node, err)
			}

			if c.lastInstructionIsPop() {
				c.removeLastPop()
			}

			afterAlternativePos := len(c.instructions)
			c.changeOperand(jumpPos, int64(afterAlternativePos))
		}
	case *ast.PrefixExpression:
		if err := c.Compile(node.Right); err != nil {
			return fmt.Errorf("c.Compile(%T): %w", node, err)
		}
		if _, err := c.emitPrefixOp(node.Operator); err != nil {
			return fmt.Errorf("c.emitPrefixOp: %w", err)
		}
	case *ast.InfixExpression:
		switch node.Operator {
		case "<":
			if err := c.Compile(node.Right); err != nil {
				return fmt.Errorf("c.Compile(%T): %w", node, err)
			}
			if err := c.Compile(node.Left); err != nil {
				return fmt.Errorf("c.Compile(%T): %w", node, err)
			}
			if _, err := c.emit(code.OpGreaterThan); err != nil {
				return fmt.Errorf("c.emit: %w", err)
			}
		default:
			if err := c.Compile(node.Left); err != nil {
				return fmt.Errorf("c.Compile(%T): %w", node, err)
			}
			if err := c.Compile(node.Right); err != nil {
				return fmt.Errorf("c.Compile(%T): %w", node, err)
			}
			if _, err := c.emitInfixOp(node.Operator); err != nil {
				return fmt.Errorf("c.emitInfixOp: %w", err)
			}
		}
	case *ast.IntegerLiteral:
		if _, err := c.emit(code.OpConstant, c.addConstant(object.Integer{Value: node.Value})); err != nil {
			return fmt.Errorf("c.emit: %w", err)
		}
	case *ast.BooleanLiteral:
		switch node.Value {
		case true:
			if _, err := c.emit(code.OpTrue); err != nil {
				return fmt.Errorf("c.emit: %w", err)
			}
		case false:
			if _, err := c.emit(code.OpFalse); err != nil {
				return fmt.Errorf("c.emit: %w", err)
			}
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
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)

	return pos, nil
}

func (c *Compiler) emitPrefixOp(op string) (int, error) {
	switch op {
	case "!":
		return c.emit(code.OpBang)
	case "-":
		return c.emit(code.OpMinus)
	default:
		return 0, fmt.Errorf("unknown operator: %s", op)
	}
}

func (c *Compiler) emitInfixOp(op string) (int, error) {
	switch op {
	case "+":
		return c.emit(code.OpAdd)
	case "-":
		return c.emit(code.OpSub)
	case "*":
		return c.emit(code.OpMul)
	case "/":
		return c.emit(code.OpDiv)
	case ">":
		return c.emit(code.OpGreaterThan)
	case "==":
		return c.emit(code.OpEqual)
	case "!=":
		return c.emit(code.OpNotEqual)
	default:
		return 0, fmt.Errorf("unknown operator: %s", op)
	}
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	c.previousInstruction = previous
	c.lastInstruction = last
}

func (c *Compiler) lastInstructionIsPop() bool {
	return c.lastInstruction.Opcode == code.OpPop
}

func (c *Compiler) removeLastPop() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
}

func (c *Compiler) replaceInstruction(pos int, newInstruction code.Instructions) {
	for i := range newInstruction {
		c.instructions[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) changeOperand(opPos int, operand int64) error {
	op := code.Opcode(c.instructions[opPos])
	newInstruction, err := code.Make(op, operand)
	if err != nil {
		return fmt.Errorf("code.Make: %w", err)
	}
	c.replaceInstruction(opPos, newInstruction)
	return nil
}
