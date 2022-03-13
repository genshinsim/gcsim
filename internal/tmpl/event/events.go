package event

import (
	"github.com/genshinsim/gcsim/pkg/coretype"
)

type Ctrl struct {
	c      core
	events [][]ehook
}

type core interface {
	coretype.Framer
	coretype.Logger
}

type ehook struct {
	f   coretype.EventHook
	key string
	src int
}

func NewCtrl(c core) *Ctrl {
	h := &Ctrl{c: c}

	h.events = make([][]ehook, coretype.EndEventTypes)

	for i := range h.events {
		h.events[i] = make([]ehook, 0, 10)
	}

	return h
}

func (h *Ctrl) Subscribe(e coretype.EventType, f coretype.EventHook, key string) {
	a := h.events[e]

	//check if override first
	ind := -1
	for i, v := range a {
		if v.key == key {
			ind = i
		}
	}
	if ind > -1 {
		h.c.NewEvent("hook added", coretype.LogHookEvent, -1, "overwrite", true, "key", key, "type", e)
		a[ind] = ehook{
			f:   f,
			key: key,
			src: h.c.F(),
		}
	} else {
		a = append(a, ehook{
			f:   f,
			key: key,
			src: h.c.F(),
		})
		h.c.NewEvent("hook added", coretype.LogHookEvent, -1, "overwrite", true, "key", key, "type", e)
	}
	h.events[e] = a
}

func (h *Ctrl) Unsubscribe(e coretype.EventType, key string) {
	n := 0
	for _, v := range h.events[e] {
		if v.key != key {
			h.events[e][n] = v
			n++
		}
	}
	h.events[e] = h.events[e][:n]
}

func (h *Ctrl) Emit(e coretype.EventType, args ...interface{}) {
	n := 0
	for i, v := range h.events[e] {
		if v.f(args...) {
			h.c.NewEvent("event hook ended", coretype.LogHookEvent, -1, "key", i, "src", v.src)
		} else {
			h.events[e][n] = v
			n++
		}
	}
	h.events[e] = h.events[e][:n]
}
