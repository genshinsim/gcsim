package character

import "github.com/genshinsim/gcsim/pkg/core/glog"

type mod interface {
	Key() string
	Expiry() int
	Event() glog.Event
	SetEvent(glog.Event)
}

type modTmpl struct {
	key    string
	expiry int
	event  glog.Event
}

func (m *modTmpl) Key() string             { return m.key }
func (m *modTmpl) Expiry() int             { return m.expiry }
func (m *modTmpl) Event() glog.Event       { return m.event }
func (m *modTmpl) SetEvent(evt glog.Event) { m.event = evt }

func deleteMod[K mod](c *CharWrapper, slice *[]K, key string) {
	n := 0
	for _, v := range *slice {
		if v.Key() == key {
			v.Event().SetEnded(*c.f)
			c.log.NewEvent("mod deleted", glog.LogStatusEvent, c.Index).
				Write("key", key)
		} else {
			(*slice)[n] = v
			n++
		}
	}
	//BUG: i think this needs to be *K here. otherwise delete wont work?
	*slice = (*slice)[:n]
}

func addMod[K mod](c *CharWrapper, slice *[]K, mod K) {
	ind := findMod(slice, mod.Key())

	//if does not exist, make new and add
	if ind == -1 {
		evt := c.log.NewEvent("mod added", glog.LogStatusEvent, c.Index).
			Write("overwrite", false).
			Write("key", mod.Key()).
			Write("expiry", mod.Expiry())
		evt.SetEnded(mod.Expiry())
		mod.SetEvent(evt)
		//BUG: i think this needs to be *K here. otherwise delete wont work?
		*slice = append(*slice, mod)
		return
	}

	//otherwise check not expired
	var evt glog.Event
	if (*slice)[ind].Expiry() > *c.f || (*slice)[ind].Expiry() == -1 {
		evt = c.log.NewEvent("mod refreshed", glog.LogStatusEvent, c.Index).
			Write("overwrite", true).
			Write("key", mod.Key()).
			Write("expiry", mod.Expiry())

	} else {
		//if expired overide the event
		evt = c.log.NewEvent("mod added", glog.LogStatusEvent, c.Index).
			Write("overwrite", false).
			Write("key", mod.Key()).
			Write("expiry", mod.Expiry())
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
