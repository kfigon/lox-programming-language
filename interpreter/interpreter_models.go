package interpreter

import (
	"fmt"
	"lox/lexer"
	"lox/parser"
)

type LoxObject struct {
	v *any
}

func toLoxObj(v any) LoxObject {
	return LoxObject{v: &v}
}

type LoxFunction struct {
	body parser.BlockStatement
	args []string
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
	return getFromLoxObj[T](loxObj)
}

func getFromLoxObj[T any](loxObj LoxObject) (T, bool) {
	val, ok := (*loxObj.v).(T)
	return val, ok
}