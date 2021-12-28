package ast

import (
	"github.com/masp/hoser/token"
)

// A hoser module is a set of pipe definitions. Pipes are componsed of blocks. A block can be either:
// - Pure (no inputs, no outputs)
// - Sources (no inputs, only outputs)
// - Sinks (only inputs, no outputs)
// - Mixed (inputs, outputs)
//
// A block is defined by a name, a set of inputs, and a set of outputs.
//
// An input and output are defined by a name and an optional type (default any).
//
// A pipe block is composed of 1 or more expressions. An expression can be either:
// - An assignment statement e.g. `a = 5`
// - Another block call e.g. `a(value: 10)`
//

type Node interface {
	Pos() token.Pos // the offset where this node starts
	End() token.Pos // the offset where this node ends
}

// Decl is only a pipe block declaration, e.g. `block()`
type Decl interface {
	Node
	declNode()
}

type Stmt interface {
	Node
	stmtNode()
}

type Expr interface {
	Node
	exprNode()
}

// ----------------------------------------------------------------------------
// Declarations
//

// BlockDecl is either a stub or a pipe
type BlockDecl interface {
	Decl
	BlockName() string
	BlockInputs() *FieldList
	BlockOutputs() *FieldList
}

// Module represents all the contents of a single file, including all defined blocks and all referenced blocks.
type Module struct {
	ModulePos     token.Pos     // position of module keyword
	Name          *LiteralExpr  // name of module identifier as a string literal
	Imports       []*ImportDecl // list of imported modules
	DefinedBlocks []BlockDecl
}

func (m *Module) Pos() token.Pos {
	return m.ModulePos
}

func (m *Module) End() token.Pos {
	if len(m.DefinedBlocks) > 0 {
		return m.DefinedBlocks[len(m.DefinedBlocks)-1].End()
	} else {
		return m.Name.End()
	}
}

func (m *Module) Lookup(name string) BlockDecl {
	for _, decl := range m.DefinedBlocks {
		if decl.BlockName() == name {
			return decl
		}
	}
	return nil
}

type ImportDecl struct {
	Keyword    token.Pos
	ModuleName *LiteralExpr // import "ModuleName"
}

func (b *ImportDecl) Pos() token.Pos {
	return b.Keyword
}

func (b *ImportDecl) End() token.Pos {
	return b.ModuleName.End()
}

type StubDecl struct {
	Name    *Ident
	Inputs  FieldList
	Outputs FieldList
}

func (b *StubDecl) Pos() token.Pos {
	return b.Name.Pos()
}

func (b *StubDecl) End() token.Pos {
	return b.Outputs.End()
}

type PipeDecl struct {
	StubDecl
	BegLBrack token.Pos
	Body      []Stmt
	EndRBrack token.Pos

	BodyDAG *Graph // nil by default, added in by tracer
}

func (b *PipeDecl) Pos() token.Pos {
	return b.Name.Pos()
}

func (b *PipeDecl) End() token.Pos {
	return b.EndRBrack
}

func (b *PipeDecl) BlockName() string        { return b.Name.V }
func (b *PipeDecl) BlockInputs() *FieldList  { return &b.Inputs }
func (b *PipeDecl) BlockOutputs() *FieldList { return &b.Outputs }
func (b *StubDecl) BlockName() string        { return b.Name.V }
func (b *StubDecl) BlockInputs() *FieldList  { return &b.Inputs }
func (b *StubDecl) BlockOutputs() *FieldList { return &b.Outputs }

func (m *Module) declNode()     {}
func (m *ImportDecl) declNode() {}
func (m *PipeDecl) declNode()   {}
func (m *StubDecl) declNode()   {}

// ----------------------------------------------------------------------------
// Expressions
//

// Field is a key-value combination like 'key: value' that shows up in pipe definitions and pattern
// matching.
type Field struct {
	Key   *Ident
	Colon token.Pos
	Value Expr
}

func (f *Field) Pos() token.Pos {
	return f.Key.Pos()
}

func (f *Field) End() token.Pos {
	return f.Value.End()
}

type FieldList struct {
	Opener token.Pos // Opening { or NoPos if none
	Fields []*Field
	Closer token.Pos // Closing } or NoPos if none
}

func (f *FieldList) Pos() token.Pos {
	if f.Opener.IsValid() {
		return f.Opener
	} else {
		if f.Len() > 0 {
			return f.Fields[0].Pos()
		} else {
			panic("impossible to have empty fields without surrounding {/(")
		}
	}
}

func (f *FieldList) End() token.Pos {
	if f.Closer.IsValid() {
		return f.Opener
	} else {
		if f.Len() > 0 {
			return f.Fields[len(f.Fields)-1].Pos()
		} else {
			panic("impossible to have empty fields without surrounding }/)")
		}
	}
}

func (f *FieldList) Len() int {
	return len(f.Fields)
}

type CallExpr struct {
	Name   *Ident
	Lparen token.Pos
	Args   []Expr
	Rparen token.Pos
}

func (c *CallExpr) Pos() token.Pos {
	return c.Name.Pos()
}

func (c *CallExpr) End() token.Pos {
	return c.Rparen
}

type ParenExpr struct {
	X Expr
}

func (p *ParenExpr) Pos() token.Pos {
	return p.X.Pos() - token.Pos(1)
}

func (p *ParenExpr) End() token.Pos {
	return p.X.Pos() + token.Pos(1)
}

type Ident struct {
	V         string // V is the string value of the ident (= Run in mod.Run())
	NamePos   token.Pos
	Module    string // Module is mod in mod.Run(), if unscoped Module="" and ModulePos=NoPos
	ModulePos token.Pos
}

func (i *Ident) Pos() token.Pos {
	return i.NamePos
}

func (i *Ident) End() token.Pos {
	return i.NamePos + token.Pos(len(i.V))
}

func (i *Ident) Local() bool {
	return !i.ModulePos.IsValid()
}

func (i *Ident) FullName() string {
	if i.Local() {
		return i.V
	}
	return i.Module + "." + i.V
}

type LiteralExpr struct {
	Start     token.Pos
	Type      token.Token // e.g. token.String, Integer or Float
	Value     string
	ParsedVal interface{}
}

func (lit *LiteralExpr) Pos() token.Pos {
	return lit.Start
}

func (lit *LiteralExpr) End() token.Pos {
	if lit.Type == token.String {
		return lit.Start + token.Pos(len(lit.Value)+2) // +2 for ""
	} else {
		return lit.Start + token.Pos(len(lit.Value))
	}
}

// AssignExpr is an expression separated by an =
type AssignExpr struct {
	Lhs   Expr
	EqPos token.Pos
	Rhs   Expr
}

func (a *AssignExpr) Pos() token.Pos { return a.Lhs.Pos() }
func (a *AssignExpr) End() token.Pos { return a.Rhs.Pos() }

func (*Field) exprNode()       {}
func (*FieldList) exprNode()   {}
func (*Ident) exprNode()       {}
func (*CallExpr) exprNode()    {}
func (*LiteralExpr) exprNode() {}
func (*ParenExpr) exprNode()   {}
func (*AssignExpr) exprNode()  {}

// ----------------------------------------------------------------------------
// Statements
//

type ExprStmt struct {
	X Expr
}

func (e *ExprStmt) Pos() token.Pos { return e.X.Pos() }
func (e *ExprStmt) End() token.Pos { return e.X.End() }

func (e *ExprStmt) stmtNode() {}
