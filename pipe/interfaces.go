package pipe

import "sync"

type inputter[I any] interface {
	input() chan<- I
}

type nexter[O any] interface {
	next(nexter[O])
	inputter[O]
	starter[O]
}

type starter[I any] interface {
	start(*sync.WaitGroup, I, chan<- error)
}
