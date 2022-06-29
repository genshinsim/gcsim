package sara

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// Implements C1 CD reduction. Waits until delay (when it hits the enemy), then procs the effect
// Triggers on her E and Q
func (c *char) c1() {
	if c.Core.F < c.c1LastProc {
		return
	}
	c.c1LastProc = c.Core.F + 180
	c.ReduceActionCooldown(action.ActionSkill, 60)
	c.Core.Log.NewEvent("c1 reducing skill cooldown", glog.LogCharacterEvent, c.Index, "new_cooldown", c.Cooldown(action.ActionSkill))
}

// The Electro DMG of characters who have had their ATK increased by Tengu Juurai has its Crit DMG increased by 60%.
func (c *char) c6(char *character.CharWrapper) {
	m := make([]float64, attributes.EndStatType)
	m[attributes.CD] = 0.6

	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("sara-c6", 360),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			if atk.Info.Element != attributes.Electro {
				return nil, false
			}
			return m, true
		},
	})
}
