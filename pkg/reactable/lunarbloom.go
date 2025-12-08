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
)

func (r *Reactable) lunarBloom() bool {
	if r.dewCount >= maxVerdantDew {
		return false
	}
	r.core.Tasks.Add(func() {
		r.addDew(verdantDewEndFrame - verdantDewStartFrame)
		r.core.Log.NewEvent("verdant dew generation extended", glog.LogElementEvent, -1).
			Write("nextDewFrame", r.nextDewFrame)
	}, verdantDewStartFrame)
	return true
}

func (r *Reactable) addDew(amount int) {
	if r.dewCount >= maxVerdantDew {
		return
	}
	currentFrame := r.core.F
	r.nextDewFrame -= min(amount, currentFrame-r.lastLBTick)
	if r.nextDewFrame <= 0 {
		r.core.Tasks.Add(func() {
			r.dewCount++
			if r.dewCount >= maxVerdantDew {
				r.nextDewFrame = nextDewFrameRefresh
			}
			r.core.Log.NewEvent(fmt.Sprintf("lunar bloom dew gained: %v/%v", r.dewCount, maxVerdantDew), glog.LogElementEvent, -1)
		}, amount+r.nextDewFrame)
		r.nextDewFrame += nextDewFrameRefresh
	}

	r.lastLBTick = currentFrame

	if r.nextDewFrame < 0 {
		r.nextDewFrame = nextDewFrameRefresh
	}
}
