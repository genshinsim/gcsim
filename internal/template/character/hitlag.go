package character

//ApplyHitlag adds hitlag to the character for specified duration
func (c *Character) ApplyHitlag(factor float64, dur int) {
	c.hitlagFactor = factor
	c.hitlagUntil = c.Core.F + dur
}

func (c *Character) FramePausedOnHitlag() bool {
	return c.hitlagUntil > c.Core.F
}
