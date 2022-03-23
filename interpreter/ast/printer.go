package ast

import (
	"fmt"
	"strings"
)

var (
	_ ExprVisitor = (*AstPrinter)(nil)
)

type AstPrinter struct{}

func (p *AstPrinter) Print(expr Expr) string {
	return expr.Accept(p).(string)
}

func (p *AstPrinter) VisitAssign(assign *Assign) interface{} {
	return p.parenthesize("=", &Variable{Name: assign.Name}, assign.Value)
}

func (p *AstPrinter) VisitBinary(binary *Binary) interface{} {
	return p.parenthesize(binary.Operator.Lexeme, binary.Left, binary.Right)
}

func (p *AstPrinter) VisitCall(call *Call) interface{} {
	callee := p.Print(call.Callee)
	return p.parenthesize(callee, call.Arguments...)
}

func (p *AstPrinter) VisitGet(get *Get) any {
	return p.parenthesize("."+get.Name.Lexeme, get.Object)
}

func (p *AstPrinter) VisitGrouping(grouping *Grouping) interface{} {
	return p.parenthesize("group", grouping.Expression)
}

func (p *AstPrinter) VisitLiteral(literal *Literal) interface{} {
	if literal.Value == nil {
		return "nil"
	}
	return fmt.Sprint(literal.Value)
}

func (p *AstPrinter) VisitLogical(logical *Logical) interface{} {
	return p.parenthesize(logical.Operator.Lexeme, logical.Left, logical.Right)
}

func (p *AstPrinter) VisitSet(set *Set) any {
	return p.parenthesize(fmt.Sprintf(".%s=", set.Name.Lexeme), set.Object, set.Value)
}

func (p *AstPrinter) VisitThis(this *This) any {
	return this.Keyword.Lexeme
}

func (p *AstPrinter) VisitUnary(unary *Unary) interface{} {
	return p.parenthesize(unary.Operator.Lexeme, unary.Right)
}

func (p *AstPrinter) VisitVariable(variable *Variable) interface{} {
	return variable.Name.Lexeme
}

func (p *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	exprRepresentations := []string{}
	for _, expr := range exprs {
		exprRepresentations = append(exprRepresentations, expr.Accept(p).(string))
	}

	return fmt.Sprintf("(%s %s)", name, strings.Join(exprRepresentations, " "))
}
