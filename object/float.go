package object

import (
	"fmt"
)

type Float struct {
	Value float64
}

func (f *Float) String() string { return fmt.Sprint(f.Value) }

func (f *Float) Type() ObjectType { return FLOAT }

func (f *Float) TypeS() string { return f.Type().Stringf(f.String()) }

func (f *Float) Is(t ObjectType) bool { return f.Type() == t }

func ToFloat(o Object) (*Float, bool) {
	inte, ok := o.(*Float)
	return inte, ok && o.Type() == FLOAT
}
