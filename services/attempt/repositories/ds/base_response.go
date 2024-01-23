package ds

import "context"

type Response interface {
	Check(ctx context.Context, epId string, msgIds []string) (map[string][]int, error)
}
