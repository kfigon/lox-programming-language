package interpreter

import (
	"fmt"
	"lox/lexer"
)

type LoxObject struct {
	v *any
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