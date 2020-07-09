package parser

import (
	"fmt"
	"github.com/DarkoKlisuric/interpreter/ast"
	"github.com/DarkoKlisuric/interpreter/lexer"
	"github.com/DarkoKlisuric/interpreter/token"
	"strconv"
)

const (
	_ int = iota
	Lowest
	Equals      // ==
	LessGreater // > or <
	Sum         // +
	Product     // *
	Prefix      // -x or !x
	Call        // myFunction(x)

)

// ( token.PLUS ) and - ( token.MINUS ) have the same precedence,
//  which is lower than the precedence of * ( token.ASTERISK ) and / ( token.SLASH ), for example.
var precedences = map[token.TokenType]int{
	token.Eq:       Equals,
	token.NotEq:    Equals,
	token.Lt:       LessGreater,
	token.Gt:       LessGreater,
	token.Plus:     Sum,
	token.Minus:    Sum,
	token.Slash:    Product,
	token.Asterisk: Product,
	token.LParen:   Call,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
	errors    []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.Ident, p.parseIdentifier)
	p.registerPrefix(token.Int, p.parseIntegerLiteral)
	p.registerPrefix(token.Bang, p.parsePrefixExpression)
	p.registerPrefix(token.Minus, p.parsePrefixExpression)
	p.registerPrefix(token.True, p.parseBoolean)
	p.registerPrefix(token.False, p.parseBoolean)
	p.registerPrefix(token.LParen, p.parseGroupedExpression)
	p.registerPrefix(token.If, p.parseIfExpression)
	p.registerPrefix(token.Function, p.parseFunctionLiteral)

	// Register one infix parse function for all of our infix operators
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.Plus, p.parseInfixExpression)
	p.registerInfix(token.Minus, p.parseInfixExpression)
	p.registerInfix(token.Slash, p.parseInfixExpression)
	p.registerInfix(token.Asterisk, p.parseInfixExpression)
	p.registerInfix(token.Eq, p.parseInfixExpression)
	p.registerInfix(token.NotEq, p.parseInfixExpression)
	p.registerInfix(token.Lt, p.parseInfixExpression)
	p.registerInfix(token.Gt, p.parseInfixExpression)
	p.registerInfix(token.LParen, p.parseCallExpression)

	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

// Return errors from parser
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s insted",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// The first thing ParseProgram does is construct the root node of the AST, an *ast.Program . It
// then iterates over every token in the input until it encounters an token.EOF token.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.Eof {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.Let:
		return p.parseLetStatement()
	case token.Return:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// Method tries to find infixParseFn s for the next token. If it finds such a function, it calls it, passing
// in the expression returned by a prefixParseFn as an argument. And it does all this again and
// again until it encounters a token that has a higher precedence
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.Semicolon) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]

		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.expectPeek(token.Ident) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(token.Assign) {
		return nil
	}

	// Skipping the expressions until encounter a semicolon
	p.nextToken()

	stmt.Value = p.parseExpression(Lowest)

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	// Skipping the expressions until encounter semicolon
	stmt.ReturnValue = p.parseExpression(Lowest)

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(Lowest)

	if p.peekTokenIs(token.Semicolon) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)

		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(Prefix)

	return expression
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)

	p.errors = append(p.errors, msg)
}

// Method returns the precedence associated with the token type of peekToken
// If it doesn’t find a precedence for p.peekToken it defaults to LOWEST , the lowest
// possible precedence any operator can have.
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return Lowest
}

// This method does the same thing as peekPrecedence, but for p.curToken.
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return Lowest
}

// It uses ast.Expression argument to construct an
// *ast.InfixExpression node, with left being in the Left field.
// Then it assigns the precedence of the current token
// (which is the operator of the infix expression) to the local variable precedence ,
// before advancing the tokens by calling nextToken and filling the Right field of the node with
// another call to parseExpression - this time passing in the precedence of the operator token
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expession := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()

	p.nextToken()

	expession.Right = p.parseExpression(precedence)

	return expession
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.curToken, Value: p.curTokenIs(token.True)}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(Lowest)

	if !p.expectPeek(token.RParen) {
		return nil
	}

	return exp
}

// If have a token.ELSE, advance the tokens two times.
func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}

	if !p.expectPeek(token.LParen) {
		return nil
	}

	p.nextToken()

	expression.Condition = p.parseExpression(Lowest)

	if !p.expectPeek(token.RParen) {
		return nil
	}

	if !p.expectPeek(token.LBrace) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.Else) {
		p.nextToken()

		if !p.expectPeek(token.LBrace) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBrace) && !p.curTokenIs(token.Eof) {
		stmt := p.parseStatement()

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if !p.expectPeek(token.LParen) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBrace) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

// method parse the literal’s parameters
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var identifiers []*ast.Identifier

	if p.peekTokenIs(token.RParen) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RParen) {
		return nil
	}

	return identifiers
}

// receives the already parsed function as argument and uses it to construct
// an *ast.CallExpression node
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseCallArguments()

	return exp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	var args []ast.Expression

	if p.peekTokenIs(token.RParen) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(Lowest))

	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(Lowest))
	}

	if !p.expectPeek(token.RParen) {
		return nil
	}

	return args
}
