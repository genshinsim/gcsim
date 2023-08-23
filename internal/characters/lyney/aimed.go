package lyney

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var (
	aimedFrames     []int
	aimedPropFrames []int
)

const (
	// TODO: proper frames, currently using tighnari
	aimedRelease     = 86
	aimedPropRelease = 175

	skillAlignedICDKey = "lyney-aligned-icd"
	skillAlignedICD    = 6 * 60

	grinMalkinHatKey      = "lyney-grinmalkinhat"
	grinMalkinHatDuration = 4 * 60

	propSurplusHPDrainThreshold = 0.6
	propSurplusHPDrainRatio     = 0.2
)

func init() {
	// TODO: proper frames, currently using tighnari
	aimedFrames = frames.InitAbilSlice(94)
	aimedFrames[action.ActionDash] = aimedRelease
	aimedFrames[action.ActionJump] = aimedRelease

	aimedPropFrames = frames.InitAbilSlice(183)
	aimedPropFrames[action.ActionDash] = aimedPropRelease
	aimedPropFrames[action.ActionJump] = aimedPropRelease
}

func (c *char) Aimed(p map[string]int) action.ActionInfo {
	level, ok := p["level"]
	if !ok {
		level = 1
	}
	if level == 1 {
		return c.PropAimed(p)
	}

	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	weakspot := p["weakspot"]

	ai := combat.AttackInfo{
		ActorIndex:           c.Index,
		Abil:                 "Aim (Charged)",
		AttackTag:            attacks.AttackTagExtra,
		ICDTag:               attacks.ICDTagNone,
		ICDGroup:             attacks.ICDGroupDefault,
		StrikeType:           attacks.StrikeTypePierce,
		Element:              attributes.Pyro,
		Durability:           25,
		Mult:                 fullaim[c.TalentLvlAttack()],
		HitWeakPoint:         weakspot == 1,
		HitlagHaltFrames:     0.12 * 60,
		HitlagFactor:         0.01,
		HitlagOnHeadshotOnly: true,
		IsDeployable:         true,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewBoxHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			geometry.Point{Y: -0.5},
			0.1,
			1,
		),
		aimedRelease,
		aimedRelease+travel,
		c.makeC4CB(),
	)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(aimedFrames),
		AnimationLength: aimedFrames[action.InvalidAction],
		CanQueueAfter:   aimedRelease,
		State:           action.AimState,
	}
}

func (c *char) PropAimed(p map[string]int) action.ActionInfo {
	travel, ok := p["travel"]
	if !ok {
		travel = 10
	}
	c6Travel, ok := p["c6_travel"]
	if !ok {
		c6Travel = 10
	}
	weakspot := p["weakspot"]

	propAI := combat.AttackInfo{
		ActorIndex:           c.Index,
		Abil:                 "Prop Arrow",
		AttackTag:            attacks.AttackTagExtra,
		ICDTag:               attacks.ICDTagNone,
		ICDGroup:             attacks.ICDGroupDefault,
		StrikeType:           attacks.StrikeTypePierce,
		Element:              attributes.Pyro,
		Durability:           25,
		Mult:                 prop[c.TalentLvlAttack()],
		HitWeakPoint:         weakspot == 1,
		HitlagHaltFrames:     0.12 * 60,
		HitlagFactor:         0.01,
		HitlagOnHeadshotOnly: true,
		IsDeployable:         true,
	}
	c.QueueCharTask(func() {
		hpDrained := c.propSurplus()
		c.c6(c6Travel)
		target := c.Core.Combat.PrimaryTarget()
		c.Core.QueueAttack(
			propAI,
			combat.NewBoxHit(
				c.Core.Combat.Player(),
				target,
				geometry.Point{Y: -0.5},
				0.1,
				1,
			),
			0,
			travel,
			c.makeGrinMalkinHatCB(hpDrained),
			c.makeC4CB(),
		)
		// TODO: proper frames
		c.QueueCharTask(c.skillAligned(target.Pos()), travel)
	}, aimedPropRelease)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(aimedPropFrames),
		AnimationLength: aimedPropFrames[action.InvalidAction],
		CanQueueAfter:   aimedPropRelease,
		State:           action.AimState,
	}
}

// not implemented: The effect will be removed after the character spends 30s out of combat.
func (c *char) propSurplus() bool {
	// When firing the Prop Arrow, and when Lyney has more than 60% HP,
	// he will consume a portion of his HP to obtain 1 Prop Surplus stack.
	if c.CurrentHPRatio() <= propSurplusHPDrainThreshold {
		return false
	}

	currentHP := c.CurrentHP()
	maxHP := c.MaxHP()
	hpdrain := propSurplusHPDrainRatio * maxHP
	// The lowest Lyney can drop to through this method is 60% of his Max HP.
	if (currentHP-hpdrain)/maxHP <= propSurplusHPDrainThreshold {
		hpdrain = currentHP - propSurplusHPDrainThreshold*maxHP
	}
	c.Core.Player.Drain(player.DrainInfo{
		ActorIndex: c.Index,
		Abil:       "Prop Surplus",
		Amount:     hpdrain,
	})

	c.increasePropSurplusStacks(c.c1StackIncrease())
	return true
}

func (c *char) increasePropSurplusStacks(increase int) {
	c.propSurplusStacks += increase
	if c.propSurplusStacks > 5 {
		c.propSurplusStacks = 5
	}
	c.Core.Log.NewEvent("Lyney Prop Surplus stack added", glog.LogCharacterEvent, c.Index).Write("prop_surplus_stacks", c.propSurplusStacks)
}

func (c *char) skillAligned(pos geometry.Point) func() {
	return func() {
		if c.StatusIsActive(skillAlignedICDKey) {
			return
		}
		c.AddStatus(skillAlignedICDKey, skillAlignedICD, true)

		propAlignedAI := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Spiritbreath Thorn (" + c.Base.Key.Pretty() + ")",
			AttackTag:  attacks.AttackTagExtra,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypePierce,
			Element:    attributes.Pyro,
			Durability: 0,
			Mult:       propAligned[c.TalentLvlAttack()],
		}
		c.Core.QueueAttack(
			propAlignedAI,
			combat.NewCircleHitOnTarget(pos, nil, 2),
			// TODO: proper frames
			0.7*60,
			0.7*60,
			c.makeC4CB(),
		)
	}
}

func (c *char) makeGrinMalkinHatCB(hpDrained bool) combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true

		hatIncrease := c.c1HatIncrease()
		for i := 0; i < hatIncrease; i++ {
			// kill existing hat if reached limit
			if len(c.hats) == c.maxHatCount {
				c.hats[0].Kill()
			}
			g := c.newGrinMalkinHat(a.Target.Pos(), hpDrained)
			c.hats = append(c.hats, g)
			c.Core.Combat.AddGadget(g)
		}
	}
}
