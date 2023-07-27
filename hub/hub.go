package hub

type Hub struct {
	channels []channel
	handlers []handler
}

func On(h Hub, name string, handlerFn func(any)) {
	var finding chan any
	for _, c := range h.channels {
		if c.name == name {
			finding = c.channel
			break
		}
	}
	if finding == nil {
		finding = make(chan any)
		h.channels = append(h.channels, channel{name: name, channel: finding})
		go func() {
			ok := true
			for ok {
				var msg any
				msg, ok = <-finding
				if ok {
					for _, h := range h.handlers {
						if h.name == name {
							h.handlerFn(msg)
						}
					}
				}
			}
		}()
	}
	h.handlers = append(h.handlers, handler{name: name, handlerFn: handlerFn})
}

func Send(h Hub, name string, msg any) {
	for _, c := range h.channels {
		if c.name == name {
			c.channel <- msg
			break
		}
	}
}

type channel struct {
	name    string
	channel chan any
}

type handler struct {
	name      string
	handlerFn func(any)
}
