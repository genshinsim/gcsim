package sucrose

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// Handles C4: Every 7 Normal and Charged Attacks, Sucrose will reduce the CD of Astable Anemohypostasis Creation-6308 by 1-7s
func (c *char) makeC4Callback() func(combat.AttackCB) {
	if c.Base.Cons < 4 {
		return nil
	}
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true

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
}

func (c *char) c6() {
	stat := attributes.EleToDmgP(c.qAbsorb)
	c.c6buff[stat] = .20

	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("sucrose-c6", 60*10),
			AffectedStat: stat,
			Amount: func() ([]float64, bool) {
				return c.c6buff, true
			},
		})
	}
}
