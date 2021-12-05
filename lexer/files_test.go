package lexer

import (
	"bytes"
	"reflect"
	"testing"
)

func TestMarshalTo(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
		wantErr bool
	}{
		{"One Token", args{"main"}, `{"Kind":"Ident","Value":"main","Line":1,"Col":1}
`, false},
		{"Many Tokens", args{"a b c"},
			`{"Kind":"Ident","Value":"a","Line":1,"Col":1}
{"Kind":"Ident","Value":"b","Line":1,"Col":3}
{"Kind":"Ident","Value":"c","Line":1,"Col":5}
`, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			tokens, err := ScanAll(tt.args.input)
			if err != nil {
				t.Errorf("ScanAll() error = %v", err)
				return
			}

			if err := MarshalTo(tokens, out); (err != nil) != tt.wantErr {
				t.Errorf("MarshalTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("MarshalTo() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func TestUnmarshalFrom(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name    string
		args    args
		want    []Token
		wantErr bool
	}{
		{
			"One Token",
			args{`{"Kind":"Ident","Value":"main","Line":1,"Col":1}
`},
			[]Token{{Kind: Ident, Value: "main", Line: 1, Col: 1}},
			false,
		},
		{
			"Many Tokens",
			args{`{"Kind":"Ident","Value":"a","Line":1,"Col":1}
{"Kind":"Ident","Value":"b","Line":1,"Col":3}
{"Kind":"Ident","Value":"c","Line":1,"Col":5}
`},
			[]Token{
				{Kind: Ident, Value: "a", Line: 1, Col: 1},
				{Kind: Ident, Value: "b", Line: 1, Col: 3},
				{Kind: Ident, Value: "c", Line: 1, Col: 5},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := bytes.NewBufferString(tt.args.in)
			got, err := UnmarshalFrom(buf)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalFrom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}
