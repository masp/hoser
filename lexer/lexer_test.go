package lexer

import (
	"reflect"
	"testing"
)

func TestReadAll(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name    string
		args    args
		want    []Token
		wantErr bool
	}{
		{"All whitespace", args{"	 \r\n\r\n"}, []Token{}, false},
		{"Standard identifier", args{"A5_t"}, []Token{
			{Kind: Ident, Value: "A5_t", Line: 1, Col: 1},
		}, false},
		{"Invalid identifier", args{"\\5A"}, []Token{}, true},
		{"Columns & Lines Correct", args{"t1 t2\nt3"}, []Token{
			{Kind: Ident, Value: "t1", Line: 1, Col: 1},
			{Kind: Ident, Value: "t2", Line: 1, Col: 4},
			{Kind: Semicolon, Value: "", Line: 1, Col: 6},
			{Kind: Ident, Value: "t3", Line: 2, Col: 1},
		}, false},
		{"Operators supported", args{"{}()=;:,"}, []Token{
			{Kind: LCurlyBrack, Value: "{", Line: 1, Col: 1},
			{Kind: RCurlyBrack, Value: "}", Line: 1, Col: 2},
			{Kind: LParen, Value: "(", Line: 1, Col: 3},
			{Kind: RParen, Value: ")", Line: 1, Col: 4},
			{Kind: Equals, Value: "=", Line: 1, Col: 5},
			{Kind: Semicolon, Value: ";", Line: 1, Col: 6},
			{Kind: Colon, Value: ":", Line: 1, Col: 7},
			{Kind: Comma, Value: ",", Line: 1, Col: 8},
		}, false},
		{"Semicolon inserts", args{"}\n)\nA\n"}, []Token{
			{Kind: RCurlyBrack, Value: "}", Line: 1, Col: 1},
			{Kind: Semicolon, Value: "", Line: 1, Col: 2},
			{Kind: RParen, Value: ")", Line: 2, Col: 1},
			{Kind: Semicolon, Value: "", Line: 2, Col: 2},
			{Kind: Ident, Value: "A", Line: 3, Col: 1},
			{Kind: Semicolon, Value: "", Line: 3, Col: 2},
		}, false},
		{"Integer", args{"13509185"}, []Token{{Kind: Integer, Value: "13509185", Line: 1, Col: 1}}, false},
		{"Float", args{"1.2"}, []Token{{Kind: Float, Value: "1.2", Line: 1, Col: 1}}, false},
		{"Float Exponential", args{"1.2e10"}, []Token{{Kind: Float, Value: "1.2e10", Line: 1, Col: 1}}, false},
		{"String", args{`"hello\n\"there"`}, []Token{{Kind: String, Value: "hello\n\"there", Line: 1, Col: 1}}, false},
		{"Return", args{"return"}, []Token{{Kind: Return, Value: "return", Line: 1, Col: 1}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ScanAll(tt.args.text)
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
