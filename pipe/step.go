package pipe

import "sync"

func Step[I, O any](runFn func(I, chan<- O, chan<- error)) Nexter[O] {
	return &step[I, O]{
		input: make(chan I),
		runFn: runFn,
	}
}

type step[I, O any] struct {
	input  chan I
	runFn  func(I, chan<- O, chan<- error)
	output chan<- O
}

func (p *step[I, O]) Start(wg *sync.WaitGroup, errCh chan<- error) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for in := range p.input {
			p.runFn(in, p.output, errCh)
		}
		close(p.output)
	}()
}

func (p *step[I, O]) Input() <-chan I {
	return p.input
}

func (p *step[I, O]) Next(in inputter[O]) {
	p.output = in.input()
}
