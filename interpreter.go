package main

import (
	"fmt"
	"strconv"
)

type loxObject struct {
	v *any
}

type interpreter struct{}

func interpret(expr []expression) ([]loxObject, error) {
	var out []loxObject
	i := &interpreter{}
	for _, e := range expr {
		v, err := e.visitExpr(i)
		if err != nil {
			return nil, err
		} else {
			out = append(out, v.(loxObject))
		}
	}
	return out, nil
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
	}
	return nil, fmt.Errorf("invalid literal %v, line %v", li, li.line)
}

func (i *interpreter) visitUnary(u unary) (any, error) {
	op := u.op.lexeme

	exp, err := u.ex.visitExpr(&interpreter{})
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
	leftV, leftErr := b.left.visitExpr(&interpreter{})
	rightV, rightErr := b.right.visitExpr(&interpreter{})

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
	loxObj, ok := (*v).(loxObject)
	if !ok {
		var out T
		return out, fmt.Errorf("invalid lox type: %v value not found %v, line %v", t.tokType, t, t.line)
	}
	val, ok := (*loxObj.v).(T)
	if !ok {
		return val, fmt.Errorf("%v value not found %v, line %v", t.tokType, t, t.line)
	}
	return val, nil
}
