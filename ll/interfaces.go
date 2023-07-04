package ll

import "sync"

type starter interface {
	start(*sync.WaitGroup, chan<- error)
}

type inputter[I any] interface {
	input() chan<- I
}

type StartInputter[I any] interface {
	starter
	inputter[I]
}

type Piper[I, O any] interface {
	starter
	inputter[I]
	Next(StartInputter[O])
}
