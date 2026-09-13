package combat

import (
	"math"

	"github.com/genshinsim/gcsim/pkg/core/info"
)

var gadgetLimits []int

func init() {
	gadgetLimits = make([]int, info.EndGadgetTyp)
	gadgetLimits[info.GadgetTypDendroCore] = 5
	gadgetLimits[info.GadgetTypTest] = 2
	gadgetLimits[info.GadgetTypLeaLotus] = 1
	gadgetLimits[info.GadgetTypYueguiThrowing] = 2
	gadgetLimits[info.GadgetTypYueguiJumping] = 3
	gadgetLimits[info.GadgetTypSourcewaterDropletHydroTrav] = 12
	gadgetLimits[info.GadgetTypSourcewaterDropletNeuv] = 12
	gadgetLimits[info.GadgetTypSourcewaterDropletSigewinne] = 12
	gadgetLimits[info.GadgetTypCrystallizeShard] = 3
}

func (h *Handler) RemoveGadget(key info.TargetKey) {
	h.ReplaceGadget(key, nil)
}

func (h *Handler) AddGadget(t info.Gadget) {
	// check for hard coded limit
	if gadgetLimits[t.GadgetTyp()] > 0 {
		// should kill oldest one if > limit
		f := math.MaxInt
		oldest := -1
		count := 0
		for i, v := range h.gadgets {
			if v == nil || v.GadgetTyp() != t.GadgetTyp() {
				continue
			}
			count++
			if v.Src() < f {
				f = v.Src()
				oldest = i
			}
		}
		if count == gadgetLimits[t.GadgetTyp()] {
			h.gadgets[oldest].Kill()
		}
	}
	h.gadgets = append(h.gadgets, t)
	t.SetKey(h.nextkey())
}

func (h *Handler) ReplaceGadget(key info.TargetKey, t info.Gadget) {
	// do nothing if not found
	for i, v := range h.gadgets {
		if v != nil && v.Key() == key {
			h.gadgets[i] = t
		}
	}
}

func (h *Handler) Gadget(i int) info.Gadget {
	return h.gadgets[i]
}

func (h *Handler) Gadgets() []info.Gadget {
	return h.gadgets
}

func (h *Handler) GadgetCount() int {
	count := 0
	for _, v := range h.gadgets {
		if v != nil {
			count++
		}
	}

	return count
}
