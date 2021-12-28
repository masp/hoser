package parser

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/masp/hoser/ast"
	"github.com/masp/hoser/token"
)

func TestValid(t *testing.T) {
	tests := []struct {
		src string
	}{
		{`module "test"; import "a"; pipe main() () {}`},
		{`module "test"; pipe B1(a: b) {}; pipe B2() (v: d) {}`},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q", tt.src), func(t *testing.T) {
			file := token.NewFile("<test>", len(tt.src))
			_, err := ParseModule(&file, []byte(tt.src))
			if err != nil {
				t.Errorf("Scan() error = %v", err)
				return
			}
		})
	}
}

func Test_LiteralExpr(t *testing.T) {
	tests := []struct {
		input string
		want  token.Token
	}{
		{"\"hello\";", token.String},
		{"12;", token.Integer},
		{"123.5;", token.Float},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			file := token.NewFile("", len(tt.input))
			expr, err := ParseExpression(&file, []byte(tt.input))
			if err != nil {
				t.Errorf("ParseExpression() error = %v", err)
				return
			}

			if lit, ok := expr.(*ast.LiteralExpr); ok {
				if lit.Type != tt.want {
					t.Errorf("ParseExpression() literal = %v, want %v", lit.Type, tt.want)
				}
			} else {
				t.Errorf("ParseExpression() = %T, not LiteralExpr", expr)
			}
		})
	}
}

func Test_parseExpression(t *testing.T) {
	type args struct {
		program string
	}
	tests := []struct {
		name string
		args args
		want ast.Expr
	}{
		{"Assignment", args{"a = a.b;"}, &ast.AssignExpr{
			Lhs:   &ast.Ident{V: "a", NamePos: 1},
			Rhs:   &ast.Ident{V: "b", NamePos: 7, Module: "a", ModulePos: 5},
			EqPos: 3,
		}},
		{"Empty Call Expr", args{"a();"}, &ast.CallExpr{
			Name:   &ast.Ident{V: "a", NamePos: 1},
			Lparen: 2,
			Args:   nil,
			Rparen: 3,
		}},
		{"Call Expr 2 Args", args{"a(c,d,a:b);"}, &ast.CallExpr{
			Name:   &ast.Ident{V: "a", NamePos: 1},
			Lparen: 2,
			Args: []ast.Expr{
				&ast.Ident{V: "c", NamePos: 3},
				&ast.Ident{V: "d", NamePos: 5},
				&ast.Field{
					Key:   &ast.Ident{V: "a", NamePos: 7},
					Colon: 8,
					Value: &ast.Ident{V: "b", NamePos: 9},
				},
			},
			Rparen: 10,
		}},
		{"Assignment To Expression", args{"a = B();"}, &ast.AssignExpr{
			Lhs: &ast.Ident{V: "a", NamePos: 1},
			Rhs: &ast.CallExpr{
				Name:   &ast.Ident{V: "B", NamePos: 5},
				Lparen: 6,
				Args:   nil,
				Rparen: 7,
			},
			EqPos: token.Pos(3),
		}},
		{"Single unbound map", args{"a: b;"}, &ast.Field{
			Key:   &ast.Ident{V: "a", NamePos: 1},
			Colon: 2,
			Value: &ast.Ident{V: "b", NamePos: 4},
		}},
		{"Single map", args{"{a: b};"}, &ast.FieldList{
			Opener: 1,
			Fields: []*ast.Field{
				{
					Key:   &ast.Ident{V: "a", NamePos: 2},
					Colon: 3,
					Value: &ast.Ident{V: "b", NamePos: 5},
				},
			},
			Closer: 6,
		}},
		{"Multi map", args{"{a: b, c: d};"}, &ast.FieldList{
			Opener: 1,
			Fields: []*ast.Field{
				{
					Key:   &ast.Ident{V: "a", NamePos: 2},
					Colon: 3,
					Value: &ast.Ident{V: "b", NamePos: 5},
				},
				{
					Key:   &ast.Ident{V: "c", NamePos: 8},
					Colon: 9,
					Value: &ast.Ident{V: "d", NamePos: 11},
				},
			},
			Closer: 12,
		}},
		{"Nested map", args{"{a:{b:c},d:e};"}, &ast.FieldList{
			Opener: 1,
			Fields: []*ast.Field{
				{
					Key:   &ast.Ident{V: "a", NamePos: 2},
					Colon: 3,
					Value: &ast.FieldList{
						Opener: 4,
						Fields: []*ast.Field{
							{
								Key:   &ast.Ident{V: "b", NamePos: 5},
								Colon: 6,
								Value: &ast.Ident{V: "c", NamePos: 7},
							},
						},
						Closer: 8,
					},
				},
				{
					Key:   &ast.Ident{V: "d", NamePos: 10},
					Colon: 11,
					Value: &ast.Ident{V: "e", NamePos: 12},
				},
			},
			Closer: 13,
		}},
		{"Nested Call Expr", args{"a(b(d:e));"}, &ast.CallExpr{
			Name:   &ast.Ident{V: "a", NamePos: 1},
			Lparen: 2,
			Args: []ast.Expr{
				&ast.CallExpr{
					Name:   &ast.Ident{V: "b", NamePos: 3},
					Lparen: 4,
					Args: []ast.Expr{
						&ast.Field{
							Key:   &ast.Ident{V: "d", NamePos: 5},
							Colon: 6,
							Value: &ast.Ident{V: "e", NamePos: 7},
						},
					},
					Rparen: 8,
				},
			},
			Rparen: 9,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := token.NewFile("<test>", len(tt.args.program))
			got, err := ParseExpression(&file, []byte(tt.args.program))
			if err != nil {
				t.Errorf("parserState.parseExpression() error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parserState.parseExpression() mismatch")
				t.Errorf("want:\n")
				wantTree, err := ast.PrintString(&file, tt.want)
				if err != nil {
					t.Error(err)
				}
				t.Errorf("\n%s", wantTree)

				t.Errorf("got:\n")
				gotTree, err := ast.PrintString(&file, got)
				if err != nil {
					t.Error(err)
				}
				t.Errorf("\n%s", gotTree)
			}
		})
	}
}
