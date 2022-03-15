package interpreter

import (
	"fmt"

	"github.com/DanielleB-R/golox/interpreter/ast"
	"github.com/DanielleB-R/golox/interpreter/token"
)

type Parser struct {
	tokens  []*token.Token
	current int
}

func NewParser(tokens []*token.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() (ast.Expr, error) {
	expr, err := p.expression()
	if err == nil {
		return expr, nil
	}

	_, ok := err.(*ParseError)
	if !ok {
		return nil, err
	}
	fmt.Println(err.Error())
	return nil, nil
}

// Grammar rules

func (p *Parser) expression() (ast.Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (ast.Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) comparison() (ast.Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) term() (ast.Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) factor() (ast.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = &ast.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) unary() (ast.Expr, error) {
	if p.match(token.BANG, token.MINUS) {
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return &ast.Unary{
			Operator: p.previous(),
			Right:    right,
		}, nil
	}
	return p.primary()
}

func (p *Parser) primary() (ast.Expr, error) {
	if p.match(token.FALSE) {
		return &ast.Literal{
			Value: false,
		}, nil
	}
	if p.match(token.TRUE) {
		return &ast.Literal{
			Value: true,
		}, nil
	}
	if p.match(token.NIL) {
		return &ast.Literal{
			Value: nil,
		}, nil
	}

	if p.match(token.NUMBER, token.STRING) {
		return &ast.Literal{
			Value: p.previous().Literal,
		}, nil
	}

	if p.match(token.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return &ast.Grouping{
			Expression: expr,
		}, nil
	}

	return nil, &ParseError{
		token:   p.peek(),
		message: "Expect expression.",
	}
}

// Helpers

func (p *Parser) peek() *token.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() *token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().TokenType == token.EOF
}

func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.current += 1
	}
	return p.previous()
}

func (p *Parser) check(tokenType int) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().TokenType == tokenType
}

func (p *Parser) match(types ...int) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(tokenType int, message string) (*token.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}

	return nil, &ParseError{token: p.peek(), message: message}
}
