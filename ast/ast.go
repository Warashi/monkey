package ast

import (
	"strings"

	"github.com/Warashi/implement-interpreter-with-go/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var b strings.Builder
	for _, s := range p.Statements {
		b.WriteString(s.String())
	}
	return b.String()
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (s *LetStatement) statementNode()       {}
func (s *LetStatement) TokenLiteral() string { return s.Token.Literal }
func (s *LetStatement) String() string {
	var b strings.Builder
	b.WriteString(s.TokenLiteral())
	b.WriteString(" ")
	b.WriteString(s.Name.String())
	b.WriteString(" = ")
	if s.Value != nil {
		b.WriteString(s.Value.String())
	}
	b.WriteString(";\n")
	return b.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (e *Identifier) expressionNode()      {}
func (e *Identifier) TokenLiteral() string { return e.Token.Literal }
func (e *Identifier) String() string       { return e.Value }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (e *IntegerLiteral) expressionNode()      {}
func (e *IntegerLiteral) TokenLiteral() string { return e.Token.Literal }
func (e *IntegerLiteral) String() string       { return e.Token.Literal }

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (e *BooleanLiteral) expressionNode()      {}
func (e *BooleanLiteral) TokenLiteral() string { return e.Token.Literal }
func (e *BooleanLiteral) String() string       { return e.Token.Literal }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (e *PrefixExpression) expressionNode()      {}
func (e *PrefixExpression) TokenLiteral() string { return e.Token.Literal }
func (e *PrefixExpression) String() string {
	var b strings.Builder
	b.WriteString("(")
	b.WriteString(e.Operator)
	b.WriteString(e.Right.String())
	b.WriteString(")")
	return b.String()
}

type InfixExpression struct {
	Token       token.Token
	Operator    string
	Left, Right Expression
}

func (e *InfixExpression) expressionNode()      {}
func (e *InfixExpression) TokenLiteral() string { return e.Token.Literal }
func (e *InfixExpression) String() string {
	var b strings.Builder
	b.WriteString("(")
	b.WriteString(e.Left.String())
	b.WriteString(" ")
	b.WriteString(e.Operator)
	b.WriteString(" ")
	b.WriteString(e.Right.String())
	b.WriteString(")")
	return b.String()
}

type ReturnStatement struct {
	Token token.Token
	Value Expression
}

func (s *ReturnStatement) statementNode()       {}
func (s *ReturnStatement) TokenLiteral() string { return s.Token.Literal }
func (s *ReturnStatement) String() string {
	var b strings.Builder
	b.WriteString(s.TokenLiteral())
	b.WriteString(" ")
	if s.Value != nil {
		b.WriteString(s.Value.String())
	}
	b.WriteString(";\n")
	return b.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (s *ExpressionStatement) statementNode()       {}
func (s *ExpressionStatement) TokenLiteral() string { return s.Token.Literal }
func (s *ExpressionStatement) String() string {
	if s.Expression != nil {
		return s.Expression.String()
	}
	return ""
}
