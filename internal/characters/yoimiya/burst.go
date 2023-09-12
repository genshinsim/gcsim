package yoimiya

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/enemy"
)

var burstFrames []int

const burstHitmark = 75
const abDebuff = "aurous-blaze"
const abIcdKey = "aurous-blaze-icd"

func init() {
	burstFrames = frames.InitAbilSlice(113) // Q -> N1
	burstFrames[action.ActionSkill] = 112   // Q -> E
	burstFrames[action.ActionDash] = 111    // Q -> D
	burstFrames[action.ActionJump] = 112    // Q -> J
	burstFrames[action.ActionSwap] = 109    // Q -> Swap
}

func (c *char) Burst(p map[string]int) action.Info {
	// assume it does skill dmg at end of it's animation
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Aurous Blaze",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeBlunt,
		Element:    attributes.Pyro,
		Durability: 50,
		Mult:       burst[c.TalentLvlBurst()],
	}

	if c.Base.Ascension >= 4 {
		c.Core.Tasks.Add(c.a4, burstHitmark)
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 6),
		0,
		burstHitmark,
		c.applyAB, // callback to apply Aurous Blaze
		c.makeC2CB(),
	)

	// add cooldown to sim
	c.SetCD(action.ActionBurst, 15*60)
	// use up energy
	c.ConsumeEnergy(5)

	c.abApplied = false

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}
}

func (c *char) applyAB(a combat.AttackCB) {
	// marker an opponent after first hit
	// ignore the bouncing around for now (just assume it's always target 0)
	// icd of 2s, removed if down

	// do nothing if ab already applied on enemy
	if c.abApplied {
		return
	}

	trg, ok := a.Target.(*enemy.Enemy)
	// do nothing if not an enemy
	if !ok {
		return
	}
	c.abApplied = true

	duration := 600
	if c.Base.Cons >= 1 {
		duration = 840
	}
	trg.AddStatus(abDebuff, duration, true) // apply Aurous Blaze
}

func (c *char) burstHook() {
	// check on attack landed for target 0
	// if aurous active then trigger dmg if not on cd
	c.Core.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
		ae := args[1].(*combat.AttackEvent)
		trg, ok := args[0].(*enemy.Enemy)
		// ignore if not an enemy
		if !ok {
			return false
		}
		// ignore if debuff not on enemy
		if !trg.StatusIsActive(abDebuff) {
			return false
		}
		// ignore for self
		if ae.Info.ActorIndex == c.Index {
			return false
		}
		// ignore if on icd
		if trg.StatusIsActive(abIcdKey) {
			return false
		}
		// ignore if wrong tags
		switch ae.Info.AttackTag {
		case attacks.AttackTagNormal:
		case attacks.AttackTagExtra:
		case attacks.AttackTagPlunge:
		case attacks.AttackTagElementalArt:
		case attacks.AttackTagElementalArtHold:
		case attacks.AttackTagElementalBurst:
		default:
			return false
		}
		// do explosion, set icd
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Aurous Blaze (Explode)",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagElementalBurst,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeBlunt,
			Element:    attributes.Pyro,
			Durability: 25,
			Mult:       burstExplode[c.TalentLvlBurst()],
		}
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(trg, nil, 3), 0, 1, c.makeC2CB())

		trg.AddStatus(abIcdKey, 120, true) // trigger Aurous Blaze ICD

		// C4
		if c.Base.Cons >= 4 {
			c.ReduceActionCooldown(action.ActionSkill, 72)
		}

		return false
	}, "yoimiya-burst-check")

	if c.Core.Flags.DamageMode {
		// add check for if yoimiya dies
		c.Core.Events.Subscribe(event.OnPlayerHPDrain, func(_ ...interface{}) bool {
			if c.CurrentHPRatio() <= 0 {
				// remove Aurous Blaze from target
				for _, x := range c.Core.Combat.Enemies() {
					trg := x.(*enemy.Enemy)
					if trg.StatusIsActive(abDebuff) {
						trg.DeleteStatus(abDebuff)
					}
				}
			}
			return false
		}, "yoimiya-died")
	}
}
