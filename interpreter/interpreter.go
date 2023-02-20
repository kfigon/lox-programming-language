package interpreter

import (
	"fmt"
	"lox/lexer"
	"lox/parser"
	"strconv"
)

type LoxObject struct {
	v *any
}

type Interpreter struct {
	env environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		env: environment{},
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
	return i.doAssignment(let.AssignmentStatement)
}

func (i *Interpreter) VisitAssignmentStatement(assign parser.AssignmentStatement) error {
	if _, ok := i.env.get(assign.Name); !ok {
		return fmt.Errorf("unknown variable %v found during assignment", assign.Name)
	}
	return i.doAssignment(assign)
}

func (i *Interpreter) doAssignment(assign parser.AssignmentStatement) error {
	v, err := assign.Expression.AcceptExpr(i)
	if err != nil {
		return err
	}

	do := func(lo LoxObject) error {
		i.env.put(assign.Name, lo)
		return nil
	}

	if boolExp, ok := canCast[bool](&v); ok {
		return do(toLoxObj(boolExp))
	} else if intExp, ok := canCast[int](&v); ok {
		return do(toLoxObj(intExp))
	} else if strExp, ok := canCast[string](&v); ok {
		return do(toLoxObj(strExp))
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
	rightV, rightErr := b.Right.AcceptExpr(i)

	if leftErr != nil {
		return nil, leftErr
	} else if rightErr != nil {
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

func toLoxObj(v any) LoxObject {
	return LoxObject{v: &v}
}

func castTo[T any](t lexer.Token, v *any) (T, error) {
	val, ok := canCast[T](v)
	if !ok {
		return val, fmt.Errorf("invalid lox type: %v value not found %v, line %v", t.TokType, t, t.Line)
	}
	return val, nil
}

func canCast[T any](v *any) (T, bool) {
	loxObj, ok := (*v).(LoxObject)
	if !ok {
		var out T
		return out, false
	}
	val, ok := (*loxObj.v).(T)
	return val, ok
}
