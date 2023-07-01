package ito

import (
	"sort"
)

func Sort[T any](ito Ito[T], less func(T, T) bool) Ito[T] {
	values := make([]T, 0)
	for ito.Next() {
		values = append(values, ito.Item())
	}
	sort.SliceStable(values, func(i, j int) bool {
		return less(values[i], values[j])
	})
	return &sortedIto[T]{values: values, index: -1}
}

type sortedIto[T any] struct {
	values []T
	index  int
}

func (s *sortedIto[T]) Next() bool {
	s.index++
	return s.index < len(s.values)
}

func (s *sortedIto[T]) Item() T {
	return s.values[s.index]
}
