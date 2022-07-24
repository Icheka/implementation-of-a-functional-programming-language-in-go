package ast

import (
	"bytes"
	"monkey/token"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
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
	}
	return ""
}

func (p *Program) String() string {
	var str bytes.Buffer

	for _, stmt := range p.Statements {
		str.WriteString(stmt.String())
	}

	return str.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Value
}
func (i *Identifier) String() string {
	return i.Value
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Value
}
func (ls *LetStatement) String() string {
	var str bytes.Buffer

	str.WriteString(ls.TokenLiteral() + " ")
	str.WriteString(ls.Name.String() + " = ")

	if ls.Value != nil {
		str.WriteString(ls.Value.String())
	}

	str.WriteString(";")

	return str.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Value
}

func (rs *ReturnStatement) String() string {
	var str bytes.Buffer

	str.WriteString(rs.TokenLiteral())
	if rs.ReturnValue != nil {
		str.WriteString(" " + rs.ReturnValue.String())
	}

	str.WriteString(";")

	return str.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Value
}
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) statementNode()       {}
func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Value }
func (il *IntegerLiteral) String() string       { return il.Token.Value }

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) statementNode()       {}
func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Value }
func (pe *PrefixExpression) String() string {
	var str bytes.Buffer

	str.WriteString("(")
	str.WriteString(pe.Operator)
	str.WriteString(pe.Right.String())
	str.WriteString(")")

	return str.String()
}

type InfixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
	Left     Expression
}

func (ie *InfixExpression) statementNode()       {}
func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Value }
func (ie *InfixExpression) String() string {
	var str bytes.Buffer

	str.WriteString("(")
	str.WriteString(ie.Left.String() + " ")
	str.WriteString(ie.Operator + " ")
	str.WriteString(ie.Right.String())
	str.WriteString(")")

	return str.String()
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) statementNode()       {}
func (b *Boolean) expressionNode()      {}
func (b *Boolean) TokenLiteral() string { return b.Token.Value }
func (b *Boolean) String() string       { return b.Token.Value }

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) statementNode()       {}
func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Value }
func (ie *IfExpression) String() string {
	var str bytes.Buffer

	str.WriteString("if ")
	str.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		str.WriteString("else ")
		str.WriteString(ie.Alternative.String())
	}

	return str.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) expressionNode()      {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Value }
func (bs *BlockStatement) String() string {
	var str bytes.Buffer

	for _, s := range bs.Statements {
		str.WriteString(s.String())
	}

	return str.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) statementNode()       {}
func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Value }
func (fl *FunctionLiteral) String() string {
	var str bytes.Buffer
	str.WriteString(fl.TokenLiteral() + " (")

	parameters := []string{}
	for _, p := range fl.Parameters {
		parameters = append(parameters, p.String())
	}
	str.WriteString(strings.Join(parameters, ", "))
	str.WriteString(") ")
	str.WriteString(fl.Body.String())

	return str.String()
}

type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Value }
func (ce *CallExpression) String() string {
	var str bytes.Buffer

	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	str.WriteString(ce.Function.String())
	str.WriteString("(")
	str.WriteString(strings.Join(args, ", "))
	str.WriteString(")")

	return str.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Value }
func (sl *StringLiteral) String() string       { return sl.Token.Value }

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Value }
func (al *ArrayLiteral) String() string {
	var str bytes.Buffer

	str.WriteString("[")

	els := []string{}
	for _, el := range al.Elements {
		els = append(els, el.String())
	}

	str.WriteString(strings.Join(els, ", "))
	str.WriteString("]")

	return str.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Value }
func (ie *IndexExpression) String() string {
	var str bytes.Buffer

	str.WriteString(ie.Left.String())
	str.WriteString("[")
	str.WriteString(ie.Index.String())
	str.WriteString("]")

	return str.String()
}
