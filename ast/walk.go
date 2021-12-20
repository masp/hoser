package ast

type Visitor func(Node) bool

// Walk does a depth first traverse over all the nodes in the AST calling visitor for each node.
// If visitor returns false, the walk stops, otherwise it continues.
func Walk(node Node, v Visitor) {
	if ok := v(node); !ok {
		return
	}

	switch n := node.(type) {
	case *Module:
		Walk(n.Name, v)
		for _, block := range n.DefinedPipes {
			Walk(block, v)
		}
	case *PipeDecl:
		Walk(n.Name, v)
		Walk(n.Inputs, v)
		Walk(n.Outputs, v)
		for _, stmt := range n.Body {
			Walk(stmt, v)
		}
	case *Field:
		Walk(n.Key, v)
		Walk(n.Value, v)
	case *FieldList:
		for _, field := range n.Fields {
			Walk(field, v)
		}
	case *CallExpr:
		Walk(n.Name, v)
		for _, arg := range n.Args {
			Walk(arg, v)
		}
	case *ParenExpr:
		Walk(n.X, v)
	case *ExprStmt:
		Walk(n.X, v)
	case *AssignExpr:
		Walk(n.Lhs, v)
		Walk(n.Rhs, v)
	case *Ident, *LiteralExpr:
	default:
	}
	v(nil)
}
