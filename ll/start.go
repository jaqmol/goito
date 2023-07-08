package ll

import (
	"sync"
)

type Starter[I, O any] interface {
	Piper[I, O]
	Start() Starter[I, O]
}

type start[I, O any] struct {
	piper Piper[I, O]
}

func Start[I, O any](runFn PipeRunFn[I, O]) Starter[I, O] {
	s := &start[I, O]{
		piper: Pipe[I, O](runFn),
	}
	return s
}

func (s *start[I, O]) waitGroup(wg *sync.WaitGroup) {
	s.piper.waitGroup(wg)
}

func (s *start[I, O]) Start() Starter[I, O] {
	var wg sync.WaitGroup
	s.waitGroup(&wg)
	return s
}

func (s *start[I, O]) Next(next Sink[O]) {
	s.piper.Next(next)
}

func (s *start[I, O]) Write(item I) {
	s.piper.Write(item)
}

func (s *start[I, O]) Error(err error) {
	s.piper.Error(err)
}

func (s *start[I, O]) Done() {
	s.piper.Done()
}
