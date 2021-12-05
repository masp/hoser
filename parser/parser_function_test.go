package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/masp/hoser/lexer"
)

func Test_parserState_parseFunction(t *testing.T) {
	tests := []struct {
		name    string
		program string
		want    *Block
		wantErr bool
	}{
		{"1 input 0 outputs", "main(v: int) {}", &Block{
			Name: &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "main", Line: 1, Col: 1}},
			Inputs: &Map{
				StartToken: lexer.Token{Kind: lexer.Ident, Value: "v", Line: 1, Col: 6},
				Entries: []Entry{
					{
						Key: Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "v", Line: 1, Col: 6}},
						Val: &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "int", Line: 1, Col: 9}},
					},
				},
			},
			Outputs: &Map{StartToken: lexer.Token{Kind: lexer.LCurlyBrack, Value: "{", Line: 1, Col: 14}},
			Body:    EmptyFnBody,
		}, false},
		{"2 inputs 1 output", "main(a: int, b: int) (o: int) {}", &Block{
			Name: &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "main", Line: 1, Col: 1}},
			Inputs: &Map{
				StartToken: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 6},
				Entries: []Entry{
					{
						Key: Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 6}},
						Val: &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "int", Line: 1, Col: 9}},
					},
					{
						Key: Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "b", Line: 1, Col: 14}},
						Val: &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "int", Line: 1, Col: 17}},
					},
				},
			},
			Outputs: &Map{
				StartToken: lexer.Token{Kind: lexer.Ident, Value: "o", Line: 1, Col: 23},
				Entries: []Entry{
					{
						Key: Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "o", Line: 1, Col: 23}},
						Val: &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "int", Line: 1, Col: 26}},
					},
				},
			},
			Body: EmptyFnBody,
		}, false},
		{"Many statements", "main() {a = b; c = d}", &Block{
			Name:    &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "main", Line: 1, Col: 1}},
			Inputs:  &Map{StartToken: lexer.Token{Kind: lexer.RParen, Value: ")", Line: 1, Col: 6}},
			Outputs: &Map{StartToken: lexer.Token{Kind: lexer.LCurlyBrack, Value: "{", Line: 1, Col: 8}},
			Body: []Expression{
				&AssignmentExpr{
					Left:  &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 9}},
					Right: &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "b", Line: 1, Col: 13}},
				},
				&AssignmentExpr{
					Left:  &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "c", Line: 1, Col: 16}},
					Right: &Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "d", Line: 1, Col: 20}},
				},
			},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newTestParser(tt.program)
			got, err := s.parseFunction()
			if (err != nil) != tt.wantErr {
				t.Errorf("parserState.parseFunction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parserState.parseFunction() = %v, want %v", got, tt.want)
				fmt.Printf("got token tree:\n%v\n", PrintTokenAst(got, 0))
				fmt.Printf("want token tree:\n%v\n", PrintTokenAst(tt.want, 0))
			}
		})
	}
}
