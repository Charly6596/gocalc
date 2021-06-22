package object

import "bytes"

type List struct {
	Values []Object
}

func (l *List) Type() ObjectType { return LIST }
func (l *List) TypeS() string {
	return LIST.Stringf(l.String())
}
func (l *List) String() string {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i, obj := range l.Values {
		buf.WriteString(obj.String())
		if i != len(l.Values)-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString("]")
	return buf.String()
}
