package graph

import (
	"fmt"

	"github.com/masp/hoser/parser"
)

// graph is a semantic representation of the parsed hoser syntax as a DAG that is more suitable for
// the scheduler and visualization tools.
//
// It condenses the following program:
// b(b0: int) (bo0: int) {}
// d(d0: int, d1: int) (do0: int) {}
// e(e0: int) (eo0: int) {}
// main(input: int) (c: int) {
// 	bo0: a = b(input)
//  do0: c = d(d0: a, d1: e(e0: 5))
// }
//
// into a series of nodes with ports where the inputs/outputs are connected
//
// Key functions are:
// - Replacing names of blocks and ports (inputs or outputs) with numerical indices
// - Flattening structure so that nesting is removed (a(b()) -> a, b)
// - Removing unused symbols/lines that are disconnected (c in the above example program)
//
// the output of the above program would:
//
// -1 refers to the inputs and outputs of the block as a whole
// Program: [
//  (id, inputs, outputs) (idx),
// 	(b, (0 <- -1.0), (0 -> 3.0)), (0)
// 	(5, (), (0 -> 2.0)), (1)
// 	(e, (0 <- 1.0), (0 -> 3.1)), (2)
// 	(d, (0 <- 0.0, 1 <- 2.0), (0 -> -1.0)), (3)
// ]
//

type PortIdx int
type NodeIdx int

// Definition is a DAG representation of a block's body where each node represents a "block" to be executed with an ordered set of input and output ports. These ports
// are labeled 0-N. Each node contains information about all its inputs and outputs, which are referenced by their indices. A graph component
// is created for each block defined in the module.
//
// Only numeric indices are used rather than names to keep the representation concise and avoid circular references which can create complexity.
// The downside is that the program graph is difficult to modify in place (all indices change if anything is added or removed). Since most programs are
// smaller, though, it is not too expensive to recalculate the whole graph each time.
type Definition struct {
	Nodes []Node // the sequence of blocks. the indices of each node is used to lookup
	Edges []Edge // edges connecting two nodes
}

// Module is a collection of subroutine definitions. It also contains the definitions of all blocks in the module so that
// information (like types of inputs, names of ports) can be referenced.
type Module struct {
	Definitions map[string]Definition // a definition for every block defined in a module by block name
}

type Node interface {
	Name() string
}

// BlockNode is a node that is a block
type BlockNode struct {
	Block string // block name that is executed by this node, can be used to look up definition.
}

func (b BlockNode) Name() string {
	return b.Block
}

// ConstNode is a node that is a constant literal like 10 or "hello" that outputs a single value on a single port always
type ConstNode struct {
	Value parser.Expression
}

func (c ConstNode) Name() string {
	return c.Value.String()
}

type Edge struct {
	DstNode, SrcNode NodeIdx
	DstPort, SrcPort PortIdx
}

// value represents possible outputs from expressions in hoser.
// a = b: 10, c: 5 is a valid expression where (b: 10, c: 5) is two blocks (10) and (5)
// that are connected to other blocks through the symbol 'a'. If later b(a) is called, we
// need to map the symbol's outputs to b's inputs and create edges between the blocks.
//
// A value then is a collection of outputs from blocks. An output from a block is identified by its node
// idx and port idx associated with a name (b and c in the above example).
//
// Examples of possible values:
// a = 10 -> "a": (n: 0, p: 0)
// a = b: 10 -> "a": [(name: "b", n: 0, p: 0)]
// b: a = b: 10 -> "a": (n: 0, p: 0)
type singleValue struct {
	SrcNode NodeIdx
	SrcPort PortIdx
}

func (s singleValue) dummy() {}

type mapValue struct {
	Values map[string]value
}

func (m mapValue) dummy() {}

// either singleValue or mapValue
type value interface {
	dummy()
}

type grapher struct {
	blockReference   *parser.Module
	symbolTable      map[string]value
	blockBeingParsed *parser.Block
}

func grapherErr(expr parser.Expression, err error) error {
	start, _ := expr.Span()
	return fmt.Errorf("syntax error: %d:%d (%v) %w", start.Line, start.Col, expr, err)
}

