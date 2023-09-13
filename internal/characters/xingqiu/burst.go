package xingqiu

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var burstFrames []int

const (
	burstHitmark = 18
	burstKey     = "xingqiuburst"
	burstICDKey  = "xingqiu-burst-icd"
)

func init() {
	burstFrames = frames.InitAbilSlice(40)
	burstFrames[action.ActionAttack] = 33
	burstFrames[action.ActionSkill] = 33
	burstFrames[action.ActionDash] = 33
	burstFrames[action.ActionJump] = 33
}

/**
The number of Hydro Swords summoned per wave follows a specific pattern, usually alternating between 2 and 3 swords.
At C6, this is upgraded and follows a pattern of 2 → 3 → 5… which then repeats.

There is an approximately 1 second interval between summoned Hydro Sword waves, so that means a theoretical maximum of 15 or 18 waves.

Each wave of Hydro Swords is capable of applying one (1) source of Hydro status, and each individual sword is capable of getting a crit.
**/

func (c *char) Burst(p map[string]int) action.Info {
	// apply hydro every 3rd hit
	// triggered on normal attack
	// also applies hydro on cast if p=1
	// how we doing that?? trigger 0 dmg?

	/** c2
	Extends the duration of Guhua Sword: Raincutter by 3s.
	Decreases the Hydro RES of opponents hit by sword rain attacks by 15% for 4s.
	**/
	dur := 15
	if c.Base.Cons >= 2 {
		dur += 3
	}
	dur *= 60
	c.AddStatus(burstKey, dur+33, true) // add 33f for anim
	c.applyOrbital(dur, burstHitmark)

	c.burstCounter = 0
	c.numSwords = 2
	c.nextRegen = false

	// c.CD[combat.BurstCD] = c.S.F + 20*60
	c.SetCD(action.ActionBurst, 20*60)
	c.ConsumeEnergy(3)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstHitmark,
		State:           action.BurstState,
	}
}

func (c *char) summonSwordWave() {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Guhua Sword: Raincutter",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypePierce,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	// only if c.nextRegen is true and first sword
	var c2cb, c6cb func(a combat.AttackCB)
	if c.nextRegen {
		done := false
		c6cb = func(a combat.AttackCB) {
			if a.Target.Type() != targets.TargettableEnemy {
				return
			}
			if done {
				return
			}
			c.AddEnergy("xingqiu-c6", 3)
			done = true
		}
	}
	if c.Base.Cons >= 2 {
		icd := -1
		c2cb = func(a combat.AttackCB) {
			if c.Core.F < icd {
				return
			}

			e, ok := a.Target.(*enemy.Enemy)
			if !ok {
				return
			}

			icd = c.Core.F + 1
			c.Core.Tasks.Add(func() {
				e.AddResistMod(combat.ResistMod{
					Base:  modifier.NewBaseWithHitlag("xingqiu-c2", 4*60),
					Ele:   attributes.Hydro,
					Value: -0.15,
				})
			}, 1)
		}
	}

	for i := 0; i < c.numSwords; i++ {
		//TODO: this snapshot timing is off? perhaps should be 0?
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHit(
				c.Core.Combat.Player(),
				c.Core.Combat.PrimaryTarget(),
				nil,
				0.5,
			),
			20,
			20,
			c2cb,
			c6cb,
		)
		c6cb = nil
		c.burstCounter++
	}

	// figure out next wave # of swords
	switch c.numSwords {
	case 2:
		c.numSwords = 3
		c.nextRegen = false
	case 3:
		if c.Base.Cons >= 6 {
			c.numSwords = 5
			c.nextRegen = true
		} else {
			c.numSwords = 2
			c.nextRegen = false
		}
	case 5:
		c.numSwords = 2
		c.nextRegen = false
	}
}
