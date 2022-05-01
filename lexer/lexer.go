package lexer

import "github.com/Warashi/implement-interpreter-with-go/token"

type Lexer struct {
	input         string
	position      int
	readPosisiton int
	ch            byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosisiton >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosisiton]
	}
	l.position = l.readPosisiton
	l.readPosisiton++
}

func (l *Lexer) NextToken() token.Token {
	defer l.readChar()
	switch l.ch {
	case '=':
		return newToken(token.ASSIGN, l.ch)
	case ';':
		return newToken(token.SEMICOLON, l.ch)
	case '(':
		return newToken(token.LPAREN, l.ch)
	case ')':
		return newToken(token.RPAREN, l.ch)
	case '{':
		return newToken(token.LBRACE, l.ch)
	case '}':
		return newToken(token.RBRACE, l.ch)
	case ',':
		return newToken(token.COMMA, l.ch)
	case '+':
		return newToken(token.PLUS, l.ch)
	case '0':
		return token.Token{Type: token.EOF}
	default:
		return token.Token{Type: token.ILLEGAL}
	}
}

func newToken(t token.Type, ch byte) token.Token {
	return token.Token{Type: t, Literal: string(ch)}
}
