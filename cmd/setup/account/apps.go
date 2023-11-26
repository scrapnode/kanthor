package account

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/scrapnode/kanthor/data/interchange"
	"github.com/scrapnode/kanthor/internal/domain/entities"
	"github.com/scrapnode/kanthor/pkg/safe"
	"github.com/scrapnode/kanthor/services/portal/usecase"
)

func apps(uc usecase.Portal, ctx context.Context, ws *entities.Workspace, file string, p *printing) error {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	data, err := interchange.Unmarshal(bytes)
	if err != nil {
		return err
	}

	apps, eps, eprs, tree := mapping(ws, data)
	in := &usecase.WorkspaceSetupIn{
		Workspace:     ws,
		Applications:  apps,
		Endpoints:     eps,
		EndpointRules: eprs,
	}
	if err := in.Validate(); err != nil {
		return err
	}
	out, err := uc.Workspace().Setup(ctx, in)
	if err != nil {
		return err
	}

	p.AddStdout(appsOutput(tree, out.Status))
	p.AddJson("applications", out.ApplicationIds)

	return nil
}

func mapping(doc *entities.Workspace, ws *interchange.Workspace) ([]entities.Application, []entities.Endpoint, []entities.EndpointRule, *safe.Node[string]) {
	now := time.Now().UTC()
	applications := []entities.Application{}
	endpoints := []entities.Endpoint{}
	rules := []entities.EndpointRule{}

	tree := &safe.Node[string]{
		Value:    doc.Id,
		Children: []safe.Node[string]{},
	}

	for _, app := range ws.Applications {
		application := entities.Application{WsId: doc.Id, Name: app.Name}
		application.GenId()
		application.SetAT(now)
		applications = append(applications, application)

		appNode := safe.Node[string]{
			Value:    application.Id,
			Children: []safe.Node[string]{},
		}

		for _, ep := range app.Endpoints {
			endpoint := entities.Endpoint{
				AppId:  application.Id,
				Name:   ep.Name,
				Method: ep.Method,
				Uri:    ep.Uri,
			}
			endpoint.GenId()
			endpoint.SetAT(now)
			endpoint.GenSecretKey()
			endpoints = append(endpoints, endpoint)

			epNode := safe.Node[string]{
				Value:    endpoint.Id,
				Children: []safe.Node[string]{},
			}

			for _, epr := range ep.Rules {
				rule := entities.EndpointRule{
					EpId:                endpoint.Id,
					Name:                epr.Name,
					Priority:            epr.Priority,
					Exclusionary:        epr.Exclusionary,
					ConditionSource:     epr.ConditionSource,
					ConditionExpression: epr.ConditionExpression,
				}
				rule.GenId()
				rule.SetAT(now)
				rules = append(rules, rule)

				eprNode := safe.Node[string]{
					Value:    rule.Id,
					Children: []safe.Node[string]{},
				}
				epNode.Children = append(epNode.Children, eprNode)
			}

			appNode.Children = append(appNode.Children, epNode)
		}

		tree.Children = append(tree.Children, appNode)
	}

	return applications, endpoints, rules, tree
}

func appsOutput(tree *safe.Node[string], status map[string]bool) string {
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
