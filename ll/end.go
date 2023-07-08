package ll

import (
	"bytes"
	"sync"
)

type Ender[I any] interface {
	Sink[I]
	Wait() error
}

type end[I any] struct {
	wg     *sync.WaitGroup
	runFn  EndRunFn[I]
	input  chan I
	error  chan error
	mErrs  MultiError
	isDone bool
}

func End[I any](runFn EndRunFn[I]) Ender[I] {
	e := &end[I]{
		runFn: runFn,
		input: make(chan I),
		error: make(chan error),
	}
	go func() {
		inputIsOpen := true
		for inputIsOpen {
			select {
			case item, ok := <-e.input:
				if ok {
					e.runFn(item, e)
				} else {
					inputIsOpen = false
				}
			case err := <-e.error:
				e.mErrs.Errors = append(e.mErrs.Errors, err)
			}
		}
		e.wg.Done()
	}()
	return e
}

func (e *end[I]) waitGroup(wg *sync.WaitGroup) {
	if e.wg == nil {
		e.wg = wg
		e.wg.Add(1)
	} else {
		panic("Wait group already set")
	}
}

func (e *end[I]) Write(item I) {
	e.input <- item
}

func (e *end[I]) Error(err error) {
	e.error <- err
}

func (e *end[I]) Done() {
	if !e.isDone {
		close(e.input)
		e.isDone = true
	}
}

func (e *end[I]) Wait() error {
	e.wg.Wait()
	if len(e.mErrs.Errors) > 0 {
		return &e.mErrs
	}
	return nil
}

type MultiError struct {
	Errors []error
}

func (m *MultiError) Error() string {
	var buff bytes.Buffer
	li := len(m.Errors) - 1
	for i, err := range m.Errors {
		buff.WriteString(err.Error())
		if i < li {
			buff.WriteString(", ")
		}
	}
	return buff.String()
}
