package testutil

import (
	"github.com/Warashi/implement-interpreter-with-go/ast"
	"github.com/Warashi/implement-interpreter-with-go/object"
)

func IntegerObject(val int64) object.Integer {
	return object.Integer{Value: val}
}

func StringObject(val string) object.String {
	return object.String{Value: val}
}

func BooleanObject(val bool) object.Boolean {
	return object.Boolean{Value: val}
}

func NullObject() object.Null {
	return object.Null{}
}

func ErrorObject(message string) object.Error {
	return object.Error{Message: message}
}

func FunctionObject(env object.Environment, body *ast.BlockStatement, params ...*ast.Identifier) object.Function {
	return object.Function{
		Parameters: params,
		Body:       body,
		Env:        env,
	}
}

func ArrayObject(elements ...object.Object) object.Array {
	return object.Array{Elements: elements}
}

func HashObject(pairs map[object.Hashable]object.Object) object.Hash {
	return object.Hash{Pairs: pairs}
}
