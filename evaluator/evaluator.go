package evaluator

import (
	"github.com/DarkoKlisuric/interpreter/ast"
	"github.com/DarkoKlisuric/interpreter/object"
)

var (
	True = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
)

// should always start at the top of the tree, receiving an *ast.Program,
// and then traverse every node in it
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		// No need for creating a new object.Boolean every time when encounter a true or false
		return nativeBoolToBooleanObject(node.Value)
	}

	return nil
}

// evaluate every statement
// if the statement is an *ast.ExpressionStatement we evaluate its expression
func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return True
	}

	return False
}