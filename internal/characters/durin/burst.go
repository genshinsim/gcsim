package durin

import (
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var (
	burstFrames      []int
	burstInitHitmark = []int{97, 97 + 24, 154} // Initial Hit
)

const (
	burstTicksWhite          = 20
	burstTicksBlack          = 16
	burstIntervalWhite       = 58.8
	burstIntervalBlack       = 73.6
	burstFirstTickDelayWhite = 154 + 60 - burstIntervalWhite
	burstFirstTickDelayBlack = 154 + 95 - burstIntervalBlack
	burstCD                  = 18 * 60
	burstDuration            = 20.5 * 60 // starts on burstInitialHitmark[0]

	burstKeyWhite = "durin-burst-white"
	burstKeyBlack = "durin-burst-black"
)

func init() {
	burstFrames = frames.InitAbilSlice(104)
	burstFrames[action.ActionAttack] = 103
	burstFrames[action.ActionSkill] = 103
	burstFrames[action.ActionDash] = 102
	burstFrames[action.ActionJump] = 104
	burstFrames[action.ActionWalk] = 102
	burstFrames[action.ActionSwap] = 102
}

func ceil(x float64) int {
	return int(math.Ceil(x))
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(blackKey) {
		travel, ok := p["travel"]
		if !ok {
			travel = 9
		}
		return c.burstBlack(travel)
	}

	return c.burstWhite()
}

func (c *char) burstWhite() (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex:       c.Index(),
		Abil:             "As the Light Shifts",
		AttackTag:        attacks.AttackTagElementalBurst,
		ICDTag:           attacks.ICDTagDurinBurstTornado,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Pyro,
		Durability:       25,
		IgnoreDefPercent: c.c6DefIgnore(true),
	}

	// TODO: if target is out of range then pos should be player pos + Y: 10 offset
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 6.5)
	for i, mult := range burstWhiteInitial {
		ai.Mult = mult[c.TalentLvlBurst()]
		c.Core.QueueAttack(
			ai,
			ap,
			burstInitHitmark[i],
			burstInitHitmark[i],
			c.c6WhiteMakeCB(),
		)
	}

	c.burstSrc = c.Core.F
	c.DeleteStatus(burstKeyBlack)
	for i := 0.0; i < burstTicksWhite; i++ {
		c.Core.Tasks.Add(c.burstTickWhite(c.burstSrc), ceil(burstFirstTickDelayWhite+burstIntervalWhite*i))
	}

	c.QueueCharTask(func() {
		c.a1OnBurst(true)
		c.a4OnBurst()
		c.c1OnBurst(true)
		c.AddStatus(burstKeyWhite, burstDuration, false)
	}, burstInitHitmark[0])

	c.SetCD(action.ActionBurst, burstCD)
	c.ConsumeEnergy(10)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) burstTickWhite(src int) func() {
	return func() {
		if src != c.burstSrc {
			return
		}

		ai := info.AttackInfo{
			ActorIndex:       c.Index(),
			Abil:             "Dragon of White Flame",
			AttackTag:        attacks.AttackTagElementalBurst,
			ICDTag:           attacks.ICDTagDurinBurstWhite,
			ICDGroup:         attacks.ICDGroupDurinBurstWhite,
			StrikeType:       attacks.StrikeTypeDefault,
			Element:          attributes.Pyro,
			Durability:       25,
			Mult:             burstWhiteDoT[c.TalentLvlBurst()] * c.a4Dmg(),
			IgnoreDefPercent: c.c6DefIgnore(true),
		}

		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3.5),
			0,
			0,
			c.c6WhiteMakeCB(),
		)
	}
}

func (c *char) burstBlack(travel int) (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex:       c.Index(),
		Abil:             "As the Stars Smolder",
		AttackTag:        attacks.AttackTagElementalBurst,
		ICDTag:           attacks.ICDTagDurinBurstTornado,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Pyro,
		Durability:       25,
		IgnoreDefPercent: c.c6DefIgnore(false),
	}

	// TODO: if target is out of range then pos should be player pos + Y: 10 offset
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 6.5)
	for i, mult := range burstBlackInitial {
		ai.Mult = mult[c.TalentLvlBurst()]
		c.Core.QueueAttack(ai, ap, burstInitHitmark[i], burstInitHitmark[i])
	}

	c.burstSrc = c.Core.F
	c.DeleteStatus(burstKeyWhite)
	for i := 0.0; i < burstTicksBlack; i++ {
		c.Core.Tasks.Add(c.burstTickBlack(c.burstSrc, travel), ceil(burstFirstTickDelayBlack+burstIntervalBlack*i))
	}

	c.QueueCharTask(func() {
		c.a1OnBurst(false)
		c.a4OnBurst()
		c.c1OnBurst(false)
		c.AddStatus(burstKeyBlack, burstDuration, false)
	}, burstInitHitmark[0])

	c.SetCD(action.ActionBurst, burstCD)
	c.ConsumeEnergy(10)
	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) burstTickBlack(src, travel int) func() {
	return func() {
		if src != c.burstSrc {
			return
		}

		ai := info.AttackInfo{
			ActorIndex:       c.Index(),
			Abil:             "Dragon of Dark Decay",
			AttackTag:        attacks.AttackTagElementalBurst,
			ICDTag:           attacks.ICDTagDurinBurstBlack,
			ICDGroup:         attacks.ICDGroupDurinBurstBlack,
			StrikeType:       attacks.StrikeTypeDefault,
			Element:          attributes.Pyro,
			Durability:       25,
			Mult:             burstBlackDoT[c.TalentLvlBurst()] * c.a4Dmg(),
			IgnoreDefPercent: c.c6DefIgnore(false),
		}

		c.Core.QueueAttack(ai, combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key()), 0, travel)
	}
}
