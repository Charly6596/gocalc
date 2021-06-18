package environment

import (
	"bytes"
	"gocalc/object"
)

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

func (e *Environment) String() string {
	var buff bytes.Buffer

	buff.WriteString("Environment inspection")
	for name, value := range e.store {
		buff.WriteString("\n\t")
		buff.WriteString(name)
		buff.WriteString(": ")
		buff.WriteString(value.String())
		buff.WriteString(" <")
		buff.WriteString(value.Type().String())
		buff.WriteString(">")
	}
	return buff.String()
}
