package ito

type Ito[T any] interface {
	Next() bool
	Item() T
}
