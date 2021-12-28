package tracer

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/masp/hoser/ast"
	"github.com/masp/hoser/token"
)

func encodeBlock(block ast.Block) (value string) {
	switch b := block.(type) {
	case *ast.PipeBlock:
		value = b.Decl.BlockName()
	case *ast.LiteralBlock:
		value = b.Lit.Value
	case *ast.StubBlock:
		value = b.Decl.BlockName() + "*"
	default:
		panic(fmt.Errorf("invalid block type: %T", block))
	}
	return
}

func encodeBlocks(blocks []ast.Block) (result []string) {
	for _, block := range blocks {
		result = append(result, encodeBlock(block))
	}
	return
}

func encodeEdge(edge ast.Edge, graph *ast.Graph) (value string) {
	return fmt.Sprintf("%s[%d]->%s[%d]",
		encodeBlock(graph.Blocks[edge.Src.Block]), edge.Src.Port,
		encodeBlock(graph.Blocks[edge.Dst.Block]), edge.Dst.Port)
}

func encodeEdges(graph *ast.Graph) (result []string) {
	for _, edge := range graph.Edges {
		result = append(result, encodeEdge(edge, graph))
	}
	return
}

func Test_TracePipe(t *testing.T) {
	tests := []struct {
		name       string
		src        string
		wantBlocks []string
		wantEdges  []string
	}{
		{
			"Single call",
			`
module "a"
B(a: int) {}
main() {
	B(a: 10)
}
`,
			[]string{"10", "B"},
			[]string{"10[0]->B[0]"},
		},
		{
			"Single call with single arg",
			`
module "a"
B(a: int) {}
main() {
	B(10)
}
`,
			[]string{"10", "B"},
			[]string{"10[0]->B[0]"},
		},
		{
			"Multiarg",
			`
module "a"
B(a: int, b: int) {}
main() {
	B(10, 12)
}
`,
			[]string{"10", "12", "B"},
			[]string{"10[0]->B[0]", "12[0]->B[1]"},
		},
		{
			"Nested",
			`
module "a"
B(a: int, b: int) {}
C() (c: int) {}
main() {
	B(10, C())
}
`,
			[]string{"10", "C", "B"},
			[]string{"10[0]->B[0]", "C[0]->B[1]"},
		},
		{
			"Symbols",
			`
module "a"
B(a: int, b: int) {}
C() (c: int) {}
main() {
	c = C()
	B(10, c)
}
`,
			[]string{"C", "10", "B"},
			[]string{"10[0]->B[0]", "C[0]->B[1]"},
		},
		{
			"Unify multiresult",
			`
module "a"
B(a: int, b: int) {}
C() (c1: int, c2: int) {}
main() {
	{c1: c1, c2: c2} = C()
	B(a: c2, b: c2)
}
`,
			[]string{"C", "B"},
			[]string{"C[1]->B[0]", "C[1]->B[1]"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := token.NewFile("", len(tt.src))
			tr := NewTracer()
			module, err := tr.TraceModule(&file, []byte(tt.src))
			if err != nil {
				t.Fatal(err)
			}

			if mainBlock := module.Lookup("main"); mainBlock != nil {
				if got, ok := mainBlock.(*ast.PipeDecl); ok {
					gotBlocks := encodeBlocks(got.BodyDAG.Blocks)
					if !reflect.DeepEqual(gotBlocks, tt.wantBlocks) {
						t.Errorf("got blocks %v, want %v", gotBlocks, tt.wantBlocks)
					}

					gotEdges := encodeEdges(got.BodyDAG)
					if !reflect.DeepEqual(gotEdges, tt.wantEdges) {
						t.Errorf("got edges %v, want %v", gotEdges, tt.wantEdges)
					}
				} else {
					t.Errorf("main is not a pipe definition")
				}
			} else {
				t.Error("missing 'main' definition")
			}
		})
	}
}

func Test_TracePipeFail(t *testing.T) {
	tests := []struct {
		name string
		src  string
	}{
		{
			"Mismatch type",
			`
module "a"
B(a: int) {}
main() { B(a: "test") }
`,
		},
		{
			"Missing declaration",
			`
module "a"
main() { B() }
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := token.NewFile("", len(tt.src))
			tr := NewTracer()
			_, err := tr.TraceModule(&file, []byte(tt.src))
			if err == nil {
				t.Errorf("expected TraceModule() err = %v", err)
			}
		})
	}
}
