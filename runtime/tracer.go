package runtime

// node is a call like grep.Filter which represents either another pipe definition or a process to be executed
type node struct {
}

// tracer will go through a pipe definition and create node's which represent processes to execute and
// edges which represent the connections between processes.
type tracer struct {
}
