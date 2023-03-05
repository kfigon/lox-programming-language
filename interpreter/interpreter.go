package interpreter

import (
	"fmt"
	"lox/lexer"
	"lox/parser"
	"strconv"
)

type Interpreter struct {
	env *environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		env: newEnv(),
	}
}

func Interpret(stms []parser.Statement) error {
	i := NewInterpreter()
	for _, stmt := range stms {
		err := stmt.AcceptStatement(i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *Interpreter) VisitStatementExpression(s parser.StatementExpression) error {
	v, err := s.Expression.AcceptExpr(i)
	_ = v // todo
	return err
}

func (i *Interpreter) VisitLetStatement(let parser.LetStatement) error {
	return i.doAssignment(let.AssignmentStatement, func(name string, lo LoxObject) error {
		i.env.create(name, lo)
		return nil
	})
}

func (i *Interpreter) VisitAssignmentStatement(assign parser.AssignmentStatement) error {
	return i.doAssignment(assign, func(name string, lo LoxObject) error {
		return i.env.put(name, lo)
	})
}

func (i *Interpreter) doAssignment(assign parser.AssignmentStatement, do func(string, LoxObject) error) error {
	v, err := assign.Expression.AcceptExpr(i)
	if err != nil {
		return err
	}

	if boolExp, ok := canCast[bool](&v); ok {
		return do(assign.Name, toLoxObj(boolExp))
	} else if intExp, ok := canCast[int](&v); ok {
		return do(assign.Name, toLoxObj(intExp))
	} else if strExp, ok := canCast[string](&v); ok {
		return do(assign.Name, toLoxObj(strExp))
	}

	return fmt.Errorf("unknown type of variable %v", assign.Name)
}

func (i *Interpreter) VisitLiteral(li parser.Literal) (any, error) {
	tok := lexer.Token(li)
	if lexer.CheckTokenType(tok, lexer.Number) {
		v, err := strconv.Atoi(li.Lexeme)
		if err != nil {
			return nil, fmt.Errorf("invalid number %v, line %v, error: %w", li, li.Line, err)
		}
		return toLoxObj(v), nil
	} else if lexer.CheckTokenType(tok, lexer.StringLiteral) {
		return toLoxObj(li.Lexeme), nil
	} else if lexer.CheckTokenType(tok, lexer.Boolean) {
		v, err := strconv.ParseBool(li.Lexeme)
		if err != nil {
			return nil, fmt.Errorf("invalid boolean %v, line %v, error: %w", li, li.Line, err)
		}
		return toLoxObj(v), nil
	} else if lexer.CheckTokenType(tok, lexer.Identifier) {
		v, ok := i.env.get(tok.Lexeme)
		if !ok {
			return nil, fmt.Errorf("unknown variable %v, line %v", tok.Lexeme, li.Line)
		}
		return v, nil
	}
	return nil, fmt.Errorf("invalid literal %v, line %v", li, li.Line)
}

func (i *Interpreter) VisitUnary(u parser.Unary) (any, error) {
	op := u.Op.Lexeme

	exp, err := u.Ex.AcceptExpr(i)
	if err != nil {
		return nil, err
	}

	if op == "!" {
		v, err := castTo[bool](u.Op, &exp)
		if err != nil {
			return nil, err
		}
		return toLoxObj(!v), nil
	} else if op == "-" {
		v, err := castTo[int](u.Op, &exp)
		if err != nil {
			return nil, err
		}
		return toLoxObj(-v), nil
	}
	return nil, fmt.Errorf("invalid unary operator %v, line %v", u.Op, u.Op.Line)
}

func (i *Interpreter) VisitBinary(b parser.Binary) (any, error) {
	leftV, leftErr := b.Left.AcceptExpr(i)
	if leftErr != nil {
		return nil, leftErr
	}

	// short curcuit optimisation
	if b.Op.Lexeme == "||" || b.Op.Lexeme == "&&" {
		leftBool, leftErr := castTo[bool](b.Op, &leftV)
		if leftErr == nil {
			if b.Op.Lexeme == "||" && leftBool {
				return toLoxObj(true), nil
			} else if b.Op.Lexeme == "&&" && !leftBool {
				return toLoxObj(false), nil
			}
		}
	}

	rightV, rightErr := b.Right.AcceptExpr(i)
 	if rightErr != nil {
		return nil, rightErr
	}

	leftBool, leftErr := castTo[bool](b.Op, &leftV)
	rightBool, rightErr := castTo[bool](b.Op, &rightV)
	if leftErr == nil && rightErr == nil {
		switch b.Op.Lexeme {
		case "!=":
			return toLoxObj(leftBool != rightBool), nil
		case "==":
			return toLoxObj(leftBool == rightBool), nil
		case "||":
			return toLoxObj(leftBool || rightBool), nil
		case "&&": 
			return toLoxObj(leftBool && rightBool), nil
		}
		return nil, fmt.Errorf("unsupported binary operator boolean strings %v, line %v", b.Op, b.Op.Line)
	}

	leftStr, leftErr := castTo[string](b.Op, &leftV)
	rightStr, rightErr := castTo[string](b.Op, &rightV)
	if leftErr == nil && rightErr == nil {
		switch b.Op.Lexeme {
		case "+":
			return toLoxObj(leftStr + rightStr), nil
		case "==":
			return toLoxObj(leftStr == rightStr), nil
		case "!=":
			return toLoxObj(leftStr != rightStr), nil
		}
		return nil, fmt.Errorf("unsupported binary operator on strings %v, line %v", b.Op, b.Op.Line)
	}

	leftI, leftErr := castTo[int](b.Op, &leftV)
	rightI, rightErr := castTo[int](b.Op, &rightV)
	if leftErr == nil && rightErr == nil {
		switch b.Op.Lexeme {
		case "+":
			return toLoxObj(leftI + rightI), nil
		case "-":
			return toLoxObj(leftI - rightI), nil
		case "*":
			return toLoxObj(leftI * rightI), nil
		case "/":
			return toLoxObj(leftI / rightI), nil
		case "%":
			return toLoxObj(leftI % rightI), nil
		case ">":
			return toLoxObj(leftI > rightI), nil
		case ">=":
			return toLoxObj(leftI >= rightI), nil
		case "<":
			return toLoxObj(leftI < rightI), nil
		case "<=":
			return toLoxObj(leftI <= rightI), nil
		case "!=":
			return toLoxObj(leftI != rightI), nil
		case "==":
			return toLoxObj(leftI == rightI), nil
		}
		return nil, fmt.Errorf("unsupported binary operator on int %v, line %v", b.Op, b.Op.Line)
	}
	return nil, fmt.Errorf("unsupported binary operator, unknown type %v, line %v", b.Op, b.Op.Line)
}

func (i *Interpreter) VisitBlockStatement(b parser.BlockStatement) error {
	previous := i.env
	defer func (){
		i.env = previous
	}()

	scopeEnv := newEnv()
	scopeEnv.enclosing = previous
	i.env = scopeEnv

	for _, s := range b.Stmts {
		if err := s.AcceptStatement(i); err != nil {
			return err
		}
	}

	return nil
}

func (i *Interpreter) VisitIfStatement(ifStmt parser.IfStatement) error {
	for _, ifEl := range ifStmt.Ifs {
		v, err := ifEl.Predicate.AcceptExpr(i)
		if err != nil {
			return fmt.Errorf("error during evaluating if predicate: %w", err)
		}

		boolExp, ok := canCast[bool](&v)
		if !ok {
			return fmt.Errorf("non boolean expression in if statement")
		} else if boolExp {
			return ifEl.Body.AcceptStatement(i)
		}
	}
	return nil
}

func (i *Interpreter) VisitWhileStatement(whileStmt parser.WhileStatement) error {
	for {
		v, err := whileStmt.Predicate.AcceptExpr(i)
		if err != nil {
			return fmt.Errorf("error during evaluating while predicate: %w", err)
		}

		boolExp, ok := canCast[bool](&v)
		if !ok {
			return fmt.Errorf("non boolean expression in while statement")
		}

		if !boolExp {
			break
		}
		if err = whileStmt.Body.AcceptStatement(i); err != nil {
			return fmt.Errorf("error during processing while block: %w", err)
		}
	}
	return nil
}

func (i *Interpreter) VisitFunctionCall(fn parser.FunctionCall) (any, error) {
	return nil,nil
}
