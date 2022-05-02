package testdata

import _ "embed"

//go:embed let.monkey
var Let string

//go:embed return.monkey
var Return string

//go:embed identifier-expression.monkey
var IdentifierExpression string
