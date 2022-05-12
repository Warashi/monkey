package evaluator

import "github.com/Warashi/monkey/object"

var builtins = map[string]object.Builtin{
	"len":   {Fn: builtinLen},
	"first": {Fn: builtinFirst},
	"last":  {Fn: builtinLast},
	"rest":  {Fn: builtinRest},
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

func builtinFirst(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newErrorf("wrong number of arguments. got=%d, want=%d", len(args), 1)
	}
	if args[0].Type() != object.TypeArray {
		return newErrorf("argument to `first` not supported, got %s", args[0].Type())
	}
	arr := args[0].(object.Array)
	if len(arr.Elements) == 0 {
		return NULL
	}
	return arr.Elements[0]
}

func builtinLast(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newErrorf("wrong number of arguments. got=%d, want=%d", len(args), 1)
	}
	if args[0].Type() != object.TypeArray {
		return newErrorf("argument to `first` not supported, got %s", args[0].Type())
	}
	arr := args[0].(object.Array)
	if len(arr.Elements) == 0 {
		return NULL
	}
	return arr.Elements[len(arr.Elements)-1]
}

func builtinRest(args ...object.Object) object.Object {
	if len(args) != 1 {
		return newErrorf("wrong number of arguments. got=%d, want=%d", len(args), 1)
	}
	if args[0].Type() != object.TypeArray {
		return newErrorf("argument to `first` not supported, got %s", args[0].Type())
	}
	arr := args[0].(object.Array)
	if len(arr.Elements) == 0 {
		return NULL
	}
	elements := make([]object.Object, len(arr.Elements)-1)
	copy(elements, arr.Elements[1:])
	return object.Array{Elements: elements}
}

func builtinPush(args ...object.Object) object.Object {
	if len(args) != 2 {
		return newErrorf("wrong number of arguments. got=%d, want=%d", len(args), 2)
	}
	if args[0].Type() != object.TypeArray {
		return newErrorf("argument to `first` not supported, got %s", args[0].Type())
	}
	arr := args[0].(object.Array)
	if len(arr.Elements) == 0 {
		return NULL
	}
	elements := make([]object.Object, len(arr.Elements), len(arr.Elements)+1)
	copy(elements, arr.Elements)
	elements = append(elements, args[1])
	return object.Array{Elements: elements}
}
