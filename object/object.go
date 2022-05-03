package object

import "strconv"

//go:generate go run golang.org/x/tools/cmd/stringer -type Type -trimprefix Type
type Type int

const (
	_ Type = iota
	TypeInteger
	TypeBoolean
	TypeNull
	TypeReturn
	TypeError
)

type Object interface {
	Type() Type
	Inspect() string
}

type Integer struct {
	Value int64
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

func (o Integer) Type() Type      { return TypeInteger }
func (o Integer) Inspect() string { return strconv.FormatInt(o.Value, 10) }

func (o Boolean) Type() Type      { return TypeBoolean }
func (o Boolean) Inspect() string { return strconv.FormatBool(o.Value) }

func (o Null) Type() Type      { return TypeNull }
func (o Null) Inspect() string { return "null" }

func (o Return) Type() Type      { return TypeReturn }
func (o Return) Inspect() string { return o.Value.Inspect() }

func (o Error) Type() Type      { return TypeError }
func (o Error) Inspect() string { return "ERROR: " + o.Message }
