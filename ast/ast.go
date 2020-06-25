package ast

import "github.com/DarkoKlisuric/interpreter/token"

type Node interface {
	// TokenLiteral() method that returns
	// the literal value of the token it’s associated with
	// TokenLiteral() will be used only for debugging and testing
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier // Name hold identifier of the binding
	Value Expression  // Value for the expression that produces the value
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
