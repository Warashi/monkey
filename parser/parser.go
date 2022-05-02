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

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, prefixParseFns: make(map[token.Type]prefixParseFn), infixParseFns: make(map[token.Type]infixParseFn)}

	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)

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
	if prefix, ok := p.prefixParseFns[p.current.Type]; ok {
		return prefix()
	}
	return nil
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
