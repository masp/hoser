package parser

// tracer will traverse a set of modules and try to resolve for a given module a fully connected DAG
//
// For example:
// A() (v: int) { v = 10 }
// main() {
// 	shell.Run(A())
// }
//
// tracer needs to resolve the symbols Run() and A() and create DAGs representing their bodies. shell.Run is defined in another module (shell)
// and so the tracer needs to find the file in the include path, parse and recursively trace its body, and then return the graph here.
//
// tracer also needs to verify that edges between ports are never incorrectly typed or that ports are missing connections or unused.
//
// The end product is a fully connected DAG with the only terminal blocks being stubs (defined in Go) and literal blocks.
type resolver struct {
}
