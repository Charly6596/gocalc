package object

type Boolean struct {
	Value bool
}

func (b *Boolean) String() string {
	if b.Value {
		return "True"
	}

	return "False"
}

func (b *Boolean) Type() ObjectType { return BOOLEAN }

func (b *Boolean) TypeS() string { return b.Type().Stringf(b.String()) }
