package routing_test

import (
	"fmt"
	"testing"

	"github.com/jaswdr/faker"
	"github.com/scrapnode/kanthor/internal/entities"
	"github.com/scrapnode/kanthor/internal/routing"
	"github.com/stretchr/testify/assert"
)

func TestConditionExpression(t *testing.T) {
	fake := faker.New()

	t.Run("malformed error", func(st *testing.T) {
		_, err := routing.ConditionExpression(&entities.EndpointRule{ConditionExpression: string(routing.ConditionExpressionDivider[0])})
		assert.NotNil(st, err)
		assert.ErrorIs(st, err, routing.ErrConditionExpressionMalformed)
	})

	t.Run("unknown error", func(st *testing.T) {
		_, err := routing.ConditionExpression(&entities.EndpointRule{ConditionExpression: routing.ConditionExpressionDivider})
		assert.NotNil(st, err)
		assert.ErrorIs(st, err, routing.ErrConditionExpressionUnknown)
	})

	t.Run("match any", func(st *testing.T) {
		expression := fmt.Sprintf("%s%s", routing.ConditionExpressionAny, routing.ConditionExpressionDivider)

		match, err := routing.ConditionExpression(&entities.EndpointRule{ConditionExpression: expression})
		assert.Nil(st, err)
		assert.True(st, match(fake.App().Name()))
		assert.True(st, match(fake.Blood().Name()))
	})

	t.Run("match equal", func(st *testing.T) {
		target := fake.Blood().Name()

		expression := fmt.Sprintf("%s%s%s", routing.ConditionExpressionEqual, routing.ConditionExpressionDivider, target)

		match, err := routing.ConditionExpression(&entities.EndpointRule{ConditionExpression: expression})
		assert.Nil(st, err)
		assert.True(st, match(target))
		assert.False(st, match(fake.App().Name()))
	})

	t.Run("match prefix", func(st *testing.T) {
		begin := fake.RandomStringWithLength(5)
		end := fake.RandomStringWithLength(5)
		target := fmt.Sprintf("%s.%s", begin, end)

		expression := fmt.Sprintf("%s%s%s", routing.ConditionExpressionPrefix, routing.ConditionExpressionDivider, begin)

		match, err := routing.ConditionExpression(&entities.EndpointRule{ConditionExpression: expression})
		assert.Nil(st, err)
		assert.True(st, match(target))
		assert.True(st, match(begin))
		assert.False(st, match(end))
	})
}
