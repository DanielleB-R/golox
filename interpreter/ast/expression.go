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
	_ Expr = (*Super)(nil)
	_ Expr = (*This)(nil)
	_ Expr = (*Unary)(nil)
	_ Expr = (*Variable)(nil)
)

type Expr interface {
	expression()
	Accept(visitor ExprVisitor) any
}

type ExprVisitor interface {
	VisitAssign(assign *Assign) any
	VisitBinary(binary *Binary) any
	VisitCall(call *Call) any
	VisitGet(get *Get) any
	VisitGrouping(grouping *Grouping) any
	VisitLiteral(literal *Literal) any
	VisitLogical(logical *Logical) any
	VisitSet(set *Set) any
	VisitSuper(super *Super) any
	VisitThis(this *This) any
	VisitUnary(unary *Unary) any
	VisitVariable(variable *Variable) any
}

type Assign struct {
	Name  *token.Token
	Value Expr
}

func (*Assign) expression() {}
func (a *Assign) Accept(visitor ExprVisitor) any {
	return visitor.VisitAssign(a)
}

type Binary struct {
	Left     Expr
	Operator *token.Token
	Right    Expr
}

func (*Binary) expression() {}
func (b *Binary) Accept(visitor ExprVisitor) any {
	return visitor.VisitBinary(b)
}

type Call struct {
	Callee    Expr
	Paren     *token.Token
	Arguments []Expr
}

func (*Call) expression() {}
func (c *Call) Accept(visitor ExprVisitor) any {
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
func (g *Grouping) Accept(visitor ExprVisitor) any {
	return visitor.VisitGrouping(g)
}

type Literal struct {
	Value any
}

func (*Literal) expression() {}
func (l *Literal) Accept(visitor ExprVisitor) any {
	return visitor.VisitLiteral(l)
}

type Logical struct {
	Left     Expr
	Operator *token.Token
	Right    Expr
}

func (*Logical) expression() {}
func (l *Logical) Accept(visitor ExprVisitor) any {
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

type Super struct {
	Keyword *token.Token
	Method  *token.Token
}

func (*Super) expression() {}
func (s *Super) Accept(visitor ExprVisitor) any {
	return visitor.VisitSuper(s)
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
func (u *Unary) Accept(visitor ExprVisitor) any {
	return visitor.VisitUnary(u)
}

type Variable struct {
	Name *token.Token
}

func (*Variable) expression() {}
func (v *Variable) Accept(visitor ExprVisitor) any {
	return visitor.VisitVariable(v)
}
