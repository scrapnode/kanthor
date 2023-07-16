package structure

import "github.com/scrapnode/kanthor/domain/entities"

type ListReq struct {
	Cursor string
	Search string
	Limit  int
	Ids    []string
}

type ListOps func(req *ListReq)

func WithListIds(ids []string) ListOps {
	return func(req *ListReq) {
		req.Ids = ids
	}
}

func WithListCursor(cursor string) ListOps {
	return func(req *ListReq) {
		req.Cursor = cursor
	}
}

func WithListSearch(search string) ListOps {
	return func(req *ListReq) {
		req.Search = search
	}
}

func WithListLimit(limit int) ListOps {
	return func(req *ListReq) {
		req.Limit = limit
	}
}

func ListReqBuild(opts []ListOps) ListReq {
	req := ListReq{}
	for _, opt := range opts {
		opt(&req)
	}
	return req
}

func ListResBuild[T any](res *ListRes[T]) *ListRes[T] {
	if len(res.Data) == 0 {
		return res
	}

	latest := any(res.Data[len(res.Data)-1])
	res.Cursor = latest.(entities.Entity).Id
	return res
}

type ListRes[T any] struct {
	Cursor string
	Data   []T
}
