package zibai

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c1Key = "zibai-c1"
	c2Key = "zibai-c2"
	c4Key = "zibai-c4"
	c6Key = "zibai-c6"

	c2Mult = 5.5
)

// After using the Elemental Skill Heaven and Earth Made Manifest, Zibai will immediately gain 100
// Phase Shift Radiance, and the max number of Spirit Steed's Stride usages per Lunar Phase Shift
// mode is increased to 5 times.
// Additionally, each time you switch to the Lunar Phase Shift mode, the first Spirit Steed's
// Stride's 2nd-hit Lunar-Crystallize Reaction DMG is increased by 220%.
func (c *char) c1Init() {
	if c.Base.Cons < 1 {
		return
	}

	c.AddReactBonusMod(character.ReactBonusMod{
		Base: modifier.NewBase(c1Key+"-buff", -1),
		Amount: func(ai info.AttackInfo) float64 {
			if ai.ActorIndex != c.Index() {
				return 0
			}

			if !c.StatusIsActive(c1Key) {
				return 0
			}

			if ai.Abil != skillAbil2 {
				return 0
			}
			if c.Core.Flags.LogDebug {
				c.Core.Log.NewEvent("Adding C1 react bonus", glog.LogCharacterEvent, c.Index())
			}

			c.QueueCharTask(func() { c.DeleteStatus(c1Key) }, 1)
			return 2.20
		},
	})
}

func (c *char) c1MaxSkillsPerSkill() int {
	if c.Base.Cons < 1 {
		return 4
	}
	return 5
}

func (c *char) c1OnSkill() {
	if c.Base.Cons < 1 {
		return
	}
	c.addRadiance(100)
	c.AddStatus(c1Key, 15*60, true)
}

// When in the Lunar Phase Shift mode, all nearby party members' Lunar-Crystallize Reaction DMG is
// increased by 30%.
// Moonsign: Ascendant Gleam: Ascension Talent The Selenic Adeptus Descends is enhanced: the DMG
// dealt by the 2nd hit of Spirit Steed's Stride is further increased by 550% of Zibai's DEF. You
// must first unlock the Ascension Talent "The Selenic Adeptus Descends."
func (c *char) c2Init() {
	if c.Base.Cons < 2 {
		return
	}

	for _, char := range c.Core.Player.Chars() {
		char.AddReactBonusMod(character.ReactBonusMod{
			Base: modifier.NewBase(c2Key+"-buff", -1),
			Amount: func(ai info.AttackInfo) float64 {
				if !c.StatusIsActive(skillKey) {
					return 0
				}

				switch ai.AttackTag {
				case attacks.AttackTagReactionLunarCrystallize:
				case attacks.AttackTagDirectLunarCrystallize:
				default:
					return 0
				}
				return 0.3
			},
		})
	}
}

func (c *char) c2A4Mult() float64 {
	if c.Base.Cons < 2 {
		return 0
	}

	if c.Core.Player.GetMoonsignLevel() < 2 {
		return 0
	}

	return c2Mult - a4Mult
}

func (c *char) c4ResetNormalCount() {
	if c.Base.Cons < 4 {
		c.Character.ResetNormalCounter()
		return
	}

	if !c.StatusIsActive(skillKey) {
		c.Character.ResetNormalCounter()
		return
	}
}

// While in the Lunar Phase Shift mode, Zibai's Normal Attack sequence will not reset, and when
// Spirit Steed's Stride hits opponents, Zibai will gain the "Scattermoon Splendor" effect: The next
// time she uses Normal Attacks, the additional attack from her 4th hit will deal 250% of the
// original damage as Lunar-Crystallize Reaction DMG.
func (c *char) c4SkillCB(ac info.AttackCB) {
	if c.Base.Cons < 4 {
		return
	}

	if _, ok := ac.Target.(*enemy.Enemy); !ok {
		return
	}
	c.AddStatus(c4Key, 30*60, false)
}

func (c *char) c4N4Bonus() float64 {
	if c.Base.Cons < 4 {
		return 0
	}

	if !c.StatusIsActive(c4Key) {
		return 0
	}
	return 2.5 - 1.0
}

// While Zibai is in the Lunar Phase Shift mode, her Phase Shift Radiance gain rate is increased by
// 50%.
// Additionally, Spirit Steed's Stride will change such that it will consume all Phase Shift
// Radiance. This will elevate the DMG dealt by this instance of Spirit Steed's Stride and the
// Lunar-Crystallize Reaction DMG dealt by Zibai within the next 3s by 1.6% for every point consumed
// above 70. This effect cannot stack.
func (c *char) c6Init() {
	if c.Base.Cons < 6 {
		return
	}

	makeHook := func(aeInd int, tag attacks.AttackTag) func(args ...any) {
		return func(args ...any) {
			atk := args[aeInd].(*info.AttackEvent)
			if atk.Info.AttackTag != tag {
				return
			}

			if !c.StatusIsActive(c6Key) {
				return
			}

			if atk.Info.ActorIndex != c.Index() {
				return
			}

			if c.Core.Flags.LogDebug {
				c.Core.Log.NewEvent("Adding zibai c6 lunar crystallize elevation", glog.LogCharacterEvent, c.Index()).Write("amt", c.c6Elev)
			}
			atk.Info.Elevation += c.c6Elev
		}
	}

	c.Core.Events.Subscribe(event.OnApplyAttack, makeHook(0, attacks.AttackTagDirectLunarCrystallize), c6Key+"-direct")
	c.Core.Events.Subscribe(event.OnSpecialReactionAttack, makeHook(1, attacks.AttackTagReactionLunarCrystallize), c6Key+"-reaction")
}

func (c *char) c6RadianceEff() float64 {
	if c.Base.Cons < 6 {
		return 1.0
	}
	return 1.5
}

func (c *char) c6ConsumeRadiance() {
	if c.Base.Cons < 6 {
		c.consumeRadiance(70)
		return
	}

	c.c6Elev = max((c.radiance-70)*0.016, 0)
	c.AddStatus(c6Key, 3*60, true)
	c.consumeRadiance(c.radiance)
}
