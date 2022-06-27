package character

//ApplyHitlag adds hitlag to the character for specified duration
func (c *Character) ApplyHitlag(factor float64, dur int) {
	c.hitlagFactor = factor
	c.hitlagUntil = c.Core.F + dur
	// if c.Core.Flags.LogDebug {
	// 	c.Core.Log.NewEvent(fmt.Sprintf("hitlag applied: %v", dur), glog.LogHitlagEvent, c.Index, "duration", dur, "factor", factor).
	// 		SetEnded(c.Core.F + dur)
	// }
}

func (c *Character) FramePausedOnHitlag() bool {
	return c.hitlagUntil >= c.Core.F
}
