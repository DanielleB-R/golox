package ast

import "github.com/DanielleB-R/golox/interpreter/token"

var (
	_ Expr = (*Assign)(nil)
	_ Expr = (*Binary)(nil)
	_ Expr = (*Call)(nil)
	_ Expr = (*Get)(nil)
	_ Expr = (*Grouping)(nil)
	_ Expr = (*Literal)(nil)
	_ Expr = (*Logical)(nil)
	_ Expr = (*Set)(nil)
	_ Expr = (*This)(nil)
	_ Expr = (*Unary)(nil)
	_ Expr = (*Variable)(nil)
)

type Expr interface {
	expression()
	Accept(visitor ExprVisitor) interface{}
}

type ExprVisitor interface {
	VisitAssign(assign *Assign) interface{}
	VisitBinary(binary *Binary) interface{}
	VisitCall(call *Call) interface{}
	VisitGet(get *Get) any
	VisitGrouping(grouping *Grouping) interface{}
	VisitLiteral(literal *Literal) interface{}
	VisitLogical(logical *Logical) interface{}
	VisitSet(set *Set) any
	VisitThis(this *This) any
	VisitUnary(unary *Unary) interface{}
	VisitVariable(variable *Variable) interface{}
}

type Assign struct {
	Name  *token.Token
	Value Expr
}

func (*Assign) expression() {}
func (a *Assign) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitAssign(a)
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

type Call struct {
	Callee    Expr
	Paren     *token.Token
	Arguments []Expr
}

func (*Call) expression() {}
func (c *Call) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitCall(c)
}

type Get struct {
	Object Expr
	Name   *token.Token
}

func (*Get) expression() {}
func (g *Get) Accept(visitor ExprVisitor) any {
	return visitor.VisitGet(g)
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

type Logical struct {
	Left     Expr
	Operator *token.Token
	Right    Expr
}

func (*Logical) expression() {}
func (l *Logical) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLogical(l)
}

type Set struct {
	Object Expr
	Name   *token.Token
	Value  Expr
}

func (*Set) expression() {}
func (s *Set) Accept(visitor ExprVisitor) any {

	return visitor.VisitSet(s)
}

type This struct {
	Keyword *token.Token
}

func (*This) expression() {}
func (t *This) Accept(visitor ExprVisitor) any {
	return visitor.VisitThis(t)
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
