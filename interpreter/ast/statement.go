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
	VisitExpressionStmt(stmt *ExpressionStmt)
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
