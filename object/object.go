package object

import "strconv"

type Type string

const (
	INTEGER_OBJ Type = "INTEGER"
	BOOLEAN_OBJ Type = "BOOLEAN"
	NULL_OJB    Type = "NULL"
	RETURN_OBJ  Type = "RETURN"
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

func (o Integer) Type() Type      { return INTEGER_OBJ }
func (o Integer) Inspect() string { return strconv.FormatInt(o.Value, 10) }

func (o Boolean) Type() Type      { return BOOLEAN_OBJ }
func (o Boolean) Inspect() string { return strconv.FormatBool(o.Value) }

func (o Null) Type() Type      { return NULL_OJB }
func (o Null) Inspect() string { return "null" }

func (o Return) Type() Type      { return RETURN_OBJ }
func (o Return) Inspect() string { return o.Value.Inspect() }
