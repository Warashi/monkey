package parser

import (
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
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() ast.Statement {
	if !p.nextIfPeekIs(token.IDENT) {
		return nil
	}
	name := &ast.Identifier{Token: p.current, Value: p.current.Literal}
	if !p.nextIfPeekIs(token.ASSIGN) {
		return nil
	}
	for !p.currentIs(token.SEMICOLON) {
		p.nextToken()
	}
	return &ast.LetStatement{Name: name}
}

func (p *Parser) currentIs(t token.Type) bool {
	return p.current.Type == t
}

func (p *Parser) peekIs(t token.Type) bool {
	return p.peek.Type == t
}

func (p *Parser) nextIfPeekIs(t token.Type) bool {
	if p.peekIs(t) {
		p.nextToken()
		return true
	}
	return false
}
