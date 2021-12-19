package token

import (
	"testing"
)

func TestFiles(t *testing.T) {
	src := []byte("a\nb\nc")
	file := NewFile("<test>", len(src))
	tests := []struct {
		name string
		pos  Pos
		want Position
	}{
		{"a", Pos(1), Position{
			Filename: "<test>",
			Offset:   Pos(1),
			Line:     1,
			Column:   1,
		}},
		{"b", Pos(3), Position{
			Filename: "<test>",
			Offset:   Pos(3),
			Line:     2,
			Column:   1,
		}},
		{"c", Pos(5), Position{
			Filename: "<test>",
			Offset:   Pos(5),
			Line:     3,
			Column:   1,
		}},
	}

	file.AddLine(1)
	file.AddLine(3)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := file.Position(tt.pos); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
