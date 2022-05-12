package ast

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Warashi/monkey/token"
	"golang.org/x/exp/slices"
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

type StringLiteral struct {
	Token token.Token
	Value string
}

func (e *StringLiteral) expressionNode()      {}
func (e *StringLiteral) TokenLiteral() string { return e.Token.Literal }
func (e *StringLiteral) String() string       { return e.Token.Literal }

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (e *BooleanLiteral) expressionNode()      {}
func (e *BooleanLiteral) TokenLiteral() string { return e.Token.Literal }
func (e *BooleanLiteral) String() string       { return e.Token.Literal }

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (e *FunctionLiteral) expressionNode()      {}
func (e *FunctionLiteral) TokenLiteral() string { return e.Token.Literal }
func (e *FunctionLiteral) String() string {
	var b strings.Builder
	params := make([]string, 0, len(e.Parameters))
	for _, p := range e.Parameters {
		params = append(params, p.String())
	}
	b.WriteString(e.TokenLiteral())
	b.WriteString("(")
	b.WriteString(strings.Join(params, ","))
	b.WriteString(")")
	b.WriteString(e.Body.String())
	return b.String()
}

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (e *ArrayLiteral) expressionNode()      {}
func (e *ArrayLiteral) TokenLiteral() string { return e.Token.Literal }
func (e *ArrayLiteral) String() string {
	var b strings.Builder
	elements := make([]string, 0, len(e.Elements))
	for _, e := range e.Elements {
		elements = append(elements, e.String())
	}
	b.WriteString("[")
	b.WriteString(strings.Join(elements, ", "))
	b.WriteString("]")
	return b.String()
}

type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (e *HashLiteral) expressionNode()      {}
func (e *HashLiteral) TokenLiteral() string { return e.Token.Literal }
func (e *HashLiteral) String() string {
	var b strings.Builder
	pairs := make([]string, 0, len(e.Pairs))
	for _, pair := range e.pairs() {
		pairs = append(pairs, fmt.Sprintf("%s:%s", pair[0], pair[1]))
	}
	b.WriteString("{")
	b.WriteString(strings.Join(pairs, ", "))
	b.WriteString("}")
	return b.String()
}

func (e *HashLiteral) pairs() [][2]Expression {
	pairs := make([][2]Expression, 0, len(e.Pairs))
	for k, v := range e.Pairs {
		pairs = append(pairs, [2]Expression{k, v})
	}
	slices.SortFunc(pairs, func(a, b [2]Expression) bool { return a[0].String() < b[0].String() })
	return pairs
}

func (e *HashLiteral) Equal(o *HashLiteral) bool {
	return reflect.DeepEqual(e.pairs(), o.pairs())
}

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

type IndexExpression struct {
	Token       token.Token
	Left, Right Expression
}

func (e *IndexExpression) expressionNode()      {}
func (e *IndexExpression) TokenLiteral() string { return e.Token.Literal }
func (e *IndexExpression) String() string {
	var b strings.Builder
	b.WriteString("(")
	b.WriteString(e.Left.String())
	b.WriteString("[")
	b.WriteString(e.Right.String())
	b.WriteString("])")
	return b.String()
}

type IfExpression struct {
	Token                    token.Token
	Condition                Expression
	Consequence, Alternative *BlockStatement
}

func (e *IfExpression) expressionNode()      {}
func (e *IfExpression) TokenLiteral() string { return e.Token.Literal }
func (e *IfExpression) String() string {
	var b strings.Builder
	b.WriteString("if")
	b.WriteString(e.Condition.String())
	b.WriteString(" ")
	b.WriteString(e.Consequence.String())
	if e.Alternative != nil {
		b.WriteString("else ")
		b.WriteString(e.Alternative.String())
	}
	return b.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (e *CallExpression) expressionNode()      {}
func (e *CallExpression) TokenLiteral() string { return e.Token.Literal }
func (e *CallExpression) String() string {
	args := make([]string, 0, len(e.Arguments))
	for _, arg := range e.Arguments {
		args = append(args, arg.String())
	}
	var b strings.Builder
	b.WriteString(e.Function.String())
	b.WriteString("(")
	b.WriteString(strings.Join(args, ", "))
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

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (s *BlockStatement) TokenLiteral() string {
	return s.Token.Literal
}

func (s *BlockStatement) String() string {
	var b strings.Builder
	for _, s := range s.Statements {
		b.WriteString(s.String())
	}
	return b.String()
}
