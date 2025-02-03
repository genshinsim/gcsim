package chasca

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
)

var skillDashFrames []int

func init() {
	skillDashFrames = frames.InitAbilSlice(29)
	skillDashFrames[action.ActionSkill] = 6
	skillDashFrames[action.ActionBurst] = 25
}
func (c *char) Dash(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		c.reduceNightsoulPoints(13.3)
		d, e := c.Character.Dash(p)
		d.Frames = c.skillNextFrames(d.Frames)
		d.CanQueueAfter = 1 // can run out of nightsoul and immediately start falling
		return d, e
	}
	return c.Character.Dash(p)
}
