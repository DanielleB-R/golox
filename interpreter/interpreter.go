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
	globals           *Environment
	environment       *Environment
	activeReturn      bool
	activeReturnValue interface{}
}

func NewInterpreter() *Interpreter {
	globals := NewEnvironment(nil)
	globals.Define("clock", Clock)
	return &Interpreter{
		globals:           globals,
		environment:       globals,
		activeReturn:      false,
		activeReturnValue: nil,
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
		i.resetReturnValue()
	}
}

func (i *Interpreter) execute(stmt ast.Stmt) {
	stmt.Accept(i)
}

func (i *Interpreter) VisitBlock(stmt *ast.Block) {
	i.executeBlock(stmt.Statements, NewEnvironment(i.environment))
}

func (i *Interpreter) executeBlock(statements []ast.Stmt, environment *Environment) {
	previous := i.environment
	defer func() {
		i.environment = previous
	}()

	i.environment = environment

	for _, statement := range statements {
		i.execute(statement)
		if i.activeReturn {
			return
		}
	}
}

func (i *Interpreter) VisitExpressionStmt(stmt *ast.ExpressionStmt) {
	i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitFunction(stmt *ast.Function) {
	function := NewLoxFunction(stmt)
	i.environment.Define(stmt.Name.Lexeme, function)
}

func (i *Interpreter) VisitIf(stmt *ast.If) {
	if isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
}

func (i *Interpreter) VisitPrint(stmt *ast.Print) {
	value := i.evaluate(stmt.Expression)
	fmt.Println(value)
}

func (i *Interpreter) VisitReturn(stmt *ast.Return) {
	var value interface{}
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}

	i.activeReturn = true
	i.activeReturnValue = value
}

func (i *Interpreter) VisitVar(stmt *ast.Var) {
	var value interface{}
	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}

	i.environment.Define(stmt.Name.Lexeme, value)
}

func (i *Interpreter) VisitWhile(stmt *ast.While) {
	for isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
		if i.activeReturn {
			return
		}
	}
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

func (i *Interpreter) VisitAssign(expr *ast.Assign) interface{} {
	value := i.evaluate(expr.Value)
	err := i.environment.Assign(expr.Name, value)
	if err != nil {
		panic(err)
	}
	return value
}

func (i *Interpreter) VisitLogical(expr *ast.Logical) interface{} {
	left := i.evaluate(expr.Left)

	if expr.Operator.TokenType == token.OR {
		if isTruthy(left) {
			return left
		}
	} else {
		if !isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitCall(expr *ast.Call) interface{} {
	callee := i.evaluate(expr.Callee)

	arguments := []interface{}{}
	for _, argument := range expr.Arguments {
		arguments = append(arguments, i.evaluate(argument))
	}

	function, ok := callee.(Callable)
	if !ok {
		panic(&RuntimeError{token: expr.Paren, message: "Can only call functions and classes."})
	}
	if len(arguments) != function.Arity() {
		panic(&RuntimeError{token: expr.Paren, message: fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(arguments))})
	}
	return function.Call(i, arguments)
}

func (i *Interpreter) resetReturnValue() {
	i.activeReturn = false
	i.activeReturnValue = nil
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
