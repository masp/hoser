package ast

type PortIdx int  // PortIdx references the node in the Graph, range: -1 -> len(Nodes)-1
type BlockIdx int // NodeIdx references the port (in/out) for each node, range: 0 -> len(Ports)-1

// Loc represents a position in a Graph.
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

type Block interface {
	block()
}

// PipeBlock refers to a user defined pipe block, e.g. pipe would be a PipeBlock below:
// pipe() { A() }
//
// This block is divisible and expandable, for example:
// 	A() { B(); C() }
// 	pipe() { A() } => pipe() { B(); C() }
type PipeBlock struct {
	*PipeDecl // block name that is executed by this node, can be used to look up definition.
}

// LiteralBlock is a block with a single output port that evaluates constantly to the literal expression
// This block is atomic.
type LiteralBlock struct {
	*LiteralExpr
}

// StubBlock refers to a block that is stubbed and is defined by an external process or Go code
// This block is atomic.
type StubBlock struct {
	*StubDecl
}

func (b PipeBlock) block()    {}
func (b LiteralBlock) block() {}
func (b StubBlock) block()    {}

type EdgeType int

const (
	IntEdge EdgeType = iota
	FloatEdge
	StringEdge
)

// Edge connects a "Src" Loc to a "Dst" Loc using with a typed flow of values
type Edge struct {
	Type     EdgeType // Type is the type of value that flows across this edge
	Src, Dst Loc
}
