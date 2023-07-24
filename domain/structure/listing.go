package structure

import "reflect"

type ListReq struct {
	Cursor string
	Search string
	Limit  int

	Ids []string
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

func ListReqBuild(opts []ListOps) *ListReq {
	req := ListReq{}
	for _, opt := range opts {
		opt(&req)
	}
	return &req
}

func ListResBuild[T any](res *ListRes[T], req *ListReq) *ListRes[T] {
	if len(res.Data) == 0 {
		return res
	}

	if req.Limit <= 0 {
		return res
	}

	if len(res.Data) < req.Limit {
		return res
	}

	latest := any(res.Data[len(res.Data)-1])
	value := reflect.ValueOf(latest)
	field := value.FieldByName("Id")

	if field.IsValid() {
		res.Cursor = field.String()
	}
	return res
}

type ListRes[T any] struct {
	Cursor string
	Data   []T
}
