package jahoda

import (
	"errors"
	"fmt"
	"sort"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

var burstFrames []int

const (
	burstDuration = 790
	burstHitmark  = 43

	firstRobotHitmarkDelay = 41

	firstRobotAbsorbDelay   = 6
	robotAbsorptionInterval = 41 // exclude the outlier in trial 3

	firsHealTickDelay = 12
	healInterval      = 87

	burstCD  = 18 * 60
	burstKey = "jahoda-burst-dot"
)

func init() {
	burstFrames = frames.InitAbilSlice(55) // Q -> Aim/W
	burstFrames[action.ActionAttack] = 48  // Q -> N1
	burstFrames[action.ActionSkill] = 53   // Q -> Skill
	burstFrames[action.ActionDash] = 54    // Q -> D
	burstFrames[action.ActionJump] = 54    // Q -> J
	burstFrames[action.ActionSwap] = 36    // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(shadowPursuitKey) {
		return action.Info{}, errors.New("burst called in skill state")
	}

	c.burstSrc = c.Core.F
	src := c.burstSrc
	c.burstAbsorbCheckLocation = combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 5)
	c.robotHitmarkInterval = 140

	c.AddStatus(burstKey, burstDuration, false)

	// initial hit damage
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Hidden Aces: Seven Tools of the Hunter",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Anemo,
		Durability: 25,
		Mult:       burst[c.TalentLvlBurst()],
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.Player(), info.Point{Y: 1}, 5),
		0,
		burstHitmark)

	c.robotCount = 2

	// apply a1 buff
	c.a1()

	// heal ticks
	c.QueueCharTask(func() {
		for i := 0; i < burstDuration-firsHealTickDelay; i += healInterval {
			c.Core.Tasks.Add(func() {
				if src != c.burstSrc {
					return
				}

				if c.Core.Player.ActiveChar().CurrentHPRatio() > 0.7 {
					c.a4()

					low := c.lowestHPChar()
					if low >= 0 {
						healOffField := (burstAdditionalHealFlat[c.TalentLvlBurst()] + burstAdditionalHealPP[c.TalentLvlBurst()]*c.TotalAtk()) * c.a1HealMult

						c.Core.Player.Heal(info.HealInfo{
							Caller:  c.Index(),
							Target:  low,
							Message: "Additional Healing",
							Src:     healOffField,
							Bonus:   c.Stat(attributes.Heal),
						})
					}
				}

				heal := (burstHealFlat[c.TalentLvlBurst()] + burstHealPP[c.TalentLvlBurst()]*c.TotalAtk()) * c.a1HealMult
				robotHi := info.HealInfo{
					Caller:  c.Index(),
					Target:  c.Core.Player.Active(),
					Message: "Purrsonal Coordinated Assistance Robot Healing",
					Src:     heal,
					Bonus:   c.Stat(attributes.Heal),
				}

				c.Core.Player.Heal(robotHi)
			}, i)
		}
	}, firsHealTickDelay)

	// dmg ticks
	if c.Core.Player.GetMoonsignLevel() >= 2 {
		for robot := 0; robot < c.robotCount; robot++ {
			c.Core.Tasks.Add(c.absorbCheckTask(src, robot), firstRobotAbsorbDelay+robotAbsorptionInterval*robot)
		}
	}

	c.SetCDWithDelay(action.ActionBurst, burstCD, 1)
	c.ConsumeEnergy(13)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionSwap], // earliest cancel
		State:           action.BurstState,
	}, nil
}

func (c *char) lowestHPChar() int {
	lowestIdx := -1
	lowestPct := 2.0 // > 1

	for i := 0; i < len(c.Core.Player.Chars()); i++ {
		ch := c.Core.Player.Chars()[i]
		if ch.CurrentHP() <= 0 {
			continue
		}

		if ch.CurrentHPRatio() < lowestPct {
			lowestPct = ch.CurrentHPRatio()
			lowestIdx = i
		}
	}

	return lowestIdx
}

func (c *char) absorbCheckTask(src, robot int) func() {
	return func() {
		if src != c.burstSrc || !c.StatusIsActive(burstKey) {
			return
		}

		ele := c.enemyAuraInArea(c.burstAbsorbCheckLocation, elePriority)

		// define attack info
		if ele == attributes.NoElement {
			c.Core.Tasks.Add(c.absorbCheckTask(src, robot), int(c.robotHitmarkInterval/3)) // formular from dm
			return
		}

		ai := info.AttackInfo{
			ActorIndex: c.Index(),
			Abil:       "Purrsonal Coordinated Assistance Robot DMG",
			AttackTag:  attacks.AttackTagElementalBurst,
			ICDTag:     attacks.ICDTagElementalBurst,
			ICDGroup:   attacks.ICDGroupJahodaBurst, // special icd, 15s/4 hits
			StrikeType: attacks.StrikeTypeDefault,
			Element:    ele,
			Durability: 25,
			FlatDmg:    burstSkill[c.TalentLvlBurst()] * c.TotalAtk() * c.a1DmgMult,
		}

		switch ele {
		case attributes.Pyro:
			ai.ICDTag = attacks.ICDTagElementalBurstPyro
		case attributes.Hydro:
			ai.ICDTag = attacks.ICDTagElementalBurstHydro
		case attributes.Electro:
			ai.ICDTag = attacks.ICDTagElementalBurstElectro
		case attributes.Cryo:
			ai.ICDTag = attacks.ICDTagElementalBurstCryo
		default:
			ai.ICDTag = attacks.ICDTagElementalBurst
		}

		c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index(), fmt.Sprintf("jahoda robot %d absorbed %v", robot, ele.String()))

		c.c4()
		c.Core.Tasks.Add(c.robotAtkTick(src, ai), firstRobotHitmarkDelay)
	}
}

func (c *char) robotAtkTick(src int, ai info.AttackInfo) func() {
	return func() {
		if src != c.burstSrc {
			return
		}

		if !c.StatusIsActive(burstKey) {
			return
		}

		// trigger an instance of damage on 3 closest enemies
		c.queueOn3Closest(c.Core.Combat.Player().Pos(), ai, 0)

		// schedule another attack
		c.Core.Tasks.Add(c.robotAtkTick(src, ai), int(c.robotHitmarkInterval))
	}
}

// helper to sort 3 closest enemies and attack them simultaneously
func (c *char) queueOn3Closest(origin info.Point, ai info.AttackInfo, hitDelay int) {
	enemies := c.Core.Combat.Enemies()
	type cand struct {
		t info.Target
		d float64
	}
	cands := make([]cand, 0, len(enemies))

	// compute distance
	for _, e := range enemies {
		if e == nil {
			continue
		}
		if e.Type() != info.TargettableEnemy {
			continue
		}
		if !e.IsAlive() {
			continue
		}

		p := e.Pos()
		d := p.Sub(origin).MagnitudeSquared()

		cands = append(cands, cand{t: e, d: d})
	}

	// sort
	sort.Slice(cands, func(i, j int) bool { return cands[i].d < cands[j].d })

	// queue on up to 3
	n := 3
	if len(cands) < n {
		n = len(cands)
	}
	for i := 0; i < n; i++ {
		t := cands[i].t
		ap := combat.NewCircleHitOnTarget(t, nil, 1.2) // couldn't find anywhere in dm
		c.Core.QueueAttack(ai, ap, hitDelay, hitDelay)
	}
}
