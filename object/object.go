package object

type ObjectType string

// represent every value which is encounter when evaluating source code as an Object
type Object interface {
	Type() ObjectType
	Inspect() string
}
