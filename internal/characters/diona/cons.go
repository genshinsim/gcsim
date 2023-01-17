package diona

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) c2() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = .15
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase("diona-c2", -1),
		Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
			return m, atk.Info.AttackTag == combat.AttackTagElementalArt
		},
	})
}

func (c *char) c6() {
	//c6 should last for the duration of the burst
	//lasts 12.5 second, ticks every 0.5s; adds mod to active char for 2s
	for i := 30; i <= 750; i += 30 {
		c.Core.Tasks.Add(func() {
			if !c.Core.Combat.Player().IsWithinArea(c.burstBuffArea) {
				return
			}
			//add 200EM to active char
			active := c.Core.Player.ActiveChar()
			if active.HPCurrent/active.MaxHP() > 0.5 {
				active.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag("diona-c6", 120),
					AffectedStat: attributes.EM,
					Amount: func() ([]float64, bool) {
						return c.c6buff, true
					},
				})
			} else {
				//add healing bonus if hp <= 0.5
				//bonus only lasts for 120 frames
				active.AddHealBonusMod(character.HealBonusMod{
					Base: modifier.NewBaseWithHitlag("diona-c6-healbonus", 120),
					Amount: func() (float64, bool) {
						// is this log even needed?
						c.Core.Log.NewEvent("diona c6 incomming heal bonus activated", glog.LogCharacterEvent, c.Index)
						return 0.3, false
					},
				})
				c.Tags["c6bonus-"+active.Base.Key.String()] = c.Core.F + 120
			}
		}, i)
	}
}
