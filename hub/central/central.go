package central

import (
	"github.com/jaqmol/goito/hub"
)

var h *hub.Hub = &hub.Hub{}

func On[T any](name string, handlerFn func(T)) {
	hub.On[T](h, name, handlerFn)
}

func Send[T any](name string, msg T) {
	hub.Send[T](h, name, msg)
}

func Cascade(names ...string) {
	hub.Cascade(h, names...)
}

func Close(name string) {
	hub.Close(h, name)
}
