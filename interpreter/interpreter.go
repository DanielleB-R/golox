package interpreter

import (
	"fmt"

	"github.com/DanielleB-R/golox/interpreter/ast"
	"github.com/DanielleB-R/golox/interpreter/token"
)

var (
	_ ast.ExprVisitor = (*Interpreter)(nil)
	_ ast.StmtVisitor = (*Interpreter)(nil)
)

type Interpreter struct {
	environment *Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		environment: NewEnvironment(),
	}
}

func (i *Interpreter) Interpret(statements []ast.Stmt) {
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		if runtimeError, ok := err.(*RuntimeError); ok {
			fmt.Println(runtimeError.Error())
			return
		}
		panic(err)
	}()

	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) execute(stmt ast.Stmt) {
	stmt.Accept(i)
}

func (i *Interpreter) VisitExpressionStmt(stmt *ast.ExpressionStmt) {
	i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitPrint(stmt *ast.Print) {
	value := i.evaluate(stmt.Expression)
	fmt.Println(value)
}

func (i *Interpreter) VisitVar(stmt *ast.Var) {
	var value interface{}
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}

	i.environment.Define(stmt.Name.Lexeme, value)
}

func (i *Interpreter) evaluate(expr ast.Expr) interface{} {
	return expr.Accept(i)
}

func (*Interpreter) VisitLiteral(literal *ast.Literal) interface{} {
	return literal.Value
}

func (i *Interpreter) VisitGrouping(grouping *ast.Grouping) interface{} {
	return i.evaluate(grouping.Expression)
}

func (i *Interpreter) VisitUnary(unary *ast.Unary) interface{} {
	right := i.evaluate(unary.Right)

	switch unary.Operator.TokenType {
	case token.BANG:
		return !isTruthy(right)
	case token.MINUS:
		r := checkNumberOperand(unary.Operator, right)
		return -r
	}

	// Should be unreachable
	return nil
}

func (i *Interpreter) VisitBinary(binary *ast.Binary) interface{} {
	left := i.evaluate(binary.Left)
	right := i.evaluate(binary.Right)

	switch binary.Operator.TokenType {
	case token.MINUS:
		l, r := checkNumberOperands(binary.Operator, left, right)
		return l - r
	case token.SLASH:
		l, r := checkNumberOperands(binary.Operator, left, right)
		return l / r
	case token.STAR:
		l, r := checkNumberOperands(binary.Operator, left, right)
		return l * r
	case token.PLUS:
		switch l := left.(type) {
		case float64:
			r := checkNumberOperand(binary.Operator, right)
			return l + r
		case string:
			return l + right.(string)
		default:
			panic(&RuntimeError{token: binary.Operator, message: "Operands must be numbers or strings"})
		}
	case token.GREATER:
		l, r := checkNumberOperands(binary.Operator, left, right)
		return l > r
	case token.GREATER_EQUAL:
		l, r := checkNumberOperands(binary.Operator, left, right)
		return l >= r
	case token.LESS:
		l, r := checkNumberOperands(binary.Operator, left, right)
		return l < r
	case token.LESS_EQUAL:
		l, r := checkNumberOperands(binary.Operator, left, right)
		return l <= r
	case token.BANG_EQUAL:
		return !isEqual(left, right)
	case token.EQUAL_EQUAL:
		return isEqual(left, right)
	}

	// Should be unreachable
	return nil
}

func (i *Interpreter) VisitVariable(expr *ast.Variable) interface{} {
	value, err := i.environment.Get(expr.Name)
	if err != nil {
		panic(err)
	}
	return value
}

func isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	if b, ok := object.(bool); ok {
		return b
	}
	return true
}

// This should be sufficient if I understand how Go equality is implemented
func isEqual(left interface{}, right interface{}) bool {
	return left == right
}

func checkNumberOperand(operator *token.Token, operand interface{}) float64 {
	numberOperand, ok := operand.(float64)
	if !ok {
		panic(&RuntimeError{token: operator, message: "Operand must be a number"})
	}
	return numberOperand
}

func checkNumberOperands(operator *token.Token, left interface{}, right interface{}) (float64, float64) {
	leftNumber, ok := left.(float64)
	rightNumber, ok2 := right.(float64)
	if !ok || !ok2 {
		panic(&RuntimeError{token: operator, message: "Operands must be numbers"})
	}
	return leftNumber, rightNumber
}
