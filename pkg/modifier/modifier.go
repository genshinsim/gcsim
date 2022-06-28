//package modifier provides a universal way of handling a slice
//of modifiers
package modifier

import "github.com/genshinsim/gcsim/pkg/core/glog"

type Mod interface {
	Key() string
	Expiry() int
	Event() glog.Event
	SetEvent(glog.Event)
	AffectedByHitlag() bool
}

type Base struct {
	key       string
	expiry    int
	extension float64
	event     glog.Event
	hitlag    bool
}

func (t *Base) Key() string             { return t.key }
func (t *Base) Expiry() int             { return t.expiry + int(t.extension) }
func (t *Base) Event() glog.Event       { return t.event }
func (t *Base) SetEvent(evt glog.Event) { t.event = evt }
func (t *Base) AffectedByHitlag() bool  { return t.hitlag }
func (t *Base) Extend(amt float64)      { t.extension += amt }

func NewBase(key string, expiry int, affectedByHitlag bool) Base {
	return Base{
		key:    key,
		expiry: expiry,
	}
}

func Delete[K Mod](f int, log glog.Logger, slice *[]K, key string) {
	n := 0
	for _, v := range *slice {
		if v.Key() == key {
			v.Event().SetEnded(f)
			log.NewEvent("enemy mod deleted", glog.LogStatusEvent, -1, "key", key)
		} else {
			(*slice)[n] = v
			n++
		}
	}
	*slice = (*slice)[:n]
}

func Add[K Mod](f int, log glog.Logger, slice *[]K, mod K) {
	ind := Find(slice, mod.Key())

	//if does not exist, make new and add
	if ind == -1 {
		evt := log.NewEvent(
			"enemy mod added", glog.LogStatusEvent, -1,
			"overwrite", false,
			"key", mod.Key(),
			"expiry", mod.Expiry(),
		)
		evt.SetEnded(mod.Expiry())
		mod.SetEvent(evt)
		*slice = append(*slice, mod)
		return
	}

	//otherwise check not expired
	var evt glog.Event
	if (*slice)[ind].Expiry() > f || (*slice)[ind].Expiry() == -1 {
		log.NewEvent(
			"enemy mod refreshed", glog.LogStatusEvent, -1,
			"overwrite", true,
			"key", mod.Key(),
			"expiry", mod.Expiry(),
		)
		evt = (*slice)[ind].Event()
	} else {
		//if expired overide the event
		evt = log.NewEvent(
			"enemy mod added", glog.LogStatusEvent, -1,
			"overwrite", false,
			"key", mod.Key(),
			"expiry", mod.Expiry(),
		)
	}
	mod.SetEvent(evt)
	evt.SetEnded(mod.Expiry())
	(*slice)[ind] = mod
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
