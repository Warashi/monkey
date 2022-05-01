package token

var keywordTypes = [...]Type{
	FUNCTION,
	LET,
}

var keywords map[string]Type

func init() {
	keywords = make(map[string]Type, len(keywordTypes))
	for _, k := range keywordTypes {
		keywords[k.String()] = k
	}
}

func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
