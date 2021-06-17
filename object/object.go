package object

type ObjectType byte

type Object interface {
	Type() ObjectType
	String() string
}

const (
	NULL ObjectType = iota
	INTEGER
	FLOAT
	ERROR
)

var typeNames = []string{
	"Null",
	"Int",
	"Float",
	"Error",
}

func (o ObjectType) String() string { return typeNames[o] }
