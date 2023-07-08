package ll

import "sync"

type Term interface {
	waitGroup(*sync.WaitGroup)
	Error(error)
	Done()
}

type Sink[T any] interface {
	Write(T)
	Term
}

type sinkWrapper[T any] struct {
	doneTapFn func()
	sink      Sink[T]
	isDone    bool
}

func newSinkWrapper[T any](doneTapFn func(), sink Sink[T]) Sink[T] {
	return &sinkWrapper[T]{
		doneTapFn: doneTapFn,
		sink:      sink,
	}
}

func (s *sinkWrapper[T]) waitGroup(wg *sync.WaitGroup) {
	s.sink.waitGroup(wg)
}
func (s *sinkWrapper[T]) Write(item T) {
	s.sink.Write(item)
}
func (s *sinkWrapper[T]) Error(err error) {
	s.sink.Error(err)
}
func (s *sinkWrapper[T]) Done() {
	if !s.isDone {
		s.doneTapFn()
		s.sink.Done()
		s.isDone = true
	}
}
