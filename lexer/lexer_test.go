package lexer

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/masp/hoser/token"
)

func TestReadAll(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name    string
		args    args
		want    []token.Token
		wantErr bool
	}{
		{"All whitespace", args{"	 \r\n\r\n"}, nil, false},
		{"Standard identifier", args{"A5_t"}, []token.Token{
			token.Ident,
		}, false},
		{"Invalid identifier", args{"\\5A"}, nil, true},
		{"Columns & Lines Correct", args{"t1 t2\nt3"}, []token.Token{
			token.Ident,
			token.Ident,
			token.Semicolon,
			token.Ident,
		}, false},
		{"Operators supported", args{"{}()=;:,"}, []token.Token{
			token.LCurlyBrack,
			token.RCurlyBrack,
			token.LParen,
			token.RParen,
			token.Equals,
			token.Semicolon,
			token.Colon,
			token.Comma,
		}, false},
		{"Semicolon inserts", args{"}\n)\nA\n"}, []token.Token{
			token.RCurlyBrack,
			token.Semicolon,
			token.RParen,
			token.Semicolon,
			token.Ident,
			token.Semicolon,
		}, false},
		{"Semicolon inserts", args{"a=b#NO!\n"}, []token.Token{
			token.Ident,
			token.Equals,
			token.Ident,
			token.Comment,
			token.Semicolon,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := []byte(tt.args.text)
			file := token.NewFile("<unknown>", len(src))
			got, err := ScanAll(&file, src)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScanAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ScanAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokens(t *testing.T) {
	tests := []struct {
		name      string
		tokenText string
		want      token.Token
	}{
		{"Integer", "13509185", token.Integer},
		{"Float", "1.2", token.Float},
		{"Float Exponential", "1.2e10", token.Float},
		{"String", `"hello\n\"there"`, token.String},
		{"Return", "return", token.Return},
		{"Module", "module", token.Module},
		{"Period", ".", token.Period},
		{"Comments", "# This is # a comment;", token.Comment},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := []byte(tt.tokenText)
			file := token.NewFile("<unknown>", len(src))
			got, err := ScanAll(&file, src)
			if err != nil {
				t.Error(err)
				return
			}

			if got[0] != tt.want {
				t.Errorf("ScanAll() = %v, want %v", got[0], tt.want)
			}
		})
	}
}

func TestTokensPos(t *testing.T) {
	type result struct {
		pos token.Pos
		tok token.Token
		lit string
	}

	tests := []struct {
		src  string
		want []result
	}{
		{
			"a b c",
			[]result{
				{token.Pos(1), token.Ident, "a"},
				{token.Pos(3), token.Ident, "b"},
				{token.Pos(5), token.Ident, "c"},
			},
		},
		{
			"12.5  5\n7",
			[]result{
				{token.Pos(1), token.Float, "12.5"},
				{token.Pos(7), token.Integer, "5"},
				{token.Pos(9), token.Integer, "7"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q", tt.src), func(t *testing.T) {
			src := []byte(tt.src)
			file := token.NewFile("<test>", len(src))
			scanner := NewScanner(&file, src)

			for _, want := range tt.want {
				pos, tok, lit := scanner.Next()
				if want.pos != pos || want.tok != tok || want.lit != lit {
					t.Errorf("ScanAll() = {pos: %v, tok: %v, lit: %v}, want {pos: %v, tok: %v, lit: %v}", pos, tok, lit, want.pos, want.tok, want.lit)
				}
			}
		})
	}
}

func TestTokensPositions(t *testing.T) {
	tests := []struct {
		src  string
		pos  []token.Pos
		want []token.Position
	}{
		{
			"a\nb\nc",
			[]token.Pos{
				token.Pos(1),
				token.Pos(3),
				token.Pos(5),
			},
			[]token.Position{
				{
					Filename: "<test>",
					Offset:   token.Pos(1),
					Line:     1,
					Column:   1,
				},
				{
					Filename: "<test>",
					Offset:   token.Pos(3),
					Line:     2,
					Column:   1,
				},
				{
					Filename: "<test>",
					Offset:   token.Pos(5),
					Line:     3,
					Column:   1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q", tt.src), func(t *testing.T) {
			src := []byte(tt.src)
			file := token.NewFile("<test>", len(src))
			_, err := ScanAll(&file, src)
			if err != nil {
				t.Error(err)
			}

			for i, pos := range tt.pos {
				if got := file.Position(pos); got != tt.want[i] {
					t.Errorf("ScanAll() = %v, want %v", got, tt.want[i])
				}
			}

		})
	}
}
