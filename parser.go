package main

import "fmt"

type expression interface {
	visitExpr(visitor) (any, error)
}

type visitor interface {
	visitLiteral(literal) (any, error)
	visitUnary(unary) (any, error)
	visitBinary(binary) (any, error)
}

type literal token
func (l literal) visitExpr(v visitor) (any, error) {
	return v.visitLiteral(l)
}

type unary struct {
	op token
	ex expression
}
func (u unary) visitExpr(v visitor)(any, error) {
	return v.visitUnary(u)
}

type binary struct {
	op token
	left expression
	right expression
}
func (b binary) visitExpr(v visitor)(any, error){
	return v.visitBinary(b)
}

type Parser struct {
	it *iter[token]
	Errors []error
	Expressions []expression
}

func NewParser(toks []token) *Parser {
	return &Parser{it: toIter(toks)}
}

func (p *Parser) Parse() ([]expression, []error) {
	for  {
		current, ok := p.it.current()
		if !ok {
			break
		}
		v,err := p.parseExpression()
		if err != nil {
			p.Errors = append(p.Errors, err)
			p.recover()
		} else {
			p.Expressions = append(p.Expressions, v)
		}

		if checkTokenType(current, semicolon) {
			p.it.consume()
		}
	}
	return p.Expressions, p.Errors
}

func (p *Parser) parseExpression() (expression,error) {
	return p.parseEquality()
}

func (p *Parser) parseEquality() (expression,error) {
	ex,err := p.parseComparison()
	if err != nil {
		return nil, err
	}
	for  {
		current, ok := p.it.current()
		if ok && (checkToken(current, operator, "!=") || checkToken(current, operator, "==")) {
			p.it.consume()
			e,err :=p.parseComparison()
			if err != nil {
				return nil, err
			}
			ex = binary{op: current, left: ex, right: e}
		} else {
			break
		}
	}
	return ex, nil
}

func (p *Parser) parseComparison() (expression,error) {
	ex,err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	for  {
		current, ok := p.it.current()
		if ok && (checkToken(current, operator, ">") || 
			checkToken(current, operator, ">=") || 
			checkToken(current, operator, "<") ||
			checkToken(current, operator, "<=")) {
			p.it.consume()
			e, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			ex = binary{op: current, left: ex, right: e}
		} else {
			break
		}
	}
	return ex,nil
}

func (p *Parser) parseTerm() (expression,error) {
	ex,err := p.parseFactor()
	if err != nil {
		return nil, err
	}
	for  {
		current, ok := p.it.current()
		if ok && (checkToken(current, operator, "-") || checkToken(current, operator, "+")) {
			p.it.consume()
			e, err := p.parseFactor()
			if err != nil {
				return nil, err
			}
			ex = binary{op: current, left: ex, right: e}
		} else {
			break
		}
	}
	return ex,nil
}

func (p *Parser) parseFactor() (expression,error) {
	ex,err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	for  {
		current, ok := p.it.current()
		if ok && (checkToken(current, operator, "/") || checkToken(current, operator, "*")) {
			p.it.consume()
			e, err := p.parseUnary()
			if err != nil {
				return nil, err
			}
			ex = binary{op: current, left: ex, right: e}
		} else {
			break
		}
	}
	return ex,nil
}

func (p *Parser) parseUnary() (expression,error) {
	current, ok := p.it.current()
	if ok && (checkToken(current, operator, "!") || checkToken(current, operator, "-")) {
		op := current
		p.it.consume()
		e, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return unary{op: op, ex: e}, nil
	} 
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() (expression,error) {
	current, ok := p.it.current()
	if !ok {
		return nil, eofError()
	} else if checkToken(current, opening, "(") {
		p.it.consume()
		ex,err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		next, ok := p.it.current()
		if ok && checkToken(next, closing, ")") {
			p.it.consume()
			return ex, nil
		} else {
			if !ok {
				return nil, eofError()
			} else {
				return nil, makeError(next, "unmatched ')'")
			}
		}
	} else if checkTokenType(current, number) || checkTokenType(current, boolean) || checkTokenType(current, stringLiteral) {
		p.it.consume()
		return literal(current), nil
	} else if checkTokenType(current, semicolon) {
		p.it.consume()
		return p.parseExpression()
	}
	return nil, makeError(current, "unexpected token when parsing primary expression")
}

func makeError(tok token, msg string) error {
	return fmt.Errorf("%v, at line %v at token %v", msg, tok.line, tok)
}

func eofError() error {
	return fmt.Errorf("unexpected end of tokens")
}

func (p *Parser) recover() {
	for {
		current, ok := p.it.current()
		if !ok || checkTokenType(current, semicolon) { // todo: more synchronisation tokens
			break
		}
		p.it.consume()
	}
}