package ast

import (
	"fmt"
	"gocalc/object"
	"gocalc/token"
)

type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()                        {}
func (bl *BooleanLiteral) TokenLiteral() string                   { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string                         { return fmt.Sprint(bl.Value) }
func (bl *BooleanLiteral) Accept(visit NodeVisitor) object.Object { return visit.BooleanLiteral(bl) }
