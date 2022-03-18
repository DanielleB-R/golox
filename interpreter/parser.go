package interpreter

import (
	"github.com/DanielleB-R/golox/interpreter/ast"
	"github.com/DanielleB-R/golox/interpreter/token"
)

type Parser struct {
	tokens  []*token.Token
	current int
	errors  ParseErrors
}

func NewParser(tokens []*token.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
		errors:  ParseErrors{},
	}
}

func (p *Parser) Parse() ([]ast.Stmt, error) {
	statements := []ast.Stmt{}
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			parseErr, ok := err.(*ParseError)
			if !ok {
				return statements, err
			}
			p.errors = append(p.errors, parseErr)
			p.synchronize()
		}
		statements = append(statements, stmt)
	}

	return statements, nil
}

// Grammar rules

func (p *Parser) declaration() (ast.Stmt, error) {
	if p.match(token.FUN) {
		return p.function("function")
	}
	if p.match(token.VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) function(kind string) (ast.Stmt, error) {
	name, err := p.consume(token.IDENTIFIER, "Expect "+kind+" name.")
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.LEFT_PAREN, "Expect '(' after "+kind+" name.")
	if err != nil {
		return nil, err
	}

	parameters := []*token.Token{}
	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				return nil, &ParseError{token: p.peek(), message: "Can't have more than 255 parameters"}
			}
			name, err := p.consume(token.IDENTIFIER, "Expect parameter name.")
			if err != nil {
				return nil, err
			}
			parameters = append(parameters, name)
			if !p.match(token.COMMA) {
				break
			}
		}
	}
	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.LEFT_BRACE, "Expect '{' before "+kind+" body.")
	if err != nil {
		return nil, err
	}

	body, err := p.block()
	if err != nil {
		return nil, err
	}

	return &ast.Function{
		Name:   name,
		Params: parameters,
		Body:   body,
	}, nil

}

func (p *Parser) varDeclaration() (ast.Stmt, error) {
	name, err := p.consume(token.IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	var initializer ast.Expr
	if p.match(token.EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	_, err = p.consume(token.SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}
	return &ast.Var{
		Name:        name,
		Initializer: initializer,
	}, nil
}

func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(token.FOR) {
		return p.forStatement()
	}
	if p.match(token.IF) {
		return p.ifStatement()
	}
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	if p.match(token.LEFT_BRACE) {
		block, err := p.block()
		if err != nil {
			return nil, err
		}
		return &ast.Block{
			Statements: block,
		}, nil
	}
	if p.match(token.WHILE) {
		return p.whileStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) forStatement() (ast.Stmt, error) {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	var initializer ast.Stmt
	if !p.match(token.SEMICOLON) {
		if p.match(token.VAR) {
			initializer, err = p.varDeclaration()
			if err != nil {
				return nil, err
			}
		} else {
			initializer, err = p.expressionStatement()
			if err != nil {
				return nil, err
			}
		}
	}

	var condition ast.Expr
	if !p.check(token.SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.SEMICOLON, "Expect ';' after loop condition")
	if err != nil {
		return nil, err
	}

	var increment ast.Expr
	if !p.check(token.RIGHT_PAREN) {
		increment, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after for clauses")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = &ast.Block{
			Statements: []ast.Stmt{
				body,
				&ast.ExpressionStmt{
					Expression: increment,
				},
			},
		}
	}

	if condition == nil {
		condition = &ast.Literal{Value: true}
	}
	body = &ast.While{
		Condition: condition,
		Body:      body,
	}

	if initializer != nil {
		body = &ast.Block{
			Statements: []ast.Stmt{initializer, body},
		}
	}

	return body, nil
}

func (p *Parser) ifStatement() (ast.Stmt, error) {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after if condition")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}
	var elseBranch ast.Stmt
	if p.match(token.ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return &ast.If{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}, nil
}

func (p *Parser) printStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.SEMICOLON, "Expect ';' after expression.")
	if err != nil {
		return nil, err
	}

	return &ast.Print{
		Expression: expr,
	}, nil
}

func (p *Parser) whileStatement() (ast.Stmt, error) {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after while condition")
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return &ast.While{
		Condition: condition,
		Body:      body,
	}, nil
}

func (p *Parser) block() ([]ast.Stmt, error) {
	statements := []ast.Stmt{}

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		statement, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}

	_, err := p.consume(token.RIGHT_BRACE, "Expect '}' after block")
	if err != nil {
		return nil, err
	}
	return statements, nil
}

func (p *Parser) expressionStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	_, err = p.consume(token.SEMICOLON, "Expect ';' after expression.")
	if err != nil {
		return nil, err
	}

	return &ast.ExpressionStmt{
		Expression: expr,
	}, nil
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (ast.Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(token.EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}
		if variable, ok := expr.(*ast.Variable); ok {
			return &ast.Assign{
				Name:  variable.Name,
				Value: value,
			}, nil
		}
		// TODO: This should not trigger a resynchronization of the parser
		return nil, &ParseError{
			token:   equals,
			message: "Invalid assignment target",
		}
	}

	return expr, nil
}

func (p *Parser) or() (ast.Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(token.OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		expr = &ast.Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
}

func (p *Parser) and() (ast.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(token.AND) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return nil, err
		}
		expr = &ast.Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr, nil
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
	return p.call()
}

func (p *Parser) call() (ast.Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	for {
		if p.match(token.LEFT_PAREN) {
			expr, err = p.finishCall(expr)
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}

	return expr, nil
}

func (p *Parser) finishCall(callee ast.Expr) (ast.Expr, error) {
	arguments := []ast.Expr{}
	if !p.check(token.RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				// NOTE: this should be non-resynchronizing
				return nil, &ParseError{
					token:   p.peek(),
					message: "Can't have more than 255 arguments",
				}
			}
			expression, err := p.expression()
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, expression)
			if !p.match(token.COMMA) {
				break
			}
		}
	}

	paren, err := p.consume(token.RIGHT_PAREN, "Expect ')' after arguments.")
	if err != nil {
		return nil, err
	}

	return &ast.Call{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}, nil
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

	if p.match(token.IDENTIFIER) {
		return &ast.Variable{
			Name: p.previous(),
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

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().TokenType == token.SEMICOLON {
			return
		}

		switch p.peek().TokenType {
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}

		p.advance()
	}
}
