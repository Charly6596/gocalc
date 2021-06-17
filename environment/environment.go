package environment

import "gocalc/object"

type Environment struct {
	store map[string]object.Object
}

func New() *Environment {
	s := make(map[string]object.Object)
	return &Environment{s}
}

func (e *Environment) Get(ident string) (object.Object, bool) {
	res, ok := e.store[ident]
	return res, ok
}

func (e *Environment) Set(name string, obj object.Object) object.Object {
	e.store[name] = obj
	return obj
}
