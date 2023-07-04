package ll

import (
	"sync"
)

func New[I, O any](runFn func(I, chan<- O, chan<- error)) Piper[I, O] {
	return &pipe[I, O]{
		inputCh: make(chan I),
		runFn:   runFn,
	}
}

type pipe[I, O any] struct {
	inputCh chan I
	runFn   func(I, chan<- O, chan<- error)
	nextSi  StartInputter[O]
}

func (p *pipe[I, O]) start(wg *sync.WaitGroup, errCh chan<- error) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for in := range p.inputCh {
			p.runFn(in, p.nextSi.input(), errCh)
		}
		// close(p.nextSi.input())
	}()
	p.nextSi.start(wg, errCh)
}

func (p *pipe[I, O]) input() chan<- I {
	return p.inputCh
}

func (p *pipe[I, O]) Next(si StartInputter[O]) {
	p.nextSi = si
}
