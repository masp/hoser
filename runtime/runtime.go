package runtime

import (
	"fmt"

	"github.com/masp/hoser/graph"
)

// runtime is responsible for taking a read-only state and applying a dataflow computation to produce a new state.

type State struct {
	intVars map[string]int
}

func (s State) Copy() State {
	newState := State{
		intVars: make(map[string]int),
	}

	for k, v := range s.intVars {
		newState.intVars[k] = v
	}
	return newState
}

type machineFn func(*Machine, State, *State)

type Machine struct {
	blocks map[string]machineFn
	args   [16]int // max 16 possible ports
}

func NewMachine() *Machine {
	m := &Machine{
		blocks: map[string]machineFn{
			"Fetch": fetchImpl,
			"Store": storeImpl,
		},
	}
	return m
}

func fetchImpl(m *Machine, in State, out *State) {
	a := m.IntArg(0)

}

func storeImpl(m *Machine, in State, out *State) {

}

func (m *Machine) DefineBlock(name string, fn machineFn) *Machine {
	m.blocks[name] = fn
	return m
}

func (m *Machine) Run(in State, graph graph.Definition) (State, error) {
	out := in.Copy()

	for _, node := range graph.Nodes {
		if block, ok := m.blocks[node.Block]; ok {
			block(m, in, &out)
		} else {
			return in, fmt.Errorf("no block found with name ")
		}
	}
}
