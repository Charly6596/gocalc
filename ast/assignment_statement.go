package ast

import (
	"bytes"
	"gocalc/token"
)

type AssignmentStatement struct {
	Token token.Token // token.ASSIGNMENT
	Name  *Identifier
	Value Expression
}

func (ae *AssignmentStatement) statementNode()       {}
func (ae *AssignmentStatement) TokenLiteral() string { return ae.Token.Literal }
func (ae *AssignmentStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ae.Name.String())
	out.WriteString(" = ")
	if ae.Value != nil {
		out.WriteString(ae.Value.String())
	}
	out.WriteString(";")
	return out.String()
}
