package yaemiko

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1Key          = "yae-c1"
	c2Key          = "yae-c2"
	c4Key          = "yae-c4"
	c4EnergyIcdKey = c4Key + "-icd"
	c6Key          = "yae-c6"
)

var c2BuffVal = []float64{0, 60, 90, 120, 200}

// After Yae Miko triggers a Superconduct or Stellar-Conduct reaction, or deals Stellar-Conduct DMG,
// nearby party members will gain a 50% Electro DMG and Stellar-Conduct DMG Bonus for 10s. Triggering
// these reactions again will refresh the duration of the DMG bonuses.
func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}

	if !c.revelation {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.ElectroP] = 0.5

	electroStatMod := character.StatMod{
		Base:         modifier.NewBaseWithHitlag(c1Key+"-electro", 10*60),
		AffectedStat: attributes.ElectroP,
		Amount: func() []float64 {
			return m
		},
	}

	sscReactMod := character.ReactBonusMod{
		Base: modifier.NewBaseWithHitlag(c1Key+"-ssc", 10*60),
		Amount: func(ai info.AttackInfo) float64 {
			if ai.AttackTag == attacks.AttackTagDirectStellarConduct {
				return 0.5
			}
			return 0
		},
	}

	gainBuffs := func() {
		for _, char := range c.Core.Player.Chars() {
			char.AddStatMod(electroStatMod)
			char.AddReactBonusMod(sscReactMod)
		}
	}

	hook := func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != c.Index() {
			return
		}

		gainBuffs()
	}

	hookDmg := func(args ...any) {
		if _, ok := args[0].(*enemy.Enemy); !ok {
			return
		}

		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != c.Index() {
			return
		}

		if atk.Info.AttackTag != attacks.AttackTagDirectStellarConduct {
			return
		}

		gainBuffs()
	}

	c.Core.Events.Subscribe(event.OnSuperconduct, hook, c1Key)
	c.Core.Events.Subscribe(event.OnStellarConduct, hook, c1Key)
	c.Core.Events.Subscribe(event.OnEnemyDamage, hookDmg, c1Key)
}

// Sesshou Sakura start at Level 2 when created, their max level is increased to 4, and their attack
// range is increased by 60%.
//
// Additionally, when there are Sesshou Sakura in the field, Yae Miko's and your current active
// character's Elemental Mastery will also be increased by 60/90/120/200 points, depending on
// the level of the Sesshou Sakura.
func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		c.kitsuneDetectionRadius = 12.5
		return
	}

	c.kitsuneDetectionRadius = 20

	if !c.revelation {
		return
	}
	m := make([]float64, attributes.EndStatType)

	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(c2Key, -1),
			AffectedStat: attributes.EM,
			Amount: func() []float64 {
				if c.Index() != char.Index() && c.Core.Player.Active() != char.Index() {
					return nil
				}

				level := c.sakuraLevel()
				if level == 0 {
					return nil
				}
				m[attributes.EM] = c2BuffVal[level]
				return m
			},
		})
	}
}

func (c *char) c2SakuraLevelBonus() int {
	if c.Base.Cons < 2 {
		return 0
	}
	return 1
}

// Additionally, Yae Miko's Elemental Burst DMG is increased by 100%.
func (c *char) c4Init() {
	if c.Base.Cons < 4 {
		return
	}

	c.c4buff = make([]float64, attributes.EndStatType)
	c.c4buff[attributes.ElectroP] = .20

	if !c.revelation {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.DmgP] = 1
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(c4Key, -1),
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			if atk.Info.AttackTag == attacks.AttackTagElementalBurst {
				return m
			}
			return nil
		},
	})
}

// When Sesshou Sakura lightning hits opponents, the Electro DMG Bonus of all nearby party members is increased by 20% for 5s,
// and Yae Miko regenerates 8 Elemental Energy.
// The aforementioned Elemental Energy recovery effect can trigger once every 5s.
func (c *char) c4MakeCB() info.AttackCBFunc {
	if c.Base.Cons < 4 {
		return nil
	}

	done := false
	return func(a info.AttackCB) {
		if a.Target.Type() != info.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true
		c.c4()
	}
}

func (c *char) c4() {
	// only called on c4MakeCB, so no need to check cons

	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("yaemiko-c4", 5*60),
			AffectedStat: attributes.ElectroP,
			Amount: func() []float64 {
				return c.c4buff
			},
		})
	}

	if !c.revelation {
		return
	}

	if c.StatusIsActive(c4EnergyIcdKey) {
		return
	}
	c.AddStatus(c4EnergyIcdKey, 5*60, true)
	c.AddEnergy(c4Key, 8)
}

// Sesshou Sakura's attacks will ignore 60% of the opponents' DEF. The CRIT DMG of Yae Miko's
// Stellar-Conduct DMG is increased by 200%.
func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}
	if !c.revelation {
		return
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.CD] = 2 // 200% CDMG
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBase(c6Key, -1),
		Amount: func(atk *info.AttackEvent, _ info.Target) []float64 {
			switch atk.Info.AttackTag {
			case attacks.AttackTagDirectStellarConduct:
			default:
				return nil
			}
			return m
		},
	})
}

func (c *char) c6DefIgnore() float64 {
	if c.Base.Cons < 6 {
		return 0
	}
	return 0.6
}
