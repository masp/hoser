package parser

import (
	"fmt"
	"strings"

	"github.com/masp/hoser/lexer"
)

// A hoser module is a set of function definitions. Functions can be either:
// - Pure (no inputs, no outputs)
// - Sources (no inputs, only outputs)
// - Sinks (only inputs, no outputs)
// - Mixed (inputs, outputs)
//
// A function is defined by a name, a set of inputs, and a set of outputs.
//
// An input and output are defined by a name and an optional type (default any).
//
// A function is composed of 1 or more expressions. An expression can be either:
// - An assignment statement e.g. `a = 5`
// - Function call e.g. `a(value: 10)`
//

type Node interface {
	Span() (lexer.Token, lexer.Token) // the first and last token that this node spans
	String() string
}

type Module struct {
	Blocks map[string]*Block
}

func (m *Module) String() string {
	var sb strings.Builder
	for _, fn := range m.Blocks {
		sb.WriteString(fn.String())
		sb.WriteString("; ")
	}
	return sb.String()
}

type Block struct {
	Name    *Identifier
	Inputs  *Map
	Outputs *Map
	Body    []Expression
}

func (f *Block) Span() (lexer.Token, lexer.Token) {
	return f.Name.Span()
}

func (f *Block) String() string {
	var sb strings.Builder
	sb.WriteString(f.Name.Token.Value)

	sb.WriteString("(")
	if f.Inputs != nil {
		sb.WriteString(f.Inputs.String())
	}
	sb.WriteString(")")

	if f.Outputs != nil {
		sb.WriteString(" (")
		sb.WriteString(f.Outputs.String())
		sb.WriteString(")")
	}

	if f.Body != nil {
		sb.WriteString(" {")
		for _, expr := range f.Body {
			sb.WriteString("\n")
			sb.WriteString(expr.String())
		}
		sb.WriteString("}")
	}

	return sb.String()
}

type Expression interface {
	Node
}

// An entry are two expressions separated by a colon
type Entry struct {
	Key Identifier
	Val Expression
}

func (e *Entry) Span() (lexer.Token, lexer.Token) {
	left, _ := e.Key.Span()
	_, right := e.Val.Span()
	return left, right
}

func (e *Entry) String() string {
	return fmt.Sprintf("%v: %v", e.Key.String(), e.Val.String())
}

type Map struct {
	StartToken lexer.Token
	Entries    []Entry
}

func (m *Map) Span() (lexer.Token, lexer.Token) {
	end := m.StartToken
	if len(m.Entries) > 0 {
		_, end = m.Entries[len(m.Entries)-1].Span()
	}
	return m.StartToken, end
}

func (m *Map) String() string {
	var sb strings.Builder
	for i := 0; i < len(m.Entries); i++ {
		sb.WriteString(m.Entries[i].String())
		if i < len(m.Entries)-1 {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}

type Identifier struct {
	Token lexer.Token
}

func (i *Identifier) Span() (lexer.Token, lexer.Token) {
	return i.Token, i.Token
}

func (i *Identifier) String() string {
	return i.Token.Value
}

type BlockCall struct {
	Name *Identifier
	Args *Map
}

func (e *BlockCall) Span() (lexer.Token, lexer.Token) {
	return e.Name.Span()
}

func (e *BlockCall) String() string {
	return fmt.Sprintf("%v(%v)", e.Name.String(), e.Args.String())
}

// AssignmentExpr is an expression separated by an =
type AssignmentExpr struct {
	Left  Expression
	Right Expression
}

func (e *AssignmentExpr) Span() (lexer.Token, lexer.Token) {
	left, _ := e.Left.Span()
	_, right := e.Right.Span()
	return left, right
}

func (e *AssignmentExpr) String() string {
	return fmt.Sprintf("%v = %v", e.Left.String(), e.Right.String())
}

type String struct {
	Token lexer.Token
}

func (s *String) Value() string {
	return s.Token.Value
}

func (s *String) Span() (lexer.Token, lexer.Token) {
	return s.Token, s.Token
}

func (s *String) String() string {
	return fmt.Sprintf("\"%s\"", s.Token.Value)
}

type Integer struct {
	Token lexer.Token
	Value int64
}

func (i *Integer) Span() (lexer.Token, lexer.Token) {
	return i.Token, i.Token
}

func (i *Integer) String() string {
	return fmt.Sprintf("%d", i.Value)
}

type Float struct {
	Token lexer.Token
	Value float64
}

func (f *Float) Span() (lexer.Token, lexer.Token) {
	return f.Token, f.Token
}

func (f *Float) String() string {
	return fmt.Sprintf("%f", f.Value)
}

type Return struct {
	Token lexer.Token
	Value Expression
}

func (r *Return) Span() (lexer.Token, lexer.Token) {
	_, end := r.Value.Span()
	return r.Token, end
}

func (r *Return) String() string {
	return fmt.Sprintf("return %s", r.Value.String())
}

func PrintTokenAst(node Node, indent int) string {
	var sb strings.Builder
	sb.WriteString(strings.Repeat("\t", indent))
	beg, end := node.Span()
	sb.WriteString(beg.String())
	sb.WriteString(" - ")
	sb.WriteString(end.String())
	sb.WriteString("\n")
	indent += 1
	switch v := node.(type) {
	case *AssignmentExpr:
		sb.WriteString(PrintTokenAst(v.Left, indent))
		sb.WriteString(PrintTokenAst(v.Right, indent))
	case *Map:
		for _, entry := range v.Entries {
			sb.WriteString(PrintTokenAst(&entry, indent))
		}
	case *Entry:
		sb.WriteString(PrintTokenAst(&v.Key, indent))
		sb.WriteString(PrintTokenAst(v.Val, indent))
	case *Block:
		sb.WriteString(PrintTokenAst(v.Inputs, indent))
		sb.WriteString(PrintTokenAst(v.Outputs, indent))
		for _, e := range v.Body {
			sb.WriteString(PrintTokenAst(e, indent))
		}
	case *BlockCall:
		sb.WriteString(PrintTokenAst(v.Args, indent))
	}
	return sb.String()
}
