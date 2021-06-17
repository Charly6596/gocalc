package ast

import (
	"bytes"
	"gocalc/object"
	"gocalc/token"
)

type PrefixExpression struct {
	Token    token.Token // the first token of the expression
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }

func (pe *PrefixExpression) Accept(visit NodeVisitor) object.Object {
	return visit.PrefixExpression(pe)
}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}
