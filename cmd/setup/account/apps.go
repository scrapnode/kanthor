package account

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/scrapnode/kanthor/data/interchange"
	"github.com/scrapnode/kanthor/domain/entities"
	"github.com/scrapnode/kanthor/domain/structure"
	usecase "github.com/scrapnode/kanthor/usecases/portal"
)

func apps(uc usecase.Portal, ctx context.Context, ws *entities.Workspace, file string, out *output) error {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	data, err := interchange.Unmarshal(bytes)
	if err != nil {
		return err
	}

	apps, eps, eprs, tree := mapping(ws, data)
	ucreq := &usecase.WorkspaceSetupReq{
		Workspace:     ws,
		Applications:  apps,
		Endpoints:     eps,
		EndpointRules: eprs,
	}
	if err := ucreq.Validate(); err != nil {
		return err
	}
	ucres, err := uc.Workspace().Setup(ctx, ucreq)
	if err != nil {
		return err
	}

	out.AddStdout(appsOutput(tree, ucres.Status))
	out.AddJson("applications", ucres.ApplicationIds)

	return nil
}

func mapping(doc *entities.Workspace, ws *interchange.Workspace) ([]entities.Application, []entities.Endpoint, []entities.EndpointRule, *structure.Node[string]) {
	now := time.Now().UTC()
	applications := []entities.Application{}
	endpoints := []entities.Endpoint{}
	rules := []entities.EndpointRule{}

	tree := &structure.Node[string]{
		Value:    doc.Id,
		Children: []structure.Node[string]{},
	}

	for _, app := range ws.Applications {
		application := entities.Application{WorkspaceId: doc.Id, Name: app.Name}
		application.GenId()
		application.SetAT(now)
		applications = append(applications, application)

		appNode := structure.Node[string]{
			Value:    application.Id,
			Children: []structure.Node[string]{},
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

			epNode := structure.Node[string]{
				Value:    endpoint.Id,
				Children: []structure.Node[string]{},
			}

			for _, epr := range ep.Rules {
				rule := entities.EndpointRule{
					EndpointId:          endpoint.Id,
					Name:                epr.Name,
					Priority:            epr.Priority,
					Exclusionary:        epr.Exclusionary,
					ConditionSource:     epr.ConditionSource,
					ConditionExpression: epr.ConditionExpression,
				}
				rule.GenId()
				rule.SetAT(now)
				rules = append(rules, rule)

				eprNode := structure.Node[string]{
					Value:    rule.Id,
					Children: []structure.Node[string]{},
				}
				epNode.Children = append(epNode.Children, eprNode)
			}

			appNode.Children = append(appNode.Children, epNode)
		}

		tree.Children = append(tree.Children, appNode)
	}

	return applications, endpoints, rules, tree
}

func appsOutput(tree *structure.Node[string], status map[string]bool) string {
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
