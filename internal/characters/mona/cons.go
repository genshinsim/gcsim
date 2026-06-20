package mona

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	c2icdkey                = "mona-c2-icd"
	c2HexereiPostBurstCAKey = "mona-c2-hexerei-post-burst-ca"
	c4key                   = "mona-c4"
	c6Key                   = "mona-c6"
)

// C1:
// When any of your own party members hits an opponent affected by an Omen, the effects of Hydro-related Elemental Reactions are enhanced for 8s:
// - Electro-Charged DMG increases by 15%
// - Lunar-Charged DMG increases by 15%
// - Vaporize DMG increases by 15%
// - Hydro Swirl DMG increases by 15%
// - Lunar-Crystallize DMG increases by 15%.
// - Frozen duration is extended by 15%.
// After Hex: Additionally, when your off-field party members trigger the above effect,
// the DMG Bonus to the above Hydro-related Elemental Reactions is enhanced to 160% of its original effect.
func (c *char) c1() {
	// TODO: "Frozen duration is extended by 15%." is bugged
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		// ignore if target doesn't have debuff
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return
		}
		if !t.StatusIsActive(omenKey) {
			return
		}

		atk := args[1].(*info.AttackEvent)

		char := c.Core.Player.Chars()[atk.Info.ActorIndex]

		// add c1 to party member that triggered the effect, delay by 1, because:
		// "This bonus does not apply in the triggering attack nor from the resulting Hydro DMG dealt by
		//  Illusory Bubble in Stellaris Phantasm regardless if they were from resulting reactions."
		c.Core.Tasks.Add(func() {
			char.AddReactBonusMod(character.ReactBonusMod{
				Base: modifier.NewBaseWithHitlag("mona-c1", 8*60),
				Amount: func(ai info.AttackInfo) float64 {
					m := 0.15

					// Hexerei passive
					// Additionally, when your off-field party members trigger the above effect, the DMG Bonus
					// to the above Hydro-related Elemental Reactions is enhanced to 160% of its original effect.
					if c.IsHexerei && c.Core.Player.Active() != char.Index() {
						m = 0.15 * 1.6
					}

					switch ai.AttackTag {
					// - Hydro Swirl DMG increases by 15%.
					// - Electro-Charged DMG increases by 15%.
					// - Lunar-Charged DMG increases by 15%.
					// - Lunar-Crystallize DMG increases by 15%.
					case attacks.AttackTagSwirlHydro,
						attacks.AttackTagECDamage,
						attacks.AttackTagReactionLunarCharge, attacks.AttackTagDirectLunarCharged,
						attacks.AttackTagReactionLunarCrystallize, attacks.AttackTagDirectLunarCrystallize:
						return m
					}

					// Vaporize DMG increases by 15%.
					// the only way Hydro Swirl can vape is via an AoE Hydro Swirl which doesn't do damage anyways, so this is fine

					if ai.Amped && ai.AmpType == info.ReactionTypeVaporize {
						return m
					}

					return 0
				},
			})
		}, 1)
	}, "mona-c1-check")
}

// C2:
// When a Normal Attack hits, there is a 20% chance that it will be automatically followed by a Charged Attack.
// This effect can only occur once every 5s.
// After Hex:
// Additionally, within 5s after Mona unleashes her Elemental Burst Stellaris Phantasm, her next Normal Attack that hits an enemy will automatically trigger a Charged Attack.
// This Charged Attack effect can only occur once every 5s.
// When Mona's Charged Attack hits an opponent, all nearby party members will have their Elemental Mastery increased by 80 for 12s.
func (c *char) c2() {
	if c.Base.Cons < 2 {
		return
	}
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		trg, ok := args[0].(*enemy.Enemy)
		if !ok {
			return
		}

		atk := args[1].(*info.AttackEvent)
		if atk.Info.ActorIndex != c.Index() {
			return
		}
		if atk.Info.AttackTag != attacks.AttackTagNormal {
			return
		}
		if !c.StatusIsActive(c2HexereiPostBurstCAKey) && c.Core.Rand.Float64() > .2 {
			return
		}
		if c.StatusIsActive(c2icdkey) {
			return
		}
		c.DeleteStatus(c2HexereiPostBurstCAKey)
		c.AddStatus(c2icdkey, 5*60, true)

		c.doCA(trg, 53)
	}, "mona-c2-followup")
}

