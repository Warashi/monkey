// Code generated by "stringer -type Type"; DO NOT EDIT.

package token

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ILLEGAL-0]
	_ = x[EOF-1]
	_ = x[IDENT-2]
	_ = x[INT-3]
	_ = x[STRING-4]
	_ = x[ASSIGN-5]
	_ = x[PLUS-6]
	_ = x[MINUS-7]
	_ = x[BANG-8]
	_ = x[ASTERISK-9]
	_ = x[SLASH-10]
	_ = x[LT-11]
	_ = x[GT-12]
	_ = x[EQ-13]
	_ = x[NOT_EQ-14]
	_ = x[COMMA-15]
	_ = x[SEMICOLON-16]
	_ = x[LPAREN-17]
	_ = x[RPAREN-18]
	_ = x[LBRACE-19]
	_ = x[RBRACE-20]
	_ = x[LBLACKET-21]
	_ = x[RBLACKET-22]
	_ = x[FUNCTION-23]
	_ = x[LET-24]
	_ = x[TRUE-25]
	_ = x[FALSE-26]
	_ = x[IF-27]
	_ = x[ELSE-28]
	_ = x[RETURN-29]
}

const _Type_name = "ILLEGALEOFIDENTINTSTRINGASSIGNPLUSMINUSBANGASTERISKSLASHLTGTEQNOT_EQCOMMASEMICOLONLPARENRPARENLBRACERBRACELBLACKETRBLACKETFUNCTIONLETTRUEFALSEIFELSERETURN"

var _Type_index = [...]uint8{0, 7, 10, 15, 18, 24, 30, 34, 39, 43, 51, 56, 58, 60, 62, 68, 73, 82, 88, 94, 100, 106, 114, 122, 130, 133, 137, 142, 144, 148, 154}

func (i Type) String() string {
	if i < 0 || i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
