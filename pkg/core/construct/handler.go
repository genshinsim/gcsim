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
	//if refresh, we nil out the old one if any
	ind := -1
	if refresh {
		for i, v := range h.constructs {
			if v.Type() == c.Type() {
				ind = i
			}
		}
	}
	if ind > -1 {
		h.log.NewEventBuildMsg(glog.LogConstructEvent, -1, "construct replaced - new: ", c.Type().String()).
			Write("key", h.constructs[ind].Key()).
			Write("prev type", h.constructs[ind].Type()).
			Write("next type", c.Type())
		//remove construct from list, reset order by removing nils and add construct to end
		h.constructs[ind].OnDestruct()
		h.constructs[ind] = nil
		h.cleanOutNils()
		h.constructs = append(h.constructs, c)
	}
	//add this one to the end
	h.constructs = append(h.constructs, c)
	h.log.NewEventBuildMsg(glog.LogConstructEvent, -1, "construct created: ", c.Type().String()).
		Write("key", c.Key()).
		Write("type", c.Type())

	//if length > 3, then destruct the beginning ones
	for i := 0; i < len(h.constructs)-3; i++ {
		h.constructs[i].OnDestruct()
		h.log.NewEventBuildMsg(glog.LogConstructEvent, -1, "construct destroyed: "+h.constructs[i].Type().String()).
			Write("key", h.constructs[i].Key()).
			Write("type", h.constructs[i].Type())
		h.constructs[i] = nil
	}

	h.cleanOutNils()
}

func (h *Handler) cleanOutNils() {
	//clean out any nils
	n := 0
	for _, x := range h.constructs {
		if x != nil {
			h.constructs[n] = x
			n++
		}
	}
	h.constructs = h.constructs[:n]
}

func (h *Handler) NewNoLimitCons(c Construct, refresh bool) {
	ind := -1
	if refresh {
		for i, v := range h.consNoLimit {
			if v.Key() == c.Key() {
				ind = i
			}
		}
	}
	if ind > -1 {
		//destroy the existing by setting expiry
		h.consNoLimit[ind].OnDestruct()
		h.log.NewEventBuildMsg(glog.LogConstructEvent, -1, "construct destroyed: "+h.consNoLimit[ind].Type().String()).
			Write("key", h.consNoLimit[ind].Key()).
			Write("type", h.consNoLimit[ind].Type())
		//remove construct from list, reset order by removing nils and add construct to end
		h.consNoLimit[ind].OnDestruct()
		h.consNoLimit[ind] = nil
		h.cleanOutNilsNoLimit()
		h.consNoLimit = append(h.consNoLimit, c)
	}
	//add this one to the end
	h.consNoLimit = append(h.consNoLimit, c)
	h.log.NewEventBuildMsg(glog.LogConstructEvent, -1, "construct created: ", c.Type().String()).
		Write("key", c.Key()).
		Write("type", c.Type())

	h.cleanOutNilsNoLimit()
}

func (h *Handler) cleanOutNilsNoLimit() {
	//clean out any nils
	n := 0
	for _, x := range h.consNoLimit {
		if x != nil {
			h.consNoLimit[n] = x
			n++
		}
	}
	h.consNoLimit = h.consNoLimit[:n]
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

func (h *Handler) Constructs() []Construct {
	var result []Construct
	result = append(result, h.constructs...)
	result = append(result, h.consNoLimit...)
	return result
}

func (h *Handler) ConstructsByType(t GeoConstructType) []Construct {
	var result []Construct
	for _, v := range h.constructs {
		if v.Type() == t {
			result = append(result, v)
		}
	}
	for _, v := range h.consNoLimit {
		if v.Type() == t {
			result = append(result, v)
		}
	}
	return result
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
