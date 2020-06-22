package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL" // signifies a token/character we don’t know about
	EOF     = "EOF"     // stands for “end of file”

	// Identifiers + literals
	IDENT = "IDENT" // add, foo-bar, x, y, ...
	INT   = "INT"   // 1234123

	//Operators
	ASSIGN = "="
	PLUS   = "+"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	COMMA = ","
	SEMICOLON = ";"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

var keywords = map[string]TokenType{
	"fn": FUNCTION,
	"let": LET,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

