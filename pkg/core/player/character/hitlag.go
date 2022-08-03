package character

import (
	"fmt"
	"math"

	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (c *CharWrapper) QueueCharTask(f func(), delay int) {
	c.queue = append(c.queue, charTask{
		f:     f,
		delay: c.timePassed + float64(delay),
	})
}

func (c *CharWrapper) Tick() {
	//decrement frozen time first
	c.frozenFrames -= 1.0
	left := 0.0
	if c.frozenFrames < 0 {
		left = -c.frozenFrames
		c.frozenFrames = 0
	}
	//if any left then increase time passed
	if left <= 0 {
		//do nothing this tick
		return
	}
	c.timePassed += left

	//check char queue for any executable actions
	n := 0
	for i := 0; i < len(c.queue); i++ {
		if c.queue[i].delay <= c.timePassed {
			c.queue[i].f()
		} else {
			// keep the actions that can't be executed yet
			c.queue[n] = c.queue[i]
			n++
		}
	}
	// set char queue len to the remaining elements
	c.queue = c.queue[:n]
}

func (c *CharWrapper) FramePausedOnHitlag() bool {
	return c.frozenFrames > 0
}

//ApplyHitlag adds hitlag to the character for specified duration
func (c *CharWrapper) ApplyHitlag(factor, dur float64) {
	//number of frames frozen is total duration * (1 - factor)
	ext := dur * (1 - factor)
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
		if v.AffectedByHitlag() && v.Expiry() != -1 {
			c.mods[i].Extend(ext)
			if c.debug {
				logs = append(logs, fmt.Sprintf("%v: %v", v.Key(), v.Expiry()))
			}
		}
	}

	if c.debug {
		evt.Write("mods affected", logs)
	}

}
