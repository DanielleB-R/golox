package ast

import "github.com/DanielleB-R/golox/interpreter/token"

var (
	_ Stmt = (*Block)(nil)
	_ Stmt = (*ExpressionStmt)(nil)
	_ Stmt = (*Function)(nil)
	_ Stmt = (*If)(nil)
	_ Stmt = (*Print)(nil)
	_ Stmt = (*Return)(nil)
	_ Stmt = (*Var)(nil)
	_ Stmt = (*While)(nil)
)

type Stmt interface {
	statement()
	Accept(visitor StmtVisitor)
}

type StmtVisitor interface {
	VisitBlock(stmt *Block)
	VisitExpressionStmt(stmt *ExpressionStmt)
	VisitFunction(stmt *Function)
	VisitIf(stmt *If)
	VisitPrint(stmt *Print)
	VisitReturn(stmt *Return)
	VisitVar(stmt *Var)
	VisitWhile(stmt *While)
}

type ExpressionStmt struct {
	Expression Expr
}

func (*ExpressionStmt) statement() {}
func (e *ExpressionStmt) Accept(visitor StmtVisitor) {
	visitor.VisitExpressionStmt(e)
}

type Function struct {
	Name   *token.Token
	Params []*token.Token
	Body   []Stmt
}

func (*Function) statement() {}
func (f *Function) Accept(visitor StmtVisitor) {
	visitor.VisitFunction(f)
}

type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (*If) statement() {}
func (i *If) Accept(visitor StmtVisitor) {
	visitor.VisitIf(i)
}

type Print struct {
	Expression Expr
}

func (*Print) statement() {}
func (p *Print) Accept(visitor StmtVisitor) {
	visitor.VisitPrint(p)
}

type Return struct {
	Keyword *token.Token
	Value   Expr
}

func (*Return) statement() {}
func (p *Return) Accept(visitor StmtVisitor) {
	visitor.VisitReturn(p)
}

type Var struct {
	Name        *token.Token
	Initializer Expr
}

func (*Var) statement() {}
func (v *Var) Accept(visitor StmtVisitor) {
	visitor.VisitVar(v)
}

type While struct {
	Condition Expr
	Body      Stmt
}

func (*While) statement() {}
func (w *While) Accept(visitor StmtVisitor) {
	visitor.VisitWhile(w)
}

type Block struct {
	Statements []Stmt
}

func (*Block) statement() {}
func (b *Block) Accept(visitor StmtVisitor) {
	visitor.VisitBlock(b)
}
