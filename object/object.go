package object

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/Warashi/implement-interpreter-with-go/ast"
	"golang.org/x/exp/slices"
)

type BuiltinFunction func(args ...Object) Object

//go:generate go run golang.org/x/tools/cmd/stringer -type Type -trimprefix Type
type Type int

const (
	_ Type = iota
	TypeInteger
	TypeString
	TypeBoolean
	TypeNull
	TypeReturn
	TypeError
	TypeFunction
	TypeBuiltin
	TypeArray
	TypeHash
)

type Object interface {
	Type() Type
	Inspect() string
}

type Integer struct {
	Value int64
}

type String struct {
	Value string
}

type Boolean struct {
	Value bool
}

type Null struct{}

type Return struct {
	Value Object
}

type Error struct {
	Message string
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        Environment
}

type Builtin struct {
	Fn BuiltinFunction
}

type Array struct {
	Elements []Object
}

type Hash struct {
	Pairs map[Object]Object
}

func (o Integer) Type() Type      { return TypeInteger }
func (o Integer) Inspect() string { return strconv.FormatInt(o.Value, 10) }

func (o String) Type() Type      { return TypeString }
func (o String) Inspect() string { return o.Value }

func (o Boolean) Type() Type      { return TypeBoolean }
func (o Boolean) Inspect() string { return strconv.FormatBool(o.Value) }

func (o Null) Type() Type      { return TypeNull }
func (o Null) Inspect() string { return "null" }

func (o Return) Type() Type      { return TypeReturn }
func (o Return) Inspect() string { return o.Value.Inspect() }

func (o Error) Type() Type      { return TypeError }
func (o Error) Inspect() string { return "ERROR: " + o.Message }

func (o Function) Type() Type { return TypeFunction }
func (o Function) Inspect() string {
	var b strings.Builder
	params := make([]string, 0, len(o.Parameters))
	for _, p := range o.Parameters {
		params = append(params, p.String())
	}
	b.WriteString("fn(")
	b.WriteString(strings.Join(params, ", "))
	b.WriteString(") {\n")
	b.WriteString(o.Body.String())
	b.WriteString("\n}")
	return b.String()
}

func (o Builtin) Type() Type      { return TypeBuiltin }
func (o Builtin) Inspect() string { return "builtin function" }

func (o Array) Type() Type { return TypeArray }
func (o Array) Inspect() string {
	var b strings.Builder
	elements := make([]string, 0, len(o.Elements))
	for _, e := range o.Elements {
		elements = append(elements, e.Inspect())
	}
	b.WriteString("[")
	b.WriteString(strings.Join(elements, ", "))
	b.WriteString("]")
	return b.String()
}

func (o Hash) Type() Type { return TypeHash }
func (o Hash) Inspect() string {
	var b strings.Builder
	pairs := make([]string, 0, len(o.Pairs))
	for k, v := range o.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s:%s", k.Inspect(), v.Inspect()))
	}
	b.WriteString("{")
	b.WriteString(strings.Join(pairs, ", "))
	b.WriteString("}")
	return b.String()
}
func (o Hash) Equal(other Hash) bool { return reflect.DeepEqual(o.pairs(), other.pairs()) }
func (o Hash) pairs() [][2]Object {
	pairs := make([][2]Object, 0, len(o.Pairs))
	for k, v := range o.Pairs {
		pairs = append(pairs, [2]Object{k, v})
	}
	slices.SortFunc(pairs, func(a, b [2]Object) bool { return a[0].Inspect() < b[0].Inspect() })
	return pairs
}
