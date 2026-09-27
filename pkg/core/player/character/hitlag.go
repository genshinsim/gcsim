package character

import (
	"fmt"
	"math"

	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (c *CharWrapper) QueueCharTask(f func(), delay int) {
	if delay <= c.frozenFrames {
		f()
		return
	}
	// align char tasks with mods. Mods tick down if added during hitlag, so
	// char tasks should also tick down during hitlag
	c.queue.Add(f, delay-c.frozenFrames)
}

func (c *CharWrapper) Tick() {
	if c.frozenFrames > 0 {
		// frozen for this frame, do nothing
		c.frozenFrames--
		return
	}
	c.TimePassed++

	// check char queue for any executable actions
	c.queue.Run()
}

func (c *CharWrapper) FramePausedOnHitlag() bool {
	return c.frozenFrames > 0
}

// ApplyHitlag adds hitlag to the character for specified duration
func (c *CharWrapper) ApplyHitlag(factor, dur float64) {
	// number of frames frozen is total duration * (1 - factor)
	newHitlag := int(math.Ceil(dur * (1 - factor)))

	// TODO: this is inaccurate for overlapping hitlags of different hitlag factors
	oldFrozen := c.frozenFrames
	c.frozenFrames = max(newHitlag, c.frozenFrames)
	ext := c.frozenFrames - oldFrozen

	var logs []string
	var evt glog.Event
	if c.debug {
		logs = make([]string, 0, len(c.mods))
		evt = c.log.NewEvent(
			fmt.Sprintf("hitlag applied to char: %.3f", dur),
			glog.LogHitlagEvent, c.Index(),
		).
			Write("duration", dur).
			Write("factor", factor).
			Write("frozen_frames", c.frozenFrames).
			SetEnded(*c.f + int(math.Ceil(dur)))
	}

	for i, v := range c.mods {
		if v.AffectedByHitlag() && v.Expiry() != -1 && v.Expiry() > *c.f {
			mod := c.mods[i]
			mod.Extend(mod.Key(), c.log, c.Index(), ext)
			if c.debug {
				logs = append(logs, fmt.Sprintf("%v: %v", v.Key(), v.Expiry()))
			}
		}
	}

	if c.debug {
		evt.Write("mods affected", logs)
	}
}
