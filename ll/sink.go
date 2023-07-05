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
