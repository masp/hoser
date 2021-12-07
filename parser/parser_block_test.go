package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/masp/hoser/ast"
	"github.com/masp/hoser/lexer"
)

func Test_parserState_parseFunction(t *testing.T) {
	tests := []struct {
		name    string
		program string
		want    *ast.Block
		wantErr bool
	}{
		{"Empty body", "main(v: int) (v: int)\nf2()", &ast.Block{
			Name: &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "main", Line: 1, Col: 1}},
			Inputs: &ast.Map{
				StartToken: lexer.Token{Kind: lexer.Ident, Value: "v", Line: 1, Col: 6},
				Entries: []ast.Entry{
					{
						Key: ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "v", Line: 1, Col: 6}},
						Val: &ast.Type{Token: lexer.Token{Kind: lexer.IntType, Value: "int", Line: 1, Col: 9}},
					},
				},
			},
			Outputs: &ast.Map{
				StartToken: lexer.Token{Kind: lexer.Ident, Value: "v", Line: 1, Col: 15},
				Entries: []ast.Entry{
					{
						Key: ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "v", Line: 1, Col: 15}},
						Val: &ast.Type{Token: lexer.Token{Kind: lexer.IntType, Value: "int", Line: 1, Col: 18}},
					},
				},
			},
			Body: nil,
		}, false},
		{"1 input 0 outputs", "main(v: int) {}", &ast.Block{
			Name: &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "main", Line: 1, Col: 1}},
			Inputs: &ast.Map{
				StartToken: lexer.Token{Kind: lexer.Ident, Value: "v", Line: 1, Col: 6},
				Entries: []ast.Entry{
					{
						Key: ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "v", Line: 1, Col: 6}},
						Val: &ast.Type{Token: lexer.Token{Kind: lexer.IntType, Value: "int", Line: 1, Col: 9}},
					},
				},
			},
			Outputs: &ast.Map{StartToken: lexer.Token{Kind: lexer.LCurlyBrack, Value: "{", Line: 1, Col: 14}},
			Body:    EmptyFnBody,
		}, false},
		{"2 inputs 1 output", "main(a: int, b: int) (o: int) {}", &ast.Block{
			Name: &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "main", Line: 1, Col: 1}},
			Inputs: &ast.Map{
				StartToken: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 6},
				Entries: []ast.Entry{
					{
						Key: ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 6}},
						Val: &ast.Type{Token: lexer.Token{Kind: lexer.IntType, Value: "int", Line: 1, Col: 9}},
					},
					{
						Key: ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "b", Line: 1, Col: 14}},
						Val: &ast.Type{Token: lexer.Token{Kind: lexer.IntType, Value: "int", Line: 1, Col: 17}},
					},
				},
			},
			Outputs: &ast.Map{
				StartToken: lexer.Token{Kind: lexer.Ident, Value: "o", Line: 1, Col: 23},
				Entries: []ast.Entry{
					{
						Key: ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "o", Line: 1, Col: 23}},
						Val: &ast.Type{Token: lexer.Token{Kind: lexer.IntType, Value: "int", Line: 1, Col: 26}},
					},
				},
			},
			Body: EmptyFnBody,
		}, false},
		{"Many statements", "main() {a = b; c = d}", &ast.Block{
			Name:    &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "main", Line: 1, Col: 1}},
			Inputs:  &ast.Map{StartToken: lexer.Token{Kind: lexer.RParen, Value: ")", Line: 1, Col: 6}},
			Outputs: &ast.Map{StartToken: lexer.Token{Kind: lexer.LCurlyBrack, Value: "{", Line: 1, Col: 8}},
			Body: []ast.Expression{
				&ast.AssignmentExpr{
					Left:  &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "a", Line: 1, Col: 9}},
					Right: &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "b", Line: 1, Col: 13}},
				},
				&ast.AssignmentExpr{
					Left:  &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "c", Line: 1, Col: 16}},
					Right: &ast.Identifier{Token: lexer.Token{Kind: lexer.Ident, Value: "d", Line: 1, Col: 20}},
				},
			},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseBlock(tt.program)
			if (err != nil) != tt.wantErr {
				t.Errorf("parserState.parseFunction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parserState.parseFunction() = %v, want %v", got, tt.want)
				fmt.Printf("got token tree:\n%v\n", ast.PrintTokenAst(got, 0))
				fmt.Printf("want token tree:\n%v\n", ast.PrintTokenAst(tt.want, 0))
			}
		})
	}
}
