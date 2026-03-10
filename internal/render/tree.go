package render

import (
	"fmt"
	"strings"

	"github.com/cstsortan/prm/internal/service"
)

// Tree renders a hierarchy of tree nodes as an indented tree.
func Tree(nodes []*service.TreeNode) string {
	if len(nodes) == 0 {
		return Dim.Render("No epics found.")
	}

	var b strings.Builder
	for i, node := range nodes {
		renderNode(&b, node, "", i == len(nodes)-1)
	}
	return b.String()
}

func renderNode(b *strings.Builder, node *service.TreeNode, prefix string, last bool) {
	e := node.Entity

	connector := "├── "
	if last {
		connector = "└── "
	}
	if prefix == "" {
		connector = ""
	}

	status := fmt.Sprintf("[%s]", StatusStyle(e.Status))
	typ := TypeStyle(e.Type)
	title := e.Title

	b.WriteString(fmt.Sprintf("%s%s%s %s %s %s\n",
		prefix,
		connector,
		typ,
		title,
		status,
		IDStyle.Render(e.ShortID()),
	))

	childPrefix := prefix
	if prefix != "" {
		if last {
			childPrefix += "    "
		} else {
			childPrefix += "│   "
		}
	} else {
		childPrefix = "  "
	}

	for i, child := range node.Children {
		renderNode(b, child, childPrefix, i == len(node.Children)-1)
	}
}
