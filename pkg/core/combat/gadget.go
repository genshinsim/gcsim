package combat

import (
	"math"

	"github.com/genshinsim/gcsim/pkg/core/targets"
)

type GadgetTyp int

const (
	GadgetTypUnknown GadgetTyp = iota
	StartGadgetTypEnemy
	GadgetTypDendroCore
	GadgetTypLeaLotus
	EndGadgetTypEnemy
	GadgetTypGuoba
	GadgetTypYueguiThrowing
	GadgetTypYueguiJumping
	GadgetTypBaronBunny
	GadgetTypSourcewaterDroplet
	GadgetTypTest
	EndGadgetTyp
)

var gadgetLimits []int

func init() {
	gadgetLimits = make([]int, EndGadgetTyp)
	gadgetLimits[GadgetTypDendroCore] = 5
	gadgetLimits[GadgetTypTest] = 2
	gadgetLimits[GadgetTypLeaLotus] = 1
	gadgetLimits[GadgetTypYueguiThrowing] = 2
	gadgetLimits[GadgetTypYueguiJumping] = 3
	gadgetLimits[GadgetTypSourcewaterDroplet] = 4
}

type Gadget interface {
	Target
	Src() int
	GadgetTyp() GadgetTyp
}

func (h *Handler) RemoveGadget(key targets.TargetKey) {
	h.ReplaceGadget(key, nil)
}

func (h *Handler) AddGadget(t Gadget) {
	//check for hard coded limit
	if gadgetLimits[t.GadgetTyp()] > 0 {
		//should kill oldest one if > limit
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

func (h *Handler) ReplaceGadget(key targets.TargetKey, t Gadget) {
	//do nothing if not found
	for i, v := range h.gadgets {
		if v != nil && v.Key() == key {
			h.gadgets[i] = t
		}
	}
}

func (h *Handler) Gadget(i int) Gadget {
	return h.gadgets[i]
}

func (h *Handler) Gadgets() []Gadget {
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
