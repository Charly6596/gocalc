package ast

import (
	"bytes"
	"gocalc/object"
	"gocalc/token"
)

type ListLiteral struct {
	Token  token.Token // token.ASSIGNMENT
	Values []Expression
}

func (ll *ListLiteral) expressionNode()      {}
func (ll *ListLiteral) TokenLiteral() string { return ll.Token.Literal }

func (ll *ListLiteral) Accept(visit NodeVisitor) object.Object {
	return visit.ListLiteral(ll)
}

func (ll *ListLiteral) String() string {
	var out bytes.Buffer
	out.WriteString("[")

	for i, exp := range ll.Values {
		out.WriteString(exp.String())
		if i != len(ll.Values) {
			out.WriteString(", ")
		}
	}

	out.WriteString("]")
	return out.String()
}
