package main

import (
	"fmt"
	"strconv"
)

type loxObject struct {
	v *any
}

type interpreter struct{
	env environment
}

func newInterpreter() *interpreter {
	return &interpreter{
		env: environment{},
	}
}

func interpret(stms []statement) error {
	i := newInterpreter()
	for _, stmt := range stms {
		err := stmt.acceptStatement(i)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *interpreter) visitStatementExpression(s statementExpression) error {
	v, err := s.expression.acceptExpr(i)
	_ = v // todo
	return err
}

func (i *interpreter) visitLetStatement(let letStatement) error {
	return i.visitAssignmentStatement(let.assignmentStatement)
}

func (i *interpreter) visitAssignmentStatement(assign assignmentStatement) error {
	v, err := assign.expression.acceptExpr(i)
	if err != nil {
		return err
	}

	do := func(lo loxObject) error {
		i.env[assign.name] = lo
		return nil
	}

	if boolExp, ok := canCast[bool](&v); ok {
		return do(toLoxObj(boolExp))
	} else if intExp, ok := canCast[int](&v); ok {
		return do(toLoxObj(intExp))
	} else if strExp, ok := canCast[string](&v); ok {
		return do(toLoxObj(strExp))
	}

	return fmt.Errorf("unknown type of variable %v", assign.name)
}

func (i *interpreter) visitLiteral(li literal) (any, error) {
	tok := token(li)
	if checkTokenType(tok, number) {
		v, err := strconv.Atoi(li.lexeme)
		if err != nil {
			return nil, fmt.Errorf("invalid number %v, line %v, error: %w", li, li.line, err)
		}
		return toLoxObj(v), nil
	} else if checkTokenType(tok, stringLiteral) {
		return toLoxObj(li.lexeme), nil
	} else if checkTokenType(tok, boolean) {
		v, err := strconv.ParseBool(li.lexeme)
		if err != nil {
			return nil, fmt.Errorf("invalid boolean %v, line %v, error: %w", li, li.line, err)
		}
		return toLoxObj(v), nil
	} else if checkTokenType(tok, identifier) {
		v, ok := i.env[tok.lexeme]
		if !ok {
			return nil, fmt.Errorf("unknown variable %v, line %v", tok.lexeme, li.line)
		}
		return v, nil
	}
	return nil, fmt.Errorf("invalid literal %v, line %v", li, li.line)
}

func (i *interpreter) visitUnary(u unary) (any, error) {
	op := u.op.lexeme

	exp, err := u.ex.acceptExpr(i)
	if err != nil {
		return nil, err
	}

	if op == "!" {
		v, err := castTo[bool](u.op, &exp)
		if err != nil {
			return nil, err
		}
		return toLoxObj(!v), nil
	} else if op == "-" {
		v, err := castTo[int](u.op, &exp)
		if err != nil {
			return nil, err
		}
		return toLoxObj(-v), nil
	}
	return nil, fmt.Errorf("invalid unary operator %v, line %v", u.op, u.op.line)
}

func (i *interpreter) visitBinary(b binary) (any, error) {
	leftV, leftErr := b.left.acceptExpr(i)
	rightV, rightErr := b.right.acceptExpr(i)

	if leftErr != nil {
		return nil, leftErr
	} else if rightErr != nil {
		return nil, rightErr
	}

	leftBool, leftErr := castTo[bool](b.op, &leftV)
	rightBool, rightErr := castTo[bool](b.op, &rightV)
	if leftErr == nil && rightErr == nil {
		switch b.op.lexeme {
		case "!=":
			return toLoxObj(leftBool != rightBool), nil
		case "==":
			return toLoxObj(leftBool == rightBool), nil
		}
		return nil, fmt.Errorf("unsupported binary operator boolean strings %v, line %v", b.op, b.op.line)
	}

	leftStr, leftErr := castTo[string](b.op, &leftV)
	rightStr, rightErr := castTo[string](b.op, &rightV)
	if leftErr == nil && rightErr == nil {
		switch b.op.lexeme {
		case "+":
			return toLoxObj(leftStr + rightStr), nil
		case "==":
			return toLoxObj(leftStr == rightStr), nil
		case "!=":
			return toLoxObj(leftStr != rightStr), nil
		}
		return nil, fmt.Errorf("unsupported binary operator on strings %v, line %v", b.op, b.op.line)
	}

	leftI, leftErr := castTo[int](b.op, &leftV)
	rightI, rightErr := castTo[int](b.op, &rightV)
	if leftErr == nil && rightErr == nil {
		switch b.op.lexeme {
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
		return nil, fmt.Errorf("unsupported binary operator on int %v, line %v", b.op, b.op.line)
	}
	return nil, fmt.Errorf("unsupported binary operator, unknown type %v, line %v", b.op, b.op.line)
}

func toLoxObj(v any) loxObject {
	return loxObject{v: &v}
}

func castTo[T any](t token, v *any) (T, error) {
	val, ok := canCast[T](v)
	if !ok {
		return val, fmt.Errorf("invalid lox type: %v value not found %v, line %v", t.tokType, t, t.line)
	}
	return val, nil
}

func canCast[T any](v *any) (T, bool) {
	loxObj, ok := (*v).(loxObject)
	if !ok {
		var out T
		return out, false
	}
	val, ok := (*loxObj.v).(T)
	return val, ok
}
