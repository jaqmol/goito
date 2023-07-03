package pipe

import "sync"

type Starter[I, O any] interface {
	nexter[O]
	Start() error
}

func Start[I, O any](runFn func(I, chan<- O, chan<- error)) Starter[I, O] {
	return &start[I, O]{
		runFn: runFn,
	}
}

type start[I, O any] struct {
	runFn func(I, chan<- O, chan<- error)
	nxt   nexter[O]
}

func (p *start[I, O]) Start() error {
	return nil
}

func (p *start[I, O]) start(wg *sync.WaitGroup, in I, errCh chan<- error) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		p.runFn(in, p.nxt.input(), errCh)
		close(p.nxt.input())
	}()
}

func (p *start[I, O]) input() chan<- I {
	return nil
}

func (p *start[I, O]) next(nxt nexter[O]) {
	p.nxt = nxt
}
