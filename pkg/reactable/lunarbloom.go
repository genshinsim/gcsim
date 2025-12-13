package reactable

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/glog"
)

const (
	maxVerdantDew        = 3
	verdantDewStartFrame = 13
	verdantDewEndFrame   = 140
	nextDewFrameRefresh  = 150
	LbKey                = "lunar-bloom"
	verdantDewKey        = "verdant-dew"
	nextDewFrameKey      = "verdant-dew-next"
)

func (r *Reactable) lunarBloom() {
	if r.core.Flags.Custom[verdantDewKey] >= maxVerdantDew {
		return
	}

	if r.core.Status.Duration(LbKey) == 0 {
		r.core.Tasks.Add(r.addDew(), verdantDewStartFrame)
	}

	r.core.Status.Add(LcKey, verdantDewEndFrame)
}

func (r *Reactable) addDew() func() {
	return func() {
		if r.core.Flags.Custom[verdantDewKey] >= maxVerdantDew {
			return
		}

		r.core.Flags.Custom[nextDewFrameKey] -= 1

		if r.core.Flags.Custom[nextDewFrameKey] <= 0 {
			r.core.Flags.Custom[verdantDewKey] += 1
			r.core.Log.NewEvent(fmt.Sprintf("lunar bloom dew gained: %v/%v", r.core.Flags.Custom[verdantDewKey], maxVerdantDew), glog.LogElementEvent, -1)
			r.core.Flags.Custom[nextDewFrameKey] = nextDewFrameRefresh
		}

		r.core.Tasks.Add(r.addDew(), 1)
	}
}
