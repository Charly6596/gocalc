package object

import "fmt"

const (
	UNKNOWN_INFIX_OPERATOR_ERROR  = "Unknown operator %s %s %s"
	UNKNOWN_PREFIX_OPERATOR_ERROR = "Unknown operator %s%s"
	IDENTIFIER_NOT_FOUND_ERROR    = "Identifier not found %s"
	SYNTAX_ERROR                  = "Syntax error: \n\t\t%s"
)

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR }
func (e *Error) String() string   { return fmt.Sprintf("%s", e.Message) }
