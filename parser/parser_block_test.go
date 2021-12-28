package parser

import (
	"testing"

	"github.com/masp/hoser/token"
)

func TestParseModule(t *testing.T) {
	tests := []struct {
		name       string
		program    string
		wantModule string
		wantBlocks []string
	}{
		{"Empty blocks", `module "main"; pipe main(v: int, v2: int) (v: int) {}; stub f()`, "main", []string{"main", "f"}},
		{"Filled blocks", `module "main"; pipe main(v: int) (v: int, v2: int) {a = b}; stub f()`, "main", []string{"main", "f"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := token.NewFile("<test>", len(tt.program))
			got, err := ParseModule(&file, []byte(tt.program))
			if err != nil {
				t.Errorf("parserState.parseFunction() error = %v", err)
				return
			}

			if got.Name.Value != tt.wantModule {
				t.Errorf("ParseModule() ModuleName = %v, want %v", got.Name.Value, tt.wantModule)
			}

			for i, block := range got.DefinedBlocks {
				want := tt.wantBlocks[i]
				if want != block.BlockName() {
					t.Errorf("ParseModule() Block (%d) = %v, want %v", i, block.BlockName(), want)
				}
			}
		})
	}
}
