package character

import (
	"fmt"
	"math"

	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/queue"
)

func (c *CharWrapper) QueueCharTask(f func(), delay int) {
	queue.Add(&c.queue, f, c.timePassed+delay)
}

func (c *CharWrapper) Tick() {
	// decrement frozen time first
	c.frozenFrames -= 1
	left := 0
	if c.frozenFrames < 0 {
		left = -c.frozenFrames
		c.frozenFrames = 0
	}
	// if any left then increase time passed
	if left <= 0 {
		// do nothing this tick
		return
	}
	c.timePassed += left

	// check char queue for any executable actions
	queue.Run(&c.queue, c.timePassed)
}

func (c *CharWrapper) FramePausedOnHitlag() bool {
	return c.frozenFrames > 0
}

// ApplyHitlag adds hitlag to the character for specified duration
func (c *CharWrapper) ApplyHitlag(factor, dur float64) {
	// number of frames frozen is total duration * (1 - factor)
	ext := int(math.Ceil(dur * (1 - factor)))
	c.frozenFrames += ext
	var logs []string
	var evt glog.Event
	if c.debug {
		logs = make([]string, 0, len(c.mods))
		evt = c.log.NewEvent(
			fmt.Sprintf("hitlag applied to char: %.3f", dur),
			glog.LogHitlagEvent, c.Index,
		).
			Write("duration", dur).
			Write("factor", factor).
			Write("frozen_frames", c.frozenFrames).
			SetEnded(*c.f + int(math.Ceil(dur)))
	}

	for i, v := range c.mods {
		if v.AffectedByHitlag() && v.Expiry() != -1 && v.Expiry() > *c.f {
			mod := c.mods[i]
			mod.Extend(mod.Key(), c.log, c.Index, ext)
			if c.debug {
				logs = append(logs, fmt.Sprintf("%v: %v", v.Key(), v.Expiry()))
			}
		}
	}

	if c.debug {
		evt.Write("mods affected", logs)
	}
}
