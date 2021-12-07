package lexer

import (
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
		{"All whitespace", args{"	 \r\n\r\n"}, []token.Token{}, false},
		{"Standard identifier", args{"A5_t"}, []token.Token{
			token.Ident,
		}, false},
		{"Invalid identifier", args{"\\5A"}, []token.Token{}, true},
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
		{"Int Type", "int", token.IntType},
		{"String Type", "string", token.StringType},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src := []byte(tt.tokenText)
			file := token.NewFile("<unknown>", len(src))
			got, err := ScanAll(&file, src)
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ScanAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
