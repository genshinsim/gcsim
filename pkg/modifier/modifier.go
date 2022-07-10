//package modifier provides a universal way of handling a slice
//of modifiers
package modifier

import "github.com/genshinsim/gcsim/pkg/core/glog"

type Mod interface {
	Key() string
	Expiry() int
	Event() glog.Event
	SetEvent(glog.Event)
}

type Base struct {
	key    string
	expiry int
	event  glog.Event
}

func (t *Base) Key() string             { return t.key }
func (t *Base) Expiry() int             { return t.expiry }
func (t *Base) Event() glog.Event       { return t.event }
func (t *Base) SetEvent(evt glog.Event) { t.event = evt }

func NewBase(key string, expiry int) Base {
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
			log.NewEvent("enemy mod deleted", glog.LogStatusEvent, -1).
				Write("key", key)
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
		evt := log.NewEvent("enemy mod added", glog.LogStatusEvent, -1).
			Write("overwrite", false).
			Write("key", mod.Key()).
			Write("expiry", mod.Expiry())
		evt.SetEnded(mod.Expiry())
		mod.SetEvent(evt)
		*slice = append(*slice, mod)
		return
	}

	//otherwise check not expired
	var evt glog.Event
	if (*slice)[ind].Expiry() > f || (*slice)[ind].Expiry() == -1 {
		evt = log.NewEvent("enemy mod refreshed", glog.LogStatusEvent, -1).
			Write("overwrite", true).
			Write("key", mod.Key()).
			Write("expiry", mod.Expiry())

	} else {
		//if expired overide the event
		evt = log.NewEvent("enemy mod added", glog.LogStatusEvent, -1).
			Write("overwrite", false).
			Write("key", mod.Key()).
			Write("expiry", mod.Expiry())
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
