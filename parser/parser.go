package parser

import (
	"fmt"
	"lox/lexer"
)

type Parser struct {
	it         *iter[lexer.Token]
	Errors     []error
	statements []Statement
}

func NewParser(toks []lexer.Token) *Parser {
	return &Parser{it: toIter(toks)}
}

func (p *Parser) Parse() ([]Statement, []error) {
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

func (p *Parser) parseTerminatedExpression() (Expression, error) {
	v, err := p.parseExpression()
	if err != nil {
		return nil, err
	}

	current, ok := p.it.current()
	if ok && lexer.CheckTokenType(current, lexer.Semicolon) {
		p.it.consume()
		return v, nil
	}
	return nil, fmt.Errorf("unterminated statement at line %d", current.Line)
}

func (p *Parser) parseStatement() (Statement, error) {
	current, ok := p.it.current()
	if !ok {
		return nil, eofError()
	}

	if lexer.CheckToken(current, lexer.Keyword, "let") {
		return p.parseLetStatement()
	} else if lexer.CheckTokenType(current, lexer.Identifier) {
		return p.parseAssignmentStatement()
	}

	v, err := p.parseTerminatedExpression()
	if err != nil {
		return nil, err
	}
	return StatementExpression{v}, nil
}

func (p *Parser) parseLetStatement() (Statement, error) {
	p.it.consume() // let
	if err := p.ensureCurrentTokenType(lexer.Identifier); err != nil {
		return nil, err
	}

	assingnment, err := p.parseAssignmentStatement()
	if err != nil {
		return nil, err
	}
	return LetStatement{AssignmentStatement: assingnment}, nil
}

func (p *Parser) parseAssignmentStatement() (AssignmentStatement, error) {
	current, _ := p.it.current()
	name := current.Lexeme
	p.it.consume() // identifier

	if err := p.ensureCurrentToken(lexer.Operator, "="); err != nil {
		return AssignmentStatement{}, err
	}
	p.it.consume() // =

	v, err := p.parseTerminatedExpression()
	if err != nil {
		return AssignmentStatement{}, err
	}
	return AssignmentStatement{name, v}, nil
}

func (p *Parser) parseExpression() (Expression, error) {
	return p.parseEquality()
}

func (p *Parser) parseEquality() (Expression, error) {
	ex, err := p.parseComparison()
	if err != nil {
		return nil, err
	}
	for {
		current, ok := p.it.current()
		if ok && (lexer.CheckToken(current, lexer.Operator, "!=") || lexer.CheckToken(current, lexer.Operator, "==")) {
			p.it.consume()
			e, err := p.parseComparison()
			if err != nil {
				return nil, err
			}
			ex = Binary{Op: current, Left: ex, Right: e}
		} else {
			break
		}
	}
	return ex, nil
}

func (p *Parser) parseComparison() (Expression, error) {
	ex, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	for {
		current, ok := p.it.current()
		if ok && (lexer.CheckToken(current, lexer.Operator, ">") ||
			lexer.CheckToken(current, lexer.Operator, ">=") ||
			lexer.CheckToken(current, lexer.Operator, "<") ||
			lexer.CheckToken(current, lexer.Operator, "<=")) {
			p.it.consume()
			e, err := p.parseTerm()
			if err != nil {
				return nil, err
			}
			ex = Binary{Op: current, Left: ex, Right: e}
		} else {
			break
		}
	}
	return ex, nil
}

func (p *Parser) parseTerm() (Expression, error) {
	ex, err := p.parseFactor()
	if err != nil {
		return nil, err
	}
	for {
		current, ok := p.it.current()
		if ok && (lexer.CheckToken(current, lexer.Operator, "-") || lexer.CheckToken(current, lexer.Operator, "+")) {
			p.it.consume()
			e, err := p.parseFactor()
			if err != nil {
				return nil, err
			}
			ex = Binary{Op: current, Left: ex, Right: e}
		} else {
			break
		}
	}
	return ex, nil
}

func (p *Parser) parseFactor() (Expression, error) {
	ex, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	for {
		current, ok := p.it.current()
		if ok && (lexer.CheckToken(current, lexer.Operator, "/") || lexer.CheckToken(current, lexer.Operator, "*")) {
			p.it.consume()
			e, err := p.parseUnary()
			if err != nil {
				return nil, err
			}
			ex = Binary{Op: current, Left: ex, Right: e}
		} else {
			break
		}
	}
	return ex, nil
}

func (p *Parser) parseUnary() (Expression, error) {
	current, ok := p.it.current()
	if ok && (lexer.CheckToken(current, lexer.Operator, "!") || lexer.CheckToken(current, lexer.Operator, "-")) {
		op := current
		p.it.consume()
		e, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return Unary{Op: op, Ex: e}, nil
	}
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() (Expression, error) {
	current, ok := p.it.current()
	if !ok {
		return nil, eofError()
	} else if lexer.CheckToken(current, lexer.Opening, "(") {
		p.it.consume()
		ex, err := p.parseExpression()
		if err != nil {
			return nil, err
		}

		if err = p.ensureCurrentToken(lexer.Closing, ")"); err != nil {
			return nil, err
		}
		p.it.consume()
		return ex, nil
	} else if lexer.CheckTokenType(current, lexer.Number) || lexer.CheckTokenType(current, lexer.Boolean) || lexer.CheckTokenType(current, lexer.StringLiteral) || lexer.CheckTokenType(current, lexer.Identifier) {
		p.it.consume()
		return Literal(current), nil
	}
	return nil, makeError(current, "unexpected token when parsing primary expression")
}

func makeError(tok lexer.Token, msg string) error {
	return fmt.Errorf("%v, at line %v at token %v", msg, tok.Line, tok)
}

func eofError() error {
	return fmt.Errorf("unexpected end of tokens")
}

func (p *Parser) recover() {
	for current, ok := p.it.current(); ok; current, ok = p.it.current() {
		if lexer.CheckTokenType(current, lexer.Semicolon) {
			p.it.consume()
			break
		}
		p.it.consume()
	}
}

func (p *Parser) ensureCurrentToken(tokType lexer.TokenType, lexeme string) error {
	current, ok := p.it.current()
	if !ok {
		return eofError()
	} else if !lexer.CheckToken(current, tokType, lexeme) {
		return makeError(current, fmt.Sprintf("expected %v", tokType))
	}
	return nil
}

func (p *Parser) ensureCurrentTokenType(tokType lexer.TokenType) error {
	current, ok := p.it.current()
	if !ok {
		return eofError()
	} else if !lexer.CheckTokenType(current, tokType) {
		return makeError(current, fmt.Sprintf("expected %v", tokType))
	}
	return nil
}
