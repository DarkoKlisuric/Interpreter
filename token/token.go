package token

// Defining TokenType to be a string beacuse
// that allows to use many different values as TokenType
type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	Illegal = "ILLEGAL" // signifies a token/character we don’t know about
	Eof     = "EOF"     // stands for “end of file”

	// Identifiers + literals
	Ident = "IDENT" // add, foo-bar, x, y, ...
	Int   = "INT"   // 1234123

	//Operators
	Assign   = "="
	Plus     = "+"
	Minus    = "-"
	Bang    = "!"
	Asterisk = "*"
	Slash    = "/"

	Lt = "<"
	Gt = ">"

	LParen = "("
	RParen = ")"
	LBrace = "{"
	RBrace = "}"

	Comma     = ","
	Semicolon = ";"

	Eq     = "=="
	NotEq = "!="

	// Keywords
	Function = "FUNCTION"
	Let      = "LET"
	True     = "TRUE"
	False    = "FALSE"
	If       = "IF"
	Else     = "ELSE"
	Return   = "RETURN"
)

// Keywords in this programming language
var keywords = map[string]TokenType{
	"fn":     Function,
	"let":    Let,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"return": Return,
}

// Checks the keywords table to see whether the given identifier is in fact a keyword
// If it is, it returns the keyword’s TokenType constant.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return Ident
}
