package parser

import (
	"fmt"
	"strconv"

	"github.com/Warashi/implement-interpreter-with-go/ast"
	"github.com/Warashi/implement-interpreter-with-go/lexer"
	"github.com/Warashi/implement-interpreter-with-go/token"
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
	precedence    int
)

const (
	_ precedence = iota
	LOWEST
	EQUALS
	LTGT
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.Type]precedence{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LTGT,
	token.GT:       LTGT,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, prefixParseFns: make(map[token.Type]prefixParseFn), infixParseFns: make(map[token.Type]infixParseFn)}

	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	for t := range precedences {
		p.registerInfix(t, p.parseInfixExpression)
	}

	// current, peek をセット
	p.nextToken()
	p.nextToken()
	return p
}

type Parser struct {
	l *lexer.Lexer

	current, peek token.Token
	errors        []string

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

func (p *Parser) registerPrefix(t token.Type, fn prefixParseFn) { p.prefixParseFns[t] = fn }
func (p *Parser) registerInfix(t token.Type, fn infixParseFn)   { p.infixParseFns[t] = fn }
func (p *Parser) Errors() []string                              { return p.errors }
func (p *Parser) nextToken()                                    { p.current, p.peek = p.peek, p.l.NextToken() }

func (p *Parser) Parse() *ast.Program {
	program := new(ast.Program)
	for p.current.Type != token.EOF {
		if stmt := p.parseStatement(); stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.current, Value: p.current.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	value, err := strconv.ParseInt(p.current.Literal, 10, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("could not parse %q as integer", p.current.Literal))
		return nil
	}
	return &ast.IntegerLiteral{Token: p.current, Value: value}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	e := &ast.PrefixExpression{Token: p.current, Operator: p.current.Literal}
	p.nextToken()
	e.Right = p.parseExpression(PREFIX)
	return e
}

func (p *Parser) currentPrecedence() precedence {
	if p, ok := precedences[p.current.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() precedence {
	if p, ok := precedences[p.peek.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	e := &ast.InfixExpression{Token: p.current, Operator: p.current.Literal, Left: left}
	precedence := p.currentPrecedence()
	p.nextToken()
	e.Right = p.parseExpression(precedence)
	return e
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.current.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() ast.Statement {
	stmt := &ast.LetStatement{Token: p.current}
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.current, Value: p.current.Literal}
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	stmt.Value = p.parseExpression(LOWEST)
	for !p.currentIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.current}
	stmt.Value = p.parseExpression(LOWEST)
	for !p.currentIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Token: p.current}
	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(prec precedence) ast.Expression {
	prefix, ok := p.prefixParseFns[p.current.Type]
	if !ok {
		p.errors = append(p.errors, fmt.Sprintf("no prefixParseFn found: %s", p.current.Type))
		return nil
	}
	left := prefix()
	for !p.peekIs(token.SEMICOLON) && prec < p.peekPrecedence() {
		infix, ok := p.infixParseFns[p.peek.Type]
		if !ok {
			return left
		}
		p.nextToken()
		left = infix(left)
	}
	return left
}

func (p *Parser) currentIs(t token.Type) bool {
	return p.current.Type == t
}

func (p *Parser) peekIs(t token.Type) bool {
	return p.peek.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekIs(t) {
		p.nextToken()
		return true
	}
	p.errors = append(p.errors, fmt.Sprintf("expect next token is %s but %s instead", t, p.peek.Type))
	return false
}
