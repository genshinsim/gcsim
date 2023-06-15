// package modifier provides a universal way of handling a slice
// of modifiers
package modifier

import (
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

type Mod interface {
	Key() string
	Expiry() int
	Event() glog.Event
	SetEvent(glog.Event)
	AffectedByHitlag() bool
	Extend(string, glog.Logger, int, int)
}

type Base struct {
	ModKey    string
	Dur       int
	Hitlag    bool
	ModExpiry int
	extension int
	event     glog.Event
}

func (t *Base) Key() string             { return t.ModKey }
func (t *Base) Expiry() int             { return t.ModExpiry + t.extension }
func (t *Base) Event() glog.Event       { return t.event }
func (t *Base) SetEvent(evt glog.Event) { t.event = evt }
func (t *Base) AffectedByHitlag() bool  { return t.Hitlag }
func (t *Base) Extend(key string, logger glog.Logger, index int, amt int) {
	t.extension += amt
	if t.extension < 0 {
		t.extension = 0
	}
	t.event.SetEnded(t.Expiry())
	logger.NewEvent("mod extended", glog.LogStatusEvent, -1).
		Write("key", key).
		Write("amt", amt).
		Write("expiry", t.Expiry())
}
func (t *Base) SetExpiry(f int) {
	if t.Dur == -1 {
		t.ModExpiry = -1
	} else {
		t.ModExpiry = f + t.Dur
	}
}

func NewBase(key string, dur int) Base {
	return Base{
		ModKey: key,
		Dur:    dur,
	}
}

func NewBaseWithHitlag(key string, dur int) Base {
	return Base{
		ModKey: key,
		Dur:    dur,
		Hitlag: true,
	}
}

// Delete removes a modifier. Returns true if deleted ok
func Delete[K Mod](slice *[]K, key string) (m Mod) {
	n := 0
	for i, v := range *slice {
		if v.Key() == key {
			m = (*slice)[i]
		} else {
			(*slice)[n] = v
			n++
		}
	}
	*slice = (*slice)[:n]
	return
}

// Add adds a modifier. Returns true if overwritten and the original evt (if overwritten)
// TODO: consider adding a map here to track the index to assist with faster lookups
func Add[K Mod](slice *[]K, mod K, f int) (overwrote bool, evt glog.Event) {
	ind := Find(slice, mod.Key())

	//if does not exist, make new and add
	if ind == -1 {
		*slice = append(*slice, mod)
		return
	}

	//otherwise check not expired
	if (*slice)[ind].Expiry() > f || (*slice)[ind].Expiry() == -1 {
		overwrote = true
		evt = (*slice)[ind].Event()
	}
	(*slice)[ind] = mod

	return
}

func Find[K Mod](slice *[]K, key string) int {
	ind := -1
	for i, v := range *slice {
		if v.Key() == key {
			ind = i
		}
	}
	return ind
}

func FindCheckExpiry[K Mod](slice *[]K, key string, f int) (int, bool) {
	ind := Find(slice, key)
	if ind == -1 {
		return ind, false
	}
	if (*slice)[ind].Expiry() < f && (*slice)[ind].Expiry() > -1 {
		return ind, false
	}
	return ind, true
}

// LogAdd is a helper that logs mod add events
func LogAdd[K Mod](prefix string, index int, mod K, logger glog.Logger, overwrote bool, oldEvt glog.Event) {
	var evt glog.Event
	if overwrote {
		logger.NewEventBuildMsg(
			glog.LogStatusEvent, index,
			prefix, " mod refreshed",
		).Write(
			"overwrite", true,
		).Write(
			"key", mod.Key(),
		).Write(
			"expiry", mod.Expiry(),
		)
		evt = oldEvt
	} else {
		evt = logger.NewEventBuildMsg(
			glog.LogStatusEvent, index,
			prefix, " mod added",
		).Write(
			"overwrite", false,
		).Write(
			"key", mod.Key(),
		).Write(
			"expiry", mod.Expiry(),
		)
	}
	evt.SetEnded(mod.Expiry())
	mod.SetEvent(evt)
}
