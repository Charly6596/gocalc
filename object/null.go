package object

type Null struct{}

func (n *Null) String() string { return NULL.String() }

func (n *Null) Type() ObjectType { return NULL }
