package ds

import "context"

type Request interface {
	Check(ctx context.Context, pairs []string) (map[string]bool, error)
}
