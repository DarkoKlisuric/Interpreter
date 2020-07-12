package object

import "fmt"

type ObjectType string

const (
	IntegerObj = "INTEGER"
	BooleanObj = "BOOLEAN"
	NullObj    = "NULL"
	ReturnValueObj = "RETURN_VALUE"
)

// represent every value which is encounter when evaluating source code as an Object
type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

type Boolean struct {
	Value bool
}

type ReturnValue struct {
	Value Object
}

type Null struct{}

func (i *Integer) Type() ObjectType { return IntegerObj }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

func (b *Boolean) Type() ObjectType { return BooleanObj }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

func (n *Null) Type() ObjectType { return NullObj }
func (n *Null) Inspect() string  { return "null" }

func (rv *ReturnValue) Type() ObjectType { return ReturnValueObj }
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

