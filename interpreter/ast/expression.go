package ast

import "github.com/DanielleB-R/golox/interpreter/token"

var (
	_ Expr = (*Binary)(nil)
	_ Expr = (*Grouping)(nil)
	_ Expr = (*Literal)(nil)
	_ Expr = (*Unary)(nil)
	_ Expr = (*Variable)(nil)
)

type Expr interface {
	expression()
	Accept(visitor ExprVisitor) interface{}
}

type ExprVisitor interface {
	VisitBinary(binary *Binary) interface{}
	VisitGrouping(grouping *Grouping) interface{}
	VisitLiteral(literal *Literal) interface{}
	VisitUnary(unary *Unary) interface{}
	VisitVariable(variable *Variable) interface{}
}

type Binary struct {
	Left     Expr
	Operator *token.Token
	Right    Expr
}

func (*Binary) expression() {}
func (b *Binary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitBinary(b)
}

type Grouping struct {
	Expression Expr
}

func (*Grouping) expression() {}
func (g *Grouping) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitGrouping(g)
}

type Literal struct {
	Value interface{}
}

func (*Literal) expression() {}
func (l *Literal) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLiteral(l)
}

type Unary struct {
	Operator *token.Token
	Right    Expr
}

func (*Unary) expression() {}
func (u *Unary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitUnary(u)
}

type Variable struct {
	Name *token.Token
}

func (*Variable) expression() {}
func (v *Variable) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitVariable(v)
}
