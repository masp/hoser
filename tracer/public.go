package tracer

import (
	"github.com/masp/hoser/ast"
	"github.com/masp/hoser/parser"
	"github.com/masp/hoser/token"
)

func NewTracer() *Tracer {
	return &Tracer{
		modCache: ast.EmptyModuleSet(),
	}
}

// TraceModule will trace all modules in main and referenced by main in includeDir and create ast Graphs representing
// the data pipes defined in the main file. The graphs are not expanded and will use references to other Nodes rather than
// expanding it into a single graph which is done by the runtime.
//
// If an external module is referenced, it wil lbe found under includeDir (recursively), parsed, and traced recursively
// until only stubs and literal blocks remain.
func (t *Tracer) TraceModule(file *token.File, src []byte) (module *ast.Module, err error) {
	module, err = parser.ParseModule(file, src)
	if err != nil {
		return
	}
	defer t.handleErrors(&err)
	t.traceModule(file, module)
	return
}

// recover bailout panics where we have gotten too many errors
func (t *Tracer) handleErrors(err *error) {
	if e := recover(); e != nil {
		// resume same panic if it's not a bailout
		if _, ok := e.(bailout); !ok {
			panic(e)
		}
	}
	t.errors.Sort()
	*err = t.errors.Err()
}

type bailout struct{}

func (t *Tracer) error(pos token.Pos, err error) {
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
