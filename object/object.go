package object

import "strconv"

type Type string

const (
	TypeInteger Type = "INTEGER"
	TypeBoolean Type = "BOOLEAN"
	TypeNull    Type = "NULL"
	TypeReturn  Type = "RETURN"
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

func (o Integer) Type() Type      { return TypeInteger }
func (o Integer) Inspect() string { return strconv.FormatInt(o.Value, 10) }

func (o Boolean) Type() Type      { return TypeBoolean }
func (o Boolean) Inspect() string { return strconv.FormatBool(o.Value) }

func (o Null) Type() Type      { return TypeNull }
func (o Null) Inspect() string { return "null" }

func (o Return) Type() Type      { return TypeReturn }
func (o Return) Inspect() string { return o.Value.Inspect() }
