package testutil

import (
	"github.com/Warashi/implement-interpreter-with-go/ast"
	"github.com/Warashi/implement-interpreter-with-go/object"
)

func IntegerObject(val int64) object.Object {
	return object.Integer{Value: val}
}

func StringObject(val string) object.Object {
	return object.String{Value: val}
}

func BooleanObject(val bool) object.Object {
	return object.Boolean{Value: val}
}

func NullObject() object.Object {
	return object.Null{}
}

func ErrorObject(message string) object.Object {
	return object.Error{Message: message}
}

func FunctionObject(env object.Environment, body *ast.BlockStatement, params ...*ast.Identifier) object.Object {
	return object.Function{
		Parameters: params,
		Body:       body,
		Env:        env,
	}
}
