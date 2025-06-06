package ororon

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var (
	burstFrames []int
	burstStatus = "supersonic-oculus"
)

const (
	burstHitmark     = 36 // Initial Hit
	burstTickHitmark = 60
)

func init() {
	burstFrames = frames.InitAbilSlice(62) // Q -> N1/CA/E/Walk
	burstFrames[action.ActionDash] = 59
	burstFrames[action.ActionJump] = 60
	burstFrames[action.ActionSwap] = 61
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	// first zap has no icd and hits everyone
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Ritual DMG",
		AttackTag:      attacks.AttackTagElementalBurst,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Electro,
		Durability:     25,
		Mult:           burst[c.TalentLvlBurst()],
		HitlagFactor:   0.05,
	}
	c.burstArea = combat.NewCircleHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: 2.3}, 6.5)
	c.Core.QueueAttack(
		ai,
		c.burstArea,
		burstHitmark,
		burstHitmark,
		c.makeC2cb(),
	)
	c.QueueCharTask(c.c4EnergyRestore, burstHitmark)

	c.burstSrc = c.Core.F
	c.AddStatus(burstStatus, 9*60, false)
	c.burstTick(c.burstSrc)

	c.c2OnBurst()
	c.c6OnBurst()

	c.SetCD(action.ActionBurst, 15*60)
	c.ConsumeEnergy(22)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) burstTick(src int) {
	tick := int(float64(burstTickHitmark) * c.c4BurstInterval())
	c.Core.Tasks.Add(func() {
		if c.burstSrc != src {
			return
		}
		if !c.StatusIsActive(burstStatus) {
			return
		}

		ai := combat.AttackInfo{
			ActorIndex:     c.Index,
			Abil:           "Soundwave Collision DMG",
			AttackTag:      attacks.AttackTagElementalBurst,
			AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
			ICDTag:         attacks.ICDTagOroronElmentalBurst,
			ICDGroup:       attacks.ICDGroupOroronElementalBurst,
			StrikeType:     attacks.StrikeTypeDefault,
			Element:        attributes.Electro,
			Durability:     25,
			Mult:           soundwave[c.TalentLvlBurst()],
			HitlagFactor:   0.02,
		}
		// TODO: make 3 boxes (bullets?) hits instead of 1 circle hit?
		c.Core.QueueAttack(
			ai,
			c.burstArea,
			0,
			0,
			c.makeC2cb(),
		)

		c.burstTick(src)
	}, tick)
}
