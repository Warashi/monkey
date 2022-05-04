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
	token.LPAREN:   CALL,
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, prefixParseFns: make(map[token.Type]prefixParseFn), infixParseFns: make(map[token.Type]infixParseFn)}

	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.STRING, p.parseStringerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.LBLACKET, p.parseArrayLiteral)

	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

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

func (p *Parser) parseStringerLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.current, Value: p.current.Literal}
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.current, Value: p.currentIs(token.TRUE)}
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	e := &ast.FunctionLiteral{Token: p.current}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	e.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	e.Body = p.parseBlockStatement()
	return e
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	if p.peekIs(token.RPAREN) {
		p.nextToken()
		return nil
	}
	p.nextToken()
	var ids []*ast.Identifier
	ids = append(ids, &ast.Identifier{Token: p.current, Value: p.current.Literal})
	for p.peekIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ids = append(ids, &ast.Identifier{Token: p.current, Value: p.current.Literal})
	}
	p.expectPeek(token.RPAREN)
	return ids
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	e := &ast.PrefixExpression{Token: p.current, Operator: p.current.Literal}
	p.nextToken()
	e.Right = p.parseExpression(PREFIX)
	return e
}

func (p *Parser) parseIfExpression() ast.Expression {
	e := &ast.IfExpression{Token: p.current}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.nextToken()
	e.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	e.Consequence = p.parseBlockStatement()
	if p.peekIs(token.ELSE) {
		p.nextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		e.Alternative = p.parseBlockStatement()
	}
	return e
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	s := &ast.BlockStatement{Token: p.current}
	p.nextToken()
	for !p.currentIs(token.RBRACE) && !p.currentIs(token.EOF) {
		if stmt := p.parseStatement(); stmt != nil {
			s.Statements = append(s.Statements, stmt)
		}
		p.nextToken()
	}
	return s
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

func (p *Parser) parseCallExpression(left ast.Expression) ast.Expression {
	e := &ast.CallExpression{Token: p.current, Function: left}
	e.Arguments = p.parseCallArguments()
	return e
}

func (p *Parser) parseCallArguments() []ast.Expression {
	return p.parseExpressionList(token.RPAREN)
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	e := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return e
}

func (p *Parser) parseArrayLiteral() ast.Expression {
	e := &ast.ArrayLiteral{Token: p.current}
	e.Elements = p.parseExpressionList(token.RBLACKET)
	return e
}

func (p *Parser) parseExpressionList(end token.Type) []ast.Expression {
	if p.peekIs(end) {
		p.nextToken()
		return nil
	}
	p.nextToken()

	var exps []ast.Expression
	exps = append(exps, p.parseExpression(LOWEST))
	for p.peekIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		exps = append(exps, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}
	return exps
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
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	for p.peekIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.current}
	p.nextToken()
	stmt.Value = p.parseExpression(LOWEST)
	for p.peekIs(token.SEMICOLON) {
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
