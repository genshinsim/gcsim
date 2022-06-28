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

//Delete removes a modifier. Returns true if deleted ok
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

//Add adds a modifier. Returns true if overwritten and the original evt (if overwritten)
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
