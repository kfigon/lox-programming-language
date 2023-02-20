package interpreter

type environment map[string]LoxObject

func (e environment) put(name string, obj LoxObject) {
	e[name] = obj
}

func (e environment) get(name string) (LoxObject, bool) {
	v, ok := e[name]
	return v, ok
}
