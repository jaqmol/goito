package ll

import "sync"

func End[I any](runFn func(I, chan<- error)) StartInputter[I] {
	return &end[I]{
		inputCh: make(chan I),
		runFn:   runFn,
	}
}

type end[I any] struct {
	inputCh chan I
	runFn   func(I, chan<- error)
}

func (e *end[I]) start(wg *sync.WaitGroup, errCh chan<- error) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for in := range e.inputCh {
			e.runFn(in, errCh)
		}
	}()
}

func (e *end[I]) input() chan<- I {
	return e.inputCh
}
