package event

import "github.com/genshinsim/gcsim/pkg/core"

type Ctrl struct {
	c      *core.Core
	events [][]ehook
}

type ehook struct {
	f   core.EventHook
	key string
	src int
}

func NewCtrl(c *core.Core) *Ctrl {
	h := &Ctrl{c: c}

	h.events = make([][]ehook, core.EndEventTypes)

	for i := range h.events {
		h.events[i] = make([]ehook, 0, 10)
	}

	return h
}

func (h *Ctrl) Subscribe(e core.EventType, f core.EventHook, key string) {
	a := h.events[e]

	//check if override first
	ind := -1
	for i, v := range a {
		if v.key == key {
			ind = i
		}
	}
	if ind > -1 {
		h.c.Log.NewEvent("hook added", core.LogHookEvent, -1, "overwrite", true, "key", key, "type", e)
		a[ind] = ehook{
			f:   f,
			key: key,
			src: h.c.F,
		}
	} else {
		a = append(a, ehook{
			f:   f,
			key: key,
			src: h.c.F,
		})
		h.c.Log.NewEvent("hook added", core.LogHookEvent, -1, "overwrite", false, "key", key, "type", e)
	}
	h.events[e] = a
}

func (h *Ctrl) Unsubscribe(e core.EventType, key string) {
	n := 0
	for _, v := range h.events[e] {
		if v.key != key {
			h.events[e][n] = v
			n++
		}
	}
	h.events[e] = h.events[e][:n]
}

func (h *Ctrl) Emit(e core.EventType, args ...interface{}) {
	n := 0
	for i, v := range h.events[e] {
		if v.f(args...) {
			h.c.Log.NewEvent("event hook ended", core.LogHookEvent, -1, "key", i, "src", v.src)
		} else {
			h.events[e][n] = v
			n++
		}
	}
	h.events[e] = h.events[e][:n]
}
