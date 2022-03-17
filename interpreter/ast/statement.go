package ast

import "github.com/DanielleB-R/golox/interpreter/token"

var (
	_ Stmt = (*ExpressionStmt)(nil)
	_ Stmt = (*Print)(nil)
	_ Stmt = (*Var)(nil)
)

type Stmt interface {
	statement()
	Accept(visitor StmtVisitor)
}

type StmtVisitor interface {
	VisitBlock(stmt *Block)
	VisitExpressionStmt(stmt *ExpressionStmt)
	VisitIf(stmt *If)
	VisitPrint(stmt *Print)
	VisitVar(stmt *Var)
}

type ExpressionStmt struct {
	Expression Expr
}

func (*ExpressionStmt) statement() {}
func (e *ExpressionStmt) Accept(visitor StmtVisitor) {
	visitor.VisitExpressionStmt(e)
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

type Var struct {
	Name        *token.Token
	Initializer Expr
}

func (*Var) statement() {}
func (v *Var) Accept(visitor StmtVisitor) {
	visitor.VisitVar(v)
}

type Block struct {
	Statements []Stmt
}

func (*Block) statement() {}
func (b *Block) Accept(visitor StmtVisitor) {
	visitor.VisitBlock(b)
}
