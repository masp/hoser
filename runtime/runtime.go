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
	}
}

func (rt *State) Lookup(ident *ast.Ident) NativeProc {
	return rt.NativeProcs[ident.FullName()]
}

func (rt *State) Register(module string, name string, proc NativeProc) {
	rt.NativeProcs[module+"."+name] = proc
}

func (rt *State) Push(arg interface{}) {
	rt.Args = append(rt.Args, arg)
}

func (rt *State) ClearArgs() {
	rt.Args = rt.Args[0:]
}

func (rt *State) RunProgram(program []byte) error {
	file := token.NewFile("", len(program))
	module, err := parser.ParseModule(&file, program)
	if err != nil {
		return err
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
						argv, err := rt.evalExpr(arg)
						if err != nil {
							return err
						}
						rt.Push(argv)
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

func (rt *State) evalExpr(x ast.Expr) (interface{}, error) {
	switch v := x.(type) {
	case *ast.LiteralExpr:
		return v.ParsedVal, nil
	default:
		return nil, fmt.Errorf("unsupported expression found in args %v", x.Pos())
	}
}

func findMainBlock(module *ast.Module) *ast.BlockDecl {
	for _, block := range module.Blocks {
		if block.Name.Name == "main" {
			return block
		}
	}
	return nil
}
