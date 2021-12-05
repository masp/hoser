package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/masp/hoser/lexer"
)

func TestScan(t *testing.T) {
	type args struct {
		program string
	}
	tests := []struct {
		name    string
		args    args
		want    *Module
		wantErr bool
	}{
		{
			"Single module single function",
			args{`main() () {}`},
			&Module{
				Blocks: map[string]*Block{
					"main": {
						Name:    &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "main", Line: 1, Col: 1}},
						Inputs:  &Map{StartToken: lexer.Token{Kind: lexer.RParen, Value: ")", Line: 1, Col: 6}},
						Outputs: &Map{StartToken: lexer.Token{Kind: lexer.RParen, Value: ")", Line: 1, Col: 9}},
						Body:    []Expression{},
					},
				},
			},
			false,
		},
		{
			"Single module multiple functions",
			args{"B1(a: b) {}\nB2() (v: d) {}"},
			&Module{
				Blocks: map[string]*Block{
					"B1": {
						Name: &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "B1", Line: 1, Col: 1}},
						Inputs: &Map{
							StartToken: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 4},
							Entries: []Entry{
								{
									Key: Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 4}},
									Val: &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "b", Line: 1, Col: 7}},
								},
							},
						},
						Outputs: &Map{StartToken: lexer.Token{Kind: lexer.LCurlyBrack, Value: "{", Line: 1, Col: 10}},
						Body:    []Expression{},
					},
					"B2": {
						Name:   &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "B2", Line: 2, Col: 1}},
						Inputs: &Map{StartToken: lexer.Token{Kind: lexer.RParen, Value: ")", Line: 2, Col: 4}},
						Outputs: &Map{
							StartToken: lexer.Token{Kind: lexer.Ident, Value: "v", Line: 2, Col: 7},
							Entries: []Entry{
								{
									Key: Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "v", Line: 2, Col: 7}},
									Val: &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "d", Line: 2, Col: 10}},
								},
							},
						},
						Body: []Expression{},
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
			tokens, errCh := lexer.Scan(tt.args.program)
			got, err := Scan(tokens, errCh)
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
		want    Expression
		wantErr bool
	}{
		{"Identifier", args{"a;"}, &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 1}}, false},
		{"String", args{"\"hello\";"}, &String{Token: lexer.Token{Kind: lexer.String, Value: "hello", Line: 1, Col: 1}}, false},
		{"Integer", args{"12;"}, &Integer{Value: 12, Token: lexer.Token{Kind: lexer.Integer, Value: "12", Line: 1, Col: 1}}, false},
		{"Float", args{"123.5;"}, &Float{Value: 123.5, Token: lexer.Token{Kind: lexer.Float, Value: "123.5", Line: 1, Col: 1}}, false},
		{"Return", args{"return 10;"}, &Return{
			Token: lexer.Token{Kind: lexer.Return, Value: "return", Line: 1, Col: 1},
			Value: &Integer{Value: 10, Token: lexer.Token{Kind: lexer.Integer, Value: "10", Line: 1, Col: 8}},
		}, false},
		{"Assignment", args{"a = b;"}, &AssignmentExpr{
			Left:  &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 1}},
			Right: &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "b", Line: 1, Col: 5}},
		}, false},
		{"Assignment To Expression", args{"a = B(x: yz, d: fg);"}, &AssignmentExpr{
			Left: &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 1}},
			Right: &BlockCall{
				Name: &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "B", Line: 1, Col: 5}},
				Args: &Map{StartToken: lexer.Token{Kind: lexer.Ident, Value: "x", Line: 1, Col: 7}, Entries: []Entry{
					{
						Key: Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "x", Line: 1, Col: 7}},
						Val: &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "yz", Line: 1, Col: 10}},
					},
					{
						Key: Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "d", Line: 1, Col: 14}},
						Val: &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "fg", Line: 1, Col: 17}},
					},
				}},
			},
		}, false},
		{"Single entry map", args{"a: b;"}, &Entry{
			Key: Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 1}},
			Val: &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "b", Line: 1, Col: 4}},
		}, false},
		{"Map", args{"a: bc, d: fg;"}, &Map{
			StartToken: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 1},
			Entries: []Entry{
				{
					Key: Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 1}},
					Val: &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "bc", Line: 1, Col: 4}},
				},
				{
					Key: Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "d", Line: 1, Col: 8}},
					Val: &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "fg", Line: 1, Col: 11}},
				},
			},
		}, false},
		{"Nested Map", args{"a: (x: yz), d: fg;"}, &Map{
			StartToken: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 1},
			Entries: []Entry{
				{
					Key: Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 1}},
					Val: &Entry{
						Key: Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "x", Line: 1, Col: 5}},
						Val: &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "yz", Line: 1, Col: 8}},
					},
				},
				{
					Key: Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "d", Line: 1, Col: 13}},
					Val: &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "fg", Line: 1, Col: 16}},
				},
			},
		}, false},
		{"Empty Block Call", args{"B();"}, &BlockCall{
			Name: &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "B", Line: 1, Col: 1}},
			Args: &Map{StartToken: lexer.Token{Kind: lexer.LParen, Value: "(", Line: 1, Col: 2}, Entries: nil},
		}, false},
		{"Nested Block Call", args{"D(e: B(v: b));"}, &BlockCall{
			Name: &Identifier{lexer.Token{Kind: lexer.Ident, Value: "D", Line: 1, Col: 1}},
			Args: &Map{
				StartToken: lexer.Token{Kind: lexer.Ident, Value: "e", Line: 1, Col: 3},
				Entries: []Entry{
					{
						Key: Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "e", Line: 1, Col: 3}},
						Val: &BlockCall{
							Name: &Identifier{lexer.Token{Kind: lexer.Ident, Value: "B", Line: 1, Col: 6}},
							Args: &Map{
								StartToken: lexer.Token{Kind: lexer.Ident, Value: "v", Line: 1, Col: 8},
								Entries: []Entry{
									{
										Key: Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "v", Line: 1, Col: 8}},
										Val: &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "b", Line: 1, Col: 11}},
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
				fmt.Printf("got token tree:\n%v\n", PrintTokenAst(got, 0))
				fmt.Printf("want token tree:\n%v\n", PrintTokenAst(tt.want, 0))
			}
		})
	}
}
