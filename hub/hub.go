package hub

type Hub struct {
	pipes    []pipe[any]
	handlers []handler[any]
}

func On[T any](h Hub, name string, handlerFn func(T)) {
	var thisPipe pipe[T]
	for _, c := range h.pipes {
		if c.is(name) {
			thisPipe = c.(pipe[T])
			break
		}
	}
	if thisPipe == nil {
		thisPipe = &pipeImpl[T]{name: name, channel: make(chan T)}
		h.pipes = append(h.pipes, thisPipe.(pipe[any]))
		go func() {
			ok := true
			for ok {
				var msg any
				msg, ok = thisPipe.receive()
				if ok {
					for _, h := range h.handlers {
						h.handle(name, msg)
					}
				}
			}
		}()
	}
	var thisHandler handler[T] = &handlerImpl[T]{name: name, handlerFn: handlerFn}
	h.handlers = append(h.handlers, thisHandler.(handler[any]))
}

func Send[T any](h Hub, name string, msg any) {
	for _, c := range h.pipes {
		c.send(name, msg)
	}
}

func Cascade(h Hub, names ...string) {
	pipes := make([]pipe[any], len(names))
	for i, n := range names {
		for _, p := range h.pipes {
			if p.is(n) {
				pipes[i] = p
			}
		}
	}
	lstIdx := len(pipes) - 1
	for i, p := range pipes {
		if i < lstIdx {
			p.onClose(pipes[i+1])
		}
	}
}

func Close(h Hub, name string) {
	delIdx := -1
	for i, c := range h.pipes {
		if c.close(name) {
			delIdx = i
			break
		}
	}
	if delIdx >= 0 {
		h.pipes = append(h.pipes[:delIdx], h.pipes[delIdx+1:]...)
	}
}

type closer interface {
	close(string) bool
}

type pipe[T any] interface {
	is(string) bool
	send(string, T)
	receive() (T, bool)
	close(string) bool
	onClose(closer)
}

type pipeImpl[T any] struct {
	name    string
	channel chan T
	closers []closer
}

func (p *pipeImpl[T]) is(name string) (ok bool) {
	ok = p.name == name
	return
}

func (p *pipeImpl[T]) send(name string, value T) {
	if p.name == name {
		p.channel <- value
	}
}

func (p *pipeImpl[T]) receive() (value T, ok bool) {
	value, ok = <-p.channel
	return
}

func (p *pipeImpl[T]) close(name string) (ok bool) {
	if p.name == name {
		close(p.channel)
		ok = true
		for _, c := range p.closers {
			c.close(name)
		}
	}
	return
}

func (p *pipeImpl[T]) onClose(c closer) {
	p.closers = append(p.closers, c)
}

type handler[T any] interface {
	is(string) bool
	handle(string, T)
}

type handlerImpl[T any] struct {
	name      string
	handlerFn func(T)
}

func (h *handlerImpl[T]) is(name string) (ok bool) {
	ok = h.name == name
	return
}

func (h *handlerImpl[T]) handle(name string, msg T) {
	if h.name == name {
		h.handlerFn(msg)
	}
}
