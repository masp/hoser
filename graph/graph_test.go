package graph

import (
	"reflect"
	"testing"

	"github.com/masp/hoser/lexer"
	"github.com/masp/hoser/parser"
)

func TestFromBlock(t *testing.T) {
	type args struct {
		moduleSrc string
	}
	tests := []struct {
		name    string
		args    args
		want    *Module
		wantErr bool
	}{
		{
			"Two blocks connected",
			args{
				moduleSrc: `
A() (A1: Int, A2: Int) {}
B(v: Int, v2: Int) (b: int) {}
C() {}
D(e: Int) {}
E(EI1: Int) (E1: Int) {}

main() {
	A1: a, A2: b = A()
	C()
	D(e: B(v: b, v2: E(EI1: a)))
}`,
			},
			&Module{
				Definitions: map[string]Definition{
					"main": {
						Nodes: []Node{BlockNode{"A"}, BlockNode{"C"}, BlockNode{"E"}, BlockNode{"B"}, BlockNode{"D"}},
						Edges: []Edge{
							{2, 0, 0, 0},
							{3, 0, 0, 1},
							{3, 2, 1, 0},
							{4, 3, 0, 0},
						},
					},
					"A": {},
					"B": {},
					"C": {},
					"D": {},
					"E": {},
				},
			},
			false,
		},
		{
			"Constant passed",
			args{
				moduleSrc: `
B(v: Int, t: Int) {}

main() {
	B(v: 10.2, t: "Hello")
}`,
			},
			&Module{
				Definitions: map[string]Definition{
					"B": {},
					"main": {
						Nodes: []Node{
							ConstNode{Value: &parser.Float{Token: lexer.Token{Kind: lexer.Float, Value: "10.2", Line: 5, Col: 7}, Value: 10.2}},
							ConstNode{Value: &parser.String{Token: lexer.Token{Kind: lexer.String, Value: "Hello", Line: 5, Col: 16}}},
							BlockNode{"B"},
						},
						Edges: []Edge{
							{2, 0, 0, 0},
							{2, 1, 1, 0},
						},
					},
				},
			},
			false,
		},
		{
			"Inputs",
			args{
				moduleSrc: `
B(v: int) {}

A(a: int) {
	B(v: a)
}`,
			},
			&Module{
				Definitions: map[string]Definition{
					"B": {},
					"A": {
						Nodes: []Node{
							BlockNode{"B"},
						},
						Edges: []Edge{
							{0, -1, 0, 0},
						},
					},
				},
			},
			false,
		},
		{
			"Outputs",
			args{
				moduleSrc: `
B(v: int) (v: int) {}

A(a: int) (t: Int) {
	return B(v: a)
}`,
			},
			&Module{
				Definitions: map[string]Definition{
					"B": {},
					"A": {
						Nodes: []Node{
							BlockNode{"B"},
						},
						Edges: []Edge{
							{0, -1, 0, 0},
							{-1, 0, 0, 0},
						},
					},
				},
			},
			false,
		},
		{
			"Multiple Outputs",
			args{
				moduleSrc: `
B() (v: int) {}

A(a: int) (o: int, o2: int) {
	return o: B(), o2: 10.2
}`,
			},
			&Module{
				Definitions: map[string]Definition{
					"B": {},
					"A": {
						Nodes: []Node{
							BlockNode{"B"},
							ConstNode{Value: &parser.Float{Token: lexer.Token{Kind: lexer.Float, Value: "10.2", Line: 5, Col: 21}, Value: 10.2}},
						},
						Edges: []Edge{
							{-1, 0, 0, 0},
							{-1, 1, 1, 0},
						},
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, errCh := lexer.Scan(tt.args.moduleSrc)
			module, err := parser.Scan(tokens, errCh)
			if err != nil {
				t.Error(err)
				return
			}
			got, err := TraceModule(module)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromModule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromModule() = %v, want %v", got, tt.want)
			}
		})
	}
}
