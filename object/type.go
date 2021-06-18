package object

type Type struct {
	Value ObjectType
}

func (t *Type) Type() ObjectType { return TYPE }
func (t *Type) String() string   { return t.Value.String() }
func (t *Type) TypeS() string    { return t.Type().Stringf(t.Value.String()) }
