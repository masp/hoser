package draw

import (
	"errors"
	"fmt"
	"strings"

	"github.com/masp/hoser/parser"
)

// drawer will take in a module in hoser and convert it into a dot format file that can be rendered
// using GraphViz software.

var (
	ErrNoRootBlocks = errors.New("No root blocks found (blocks with no inputs or outputs")
)

func Module(m *parser.Module) (string, error) {
	var sb strings.Builder
	sb.WriteString("digraph G {")
	roots := m.RootBlocks()
	if len(roots) == 0 {
		return "", ErrNoRootBlocks
	}

	for _, root := range roots {
		if err := rootBlock(root, &sb); err != nil {
			return "", err
		}
	}
	sb.WriteString("}")
	return sb.String(), nil
}

func rootBlock(root *parser.Block, sb *strings.Builder) error {
	sb.WriteString(fmt.Sprintf("subgraph %s {", root.Name))

	sb.WriteString(fmt.Sprintf(`%s [label=<
	<TABLE BORDER="0" CELLBORDER="1" CELLSPACING="0" CELLPADDING="4">
	  <TR>
		<TD ROWSPAN="3">hello<BR/>world</TD>
		<TD COLSPAN="3">b</TD>
		<TD ROWSPAN="3">g</TD>
		<TD ROWSPAN="3">h</TD>
	  </TR>
	  <TR>
		<TD>c</TD><TD PORT="here">d</TD><TD>e</TD>
	  </TR>
	  <TR>
		<TD COLSPAN="3">f</TD>
	  </TR>
	</TABLE>>];`, ))
	sb.WriteString("}")
}
