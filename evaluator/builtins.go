package evaluator

import "github.com/Warashi/implement-interpreter-with-go/object"

var builtins = map[string]object.Builtin{
	"len": {Fn: builtinLen},
}

func builtinLen(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newErrorf("wrong number of arguments. got=%d, want=%d", len(args), 1)
	}
	switch args[0].Type() {
	case object.TypeString:
		return object.Integer{Value: int64(len(args[0].(object.String).Value))}
	case object.TypeArray:
		return object.Integer{Value: int64(len(args[0].(object.Array).Elements))}
	default:
		return newErrorf("argument to `len` not supported, got %s", args[0].Type())
	}
}