func (c *char) c2HexereiCB(a info.AttackCB) {
	if !c.IsHexerei {
		return
	}

	if a.Target.Type() != info.TargettableEnemy {
		return
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.EM] = 80

	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBaseWithHitlag("mona-hexerei-c2-em", 60*8),
			AffectedStat: attributes.EM,
			Amount: func() []float64 {
				return m
			},
		})
	}
}

// C4:
// When any party member attacks an opponent affected by an Omen, their CRIT Rate is increased by 15%.
// After Hex: When any Hexerei party member attacks an opponent affected by an Omen, their CRIT DMG is increased by 15%.
func (c *char) c4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = 0.15

	for _, char := range c.Core.Player.Chars() {
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase(c4key, -1),
			Amount: func(_ *info.AttackEvent, t info.Target) []float64 {
				x, ok := t.(*enemy.Enemy)
				if !ok {
					return nil
				}

				if c.IsHexerei && char.IsHexerei {
					m[attributes.CD] = 0.15
				}
				// ok only if either bubble or omen is present
				if x.StatusIsActive(bubbleKey) || x.StatusIsActive(omenKey) {
					return m
				}
				return nil
			},
		})
	}

	// workaround for giving lunarcharge and lunarcrystallize the 15% CR and 15% CDMG
	c.Core.Events.Subscribe(event.OnLunarReactionAttack, func(args ...any) {
		x, ok := args[0].(*enemy.Enemy)
		if !ok {
			return
		}

		ae, ok := args[1].(*info.AttackEvent)
		if !ok {
			return
		}

		if !x.StatusIsActive(bubbleKey) && !x.StatusIsActive(omenKey) {
			return
		}

		if c.Core.Flags.LogDebug {
			var reaction string
			switch ae.Info.AttackTag {
			case attacks.AttackTagReactionLunarCharge:
				reaction = "Lunarcharged"
			case attacks.AttackTagReactionLunarCrystallize:
				reaction = "Lunarcrystallize"
			}

			c.Core.Log.NewEvent("Mona C4 CR added to "+reaction, glog.LogPreDamageMod, ae.Info.ActorIndex).
				Write("before", ae.Snapshot.Stats[attributes.CR]).
				Write("addition", 0.15)
		}

		ae.Snapshot.Stats[attributes.CR] += 0.15

		char := c.Core.Player.Chars()[ae.Info.ActorIndex]
		if c.IsHexerei && char.IsHexerei {
			ae.Snapshot.Stats[attributes.CD] += 0.15
		}
	}, c4key+"-lunar-reaction")
}

func (c *char) c6OnDash() {
	if c.Base.Cons < 6 {
		return
	}
	isHexAndOmen := c.IsHexerei && c.omenIsNearby()

	if !isHexAndOmen {
		// start a new task if an omen task isn't already ticking
		c.c6Src = c.Core.F
		c.QueueCharTask(c.c6Ticker(c.c6Src), 60)
	}
}

func (c *char) c6OnDashEnd(_ action.AnimationState) {
	if c.Base.Cons < 6 {
		return
	}
	isHexAndOmen := c.IsHexerei && c.omenIsNearby()
	if !isHexAndOmen {
		// Cancel existing dash ticker when dash ends if Mona isn't Hex, or the Nearby Omen condition isn't met
		c.c6Src = -1
	}
}

// this must be called before the actual bubble/omen status is applied
func (c *char) c6OnApplyBubble() {
	if c.Base.Cons < 6 {
		return
	}

	// start a Omen task only if a previous Omen task isn't still active
	if c.IsHexerei && !c.omenIsNearby() {
		c.c6Src = c.Core.F
		c.QueueCharTask(c.c6Ticker(c.Core.F), 60)
	}
}

