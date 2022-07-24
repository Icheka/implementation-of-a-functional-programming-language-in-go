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
	STRING_OBJECT       = "STRING"
	BUILTIN_OBJECT      = "BUILTIN"
	ARRAY_OBJECT        = "ARRAY"
)

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJECT }
func (a *Array) Inspect() string {
	var str bytes.Buffer

	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	str.WriteString("[")
	str.WriteString(strings.Join(elements, ", "))
	str.WriteString("]")

	return str.String()
}

type String struct{ Value string }

func (s *String) Type() ObjectType { return STRING_OBJECT }
func (s *String) Inspect() string  { return s.Value }

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

type BuiltinFunction func(args ...Object) Object
type Builtin struct {
	Function BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJECT }
func (b *Builtin) Inspect() string  { return "built-in function" }
