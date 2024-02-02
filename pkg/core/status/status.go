package status

import "github.com/genshinsim/gcsim/pkg/core/glog"

type status struct {
	expiry int
	evt    glog.Event
}

type Handler struct {
	status map[string]status
	f      *int
	log    glog.Logger
}

func New(f *int, log glog.Logger) *Handler {
	return &Handler{
		status: make(map[string]status),
		f:      f,
		log:    log,
	}
}

func (t *Handler) Duration(key string) int {
	a, ok := t.status[key]
	if !ok {
		return 0
	}
	if a.expiry > *t.f {
		return a.expiry - *t.f
	}
	return 0
}

func (t *Handler) Add(key string, dur int) {
	// check if exists
	a, ok := t.status[key]

	// if ok we want to reuse the old evt
	if ok && a.expiry > *t.f {
		// just reuse the old and update expiry + evt.Ended
		a.expiry = *t.f + dur
		a.evt.SetEnded(a.expiry)
		t.status[key] = a
		t.log.NewEvent("status refreshed", glog.LogStatusEvent, -1).
			Write("key", key).
			Write("expiry", *t.f+dur)
		return
	}

	// otherwise create a new event
	a.evt = t.log.NewEvent("status added", glog.LogStatusEvent, -1).
		Write("key", key).
		Write("expiry", *t.f+dur)
	a.expiry = *t.f + dur
	a.evt.SetEnded(a.expiry)

	t.status[key] = a
}

func (t *Handler) Extend(key string, amt int) {
	a, ok := t.status[key]

	// do nothing if status doesn't exist
	if !ok || a.expiry <= *t.f {
		return
	}

	a.expiry += amt
	a.evt.SetEnded(a.expiry)
	t.status[key] = a
	t.log.NewEvent("status extended", glog.LogStatusEvent, -1).
		Write("key", key).
		Write("amt", amt).
		Write("expiry", a.expiry)
}

func (t *Handler) Delete(key string) {
	// check if it exists first
	a, ok := t.status[key]
	if ok && a.expiry > *t.f {
		a.evt.SetEnded(*t.f)
	}
	delete(t.status, key)
}
