package sucrose

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

// Handles C4: Every 7 Normal and Charged Attacks, Sucrose will reduce the CD of Astable Anemohypostasis Creation-6308 by 1-7s
func (c *char) c4() {
	c.c4Count++
	if c.c4Count < 7 {
		return
	}
	c.c4Count = 0

	// Change can be in float. See this Terrapin video for example
	// https://youtu.be/jB3x5BTYWIA?t=20
	cdReduction := 60 * int(c.Core.Rand.Float64()*6+1)

	//we simply reduce the action cd
	c.ReduceActionCooldown(action.ActionSkill, cdReduction)

	c.Core.Log.NewEvent("sucrose c4 reducing E CD", glog.LogCharacterEvent, c.Index).
		Write("cd_reduction", cdReduction)
}

func (c *char) c6() {
	m := make([]float64, attributes.EndStatType)
	stat := attributes.EleToDmgP(c.qInfused)
	m[stat] = .20

	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod("sucrose-c6", 60*10, stat, func() ([]float64, bool) {
			return m, true
		})
	}
}
