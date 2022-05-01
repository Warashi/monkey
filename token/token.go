package token

//go:generate go run golang.org/x/tools/cmd/stringer -type Type -linecomment
type Type int

const (
	ILLEGAL Type = iota // ILLEGAL
	EOF                 // EOF

	// 識別子, リテラル
	IDENT // IDENT
	INT   // INT

	// 演算子
	ASSIGN // =
	PLUS   // +

	// デリミタ
	COMMA     // ,
	SEMICOLON // ;
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }

	// キーワード
	FUNCTION // fn
	LET      // let
)

type Token struct {
	Type    Type
	Literal string
}
