package runtime

import (
	"fmt"

	"github.com/masp/hoser/ast"
	"github.com/masp/hoser/graph"
)

// runtime is responsible for taking a read-only state and applying a dataflow computation to produce a new state.

// Ports is a small struct that is a list of pointers to underlying flat arrays
// For example, if a block has 2 input ports first int and second str type, then:
// Ports{
//  order: 'is'x
// 	intPorts: {&port1}
//  strPorts: {&port2}
// }
type Ports struct {
	order    [8]byte
	intPorts []*int
	strPorts []*string
}

type machineFn func(m *Machine, in Ports, out Ports)

type Machine struct {
	blockTable   map[string]*ast.Block
	nativeBlocks map[string]machineFn
}

func NewMachine() *Machine {
	m := &Machine{
		blockTable: make(map[string]*ast.Block),
		nativeBlocks: map[string]machineFn{
			"Print": printFn,
		},
	}
	return m
}

func (m *Machine) LoadModule(module *ast.Module) *Machine {
	for _, block := range module.Blocks {
		m.blockTable[block.Name.Token.Value] = block
	}
	return m
}

func printFn(m *Machine, in Args, out Args) {
	a := in.IntArg(0)

}

func (m *Machine) DefineBlock(name string, fn machineFn) *Machine {
	m.nativeBlocks[name] = fn
	return m
}

// argBuffers is the underlying backing for where blocks store and retrieve their input and output args from.
// There are two types of arrays:
// - The ports array is equal to the number of ports in the whole graph, which holds all the references for each block
// to the data array.
// - The data array is what is written to and read from by the blocks when they are doing their computation.
type argBuffers struct {
	intDataBuf []int
	intPortBuf []*int // points to intDataBuf

	strDataBuf []string
	strPortBuf []*string // points to strDataBuf
}

func (m *Machine) buildArgBuffers(program graph.Definition) (argBuffers, error) {
	var (
		numIntPorts, numStrPorts int = 0, 0
	)

	for i, node := range program.Nodes {
		switch v := node.(type) {
		case graph.BlockNode:
			blockDef, ok := m.blockTable[v.Block]
			if !ok {
				return argBuffers{}, fmt.Errorf("block name '%s' not found", v.Block)
			}
			for _, entry := range blockDef.Inputs.Entries {
				switch entry.Val.(type) {
				case lexer.
				}
			}
		}

		numIntPorts += node
	}
}

func (m *Machine) Run(program graph.Definition) error {
	var numIntPorts, numStrPorts int

	intPortBuf := make([]int, 0)
	intDataBuf := make([]int, 0)
	strArgBuf := make([]string)
	for _, edge := range graph.Edges {
		edge.SrcNode
	}

	out := in.Copy()

	for _, node := range graph.Nodes {
		if block, ok := m.blocks[node.Block]; ok {
			block(m, in, &out)
		} else {
			return in, fmt.Errorf("no block found with name ")
		}
	}
}
