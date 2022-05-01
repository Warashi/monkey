package token

import "golang.org/x/exp/maps"

var (
	keywordTypes = map[string]Type{
		"fn":     FUNCTION,
		"let":    LET,
		"true":   TRUE,
		"false":  FALSE,
		"if":     IF,
		"else":   ELSE,
		"return": RETURN,
	}
	keywords = maps.Values(keywordTypes)
)

func LookupIdent(ident string) Type {
	if tok, ok := keywordTypes[ident]; ok {
		return tok
	}
	return IDENT
}
