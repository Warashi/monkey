package parser

import (
	"fmt"

	"github.com/Warashi/implement-interpreter-with-go/ast"
	"github.com/Warashi/implement-interpreter-with-go/lexer"
	"github.com/Warashi/implement-interpreter-with-go/token"
)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// current, peek をセット
	p.nextToken()
	p.nextToken()
	return p
}

type Parser struct {
	l *lexer.Lexer

	current, peek token.Token
	errors        []string
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.current, p.peek = p.peek, p.l.NextToken()
}

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

func (p *Parser) parseStatement() ast.Statement {
	switch p.current.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
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
	stmt.Value = p.parseExpression()
	return stmt
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.current}
	stmt.Value = p.parseExpression()
	return stmt
}

func (p *Parser) parseExpression() ast.Expression {
	for !p.currentIs(token.SEMICOLON) {
		p.nextToken()
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
