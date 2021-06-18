package evaluator

import (
	"gocalc/object"
	"reflect"
)

type NativeFn func(*Evaluator, ...object.Object) object.Object

type NativeFunction struct {
	Name     string
	Function NativeFn
}

func (nf *NativeFunction) String() string {
	return reflect.ValueOf(nf.Function).String()
}
func (nf *NativeFunction) Type() object.ObjectType     { return object.NATIVE_FUNCTION }
func (nf *NativeFunction) TypeS() string               { return nf.Type().Stringf(nf.Name) }
func (nf *NativeFunction) Is(t object.ObjectType) bool { return nf.Type() == t }
