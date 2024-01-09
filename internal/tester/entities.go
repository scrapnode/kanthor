package tester

import (
	"fmt"
	"net/http"

	"github.com/jaswdr/faker"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/pkg/identifier"
	"github.com/scrapnode/kanthor/pkg/timer"
	"github.com/scrapnode/kanthor/pkg/utils"
	"github.com/scrapnode/kanthor/project"
)

var fake = faker.New()

func Application(t timer.Timer) *entities.Application {
	doc := &entities.Application{
		WsId: identifier.New(entities.IdNsWs),
		Name: fake.App().Name(),
	}
	doc.Id = identifier.New(entities.IdNsApp)
	doc.SetAT(t.Now())

	return doc
}

func EndpointOfApp(t timer.Timer, app *entities.Application) *entities.Endpoint {
	doc := &entities.Endpoint{
		AppId:     app.Id,
		SecretKey: utils.RandomString(32),
		Name:      fake.Beer().Name(),
		Method:    fake.RandomStringElement([]string{http.MethodPost, http.MethodPut}),
		Uri:       fake.Internet().URL(),
	}
	doc.Id = identifier.New(entities.IdNsEp)
	doc.SetAT(t.Now())

	return doc
}

func RuleOfEndpoint(t timer.Timer, ep *entities.Endpoint) *entities.EndpointRule {
	doc := &entities.EndpointRule{
		EpId:                ep.Id,
		Name:                fake.Blood().Name(),
		Priority:            fake.Int32Between(1, 100),
		Exclusionary:        fake.Bool(),
		ConditionSource:     "type",
		ConditionExpression: "any::",
	}
	doc.Id = identifier.New(entities.IdNsEpr)
	doc.SetAT(t.Now())

	return doc
}

func MessageOfApp(t timer.Timer, app *entities.Application) *entities.Message {
	json := fake.Json()
	doc := &entities.Message{
		AppId:    app.Id,
		Tier:     project.Tier(),
		Type:     fmt.Sprintf("%s.%s", fake.RandomStringWithLength(5), fake.RandomStringWithLength(5)),
		Body:     (&json).String(),
		Headers:  entities.Header{},
		Metadata: entities.Metadata{},
	}
	doc.Id = identifier.New(entities.IdNsMsg)
	doc.SetTS(t.Now())

	return doc
}
