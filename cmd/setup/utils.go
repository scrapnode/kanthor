package setup

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/segmentio/ksuid"
	"os"
	"time"
)

func id(prefix, payload string) string {
	t := time.Date(2022, 12, 01, 17, 01, 00, 00, time.UTC)
	var p [16]byte
	copy(p[:], payload)

	id, err := ksuid.FromParts(t, p[:])
	if err != nil {
		panic(err)
	}

	return prefix + "_" + id.String()
}

func showDemo(demo *demo) {
	t := table.NewWriter()
	t.AppendHeader(table.Row{"table", "id"})

	t.AppendRow([]interface{}{demo.Workspace.TableName(), demo.Workspace.Id})
	t.AppendRow([]interface{}{demo.WorkspaceTier.TableName(), demo.WorkspaceTier.Name})
	t.AppendRow([]interface{}{demo.Application.TableName(), demo.Application.Id})
	t.AppendRow([]interface{}{demo.Endpoint.TableName(), demo.Endpoint.Id})
	for _, rule := range demo.EndpointRules {
		t.AppendRow([]interface{}{rule.TableName(), rule.Id})
	}

	t.SetOutputMirror(os.Stdout)
	t.Render()
}
