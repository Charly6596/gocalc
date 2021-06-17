package ast

import (
	"bytes"
	"gocalc/object"
)

type NodeVisitor interface {
	Program(*Program) object.Object
	Identifier(*Identifier) object.Object
	FloatLiteral(*FloatLiteral) object.Object
	AssignmentStatement(*AssignmentStatement) object.Object
	ExpressionStatement(*ExpressionStatement) object.Object
	PrefixExpression(*PrefixExpression) object.Object
	InfixExpression(*InfixExpression) object.Object
	CallExpression(*CallExpression) object.Object
}

type Node interface {
	TokenLiteral() string
	String() string
	Accept(visit NodeVisitor) object.Object
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) Accept(visit NodeVisitor) object.Object {
	return visit.Program(p)
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}
