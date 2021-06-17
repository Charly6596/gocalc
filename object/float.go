package object

import (
	"fmt"
)

type Float struct {
	Value float64
}

func (i *Float) String() string {
	return fmt.Sprint(i.Value)
}

func (i *Float) Type() ObjectType {
	return FLOAT
}

func ToFloat(o Object) (*Float, bool) {
	inte, ok := o.(*Float)
	return inte, ok && o.Type() == FLOAT
}
