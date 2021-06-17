package ast

import (
	"gocalc/object"
	"gocalc/token"
)

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

func (es *ExpressionStatement) Accept(visit NodeVisitor) object.Object {
	return visit.ExpressionStatement(&ExpressionStatement{Token: es.Token, Expression: es.Expression})
}
