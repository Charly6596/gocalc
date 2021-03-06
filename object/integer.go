package object

import (
	"fmt"
)

const DIVIDE_BY_ZERO = "Cannot divide by zero (%v / %v)"

type Integer struct {
	Value int64
}

func (i *Integer) String() string {
	return fmt.Sprint(i.Value)
}

func (i *Integer) Type() ObjectType {
	return INTEGER
}

func (i *Integer) TypeS() string {
	return i.Type().Stringf(i.String())
}

func ToInteger(o Object) (*Integer, bool) {
	inte, ok := o.(*Integer)
	return inte, ok && o.Type() == INTEGER
}

func DivideByZeroError(x1, x2 fmt.Stringer) *Error {
	msg := fmt.Sprintf(DIVIDE_BY_ZERO, x1, x2)
	return &Error{Message: msg}
}
