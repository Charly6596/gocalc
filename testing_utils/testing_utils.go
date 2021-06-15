package testingutils

import (
	"reflect"
	"testing"
)

// assert fails the test if the condition is false.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		tb.Errorf(msg, v...)
		tb.FailNow()
	}
}

func Equals(tb testing.TB, exp, act interface{}, msg string) {
	if !reflect.DeepEqual(exp, act) {
		tb.Errorf("\n\n"+msg+"\n\nExpected: %#v\n\nGot: %#v", exp, act)
		tb.FailNow()
	}
}
