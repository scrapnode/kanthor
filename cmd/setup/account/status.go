package account

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/scrapnode/kanthor/domain/structure"
)

func status(tree *structure.Node[string], status map[string]bool) string {
	l := list.NewWriter()
	l.SetStyle(list.StyleConnectedRounded)

	l.AppendItems([]interface{}{tree.Value})
	// make this one to be simple by 3 lopps
	for _, app := range tree.Children {
		l.Indent()

		l.AppendItems([]interface{}{fmt.Sprintf("%s %s", icon(status[app.Value]), app.Value)})

		l.Indent()
		for _, ep := range app.Children {
			l.Indent()
			l.AppendItems([]interface{}{fmt.Sprintf("%s %s", icon(status[ep.Value]), ep.Value)})
			for _, epr := range ep.Children {
				l.Indent()
				l.AppendItems([]interface{}{fmt.Sprintf("%s %s", icon(status[epr.Value]), epr.Value)})
				l.UnIndent()
			}
			l.UnIndent()
		}
		l.UnIndent()

		l.UnIndent()
	}

	return l.Render()
}
