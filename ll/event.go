package ll

type Event[T any] struct {
	Close bool
	Value T
}
