package main

type environment map[string]loxObject

func(e environment) put(name string, obj loxObject) {
	e[name] = obj
}

func(e environment) get(name string) (loxObject, bool) {
	v, ok := e[name]
	return v, ok
}