package runtime

import (
	"reflect"
	"testing"
)

const NotCalled = "NEVER_CALLED"

func TestState_Run(t *testing.T) {
	tests := []struct {
		name    string
		program string
		want    []interface{}
	}{
		{"never called", "module main; main() {}", []interface{}{NotCalled}},
		{"zero args", "module main; main() { test.Pass() }", nil},
		{"two args", `module main; main() { test.Pass(10, "hello") }`, []interface{}{int64(10), "hello"}},
		{"named args", `module main; main() { test.Pass(a: 10, b: "hello") }`, []interface{}{int64(10), "hello"}},
		{"subcall", `
module main
sub(v: int) { test.Pass(v) }
main() { sub(10) }
`, []interface{}{int64(10)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rt := New()
			got := []interface{}{NotCalled}
			rt.Register("test", "Pass", func(state *State) {
				got = nil
				got = append(got, state.Args...)
			})

			if err := rt.RunProgram([]byte(tt.program)); err != nil {
				t.Errorf("Run() error = %v", err)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() called Pass() with %v, want %v", got, tt.want)
			}
		})
	}
}
