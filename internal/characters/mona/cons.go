package mona

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const c6Key = "mona-c6"

// C1:
// When any of your own party members hits an opponent affected by an Omen, the effects of Hydro-related Elemental Reactions are enhanced for 8s:
// - Electro-Charged DMG increases by 15%.
// - Vaporize DMG increases by 15%.
// - Hydro Swirl DMG increases by 15%.
// - Frozen duration is extended by 15%.
func (c *char) c1() {
	// TODO: "Frozen duration is extended by 15%." is bugged
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		//ignore if target doesn't have debuff
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		if !t.StatusIsActive(bubbleKey) && !t.StatusIsActive(omenKey) {
			return false
		}
		// add c1 to all party members, delay by 1, because:
		// "This bonus does not apply in the triggering attack nor from the resulting Hydro DMG dealt by Illusory Bubble in Stellaris Phantasm regardless if they were from resulting reactions."
		for _, x := range c.Core.Player.Chars() {
			char := x
			c.Core.Tasks.Add(func() {
				// TODO: "Vaporize DMG increases by 15%." should be getting snapshot, see https://library.keqingmains.com/evidence/characters/hydro/mona#mona-c1-snapshot-for-vape
				// requires ReactBonusMod refactor
				char.AddReactBonusMod(character.ReactBonusMod{
					Base: modifier.NewBase("mona-c1", 8*60),
					Amount: func(ai combat.AttackInfo) (float64, bool) {
						// doesn't work off-field
						if c.Core.Player.Active() != char.Index {
							return 0, false
						}
						// Electro-Charged DMG increases by 15%.
						if ai.AttackTag == attacks.AttackTagECDamage {
							return 0.15, false
						}
						// Vaporize DMG increases by 15%.
						// the only way Hydro Swirl can vape is via an AoE Hydro Swirl which doesn't do damage anyways, so this is fine
						if ai.Amped {
							return 0.15, false
						}
						// Hydro Swirl DMG increases by 15%.
						if ai.AttackTag == attacks.AttackTagSwirlHydro {
							return 0.15, false
						}
						return 0, false
					},
				})
			}, 1)
		}
		return false
	}, "mona-c1-check")
}

// C2:
// When a Normal Attack hits, there is a 20% chance that it will be automatically followed by a Charged Attack.
// This effect can only occur once every 5s.
func (c *char) c2(a combat.AttackCB) {
	trg := a.Target
	if c.Base.Cons < 2 {
		return
	}
	if a.Target.Type() != combat.TargettableEnemy {
		return
	}
	if c.Core.Rand.Float64() > .2 {
		return
	}
	if c.c2icd > c.Core.F {
		return
	}
	c.c2icd = c.Core.F + 300 //every 5 seconds
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Charge Attack",
		AttackTag:  attacks.AttackTagExtra,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       charge[c.TalentLvlAttack()],
	}

	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(trg, nil, 3), 0, 0)
}

// C4:
// When any party member attacks an opponent affected by an Omen, their CRIT Rate is increased by 15%.
func (c *char) c4() {
	m := make([]float64, attributes.EndStatType)
	m[attributes.CR] = 0.15

	for _, char := range c.Core.Player.Chars() {
		char.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase("mona-c4", -1),
			Amount: func(_ *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				x, ok := t.(*enemy.Enemy)
				if !ok {
					return nil, false
				}
				//ok only if either bubble or omen is present
				if x.StatusIsActive(bubbleKey) || x.StatusIsActive(omenKey) {
					return m, true
				}
				return nil, false
			},
		})
	}
}

// C6:
// Upon entering Illusory Torrent, Mona gains a 60% increase to the DMG of her next Charged Attack per second of movement.
// A maximum DMG Bonus of 180% can be achieved in this manner.
// The effect lasts for no more than 8s.
func (c *char) c6(src int) func() {
	return func() {
		if c.c6Src != src {
			c.Core.Log.NewEvent(fmt.Sprintf("%v stack gain check ignored, src diff", c6Key), glog.LogCharacterEvent, c.Index).
				Write("src", src).
				Write("new src", c.c6Src)
			return
		}
		// do nothing if not Mona
		if c.Core.Player.Active() != c.Index {
			return
		}
		// do nothing if we aren't dashing anymore
		if c.Core.Player.CurrentState() != action.DashState {
			return
		}

		c.c6Stacks++
		if c.c6Stacks > 3 {
			c.c6Stacks = 3
		}
		c.Core.Log.NewEvent(fmt.Sprintf("%v stack gained", c6Key), glog.LogCharacterEvent, c.Index).
			Write("c6Stacks", c.c6Stacks)

		m := make([]float64, attributes.EndStatType)
		c.AddAttackMod(character.AttackMod{
			Base: modifier.NewBase(c6Key, 8*60),
			Amount: func(atk *combat.AttackEvent, t combat.Target) ([]float64, bool) {
				if atk.Info.AttackTag != attacks.AttackTagExtra {
					return nil, false
				}
				m[attributes.DmgP] = 0.60 * float64(c.c6Stacks)
				return m, true
			},
		})

		// reset C6 stacks in 8s if we didn't use a CA
		c.Core.Tasks.Add(c.c6TimerReset, 8*60+1)
		// queue up another stack and buff refresh in 1s
		c.Core.Tasks.Add(c.c6(src), 60)
	}
}

func (c *char) makeC6CAResetCB() combat.AttackCBFunc {
	if c.Base.Cons < 6 || !c.StatusIsActive(c6Key) {
		return nil
	}
	return func(a combat.AttackCB) {
		if a.Target.Type() == combat.TargettableEnemy {
			return
		}
		if !c.StatusIsActive(c6Key) {
			return
		}
		c.DeleteStatus(c6Key)
		c.c6Stacks = 0
		c.Core.Log.NewEvent(fmt.Sprintf("%v stacks reset via charge attack", c6Key), glog.LogCharacterEvent, c.Index)
	}
}

func (c *char) c6TimerReset() {
	// handle C6 stack reset if CA not used before c6 buff expires
	if c.c6Stacks > 0 && !c.StatusIsActive(c6Key) {
		c.c6Stacks = 0
		c.Core.Log.NewEvent(fmt.Sprintf("%v stacks reset via timer", c6Key), glog.LogCharacterEvent, c.Index)
	}
}
