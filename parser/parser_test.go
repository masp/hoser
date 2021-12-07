package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/masp/hoser/ast"
	"github.com/masp/hoser/lexer"
)

func TestScan(t *testing.T) {
	type args struct {
		program string
	}
	tests := []struct {
		name    string
		args    args
		want    *ast.Module
		wantErr bool
	}{
		{
			"Single module single function",
			args{`main() () {}`},
			&ast.Module{
				Blocks: map[string]*ast.Block{
					"main": {
						Name:    &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "main", Line: 1, Col: 1}},
						Inputs:  &ast.Map{StartToken: lexer.Token{Kind: lexer.RParen, Value: ")", Line: 1, Col: 6}},
						Outputs: &ast.Map{StartToken: lexer.Token{Kind: lexer.RParen, Value: ")", Line: 1, Col: 9}},
						Body:    []ast.Expression{},
					},
				},
			},
			false,
		},
		{
			"Single module multiple functions",
			args{"B1(a: b) {}\nB2() (v: d) {}"},
			&ast.Module{
				Blocks: map[string]*ast.Block{
					"B1": {
						Name: &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "B1", Line: 1, Col: 1}},
						Inputs: &ast.Map{
							StartToken: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 4},
							Entries: []ast.Entry{
								{
									Key: ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 4}},
									Val: &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "b", Line: 1, Col: 7}},
								},
							},
						},
						Outputs: &ast.Map{StartToken: lexer.Token{Kind: lexer.LCurlyBrack, Value: "{", Line: 1, Col: 10}},
						Body:    []ast.Expression{},
					},
					"B2": {
						Name:   &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "B2", Line: 2, Col: 1}},
						Inputs: &ast.Map{StartToken: lexer.Token{Kind: lexer.RParen, Value: ")", Line: 2, Col: 4}},
						Outputs: &ast.Map{
							StartToken: lexer.Token{Kind: lexer.Ident, Value: "v", Line: 2, Col: 7},
							Entries: []ast.Entry{
								{
									Key: ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "v", Line: 2, Col: 7}},
									Val: &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "d", Line: 2, Col: 10}},
								},
							},
						},
						Body: []ast.Expression{},
					},
				},
			},
			false,
		},
		{
			"Unexpected end of input",
			args{`main`},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseModule(tt.args.program)
			if (err != nil) != tt.wantErr {
				t.Errorf("Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func newTestParser(program string) *parserState {
	tokens, errCh := lexer.Scan(program)
	s := &parserState{
		tokens: tokens,
		errCh:  errCh,
	}
	return s
}

func Test_parserState_parseExpression(t *testing.T) {
	type args struct {
		program string
	}
	tests := []struct {
		name    string
		args    args
		want    ast.Expression
		wantErr bool
	}{
		{"Identifier", args{"a;"}, &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 1}}, false},
		{"String", args{"\"hello\";"}, &ast.String{Token: lexer.Token{Kind: lexer.String, Value: "hello", Line: 1, Col: 1}}, false},
		{"Integer", args{"12;"}, &ast.Integer{Value: 12, Token: lexer.Token{Kind: lexer.Integer, Value: "12", Line: 1, Col: 1}}, false},
		{"Float", args{"123.5;"}, &ast.Float{Value: 123.5, Token: lexer.Token{Kind: lexer.Float, Value: "123.5", Line: 1, Col: 1}}, false},
		{"Return", args{"return 10;"}, &ast.Return{
			Token: lexer.Token{Kind: lexer.Return, Value: "return", Line: 1, Col: 1},
			Value: &ast.Integer{Value: 10, Token: lexer.Token{Kind: lexer.Integer, Value: "10", Line: 1, Col: 8}},
		}, false},
		{"Assignment", args{"a = b;"}, &ast.AssignmentExpr{
			Left:  &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 1}},
			Right: &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "b", Line: 1, Col: 5}},
		}, false},
		{"Assignment To Expression", args{"a = B(x: yz, d: fg);"}, &ast.AssignmentExpr{
			Left: &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 1}},
			Right: &ast.BlockCall{
				Name: &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "B", Line: 1, Col: 5}},
				Args: &ast.Map{StartToken: lexer.Token{Kind: lexer.Ident, Value: "x", Line: 1, Col: 7}, Entries: []ast.Entry{
					{
						Key: ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "x", Line: 1, Col: 7}},
						Val: &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "yz", Line: 1, Col: 10}},
					},
					{
						Key: ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "d", Line: 1, Col: 14}},
						Val: &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "fg", Line: 1, Col: 17}},
					},
				}},
			},
		}, false},
		{"Single entry map", args{"a: b;"}, &ast.Entry{
			Key: ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 1}},
			Val: &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "b", Line: 1, Col: 4}},
		}, false},
		{"Map", args{"a: bc, d: fg;"}, &ast.Map{
			StartToken: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 1},
			Entries: []ast.Entry{
				{
					Key: ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 1}},
					Val: &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "bc", Line: 1, Col: 4}},
				},
				{
					Key: ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "d", Line: 1, Col: 8}},
					Val: &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "fg", Line: 1, Col: 11}},
				},
			},
		}, false},
		{"Nested Map", args{"a: (x: yz), d: fg;"}, &ast.Map{
			StartToken: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 1},
			Entries: []ast.Entry{
				{
					Key: ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 1}},
					Val: &ast.Entry{
						Key: ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "x", Line: 1, Col: 5}},
						Val: &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "yz", Line: 1, Col: 8}},
					},
				},
				{
					Key: ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "d", Line: 1, Col: 13}},
					Val: &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "fg", Line: 1, Col: 16}},
				},
			},
		}, false},
		{"Empty Block Call", args{"B();"}, &ast.BlockCall{
			Name: &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "B", Line: 1, Col: 1}},
			Args: &ast.Map{StartToken: lexer.Token{Kind: lexer.LParen, Value: "(", Line: 1, Col: 2}, Entries: nil},
		}, false},
		{"Nested Block Call", args{"D(e: B(v: b));"}, &ast.BlockCall{
			Name: &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "D", Line: 1, Col: 1}},
			Args: &ast.Map{
				StartToken: lexer.Token{Kind: lexer.Ident, Value: "e", Line: 1, Col: 3},
				Entries: []ast.Entry{
					{
						Key: ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "e", Line: 1, Col: 3}},
						Val: &ast.BlockCall{
							Name: &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "B", Line: 1, Col: 6}},
							Args: &ast.Map{
								StartToken: lexer.Token{Kind: lexer.Ident, Value: "v", Line: 1, Col: 8},
								Entries: []ast.Entry{
									{
										Key: ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "v", Line: 1, Col: 8}},
										Val: &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "b", Line: 1, Col: 11}},
									},
								},
							},
						},
					},
				},
			},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestParser(tt.args.program)
			got, err := s.parseExpression(lexer.Invalid)
			if (err != nil) != tt.wantErr {
				t.Errorf("parserState.parseExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parserState.parseExpression() = %v, want %v", got, tt.want)
				fmt.Printf("got token tree:\n%v\n", ast.PrintTokenAst(got, 0))
				fmt.Printf("want token tree:\n%v\n", ast.PrintTokenAst(tt.want, 0))
			}
		})
	}
}
