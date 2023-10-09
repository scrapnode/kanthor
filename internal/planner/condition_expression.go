package planner

import (
	"errors"
	"strings"

	"github.com/scrapnode/kanthor/domain/entities"
)

var (
	ConditionExpressionDivider = "::"
	ConditionExpressionEqual   = "equal"
	ConditionExpressionAny     = "any"
)

func ConditionExpression(rule *entities.EndpointRule) (func(source string) bool, error) {
	expression := strings.Split(rule.ConditionExpression, ConditionExpressionDivider)
	if len(expression) != 2 {
		return nil, errors.New("PLANNER.CONDITION.EXPRESSION.INVALID")
	}

	if expression[0] == ConditionExpressionEqual {
		return func(source string) bool { return expression[1] == source }, nil
	}

	if expression[0] == ConditionExpressionAny {
		return func(source string) bool { return true }, nil
	}

	return nil, errors.New("PLANNER.CONDITION.EXPRESSION.UNKNOWN")
}
