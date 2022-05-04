package token

//go:generate go run golang.org/x/tools/cmd/stringer -type Type
type Type int

const (
	ILLEGAL Type = iota // ILLEGAL
	EOF                 // EOF

	// 識別子, リテラル
	IDENT  // IDENT
	INT    // INT
	STRING // STRING

	// 演算子
	ASSIGN   // =
	PLUS     // +
	MINUS    // -
	BANG     // !
	ASTERISK // *
	SLASH    // /

	LT     // <
	GT     // >
	EQ     // ==
	NOT_EQ // !=

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
	TRUE     // true
	FALSE    // false
	IF       // if
	ELSE     // else
	RETURN   // return
)

type Token struct {
	Type    Type
	Literal string
}
