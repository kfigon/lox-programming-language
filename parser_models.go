package main

type expression interface {
	acceptExpr(visitorExpr) (any, error)
}

type visitorExpr interface {
	visitLiteral(literal) (any, error)
	visitUnary(unary) (any, error)
	visitBinary(binary) (any, error)
}

type literal token

func (l literal) acceptExpr(v visitorExpr) (any, error) {
	return v.visitLiteral(l)
}

type unary struct {
	op token
	ex expression
}

func (u unary) acceptExpr(v visitorExpr) (any, error) {
	return v.visitUnary(u)
}

type binary struct {
	op    token
	left  expression
	right expression
}

func (b binary) acceptExpr(v visitorExpr) (any, error) {
	return v.visitBinary(b)
}

type statement interface {
	acceptStatement(visitorStatement) error
}

type visitorStatement interface {
	visitStatementExpression(statementExpression) error
	visitLetStatement(letStatement) error
	visitAssignmentStatement(assignmentStatement) error
}

type statementExpression struct {
	expression
}

func (s statementExpression) acceptStatement(v visitorStatement) error {
	return v.visitStatementExpression(s)
}

type letStatement struct {
	assignmentStatement
}

func (s letStatement) acceptStatement(v visitorStatement) error {
	return v.visitLetStatement(s)
}

type assignmentStatement struct {
	name string
	expression
}

func (a assignmentStatement) acceptStatement(v visitorStatement) error {
	return v.visitAssignmentStatement(a)
}