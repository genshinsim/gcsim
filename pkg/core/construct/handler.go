package construct

import (
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

type Handler struct {
	constructs  []Construct
	consNoLimit []Construct
	log         glog.Logger
	f           *int
}

func New(f *int, log glog.Logger) *Handler {
	return &Handler{
		constructs:  make([]Construct, 0, 3),
		consNoLimit: make([]Construct, 0, 3),
		f:           f,
		log:         log,
	}
}

func (h *Handler) New(c Construct, refresh bool) {
	h.NewConstruct(c, refresh, &h.constructs, true)
}

func (h *Handler) NewNoLimitCons(c Construct, refresh bool) {
	h.NewConstruct(c, refresh, &h.consNoLimit, false)
}

func (h *Handler) NewConstruct(c Construct, refresh bool, constructs *[]Construct, hasLimit bool) {
	//if refresh, we nil out the old one if any
	ind := -1
	if refresh {
		for i, v := range *constructs {
			if v.Type() == c.Type() {
				ind = i
			}
		}
	}
	if ind > -1 {
		h.log.NewEventBuildMsg(glog.LogConstructEvent, -1, "construct replaced - new: ", c.Type().String()).
			Write("key", (*constructs)[ind].Key()).
			Write("prev type", (*constructs)[ind].Type()).
			Write("next type", c.Type())
		//remove construct from list, reset order by removing nils and add construct to end
		(*constructs)[ind].OnDestruct()
		(*constructs)[ind] = nil
		h.cleanOutNils(constructs)
	} else {
		//add this one to the end
		(*constructs) = append((*constructs), c)
		h.log.NewEventBuildMsg(glog.LogConstructEvent, -1, "construct created: ", c.Type().String()).
			Write("key", c.Key()).
			Write("type", c.Type())
	}

	if hasLimit {
		//if length > 3, then destruct the beginning ones
		for i := 0; i < len((*constructs))-3; i++ {
			(*constructs)[i].OnDestruct()
			h.log.NewEventBuildMsg(glog.LogConstructEvent, -1, "construct destroyed: "+(*constructs)[i].Type().String()).
				Write("key", (*constructs)[i].Key()).
				Write("type", (*constructs)[i].Type())
			(*constructs)[i] = nil
		}
	}

	h.cleanOutNils(constructs)
}

func (h *Handler) cleanOutNils(constructs *[]Construct) {
	//clean out any nils
	n := 0
	for _, x := range *constructs {
		if x != nil {
			(*constructs)[n] = x
			n++
		}
	}
	(*constructs) = (*constructs)[:n]
}

func (h *Handler) Tick() {
	//clean out expired
	n := 0
	for _, v := range h.constructs {
		if v.Expiry() == *h.f {
			v.OnDestruct()
			h.log.NewEventBuildMsg(glog.LogConstructEvent, -1, "construct destroyed: "+v.Type().String()).
				Write("key", v.Key()).
				Write("type", v.Type())
		} else {
			h.constructs[n] = v
			n++
		}
	}
	h.constructs = h.constructs[:n]
	n = 0
	for i, v := range h.consNoLimit {
		if v.Expiry() == *h.f {
			h.consNoLimit[i].OnDestruct()
			h.log.NewEventBuildMsg(glog.LogConstructEvent, -1, "construct destroyed: "+v.Type().String()).
				Write("key", v.Key()).
				Write("type", v.Type())
		} else {
			h.consNoLimit[n] = v
			n++
		}
	}
	h.consNoLimit = h.consNoLimit[:n]

}

func (h *Handler) ConstructsByType(t GeoConstructType) ([]Construct, []Construct) {
	var match []Construct
	var notMatch []Construct
	for _, v := range h.constructs {
		if v.Type() == t {
			match = append(match, v)
		} else {
			notMatch = append(notMatch, v)
		}
	}
	for _, v := range h.consNoLimit {
		if v.Type() == t {
			match = append(match, v)
		} else {
			notMatch = append(notMatch, v)
		}
	}
	return match, notMatch
}

// how many of the given
func (h *Handler) Count() int {
	count := 0
	for _, v := range h.constructs {
		count += v.Count()
	}
	for _, v := range h.consNoLimit {
		count += v.Count()
	}
	return count
}

func (h *Handler) CountByType(t GeoConstructType) int {
	count := 0
	for _, v := range h.constructs {
		if v.Type() == t {
			count++
		}
	}
	for _, v := range h.consNoLimit {
		if v.Type() == t {
			count++
		}
	}
	return count
}

func (h *Handler) Has(key int) bool {
	for _, v := range h.constructs {
		if v.Key() == key {
			return true
		}
	}
	for _, v := range h.consNoLimit {
		if v.Key() == key {
			return true
		}
	}
	return false
}

func (h *Handler) Expiry(t GeoConstructType) int {
	expiry := -1
	for _, v := range h.constructs {
		if v.Type() == t {
			if expiry == -1 {
				expiry = v.Expiry()
			} else if expiry > v.Expiry() {
				expiry = v.Expiry()
			}
		}
	}
	for _, v := range h.consNoLimit {
		if v.Type() == t {
			if expiry == -1 {
				expiry = v.Expiry()
			} else if expiry > v.Expiry() {
				expiry = v.Expiry()
			}
		}
	}

	expiry = expiry - *h.f

	if expiry < 0 {
		return 0
	}

	return expiry
}

// destroy key if exist, return true if destroyed
func (h *Handler) Destroy(key int) bool {
	ok := false
	//clean out expired
	n := 0
	for _, v := range h.constructs {
		if v.Key() == key {
			v.OnDestruct()
			ok = true
			h.log.NewEventBuildMsg(glog.LogConstructEvent, -1, "construct destroyed: "+v.Type().String()).
				Write("key", v.Key()).
				Write("type", v.Type())
		} else {
			h.constructs[n] = v
			n++
		}
	}
	h.constructs = h.constructs[:n]
	if ok {
		return ok
	}
	n = 0
	for i, v := range h.consNoLimit {
		if v.Key() == key {
			h.consNoLimit[i].OnDestruct()
			ok = true
			h.log.NewEventBuildMsg(glog.LogConstructEvent, -1, "construct destroyed: "+v.Type().String()).
				Write("key", v.Key()).
				Write("type", v.Type())
		} else {
			h.consNoLimit[n] = v
			n++
		}
	}
	h.consNoLimit = h.consNoLimit[:n]
	return ok
}
