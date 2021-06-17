package ast

import (
	"fmt"
	"gocalc/object"
	"gocalc/token"
)

type FloatLiteral struct {
	Token token.Token
	Value float64
}

func (fl *FloatLiteral) expressionNode()                        {}
func (fl *FloatLiteral) TokenLiteral() string                   { return fl.Token.Literal }
func (fl *FloatLiteral) String() string                         { return fmt.Sprint(fl.Value) }
func (fl *FloatLiteral) Accept(visit NodeVisitor) object.Object { return visit.FloatLiteral(fl) }
