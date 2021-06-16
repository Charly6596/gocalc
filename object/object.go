package object

type ObjectType byte
type Object interface {
	Type() ObjectType
	String() string
}

const (
	NULL ObjectType = iota
	INTEGER
	ERROR
)

var typeNames = []string{
	"Null",
	"Int",
	"Error",
}

func (o ObjectType) String() string { return typeNames[o] }
