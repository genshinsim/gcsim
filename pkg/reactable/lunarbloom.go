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
	verdantDewKey        = "verdant-dew"
	nextDewFrameKey      = "verdant-dew-next"
	lastLBTickKey        = "last-lb-frame"
)

func (r *Reactable) lunarBloom() {
	if r.core.Flags.Custom[verdantDewKey] >= maxVerdantDew {
		return
	}
	r.core.Tasks.Add(func() {
		r.addDew(verdantDewEndFrame - verdantDewStartFrame)
		r.core.Log.NewEvent("verdant dew generation extended", glog.LogElementEvent, -1).
			Write("nextDewFrame", r.core.Flags.Custom[nextDewFrameKey])
	}, verdantDewStartFrame)
}

func (r *Reactable) addDew(amount int) {
	if r.core.Flags.Custom[verdantDewKey] >= maxVerdantDew {
		return
	}
	currentFrame := r.core.F
	r.core.Flags.Custom[nextDewFrameKey] -= float64(min(amount, currentFrame-int(r.core.Flags.Custom[lastLBTickKey])))
	if r.core.Flags.Custom[nextDewFrameKey] <= 0 {
		r.core.Tasks.Add(func() {
			r.core.Flags.Custom[verdantDewKey]++
			if r.core.Flags.Custom[verdantDewKey] >= maxVerdantDew {
				r.core.Flags.Custom[nextDewFrameKey] = nextDewFrameRefresh
			}
			r.core.Log.NewEvent(fmt.Sprintf("lunar bloom dew gained: %v/%v", r.core.Flags.Custom[verdantDewKey], maxVerdantDew), glog.LogElementEvent, -1)
		}, amount+int(r.core.Flags.Custom[nextDewFrameKey]))
		r.core.Flags.Custom[nextDewFrameKey] += nextDewFrameRefresh
	}

	r.core.Flags.Custom[lastLBTickKey] = float64(currentFrame)

	if r.core.Flags.Custom[nextDewFrameKey] < 0 {
		r.core.Flags.Custom[nextDewFrameKey] = nextDewFrameRefresh
	}
}
