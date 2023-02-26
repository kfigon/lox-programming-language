package parser

import "lox/lexer"

type Expression interface {
	AcceptExpr(VisitorExpr) (any, error)
}

type VisitorExpr interface {
	VisitLiteral(Literal) (any, error)
	VisitUnary(Unary) (any, error)
	VisitBinary(Binary) (any, error)
}

type Literal lexer.Token

func (l Literal) AcceptExpr(v VisitorExpr) (any, error) {
	return v.VisitLiteral(l)
}

type Unary struct {
	Op lexer.Token
	Ex Expression
}

func (u Unary) AcceptExpr(v VisitorExpr) (any, error) {
	return v.VisitUnary(u)
}

type Binary struct {
	Op    lexer.Token
	Left  Expression
	Right Expression
}

func (b Binary) AcceptExpr(v VisitorExpr) (any, error) {
	return v.VisitBinary(b)
}

type Statement interface {
	AcceptStatement(VisitorStatement) error
}

type VisitorStatement interface {
	VisitStatementExpression(StatementExpression) error
	VisitLetStatement(LetStatement) error
	VisitAssignmentStatement(AssignmentStatement) error
	VisitBlockStatement(BlockStatement) error
}

type StatementExpression struct {
	Expression
}

func (s StatementExpression) AcceptStatement(v VisitorStatement) error {
	return v.VisitStatementExpression(s)
}

type LetStatement struct {
	AssignmentStatement
}

func (s LetStatement) AcceptStatement(v VisitorStatement) error {
	return v.VisitLetStatement(s)
}

type AssignmentStatement struct {
	Name string
	Expression
}

func (a AssignmentStatement) AcceptStatement(v VisitorStatement) error {
	return v.VisitAssignmentStatement(a)
}


type BlockStatement struct {
	Stmts []Statement
}

func (b BlockStatement) AcceptStatement(v VisitorStatement) error {
	return v.VisitBlockStatement(b)
}
