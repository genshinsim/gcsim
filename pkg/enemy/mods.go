package enemy

import (
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

type mod interface {
	Key() string
	Expiry() int
	Event() glog.Event
	SetEvent(glog.Event)
}

type tmpl struct {
	key    string
	expiry int
	event  glog.Event
}

func (t *tmpl) Key() string             { return t.key }
func (t *tmpl) Expiry() int             { return t.expiry }
func (t *tmpl) Event() glog.Event       { return t.event }
func (t *tmpl) SetEvent(evt glog.Event) { t.event = evt }

func deleteMod[K mod](c *Enemy, slice *[]K, key string) {
	n := 0
	for _, v := range *slice {
		if v.Key() == key {
			v.Event().SetEnded(c.Core.F)
			c.Core.Log.NewEvent("enemy mod deleted", glog.LogStatusEvent, -1, "key", key)
		} else {
			(*slice)[n] = v
			n++
		}
	}
	*slice = (*slice)[:n]
}

func addMod[K mod](c *Enemy, slice *[]K, mod K) {
	ind := findMod(slice, mod.Key())

	//if does not exist, make new and add
	if ind == -1 {
		evt := c.Core.Log.NewEvent(
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
	if (*slice)[ind].Expiry() > c.Core.F || (*slice)[ind].Expiry() == -1 {
		evt = c.Core.Log.NewEvent(
			"enemy mod refreshed", glog.LogStatusEvent, -1,
			"overwrite", true,
			"key", mod.Key(),
			"expiry", mod.Expiry(),
		)

	} else {
		//if expired overide the event
		evt = c.Core.Log.NewEvent(
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

func findMod[K mod](slice *[]K, key string) int {
	ind := -1
	for i, v := range *slice {
		if v.Key() == key {
			ind = i
		}
	}
	return ind
}

func findModCheckExpiry[K mod](slice *[]K, key string, f int) (int, bool) {
	ind := findMod(slice, key)
	if ind == -1 {
		return ind, false
	}
	if (*slice)[ind].Expiry() < f && (*slice)[ind].Expiry() > -1 {
		return ind, false
	}
	return ind, true
}
