package parser

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
//  Given a block like `A() (v1: int, v2: int)` called like `v = A()`, v would be `manyOutput` with [v1, v2].
type outputBundle struct {
	Outputs map[string]oneOutput
}

func (o oneOutput) output()    {}
func (o outputBundle) output() {}

// tracer will traverse a set of modules and try to resolve for a given module a fully connected DAG
//
// For example:
// A() (v: int) { v = 10 }
// main() {
// 	shell.Run(A())
// }
//
// tracer needs to resolve the symbols Run() and A() and create DAGs representing their bodies. shell.Run is defined in another module (shell)
// and so the tracer needs to find the file in the include path, parse and recursively trace its body, and then return the graph here.
//
// tracer also needs to verify that edges between ports are never incorrectly typed or that ports are missing connections or unused.
//
// The end product is a fully connected set of DAGs with the only terminal blocks being stubs (defined in Go) and literal blocks.
type tracer struct {
	loadedMods  map[string]*ast.Module
	loadedFiles map[string]*token.File

	tracingMod  *ast.Module
	tracingFile *token.File

	errors token.ErrorList
}

func (t *tracer) error(pos token.Pos, err error) {
	epos := t.tracingFile.Position(pos)

	// If AllErrors is not set, discard errors reported on the same line
	// as the last recorded error and stop parsing if there are more than
	// 10 errors.
	n := len(t.errors)
	if n > 0 && t.errors[n-1].Pos.Line == epos.Line {
		return // discard - likely a spurious error
	}
	if n > 10 {
		panic(bailout{})
	}

	t.errors.Add(epos, err)
}

func (t *tracer) expectedError(node ast.Node, msg string) {
	t.error(node.Pos(), fmt.Errorf("expected %v, got %v", node, msg))
}

func newTracer() *tracer {
	return &tracer{
		loadedMods:  make(map[string]*ast.Module),
		loadedFiles: make(map[string]*token.File),
	}
}

func (t *tracer) traceModule(file *token.File, mod *ast.Module) {
	t.loadedMods[mod.Name.V] = mod
	t.loadedFiles[mod.Name.V] = file

	t.tracingMod = mod
	t.tracingFile = file
	for _, pipe := range mod.DefinedPipes {
		pipe.BodyDAG = t.tracePipe(pipe)
	}
}

type pipeTrace struct {
	Graph       ast.Graph
	symbolTable map[string]output
}

func (t *tracer) tracePipe(pipe *ast.PipeDecl) *ast.Graph {
	trace := pipeTrace{Graph: ast.Graph{}, symbolTable: make(map[string]output)}
	for _, stmt := range pipe.Body {
		t.traceStmt(stmt, &trace)
	}
	return &trace.Graph
}

func (t *tracer) traceStmt(stmt ast.Stmt, state *pipeTrace) {
	switch st := stmt.(type) {
	case *ast.ExprStmt:
		t.traceExpr(st.X, state)
	}
}

func (t *tracer) traceExpr(expr ast.Expr, state *pipeTrace) output {
	switch x := expr.(type) {
	case *ast.AssignExpr:
		return t.traceAssign(x, state)
	case *ast.CallExpr:
		return t.traceCall(x, state)
	default:
		return NilOutput
	}
}

func (t *tracer) traceCall(call *ast.CallExpr, state *pipeTrace) output {
	if call.Name.Local() {
		decl := t.tracingMod.Lookup(call.Name.V)

	}
}

func (t *tracer) traceAssign(assign *ast.AssignExpr, state *pipeTrace) output {
	// "Assignments" are more pattern matching than actual variables storing state. We destructure the right hand side
	// according to the pattern on the left hand side. Unbinded symbol names are bound to the ports/outputs on the right hand side.
	rhs := t.traceExpr(assign.Rhs, state)
	t.unifyExpr(assign.Lhs, rhs, state)
	return rhs
}

func (t *tracer) unifyExpr(pattern ast.Expr, rhs output, state *pipeTrace) {
	switch p := pattern.(type) {
	case *ast.Ident:
		// 1. An identifier a = b()
		t.unifyOne(p, rhs, state)
	case *ast.FieldList:
		// 2. A map a: b, c: d = v()
		t.unifyBundle(p, rhs, state)
	}
	t.expectedError(pattern, "variable name or map of variables")
}

func (t *tracer) unifyOne(pattern *ast.Ident, rhs output, state *pipeTrace) {
	varName := pattern.V
	state.symbolTable[varName] = rhs
}

func (t *tracer) unifyBundle(pattern *ast.FieldList, rhs output, state *pipeTrace) error {
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
