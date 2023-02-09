package main

import "fmt"

type Parser struct {
	it         *iter[token]
	Errors     []error
	statements []statement
}

func NewParser(toks []token) *Parser {
	return &Parser{it: toIter(toks)}
}

func (p *Parser) Parse() ([]statement, []error) {
	for _, ok := p.it.current(); ok; _, ok = p.it.current() {

		v, err := p.parseStatement()
		if err != nil {
			p.Errors = append(p.Errors, err)
			p.recover()
			continue
		}

		p.statements = append(p.statements, v)
	}
	return p.statements, p.Errors
}

func (p *Parser) parseTerminatedExpression() (expression, error) {
	v, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	current, ok := p.it.current()
	if ok && checkTokenType(current, semicolon) {
		p.it.consume()
		return v, nil
	}
	return nil, fmt.Errorf("unterminated statement at line %d", current.line)
}

func (p *Parser) parseStatement() (statement, error) {
	current, ok := p.it.current()
	if !ok {
		return nil, eofError()
	}

	if checkToken(current, keyword, "let") {
		return p.parseLetStatement()
	}

	v, err := p.parseTerminatedExpression()
	if err != nil {
		return nil, err
	} 
	return statementExpression{v}, nil
}

func (p *Parser) parseLetStatement() (statement, error) {
	p.it.consume() // let
	if err := p.ensureCurrentTokenType(identifier); err != nil {
		return nil, err
	}

	current, _ := p.it.current()
	name := current.lexeme
	p.it.consume() // identifier
	
	if err := p.ensureCurrentToken(operator, "="); err != nil {
		return nil, err
	}
	p.it.consume() // =

	v, err := p.parseTerminatedExpression()
	if err != nil {
		return nil, err
	} 
	return letStatement{name, v}, nil
}

func (p *Parser) parseExpression() (expression, error) {
	return p.parseEquality()
}

func (p *Parser) parseEquality() (expression, error) {
	ex, err := p.parseComparison()
	if err != nil {
		return nil, err
	}
	for {
		current, ok := p.it.current()
		if ok && (checkToken(current, operator, "!=") || checkToken(current, operator, "==")) {
			p.it.consume()
			e, err := p.parseComparison()
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

func (p *Parser) parseComparison() (expression, error) {
	ex, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	for {
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
	return ex, nil
}

func (p *Parser) parseTerm() (expression, error) {
	ex, err := p.parseFactor()
	if err != nil {
		return nil, err
	}
	for {
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
	return ex, nil
}

func (p *Parser) parseFactor() (expression, error) {
	ex, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	for {
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
	return ex, nil
}

func (p *Parser) parseUnary() (expression, error) {
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

func (p *Parser) parsePrimary() (expression, error) {
	current, ok := p.it.current()
	if !ok {
		return nil, eofError()
	} else if checkToken(current, opening, "(") {
		p.it.consume()
		ex, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		if err = p.ensureCurrentToken(closing, ")"); err != nil {
			return nil ,err
		}
		p.it.consume()
		return ex, nil
	} else if checkTokenType(current, number) || checkTokenType(current, boolean) || checkTokenType(current, stringLiteral) || checkTokenType(current, identifier) {
		p.it.consume()
		return literal(current), nil
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
	for current, ok := p.it.current(); ok; current, ok = p.it.current() {
		if checkTokenType(current, semicolon) {
			p.it.consume()
			break
		}
		p.it.consume()
	}
}

func (p *Parser) ensureCurrentToken(tokType tokenType, lexeme string) error {
	current, ok := p.it.current()
	if !ok {
		return eofError()
	} else if !checkToken(current, tokType, lexeme) {
		return makeError(current, fmt.Sprintf("expected %v", tokType))
	}
	return nil
}

func (p *Parser) ensureCurrentTokenType(tokType tokenType) error {
	current, ok := p.it.current()
	if !ok {
		return eofError()
	} else if !checkTokenType(current, tokType) {
		return makeError(current, fmt.Sprintf("expected %v", tokType))
	}
	return nil
}