package ast

import (
	"fmt"
	"strings"
)

var (
	_ Visitor = (*AstPrinter)(nil)
)

type AstPrinter struct{}

func (p *AstPrinter) Print(expr Expr) string {
	return expr.Accept(p).(string)
}

func (p *AstPrinter) VisitBinary(binary *Binary) interface{} {
	return p.parenthesize(binary.Operator.Lexeme, binary.Left, binary.Right)
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

func (p *AstPrinter) VisitUnary(unary *Unary) interface{} {
	return p.parenthesize(unary.Operator.Lexeme, unary.Right)
}

func (p *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	exprRepresentations := []string{}
	for _, expr := range exprs {
		exprRepresentations = append(exprRepresentations, expr.Accept(p).(string))
	}

	return fmt.Sprintf("(%s %s)", name, strings.Join(exprRepresentations, " "))
}
