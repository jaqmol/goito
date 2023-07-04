package ll

import "sync"

type Generator[I, O any] interface {
	Piper[I, O]
	Start(iv I) error
}

func Gen[I, O any](runFn func(I, chan<- O, chan<- error)) Generator[I, O] {
	return &gen[I, O]{
		runFn: runFn,
	}
}

type gen[I, O any] struct {
	runFn  func(I, chan<- O, chan<- error)
	nextSi StartInputter[O]
	startV I
}

func (g *gen[I, O]) Start(iv I) error {
	g.startV = iv
	errCh := make(chan error)
	var wg sync.WaitGroup
	g.start(&wg, errCh)
	errs := multiErr{}
	go func() {
		for err := range errCh {
			errs.add(err)
		}
	}()
	wg.Wait()
	if len(errs.errs) > 0 {
		return &errs
	}
	return nil
}

func (g *gen[I, O]) start(wg *sync.WaitGroup, errCh chan<- error) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		g.runFn(g.startV, g.nextSi.input(), errCh)
		// close(g.nextSi.input())
	}()
	g.nextSi.start(wg, errCh)
}

func (g *gen[I, O]) input() chan<- I {
	return nil
}

func (g *gen[I, O]) Next(si StartInputter[O]) {
	g.nextSi = si
}
