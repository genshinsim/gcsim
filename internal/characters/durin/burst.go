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
	burstInitHitmark = []int{99, 99 + 25, 99 + 25 + 26} // Initial Hit
)

const (
	burstTicksWhite          = 20
	burstTicksBlack          = 16
	burstIntervalWhite       = 59
	burstIntervalBlack       = 73
	burstFirstTickDelayBlack = 99 + 25 + 26 + 53
	burstFirstTickDelayWhite = 99 + 25 + 26 + 22
	burstCD                  = 18 * 60
	burstKeyWhite            = "durin-burst-white"
	burstKeyBlack            = "durin-burst-black"
)

func init() {
	burstFrames = frames.InitAbilSlice(110) // E -> D/J
	burstFrames[action.ActionAttack] = 110
	burstFrames[action.ActionBurst] = 110
	burstFrames[action.ActionWalk] = 110
	burstFrames[action.ActionSwap] = 110
}

func ceil(x float64) int {
	return int(math.Ceil(x))
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(blackKey) {
		return c.burstBlack()
	}

	return c.burstWhite()
}

func (c *char) burstWhite() (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex:       c.Index(),
		Abil:             "As the Light Shifts",
		AttackTag:        attacks.AttackTagElementalBurst,
		ICDTag:           attacks.ICDTagElementalBurst,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Pyro,
		Durability:       25,
		IgnoreDefPercent: c.c6DefIgnore(true),
	}
	for i, mult := range burstWhiteInitial {
		ai.Mult = mult[c.TalentLvlBurst()]
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: -1.5}, 5),
			burstInitHitmark[i],
			burstInitHitmark[i],
			c.c6WhiteMakeCB(),
		)
	}

	c.burstSrc = c.Core.F
	for i := 0.0; i < burstTicksWhite; i++ {
		c.Core.Tasks.Add(c.burstTickWhite(c.burstSrc), burstFirstTickDelayWhite+ceil(burstIntervalWhite*i))
	}
	c.DeleteStatus(burstKeyBlack)
	c.AddStatus(burstKeyWhite, burstFirstTickDelayWhite+ceil((burstTicksWhite-1)*burstIntervalWhite), false)

	c.SetCDWithDelay(action.ActionBurst, burstCD, 22)
	c.ConsumeEnergy(10)
	c.a1OnBurst(true)
	c.a4OnBurst()
	c.c1OnBurst(true)
	c.c4OnBurst()
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
			ICDTag:           attacks.ICDTagDurinBurst,
			ICDGroup:         attacks.ICDGroupDurinBurst,
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

func (c *char) burstBlack() (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex:       c.Index(),
		Abil:             "As the Stars Smolder",
		AttackTag:        attacks.AttackTagElementalBurst,
		ICDTag:           attacks.ICDTagElementalBurst,
		ICDGroup:         attacks.ICDGroupDefault,
		StrikeType:       attacks.StrikeTypeDefault,
		Element:          attributes.Pyro,
		Durability:       25,
		IgnoreDefPercent: c.c6DefIgnore(false),
	}
	for i, mult := range burstBlackInitial {
		ai.Mult = mult[c.TalentLvlBurst()]
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: -1.5}, 5), burstInitHitmark[i], burstInitHitmark[i])
	}

	c.burstSrc = c.Core.F
	for i := 0.0; i < burstTicksBlack; i++ {
		c.Core.Tasks.Add(c.burstTickBlack(c.burstSrc), burstFirstTickDelayBlack+ceil(burstIntervalBlack*i))
	}
	c.DeleteStatus(burstKeyWhite)
	c.AddStatus(burstKeyBlack, burstFirstTickDelayBlack+ceil((burstTicksBlack-1)*burstIntervalBlack), false)

	c.SetCDWithDelay(action.ActionBurst, burstCD, 22)
	c.ConsumeEnergy(10)
	c.a1OnBurst(false)
	c.a4OnBurst()
	c.c1OnBurst(false)
	c.c4OnBurst()
	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) burstTickBlack(src int) func() {
	return func() {
		if src != c.burstSrc {
			return
		}

		ai := info.AttackInfo{
			ActorIndex:       c.Index(),
			Abil:             "Dragon of Dark Decay",
			AttackTag:        attacks.AttackTagElementalBurst,
			ICDTag:           attacks.ICDTagDurinBurst,
			ICDGroup:         attacks.ICDGroupDurinBurst,
			StrikeType:       attacks.StrikeTypeDefault,
			Element:          attributes.Pyro,
			Durability:       25,
			Mult:             burstBlackDoT[c.TalentLvlBurst()] * c.a4Dmg(),
			IgnoreDefPercent: c.c6DefIgnore(false),
		}

		c.Core.QueueAttack(ai, combat.NewSingleTargetHit(c.Core.Combat.PrimaryTarget().Key()), 0, 0)
	}
}
