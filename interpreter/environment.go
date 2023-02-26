package interpreter

import "fmt"

type environment struct {
	d map[string]LoxObject
	enclosing *environment
}

func newEnv() *environment {
	return &environment{
		d: map[string]LoxObject{},
	}
}

func (e *environment) put(name string, obj LoxObject) error {
	_, ok := e.d[name]
	if ok {
		e.d[name] = obj
		return nil
	} else if e.enclosing != nil {
		e.enclosing.put(name, obj)
		return nil
	}
	return fmt.Errorf("undeclared variable %v", name)
}

func (e *environment) create(name string, obj LoxObject) {
	e.d[name] = obj
}

func (e *environment) get(name string) (LoxObject, bool) {
	v, ok := e.d[name]
	if !ok && e.enclosing != nil {
		return e.enclosing.get(name)
	}
	return v, ok
}
