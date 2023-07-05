package ll

import "sync"

type Piper[I, O any] interface {
	Sink[I]
	Next(Sink[O])
}

type pipe[I, O any] struct {
	wg    *sync.WaitGroup
	next  Sink[O]
	runFn PipeRunFn[I, O]
	input chan I
	error chan error
}

func Pipe[I, O any](runFn PipeRunFn[I, O]) Piper[I, O] {
	p := &pipe[I, O]{
		runFn: runFn,
		input: make(chan I),
		error: make(chan error),
	}
	go func() {
		inputIsOpen := true
		for inputIsOpen {
			select {
			case item, ok := <-p.input:
				if ok {
					p.runFn(item, p.next)
				} else {
					inputIsOpen = false
				}
			case err := <-p.error:
				p.next.Error(err)
			}
		}
		p.next.Done()
		p.wg.Done()
	}()
	return p
}

func (p *pipe[I, O]) waitGroup(wg *sync.WaitGroup) {
	if p.wg == nil {
		p.wg = wg
		p.wg.Add(1)
		if p.next == nil {
			panic("Nexts must all be set before calling Start()")
		}
		p.next.waitGroup(wg)
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
	close(p.input)
}

func (p *pipe[I, O]) Next(sink Sink[O]) {
	p.next = sink
}
