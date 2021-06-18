package object

import (
	"reflect"
)

type nativeFn func(objs ...Object) Object

type NativeFunction struct {
	Name     string
	Function nativeFn
}

func (nf *NativeFunction) String() string {
	return reflect.ValueOf(nf.Function).String()
}
func (nf *NativeFunction) Type() ObjectType     { return NATIVE_FUNCTION }
func (nf *NativeFunction) TypeS() string        { return nf.Type().Stringf(nf.Name) }
func (nf *NativeFunction) Is(t ObjectType) bool { return nf.Type() == t }
