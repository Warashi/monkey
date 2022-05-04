package lexer

import (
	"github.com/Warashi/implement-interpreter-with-go/token"
)

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

func (l *Lexer) peekChar() byte {
	if l.readPosisiton >= len(l.input) {
		return 0
	}
	return l.input[l.readPosisiton]
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()
	defer l.readChar()
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			return token.Token{Type: token.EQ, Literal: literal}
		}
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
	case '-':
		return newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			return token.Token{Type: token.NOT_EQ, Literal: literal}
		}
		return newToken(token.BANG, l.ch)
	case '/':
		return newToken(token.SLASH, l.ch)
	case '*':
		return newToken(token.ASTERISK, l.ch)
	case '<':
		return newToken(token.LT, l.ch)
	case '>':
		return newToken(token.GT, l.ch)
	case '"':
		return token.Token{Type: token.STRING, Literal: l.readString()}
	case 0:
		return token.Token{Type: token.EOF}
	default:
		switch {
		case isLetter(l.ch):
			ident := l.readIdentifier()
			return token.Token{Type: token.LookupIdent(ident), Literal: ident}
		case isNumber(l.ch):
			return token.Token{Type: token.INT, Literal: l.readNumber()}
		default:
			return newToken(token.ILLEGAL, l.ch)
		}
	}
}

func (l *Lexer) readIdentifier() string {
	p := l.position
	for isLetter(l.peekChar()) {
		l.readChar()
	}
	return l.input[p:l.readPosisiton]
}

func (l *Lexer) readNumber() string {
	p := l.position
	for isNumber(l.peekChar()) {
		l.readChar()
	}
	return l.input[p:l.readPosisiton]
}

func (l *Lexer) readString() string {
	l.readChar()
	p := l.position
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	return l.input[p:l.position]
}

func (l *Lexer) skipWhitespace() {
	for isSpace(l.ch) {
		l.readChar()
	}
}

func newToken(t token.Type, ch byte) token.Token {
	return token.Token{Type: t, Literal: string(ch)}
}

func isNumber(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isSpace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}
