package dori

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// The number of After-Sales Service Rounds created by Troubleshooter Shots is increased by 1.
func (c *char) c1() {
	c.afterCount++
}

// When you are in combat and the Jinni heals the character it is connected to,
// it will fire a Jinni Toop from that character's position that deals 50% of Dori's ATK DMG.
func (c *char) c2(travel int) {
	c.Core.Events.Subscribe(event.OnHeal, func(args ...interface{}) bool {
		src := args[0].(int)
		if src != c.Index {
			return false
		}
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Special Franchise",
			AttackTag:  combat.AttackTagNone,
			ICDTag:     combat.ICDTagDoriC2,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       0.5,
		}
		c.Core.QueueAttack(ai, combat.NewDefSingleTarget(c.Core.Combat.DefaultTarget, combat.TargettableEnemy), 0, travel)
		return false
	}, "dori-c2")
}

// The character connected to the Jinni will obtain the following buffs based on their current HP and Energy:
// ·When their HP is lower than 50%, they gain 50% Incoming Healing Bonus.
// ·When their Energy is less than 50%, they gain 30% Energy Recharge.
func (c *char) c4() {
	// c4 should last for the duration of the burst
	// lasts 12.5 second, ticks every 0.5s; adds mod to active char for 2s
	for i := 30; i < 750; i += 30 {
		c.Core.Tasks.Add(func() {
			active := c.Core.Player.ActiveChar()
			if active.HPCurrent/active.MaxHP() < 0.5 {
				active.AddHealBonusMod(character.HealBonusMod{
					Base: modifier.NewBaseWithHitlag("dori-c4-healbonus", 120),
					Amount: func() (float64, bool) {
						// is this log even needed?
						c.Core.Log.NewEvent("dori c4 incomming heal bonus activated", glog.LogCharacterEvent, c.Index)
						return 0.5, false
					},
				})
				c.Tags["c4bonus-"+active.Base.Key.String()] = c.Core.F + 120
			}
			// add energy recharge
			if active.Energy/active.EnergyMax < 0.5 {
				erMod := make([]float64, attributes.EndStatType)
				erMod[attributes.ER] = 0.3
				active.AddStatMod(character.StatMod{
					Base:         modifier.NewBase("dori-c4-er-bonus", 120),
					AffectedStat: attributes.ER,
					Amount: func() ([]float64, bool) {
						return erMod, false
					},
				})
			}
		}, i)
	}
}

// Dori gains the following effects for 3s after using Spirit-Warding Lamp: Troubleshooter Cannon:
// ·Electro Infusion.
// ·When Normal Attacks hit opponents, all nearby party members will heal HP equivalent to 4% of Dori's Max HP.
// This type of healing can occur once every 0.1s.
func (c *char) c6() {
	const c6icd = "dori-c6-heal-icd"
	c.Core.Events.Subscribe(event.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		if atk.Info.ActorIndex != c.Index {
			return false
		}
		if !c.Core.Player.WeaponInfuseIsActive(c.Index, c6key) {
			return false
		}
		if c.StatusIsActive(c6icd) {
			return false
		}
		if atk.Info.AttackTag != combat.AttackTagNormal {
			return false
		}

		c.AddStatus(c6icd, 6, false) // 0.1s*60 icd
		// heal party members

		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  -1,
			Message: "Dori-c6-onNormal",
			Src:     0.04 * c.MaxHP(),
			Bonus:   c.Stat(attributes.Heal),
		})

		return false
	}, "dori-c6")
}
