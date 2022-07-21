package parser

import (
	"monkey/ast"
	"monkey/errors"
	"monkey/lexer"
	"monkey/token"
	"strconv"
)

type (
	prefixParseFunction func() ast.Expression
	infixParseFunction  func(ast.Expression) ast.Expression
)

const (
	LOWEST int = iota
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL
)

var precedences = map[token.TokenType]int{
	token.EQUAL:     EQUALS,
	token.NOT_EQUAL: EQUALS,

	token.L_THAN: LESSGREATER,
	token.G_THAN: LESSGREATER,

	token.PLUS:  SUM,
	token.MINUS: SUM,

	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,

	token.LPAREN: CALL,
}

type Parser struct {
	lexer *lexer.Lexer

	currentToken token.Token
	nextToken    token.Token

	errors []string

	prefixParseFunctions map[token.TokenType]prefixParseFunction
	infixParseFunctions  map[token.TokenType]infixParseFunction
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lexer: l, errors: []string{}}

	p.prefixParseFunctions = make(map[token.TokenType]prefixParseFunction)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)

	p.infixParseFunctions = make(map[token.TokenType]infixParseFunction)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQUAL, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.L_THAN, p.parseInfixExpression)
	p.registerInfix(token.G_THAN, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	p.advanceToNextToken()
	p.advanceToNextToken()

	return p
}

func (p *Parser) registerPrefix(tokenType token.TokenType, parseFunction prefixParseFunction) {
	p.prefixParseFunctions[tokenType] = parseFunction
}
func (p *Parser) registerInfix(tokenType token.TokenType, parseFunction infixParseFunction) {
	p.infixParseFunctions[tokenType] = parseFunction
}

func (p *Parser) advanceToNextToken() {
	p.currentToken = p.nextToken
	p.nextToken = p.lexer.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for !p.expectCurrentTokenToBe(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.advanceToNextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currentToken}

	if !p.expectNextTokenToBe(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Value,
	}
	if !p.expectNextTokenToBe(token.ASSIGN) {
		return nil
	}

	p.advanceToNextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.nextToken.Type == token.SEMICOLON {
		p.advanceToNextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}
	p.advanceToNextToken()
	stmt.ReturnValue = p.parseExpression(LOWEST)

	for !p.expectCurrentTokenToBe(token.SEMICOLON) {
		p.advanceToNextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{
		Token: p.currentToken,
	}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.nextToken.Type == token.SEMICOLON {
		p.advanceToNextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFunctions[p.currentToken.Type]
	if prefix == nil {
		p.errors = append(p.errors, errors.NoPrefixParseError(p.currentToken.Type))
		return nil
	}

	leftSideExpression := prefix()

	for p.nextToken.Type != token.SEMICOLON && precedence < p.peekPrecedence() {
		infix := p.infixParseFunctions[p.nextToken.Type]
		if infix == nil {
			return leftSideExpression
		}

		p.advanceToNextToken()

		leftSideExpression = infix(leftSideExpression)
	}

	return leftSideExpression
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Value,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	literal := &ast.IntegerLiteral{
		Token: p.currentToken,
	}
	value, err := strconv.ParseInt(p.currentToken.Value, 0, 64)
	if err != nil {
		p.errors = append(p.errors, errors.CouldNotParseInteger(p.currentToken.Value))
		return nil
	}

	literal.Value = value

	return literal
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Value,
	}

	p.advanceToNextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Value,
		Left:     left,
	}

	precedence := p.currentPrecedence()
	p.advanceToNextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.currentToken,
		Value: p.currentToken.Type == token.TRUE,
	}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.advanceToNextToken()

	expr := p.parseExpression(LOWEST)

	if !p.expectNextTokenToBe(token.RPAREN) {
		return nil
	}

	return expr
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.currentToken}

	if !p.expectNextTokenToBe(token.LPAREN) {
		return nil
	}

	p.advanceToNextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectNextTokenToBe(token.RPAREN) {
		return nil
	}

	if !p.expectNextTokenToBe(token.LBRACE) {
		return nil
	}

	expression.Consequence = p.parseBlockStatement()

	if p.nextToken.Type == token.ELSE {
		p.advanceToNextToken()

		if !p.expectNextTokenToBe(token.LBRACE) {
			return nil
		}

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token: p.currentToken,
	}
	block.Statements = []ast.Statement{}

	p.advanceToNextToken()

	for p.currentToken.Type != token.RBRACE && p.currentToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.advanceToNextToken()
	}

	return block
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	literal := &ast.FunctionLiteral{Token: p.currentToken}

	if !p.expectNextTokenToBe(token.LPAREN) {
		return nil
	}

	literal.Parameters = p.parseFunctionParameters()

	if !p.expectNextTokenToBe(token.LBRACE) {
		return nil
	}

	literal.Body = p.parseBlockStatement()

	return literal
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.nextToken.Type == token.RPAREN {
		p.advanceToNextToken()
		return identifiers
	}

	p.advanceToNextToken()

	identifier := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}
	identifiers = append(identifiers, identifier)

	for p.nextToken.Type == token.COMMA {
		p.advanceToNextToken()
		p.advanceToNextToken()

		identifier := &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Value}
		identifiers = append(identifiers, identifier)
	}

	if !p.expectNextTokenToBe(token.RPAREN) {
		return nil
	}

	return identifiers
}

func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expr := &ast.CallExpression{Token: p.currentToken, Function: function}
	expr.Arguments = p.parseCallArguments()
	return expr
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.nextToken.Type == token.RPAREN {
		p.advanceToNextToken()
		return args
	}

	p.advanceToNextToken()

	args = append(args, p.parseExpression(LOWEST))

	for p.nextToken.Type == token.COMMA {
		p.advanceToNextToken()
		p.advanceToNextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectNextTokenToBe(token.RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) expectCurrentTokenToBe(token token.TokenType) bool {
	return p.currentToken.Type == token
}

func (p *Parser) expectNextTokenToBe(token token.TokenType) bool {
	if p.nextToken.Type == token {
		p.advanceToNextToken()
		return true
	}
	p.catchPeekError(token)
	return false
}

func (p *Parser) catchPeekError(t token.TokenType) {
	p.errors = append(p.errors, errors.ExpectedNextTokenToBe(t, token.TokenType(p.nextToken.Value)))
}

func (p *Parser) tokenPrecedence(t *token.Token) int {
	if precedence, ok := precedences[t.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	return p.tokenPrecedence(&p.nextToken)
}

func (p *Parser) currentPrecedence() int {
	return p.tokenPrecedence(&p.currentToken)
}

func (p *Parser) Errors() []string {
	return p.errors
}
