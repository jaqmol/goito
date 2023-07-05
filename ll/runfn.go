package ll

type PipeRunFn[I, O any] func(I, Sink[O])

type EndRunFn[I any] func(I, Term)
