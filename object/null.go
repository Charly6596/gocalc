package object

type Null struct{}

func (i *Null) String() string { return NULL.String() }

func (o *Null) Type() ObjectType { return NULL }
