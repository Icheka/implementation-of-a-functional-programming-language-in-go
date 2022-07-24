package object

import (
	"bytes"
	"fmt"
	"monkey/ast"
	"monkey/token"
	"strings"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

const (
	INTEGER_OBJECT      = "INTEGER"
	BOOLEAN_OBJECT      = "BOOLEAN"
	NULL_OBJECT         = "NULL"
	RETURN_VALUE_OBJECT = "RETURN_VALUE"
	ERROR_OBJECT        = "ERROR"
	FUNCTION_OBJECT     = "FUNCTION"
)

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJECT }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJECT }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct {
	Value int64
}

func (n *Null) Type() ObjectType { return NULL_OBJECT }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJECT }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJECT }
func (e *Error) Inspect() string  { return "Uncaught syntax error!: " + e.Message }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJECT }
func (f *Function) Inspect() string {
	var str bytes.Buffer

	parameters := []string{}
	for _, p := range f.Parameters {
		parameters = append(parameters, p.String())
	}

	str.WriteString(token.FUNCTION)
	str.WriteString("(")
	str.WriteString(strings.Join(parameters, ", "))
	str.WriteString(") {\n")
	str.WriteString(f.Body.String())
	str.WriteString("\n}")

	return str.String()
}
