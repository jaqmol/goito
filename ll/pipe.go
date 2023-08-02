package ll

import (
	"sync"
)

type Piper[I, O any] interface {
	Sink[I]
	Next(Sink[O])
}

type Nexter[O any] interface {
	Next(Sink[O])
}

type pipe[I, O any] struct {
	wg     *sync.WaitGroup
	sink   Sink[O]
	runFn  PipeRunFn[I, O]
	input  chan I
	error  chan error
	isDone bool
}

func Pipe[I, O any](runFn PipeRunFn[I, O]) Piper[I, O] {
	p := &pipe[I, O]{
		runFn: runFn,
		input: make(chan I, 1),
		error: make(chan error, 1),
	}
	go func() {
		inputIsOpen := true
		for inputIsOpen {
			select {
			case item, ok := <-p.input:
				if ok {
					p.runFn(item, p.sink)
				} else {
					inputIsOpen = false
				}
			case err := <-p.error:
				p.sink.Error(err)
			}
		}
		p.sink.Done()
		p.wg.Done()
	}()
	return p
}

func (p *pipe[I, O]) waitGroup(wg *sync.WaitGroup) {
	if p.wg == nil {
		p.wg = wg
		p.wg.Add(1)
		if p.sink == nil {
			panic("Nexts must all be set before calling Start()")
		}
		p.sink.waitGroup(wg)
	} else {
		panic("Wait group already set")
	}
}

func (p *pipe[I, O]) Write(item I) {
	p.input <- item
}

func (p *pipe[I, O]) Error(err error) {
	p.error <- err
}

func (p *pipe[I, O]) Done() {
	if !p.isDone {
		close(p.input)
		p.isDone = true
	}
}

func (p *pipe[I, O]) Next(sink Sink[O]) {
	p.sink = newSinkWrapper[O](p.Done, sink)
}
