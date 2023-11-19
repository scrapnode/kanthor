package datastore

import (
	"errors"
	"fmt"
	"sort"
)

// use mapper to create a better insert query
// ref https://klotzandrew.com/blog/postgres-passing-65535-parameter-limit
func NewMapper[T any](fns map[string]func(doc T) any, casters map[string]string) *Mapper[T] {
	if casters == nil {
		casters = map[string]string{}
	}
	return &Mapper[T]{fns: fns, data: map[string][]any{}, casters: casters}
}

type Mapper[T any] struct {
	fns     map[string]func(doc T) any
	casters map[string]string
	data    map[string][]any
}

func (mapper *Mapper[T]) Parse(docs []T) error {
	if len(docs) == 0 {
		return nil
	}

	attributes := mapper.Names()
	if len(attributes) == 0 {
		return errors.New("no parse function was registered for documents")
	}

	// init values
	for _, attribute := range attributes {
		mapper.data[attribute] = []any{}
	}

	for _, doc := range docs {
		for _, attribute := range attributes {
			parse := mapper.fns[attribute]
			value := parse(doc)
			mapper.data[attribute] = append(mapper.data[attribute], value)
		}
	}

	return nil
}

func (mapper *Mapper[T]) Names() []string {
	var attributes []string
	for key := range mapper.fns {
		attributes = append(attributes, key)
	}
	sort.Slice(attributes, func(i, j int) bool {
		return attributes[i] < attributes[j]
	})
	return attributes
}

func (mapper *Mapper[T]) Values() []any {
	var values []any

	// it's important to return data in a persistent order
	attributes := mapper.Names()
	for _, attribute := range attributes {
		values = append(values, mapper.data[attribute])
	}

	return values
}

func (mapper *Mapper[T]) Casters() []string {
	var returning []string

	attributes := mapper.Names()
	for i, attribute := range attributes {
		cast := "varchar[]"
		if custom, ok := mapper.casters[attribute]; ok {
			cast = custom
		}

		returning = append(returning, fmt.Sprintf("$%d::%s", i+1, cast))
	}

	return returning
}
