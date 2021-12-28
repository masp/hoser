package tracer

import (
	"fmt"

	"github.com/masp/hoser/ast"
	"github.com/masp/hoser/token"
)

// output represents what an expression is evaluated to when tracing
// For example
// 	a = Run()
// assigns whatever outputs Run produces (0, 1, 2...) to the symbol 'a'.
// The outputs can then be referenced later on to be fed into other blocks.
// Outputs differ from values because they are not copyable and are inseparable from
// their sources, since outputs represent continuous streams of data and not single quantities.
type output interface {
	output()
}

var NilOutput output = nil

// oneOutput is a single output from a single port.
// Given a block like `A() (v: int)` called like `v = A()`, v would have a `portOutput` value.
type oneOutput struct {
	From ast.Loc
}

// outputBundle is a bundle of outputs, usually when a block returns multiple values.
//  Given a block like `A() (v1: int, v2: int)` called like `v = A()`, v would be `manyOutput` with {v1, v2}.
type outputBundle struct {
	Outputs map[string]output
}

func (o oneOutput) output()    {}
func (o outputBundle) output() {}

func makeOutputBundle(block ast.BlockIdx, decl ast.BlockDecl) output {
	if decl.BlockOutputs() == nil {
		return NilOutput
	}

	switch len(decl.BlockOutputs().Fields) {
	case 0:
		return NilOutput
	case 1:
		return oneOutput{From: ast.Loc{Block: block, Port: 0}}
	default:
		bundle := make(map[string]output)
		for port, field := range decl.BlockOutputs().Fields {
			bundle[field.Key.V] = oneOutput{
				From: ast.Loc{
					Block: block,
					Port:  ast.PortIdx(port),
				},
			}
		}
		return outputBundle{Outputs: bundle}
	}
}

// Tracer will traverse a set of modules and try to resolve for a given module a fully connected DAG
//
// For example:
// A() (v: int) { v = 10 }
// main() {
// 	shell.Run(A())
// }
//
// Tracer needs to resolve the symbols Run() and A() and create DAGs representing their bodies. shell.Run is defined in another module (shell)
// and so the Tracer needs to find the file in the include path, parse and recursively trace its body, and then return the graph here.
//
// Tracer also needs to verify that edges between ports are never incorrectly typed or that ports are missing connections or unused.
//
// The end product is a fully connected set of DAGs with the only terminal blocks being stubs (defined in Go) and literal blocks.
type Tracer struct {
	modCache ast.ModuleSet

	tracingMod  *ast.Module
	tracingFile *token.File

	errors token.ErrorList
}

func (t *Tracer) expectedError(node ast.Node, msg string) {
	t.error(node.Pos(), fmt.Errorf("expected %v, got %T", msg, node))
}

func (t *Tracer) connect(src ast.Loc, dst ast.Loc, graph *ast.Graph) {
	var srcType, dstType ast.EdgeType
	srcBlock := graph.Blocks[src.Block]
	dstBlock := graph.Blocks[dst.Block]
	srcType = srcBlock.OutPorts()[src.Port]
	dstType = dstBlock.InPorts()[dst.Port]
	if srcType != dstType {
		t.error(dstBlock.CreatedBy().Pos(), fmt.Errorf("type mismatch: got %v, expected %v", srcType, dstType))
		return
	}
	graph.Connect(src, dst, srcType)
}

func (t *Tracer) traceModule(file *token.File, mod *ast.Module) {
	t.tracingMod = mod
	t.tracingFile = file
	for _, decl := range mod.DefinedBlocks {
		if pipe, ok := decl.(*ast.PipeDecl); ok {
			pipe.BodyDAG = t.tracePipe(pipe)
		}
	}
}

type pipeTrace struct {
	Graph       ast.Graph
	symbolTable map[string]output
}

func (t *Tracer) tracePipe(pipe *ast.PipeDecl) *ast.Graph {
	trace := pipeTrace{Graph: ast.Graph{}, symbolTable: make(map[string]output)}
	for _, stmt := range pipe.Body {
		t.traceStmt(stmt, &trace)
	}
	return &trace.Graph
}

func (t *Tracer) traceStmt(stmt ast.Stmt, state *pipeTrace) {
	switch st := stmt.(type) {
	case *ast.ExprStmt:
		t.traceExpr(st.X, state)
	}
}

func (t *Tracer) traceExpr(expr ast.Expr, state *pipeTrace) output {
	switch x := expr.(type) {
	case *ast.AssignExpr:
		return t.traceAssign(x, state)
	case *ast.CallExpr:
		return t.traceCall(x, state)
	case *ast.Ident:
		return t.traceIdent(x, state)
	case *ast.LiteralExpr:
		return t.traceLit(x, state)
	default:
		return NilOutput
	}
}

