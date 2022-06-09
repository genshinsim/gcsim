package character

import "github.com/genshinsim/gcsim/pkg/core/glog"

//ApplyHitlag adds hitlag to the character for specified duration
func (c *Character) ApplyHitlag(factor float64, dur int) {
	c.hitlagFactor = factor
	c.hitlagUntil = c.Core.F + dur
	c.Core.Log.NewEvent("hitlag applied", glog.LogCharacterEvent, c.Index, "duration", dur, "factor", factor).
		SetEnded(c.Core.F + dur)
}

func (c *Character) FramePausedOnHitlag() bool {
	return c.hitlagUntil >= c.Core.F
}
