package routing

import (
	"errors"
	"strings"

	"github.com/scrapnode/kanthor/internal/entities"
)

var (
	ConditionExpressionDivider = "::"
	ConditionExpressionAny     = "any"
	ConditionExpressionEqual   = "equal"
	ConditionExpressionPrefix  = "prefix"
)

var (
	ErrConditionExpressionMalformed = errors.New("PLANNER.CONDITION.EXPRESSION.MALFORMED")
	ErrConditionExpressionUnknown   = errors.New("PLANNER.CONDITION.EXPRESSION.UNKNOWN")
)

func ConditionExpression(rule *entities.EndpointRule) (func(source string) bool, error) {
	expression := strings.Split(rule.ConditionExpression, ConditionExpressionDivider)
	if len(expression) != 2 {
		return nil, ErrConditionExpressionMalformed
	}

	if expression[0] == ConditionExpressionAny {
		return func(source string) bool { return true }, nil
	}

	if expression[0] == ConditionExpressionEqual {
		return func(source string) bool { return source == expression[1] }, nil
	}

	if expression[0] == ConditionExpressionPrefix {
		return func(source string) bool { return strings.HasPrefix(source, expression[1]) }, nil
	}

	return nil, ErrConditionExpressionUnknown
}
