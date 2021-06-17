package ast

import (
	"gocalc/object"
	"gocalc/token"
)

type Identifier struct {
	Token token.Token // token.IDENT
	Value string
}

func (i *Identifier) expressionNode()                        {}
func (i *Identifier) TokenLiteral() string                   { return i.Token.Literal }
func (i *Identifier) String() string                         { return i.Value }
func (i *Identifier) Accept(visit NodeVisitor) object.Object { return visit.Identifier(i) }
