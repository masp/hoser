package ast

import (
	"github.com/masp/hoser/token"
)

type PortIdx int  // PortIdx references the node in the Graph, range: -1 -> len(Nodes)-1
type BlockIdx int // NodeIdx references the port (in/out) for each node, range: 0 -> len(Ports)-1

// Loc represents a position in a Graph and is connected to the source by Pos.
type Loc struct {
	Block BlockIdx
	Port  PortIdx
}

// RootNode is the node index referring to the block this graph describes
// Example:
// 	main() { a() }
// main() is the root block that has zero outputs and zero inputs.
const RootBlock BlockIdx = -1

// Graph is a DAG representation of a pipe's body where each node represents a "block" to be executed with an ordered set of input and output ports. These ports
// are labeled 0-N. Each block contains information about all its inputs and outputs, which are referenced by their indices. A graph component
// is created for each block defined in the module and referenced externally by the tracer.
//
// Only numeric indices are used rather than names to keep the representation concise and avoid circular references.
// The downside is that the program graph is difficult to modify in place (all indices change if anything is added or removed). Since most programs are
// smaller, though, it is not too expensive to recalculate the whole graph each time.
type Graph struct {
	Blocks []Block // the sequence of blocks. the NodeIdx is used to lookup in the slice
	Edges  []Edge  // edges connecting two nodes
}

func portsFromFields(fields FieldList) (ports []EdgeType) {
	for _, field := range fields.Fields {
		ports = append(ports, EdgeType(field.Value.(*Ident).V))
	}
	return
}

func (g *Graph) AddNamedBlock(decl BlockDecl, createdBy Node) BlockIdx {
	var newBlock Block
	switch b := decl.(type) {
	case *PipeDecl:
		newBlock = &PipeBlock{
			Decl:      b,
			inPorts:   portsFromFields(b.Inputs),
			outPorts:  portsFromFields(b.Outputs),
			createdBy: createdBy,
		}
	case *StubDecl:
		newBlock = &StubBlock{
			Decl:      b,
			inPorts:   portsFromFields(b.Inputs),
			outPorts:  portsFromFields(b.Outputs),
			createdBy: createdBy,
		}
	}
	g.Blocks = append(g.Blocks, newBlock)
	return BlockIdx(len(g.Blocks) - 1)
}

func (g *Graph) AddLiteralBlock(lit *LiteralExpr) BlockIdx {
	g.Blocks = append(g.Blocks, &LiteralBlock{lit})
	return BlockIdx(len(g.Blocks) - 1)
}

func (g *Graph) Connect(src Loc, dst Loc, typ EdgeType) {
	g.Edges = append(g.Edges, Edge{Type: typ, Src: src, Dst: dst})
}

type Block interface {
	// CreatedBy is the AST node that corresponds to this Block being created (connects the graph to the AST)
	CreatedBy() Node

	// The only commonality between Blocks is they have input/output ports named by indices with types.
	InPorts() []EdgeType
	OutPorts() []EdgeType
}

// PipeBlock refers to a user defined pipe block, e.g. pipe would be a PipeBlock below:
// pipe() { A() }
//
// This block is divisible and expandable, for example:
// 	A() { B(); C() }
// 	pipe() { A() } => pipe() { B(); C() }
type PipeBlock struct {
	createdBy Node
	Decl      *PipeDecl // block name that is executed by this node, can be used to look up definition.
	inPorts   []EdgeType
	outPorts  []EdgeType
}

// LiteralBlock is a block with a single output port that evaluates constantly to the literal expression
// This block is atomic.
type LiteralBlock struct {
	Lit *LiteralExpr
}

// StubBlock refers to a block that is stubbed and is defined by an external process or Go code
// This block is atomic.
type StubBlock struct {
	createdBy Node
	Decl      *StubDecl
	inPorts   []EdgeType
	outPorts  []EdgeType
}

func (b PipeBlock) CreatedBy() Node      { return b.createdBy }
func (b PipeBlock) InPorts() []EdgeType  { return b.inPorts }
func (b PipeBlock) OutPorts() []EdgeType { return b.outPorts }

func (b StubBlock) CreatedBy() Node      { return b.createdBy }
func (b StubBlock) InPorts() []EdgeType  { return b.inPorts }
func (b StubBlock) OutPorts() []EdgeType { return b.outPorts }

func (b LiteralBlock) CreatedBy() Node     { return b.Lit }
func (b LiteralBlock) InPorts() []EdgeType { return nil }
func (b LiteralBlock) OutPorts() []EdgeType {
	switch b.Lit.Type {
	case token.Integer:
		return []EdgeType{IntEdge}
	case token.Float:
		return []EdgeType{FloatEdge}
	case token.String:
		return []EdgeType{StringEdge}
	default:
		panic("invalid edge type")
	}
}

type EdgeType string

var (
	InvalidEdge EdgeType = ""
	StringEdge  EdgeType = "string"
	IntEdge     EdgeType = "int"
	FloatEdge   EdgeType = "float"
)

// Edge connects a "Src" Loc to a "Dst" Loc using with a typed flow of values
type Edge struct {
	Type     EdgeType // Type is the type of value that flows across this edge
	Src, Dst Loc
}
