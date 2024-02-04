package dehya

import (
	"errors"

	"github.com/genshinsim/gcsim/pkg/core/action"
)

func (c *char) Jump(p map[string]int) (action.Info, error) {
	// jumping in jump kick window leads to kick being done
	if c.StatusIsActive(burstKey) && c.StatusIsActive(jumpKickWindowKey) {
		c.burstHitSrc++                   // invalidate any other punch/kick tasks
		c.DeleteStatus(jumpKickWindowKey) // delete window
		return c.burstKick(c.burstHitSrc), nil
	}

	// jumping during a kick is not allowed
	// TODO: not sure if this can even happen...
	if c.StatusIsActive(kickKey) {
		return action.Info{}, errors.New("can't jump cancel burst kick")
	}

	// if burst is active at time of jump and it was not during jump kick window
	if c.StatusIsActive(burstKey) {
		c.burstHitSrc = -1       // invalidate any other punch/kick tasks
		c.DeleteStatus(burstKey) // delete burst
		// place field
		if dur := c.sanctumSavedDur; dur > 0 {
			c.sanctumSavedDur = 0
			c.addField(dur)
		}
	}

	return c.Character.Jump(p)
}