// C6:
// Upon entering Illusory Torrent, Mona gains a 60% increase to the DMG of her next Charged Attack per second of movement.
// A maximum DMG Bonus of 180% can be achieved in this manner.
// The effect lasts for no more than 8s.
func (c *char) c6Ticker(src int) func() {
	return func() {
		// do nothing if not Mona
		if c.Core.Player.Active() != c.Index() {
			return
		}

		// cancel the task if the src was changed
		if c.c6Src != src {
			return
		}

		// do nothing if we aren't dashing anymore and we aren't hexerei and enemy with omen is nearby
		isDashing := c.Core.Player.CurrentState() == action.DashState
		isHexAndOmen := c.IsHexerei && c.omenIsNearby()

		if !isDashing && !isHexAndOmen {
			return
		}

		c.c6Stacks = min(c.c6Stacks+1, 3)

		c.Core.Log.NewEvent(c6Key+" stack gained", glog.LogCharacterEvent, c.Index()).
			Write("c6Stacks", c.c6Stacks)

		// reset C6 stacks in 8s if we didn't use a CA
		c.QueueCharTask(c.c6TimerReset(c.c6Src), 8*60+1)
		// queue up another stack and buff refresh in 1s
		c.QueueCharTask(c.c6Ticker(src), 60)
	}
}

func (c *char) c6Init() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		if !c.IsHexerei {
			return
		}

		if c.Core.Player.Active() != c.Index() {
			return
		}

		t, ok := args[0].(enemy.Enemy)

		if ok {
			return
		}

		if t.StatusIsActive(omenKey) || t.StatusIsActive(bubbleKey) {
			c.Core.Tasks.Add(c.c6Ticker(c.Core.F), 60)
		}
	}, "mona-c6-init")

	m := make([]float64, attributes.EndStatType)
	c.AddAttackMod(character.AttackMod{
		Base: modifier.NewBaseWithHitlag(c6Key, 8*60),
		Amount: func(atk *info.AttackEvent, t info.Target) []float64 {
			if atk.Info.AttackTag != attacks.AttackTagExtra {
				return nil
			}
			if c.c6Stacks == 0 {
				return nil
			}

			m[attributes.DmgP] = 0.60 * float64(c.c6Stacks)

			return m
		},
	})
}

func (c *char) c6ChargeAttackInit() {
	if !c.IsHexerei {
		return
	}

	c6HexCABuff := func(args ...any) {
		atk := args[1].(*info.AttackEvent)
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return
		}
		if atk.Info.ActorIndex != c.Index() {
			return
		}
		if atk.Info.AttackTag != attacks.AttackTagExtra {
			return
		}
		if !t.StatusIsActive(omenKey) && !t.StatusIsActive(bubbleKey) {
			return
		}
		atk.Info.Mult *= 2
	}

	c.Core.Events.Subscribe(event.OnEnemyHit, c6HexCABuff, "mona-hexerei-c6-ca-buff-%v")
}

func (c *char) makeC6CAResetCB() info.AttackCBFunc {
	if c.Base.Cons < 6 || !c.StatusIsActive(c6Key) {
		return nil
	}
	return func(a info.AttackCB) {
		if a.Target.Type() == info.TargettableEnemy {
			return
		}
		if !c.StatusIsActive(c6Key) {
			return
		}
		c.DeleteStatus(c6Key)
		c.c6Stacks = 0
		c.Core.Log.NewEvent(fmt.Sprintf("%v stacks reset via charge attack", c6Key), glog.LogCharacterEvent, c.Index())
	}
}

func (c *char) c6TimerReset(src int) func() {
	return func() {
		if c.c6Src != src {
			return
		}
		// handle C6 stack reset if CA not used before c6 buff expires
		if c.c6Stacks > 0 && !c.StatusIsActive(c6Key) {
			c.c6Stacks = 0
			c.Core.Log.NewEvent(fmt.Sprintf("%v stacks reset via timer", c6Key), glog.LogCharacterEvent, c.Index())
		}
	}
}
