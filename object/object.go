package object

import (
	"fmt"
	"strings"
)

type ObjectType byte

type Object interface {
	Type() ObjectType
	TypeS() string
	String() string
}

const (
	NULL ObjectType = iota
	INTEGER
	FLOAT
	NATIVE_FUNCTION
	ERROR
	STRING
	TYPE
)

var typeNames = []string{
	NULL:            "Nil",
	INTEGER:         "Int",
	FLOAT:           "Float",
	ERROR:           "Err",
	STRING:          "Str",
	TYPE:            "Type",
	NATIVE_FUNCTION: "NativeFn",
}

func (o ObjectType) String() string { return typeNames[o] }
func (o ObjectType) Stringf(params ...string) string {
	if len(params) == 0 {
		return fmt.Sprintf("<%s>", o)
	}

	p := strings.Join(params, ", ")
	return fmt.Sprintf("<%s: %s>", o, p)
}
