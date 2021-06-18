package object

type String struct {
	Value string
}

func (s *String) Type() ObjectType   { return STRING }
func (s *String) TypeS() string      { return s.Type().Stringf(s.String()) }
func (s *String) String() string     { return s.Value }
func NewString(value string) *String { return &String{Value: value} }