func (t *Tracer) traceCall(call *ast.CallExpr, state *pipeTrace) (out output) {
	if call.Name.Local() {
		decl := t.tracingMod.Lookup(call.Name.V)
		if decl == nil {
			t.error(call.Pos(), fmt.Errorf("unable to find local pipe or stub with name %v", call.Name.V))
			return NilOutput
		}

		if len(call.Args) != len(decl.BlockInputs().Fields) {
			t.error(call.Pos(), fmt.Errorf("wrong number of args for call to %v, expected %d, got %d",
				call.Name.V,
				len(decl.BlockInputs().Fields),
				len(call.Args)))
		}

		var (
			usedPorts []int
			argval    ast.Expr
		)
		incomingEdges := make([]ast.Loc, len(decl.BlockInputs().Fields))
		for _, arg := range call.Args {
			usedPorts, argval = t.matchArgToInput(arg, decl.BlockInputs(), usedPorts)
			if argval == nil {
				// error returned by matchArgToField
				return NilOutput
			}

			foundPort := ast.PortIdx(usedPorts[len(usedPorts)-1])
			tracedarg := t.traceExpr(argval, state)
			if inarg, ok := tracedarg.(oneOutput); ok {
				incomingEdges[foundPort] = inarg.From
			} else {
				t.error(argval.Pos(), fmt.Errorf("expected single output, got %v", tracedarg))
			}
		}

		// Add the block and edges now after all the args have added their input blocks to the graph
		thisBlock := state.Graph.AddNamedBlock(decl, call)
		for port, from := range incomingEdges {
			t.connect(
				from,
				ast.Loc{Block: thisBlock, Port: ast.PortIdx(port)},
				&state.Graph,
			)
		}

		out = makeOutputBundle(thisBlock, decl)
		return
	}
	panic("unsupported global name")
}

const namedArgUsedPort = 999 // namedArgUsedPort is in usedArgs it means a named arg was used and a positional cannot be used anymore

func (t *Tracer) matchArgToInput(argExpr ast.Expr, inputs *ast.FieldList, usedPorts []int) (ports []int, value ast.Expr) {
	switch arg := argExpr.(type) {
	case *ast.Field:
		// the arg is a named argument (field: value)
		for i, field := range inputs.Fields {
			if field.Key.V == arg.Key.V {
				for _, used := range usedPorts {
					if used == i {
						t.error(arg.Pos(), fmt.Errorf("already used argument with name %s", field.Key.V))
						return
					}
				}
				ports = append(ports, namedArgUsedPort, i)
				value = arg.Value
				return
			}
		}
		t.error(arg.Pos(), fmt.Errorf("no input found with name %s", arg.Key.V))
	default:
		// the arg is a positional argument, the field is based on the position
		nextPort := len(usedPorts)
		for _, port := range usedPorts {
			if port == namedArgUsedPort {
				t.error(arg.Pos(), fmt.Errorf("positional arg cannot be after named arg"))
				return
			}
		}

		if len(inputs.Fields) <= nextPort {
			t.error(arg.Pos(), fmt.Errorf("too many arguments, expected %d got %d", len(inputs.Fields), nextPort))
			return
		}
		ports = append(ports, nextPort)
		value = arg
	}
	return
}

func (t *Tracer) traceIdent(ident *ast.Ident, state *pipeTrace) (out output) {
	var ok bool
	if out, ok = state.symbolTable[ident.V]; !ok {
		t.error(ident.Pos(), fmt.Errorf("no symbol found with name %v", ident.V))
	}
	return
}

func (t *Tracer) traceLit(lit *ast.LiteralExpr, state *pipeTrace) oneOutput {
	idx := state.Graph.AddLiteralBlock(lit)
	return oneOutput{ast.Loc{Block: idx, Port: 0}}
}

func (t *Tracer) traceAssign(assign *ast.AssignExpr, state *pipeTrace) output {
	// "Assignments" are more pattern matching than actual variables storing state. We destructure the right hand side
	// according to the pattern on the left hand side. Unbinded symbol names are bound to the ports/outputs on the right hand side.
	rhs := t.traceExpr(assign.Rhs, state)
	t.unifyExpr(assign.Lhs, rhs, state)
	return rhs
}

func (t *Tracer) unifyExpr(pattern ast.Expr, rhs output, state *pipeTrace) {
	switch p := pattern.(type) {
	case *ast.Ident:
		// 1. An identifier a = b()
		t.unifyOne(p, rhs, state)
	case *ast.FieldList:
		// 2. A map a: b, c: d = v()
		t.unifyBundle(p, rhs, state)
	default:
		t.expectedError(pattern, "variable name or map of variables")
	}
}

func (t *Tracer) unifyOne(pattern *ast.Ident, rhs output, state *pipeTrace) {
	varName := pattern.V
	state.symbolTable[varName] = rhs
}

func (t *Tracer) unifyBundle(pattern *ast.FieldList, rhs output, state *pipeTrace) error {
	if rBundle, ok := rhs.(outputBundle); ok {
		for _, field := range pattern.Fields {
			if foundOutput, ok := rBundle.Outputs[field.Key.V]; ok {
				t.unifyOne(field.Key, foundOutput, state)
			} else {
				t.error(field.Key.Pos(), fmt.Errorf("name does not match any output on right side of assignment"))
			}
		}
	} else {
		// a: b = 10 not okay
		t.expectedError(pattern, "more than one output on right side of assignment")
	}
	return nil
}