func TraceModule(input *parser.Module) (*Module, error) {
	module := Module{
		Definitions: make(map[string]Definition, len(input.Blocks)),
	}

	for i, block := range input.Blocks {
		grapher := grapher{
			blockReference:   input,
			symbolTable:      make(map[string]value),
			blockBeingParsed: block,
		}

		for port, entry := range block.Inputs.Entries {
			name := entry.Key.Token.Value
			grapher.symbolTable[name] = singleValue{
				SrcNode: -1,
				SrcPort: PortIdx(port),
			}
		}
		def, err := grapher.traceBlock(block)
		if err != nil {
			return nil, err
		}
		module.Definitions[i] = def
	}
	return &module, nil
}

func (g *grapher) traceBlock(block *parser.Block) (Definition, error) {
	def := Definition{}
	for _, expr := range block.Body {
		if _, err := g.traceExpression(expr, &def); err != nil {
			return def, err
		}
	}
	return def, nil
}

// Convert an expression to a sequence of blocks (added to def) and return the returned outputs from
// the expression (0 or more outputs)
// For example, the expression a() could return 0 or more outputs
func (g *grapher) traceExpression(expr parser.Expression, def *Definition) (value, error) {
	switch v := expr.(type) {
	case *parser.AssignmentExpr:
		return g.traceAssignment(v, def)
	case *parser.BlockCall:
		return g.traceBlockCall(v, def)
	case *parser.Identifier:
		return g.traceIdentifier(v, def)
	case *parser.Integer, *parser.Float, *parser.String:
		return g.traceLiteral(v, def)
	case *parser.Return:
		return g.traceReturn(v, def)
	case *parser.Map:
		return g.traceMap(v, def)
	}
	return nil, nil
}

func (g *grapher) traceAssignment(assign *parser.AssignmentExpr, def *Definition) (value, error) {
	// "Assignments" are more pattern matching than actual variables storing state. We destructure the right hand side
	// according to the pattern on the left hand side. Unbinded symbol names are bound to the ports/outputs on the right hand side.
	rhs, err := g.traceExpression(assign.Right, def)
	if err != nil {
		return nil, err
	}
	return rhs, g.unifyExpr(assign.Left, rhs)
}

func (g *grapher) unifyExpr(pattern parser.Expression, rhs value) error {
	switch p := pattern.(type) {
	case *parser.Identifier:
		// 1. An identifier a = b()
		return g.unifyValue(p, rhs)
	case *parser.Map:
		// 2. A map a: b, c: d = v()
		return g.unifyMap(p, rhs)
	}
	return grapherErr(pattern, fmt.Errorf("expected left-hand side to be map or variable names only"))
}

func (g *grapher) unifyValue(pattern *parser.Identifier, rhs value) error {
	varName := pattern.Token.Value
	g.symbolTable[varName] = rhs

	if rv, ok := rhs.(mapValue); ok {
		if len(rv.Values) == 1 {
			// `a = b: 5` -> `a = 5` as a simplification
			for _, entry := range rv.Values {
				g.symbolTable[varName] = entry
			}
		}
	}
	return nil
}

func (g *grapher) unifyMap(pattern *parser.Map, rhs value) error {
	if rMap, ok := rhs.(mapValue); ok {
		for _, entry := range pattern.Entries {
			lKey := entry.Key.Token.Value
			err := g.unifyExpr(entry.Val, rMap.Values[lKey])
			if err != nil {
				return err
			}
		}
	} else {
		// a: b = 10 not okay
		return grapherErr(pattern, fmt.Errorf("mismatch between left and right hand side, value %v", rhs))
	}
	return nil
}

func (g *grapher) traceMap(v *parser.Map, def *Definition) (value, error) {
	rv := mapValue{Values: make(map[string]value)}
	for _, entry := range v.Entries {
		name := entry.Key.Token.Value
		val, err := g.traceExpression(entry.Val, def)
		if err != nil {
			return nil, err
		}
		rv.Values[name] = val
	}
	return rv, nil
}

