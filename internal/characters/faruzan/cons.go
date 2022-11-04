package faruzan

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c1ICDKey = "faruzan-c1-icd"

// Implements C1 CD reduction. Waits until delay (when it hits the enemy), then procs the effect
// Triggers on her E and Q
func (c *char) c1() {
	if c.StatusIsActive(c1ICDKey) {
		return
	}
	c.AddStatus(c1ICDKey, 180, true)
	c.ReduceActionCooldown(action.ActionSkill, 60)
	c.Core.Log.NewEvent("c1 reducing skill cooldown", glog.LogCharacterEvent, c.Index).
		Write("new_cooldown", c.Cooldown(action.ActionSkill))
}

// The Electro DMG of characters who have had their ATK increased by Tengu Juurai has its Crit DMG increased by 60%.
func (c *char) c6(char *character.CharWrapper) {
	char.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag("faruzan-c6", 360),
		Amount: func(atk *combat.AttackEvent, _ combat.Target) ([]float64, bool) {
			if atk.Info.Element != attributes.Electro {
				return nil, false
			}
			return c.c6buff, true
		},
	})
}
