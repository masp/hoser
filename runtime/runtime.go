package runtime

import (
	"fmt"

	"github.com/masp/hoser/ast"
	"github.com/masp/hoser/parser"
	"github.com/masp/hoser/token"
)

type NativeProc func(state *State)

type State struct {
	NativeProcs map[string]NativeProc
	Args        []interface{}
}

func New() *State {
	return &State{
		NativeProcs: make(map[string]NativeProc),
		Pipes:       make(map[string]ast.PipeDecl),
	}
}

func (rt *State) Lookup(ident *ast.Ident) NativeProc {
	if fn, ok := rt.NativeProcs[ident.FullName()]; ok {
		return fn
	}
}

func (rt *State) RegisterProc(module string, name string, proc NativeProc) {
	if module != "" {
		rt.NativeProcs[module+"."+name] = proc
	} else {
		rt.NativeProcs[name] = proc
	}
}

func (rt *State) RegisterPipe(module string, name string, decl *ast.PipeDecl) {
	if module != "" {
		rt.Pipes[module+"."+name] = decl
	} else {
		rt.Pipes[name] = decl
	}
}

func (rt *State) Push(arg interface{}) {
	rt.Args = append(rt.Args, arg)
}

func (rt *State) ArgInt(idx int) int64     { return rt.Args[idx].(int64) }
func (rt *State) ArgFloat(idx int) float64 { return rt.Args[idx].(float64) }
func (rt *State) ArgString(idx int) string { return rt.Args[idx].(string) }

func (rt *State) ClearArgs() {
	rt.Args = rt.Args[0:]
}

func (rt *State) RunProgram(program []byte) error {
	file := token.NewFile("", len(program))
	module, err := parser.ParseModule(&file, program)
	if err != nil {
		return err
	}

	for _, block := range module.DefinedPipes {
		rt.RegisterPipe("", block.Name.Name, block)
	}

	return rt.Run(module)
}

func (rt *State) Run(module *ast.Module) error {
	mainBlock := findMainBlock(module)
	if mainBlock == nil {
		return fmt.Errorf("missing 'main' pipe in module")
	}

	for _, stmt := range mainBlock.Body {
		switch st := stmt.(type) {
		case *ast.ExprStmt:
			switch x := st.X.(type) {
			case *ast.CallExpr:
				if proc := rt.Lookup(x.Name); proc != nil {
					for _, arg := range x.Args {
						err := rt.evalArg(arg)
						if err != nil {
							return err
						}
					}
					proc(rt)

					// Clear args for next call
					rt.ClearArgs()
				} else {
					return fmt.Errorf("no proc with name %v found", x.Name.FullName())
				}
			}
		}
	}

	return nil
}

func (rt *State) evalArg(x ast.Expr) error {
	switch v := x.(type) {
	case *ast.LiteralExpr:
		rt.Push(v.ParsedVal)
		return nil
	case *ast.Field:
		err := rt.evalArg(v.Value)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unsupported expression found in args")
	}
}

func findMainBlock(module *ast.Module) *ast.PipeDecl {
	for _, block := range module.DefinedPipes {
		if block.Name.Name == "main" {
			return block
		}
	}
	return nil
}