func (g *grapher) traceBlockCall(call *parser.BlockCall, def *Definition) (value, error) {
	blockName := call.Name.Token.Value
	if block, ok := g.blockReference.Blocks[blockName]; ok {
		args := make([]value, len(block.Inputs.Entries))
		for _, arg := range call.Args.Entries {
			inputName := arg.Key.Token.Value
			matchingPortNum := -1
			for i, in := range block.Inputs.Entries {
				if inputName == in.Key.Token.Value {
					matchingPortNum = i
				}
			}

			if matchingPortNum == -1 {
				return nil, fmt.Errorf("unknown input named '%s'", inputName)
			}

			argVal, err := g.traceExpression(arg.Val, def)
			if err != nil {
				return nil, err
			}

			args[matchingPortNum] = argVal
		}

		def.Nodes = append(def.Nodes, BlockNode{Block: blockName})
		thisNode := NodeIdx(len(def.Nodes) - 1)
		for portNum, arg := range args {
			if arg == nil {
				continue // input not provided
			}

			if argSingle, ok := arg.(singleValue); ok {
				def.Edges = append(def.Edges, Edge{
					DstNode: NodeIdx(thisNode),
					DstPort: PortIdx(portNum),
					SrcNode: argSingle.SrcNode,
					SrcPort: argSingle.SrcPort,
				})
			} else {
				return nil, grapherErr(&call.Args.Entries[portNum], fmt.Errorf("expected single value but got %v", arg))
			}
		}

		var outputs value
		if block.Outputs != nil {
			if len(block.Outputs.Entries) == 1 {
				outputs = singleValue{thisNode, PortIdx(0)}
			} else {
				r := make(map[string]value)
				for port, entry := range block.Outputs.Entries {
					r[entry.Key.Token.Value] = singleValue{
						SrcNode: thisNode,
						SrcPort: PortIdx(port),
					}
				}
				outputs = mapValue{Values: r}
			}
		}
		return outputs, nil
	} else {
		return nil, grapherErr(call, fmt.Errorf("block with name %s not found in module", blockName))
	}
}

func (g *grapher) traceIdentifier(id *parser.Identifier, def *Definition) (value, error) {
	variable := id.Token.Value
	if src, ok := g.symbolTable[variable]; ok {
		return src, nil
	} else {
		return nil, grapherErr(id, fmt.Errorf("variable with name '%s' not found", variable))
	}
}

func (g *grapher) traceLiteral(v parser.Expression, def *Definition) (value, error) {
	def.Nodes = append(def.Nodes, ConstNode{v})
	thisNode := NodeIdx(len(def.Nodes) - 1)
	return singleValue{
		SrcNode: thisNode,
		SrcPort: PortIdx(0),
	}, nil
}

func (g *grapher) traceReturn(v *parser.Return, def *Definition) (value, error) {
	retVal, err := g.traceExpression(v.Value, def)
	if err != nil {
		return nil, err
	}

	if len(g.blockBeingParsed.Outputs.Entries) == 0 {
		return nil, grapherErr(v, fmt.Errorf("cannot return values from block with no outputs defined"))
	}

	for port, entry := range g.blockBeingParsed.Outputs.Entries {
		name := entry.Key.Token.Value
		switch r := retVal.(type) {
		case mapValue:
			if storedValue, ok := r.Values[name]; ok {
				if sv, ok := storedValue.(singleValue); ok {
					def.Edges = append(def.Edges, Edge{
						DstNode: -1,
						SrcNode: sv.SrcNode,
						DstPort: PortIdx(port),
						SrcPort: sv.SrcPort,
					})
				}
			}
		case singleValue:
			if len(g.blockBeingParsed.Outputs.Entries) == 1 {
				def.Edges = append(def.Edges, Edge{
					DstNode: -1,
					SrcNode: r.SrcNode,
					DstPort: 0,
					SrcPort: r.SrcPort,
				})
			} else {
				return nil, grapherErr(v, fmt.Errorf("return expected map value not single value"))
			}
		}
	}
	return retVal, nil
}