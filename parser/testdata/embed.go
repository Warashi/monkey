package testdata

import _ "embed"

var (
	//go:embed let.monkey
	Let string
	//go:embed return.monkey
	Return string
	//go:embed identifier-expression.monkey
	IdentifierExpression string
	//go:embed integer-literal-expression.monkey
	IntegerLiteralExpression string
	//go:embed boolean-literal-expression.monkey
	BooleanLiteralExpression string
)
