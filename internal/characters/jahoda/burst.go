package jahoda

import (
	"errors"
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
	burstDuration      = 790
	absorptionInterval = 41
	firstrobotHitmark  = 45
	robotDelay         = 94
	healInterval       = 87
	firsHealTickDelay  = 12
	burstCD            = 18 * 60
	burstKey           = "jahoda-burst-dot"
)

func init() {
	burstFrames = frames.InitAbilSlice(48) // Q -> N1
	burstFrames[action.ActionSkill] = 53   // Q -> Skill
	burstFrames[action.ActionAim] = 55     // Q -> Aim
	burstFrames[action.ActionDash] = 54    // Q -> D
	burstFrames[action.ActionJump] = 54    // Q -> J
	burstFrames[action.ActionWalk] = 55    // Q -> W
	burstFrames[action.ActionSwap] = 36    // Q -> Swap
}

func (c *char) Burst(p map[string]int) (action.Info, error) {
	if c.StatusIsActive(shadowPursuitKey) {
		return action.Info{}, errors.New("burst called in skill state")
	}

	c.robotHitmarkInterval = 140
	c.burstSrc = c.Core.F
	src := c.burstSrc
	c.burstAbsorbCheckLocation = combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 1.2) // Couldn't find anywhere in dm, assume top be the same as Sayu

	c.AddStatus(burstKey, burstDuration, false)

	// Initial hit damage
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
		0)

	// Define Info
	c.robotAi = info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Purrsonal Coordinated Assistance Robot DMG",
		AttackTag:  attacks.AttackTagElementalBurst,
		ICDTag:     attacks.ICDTagElementalBurst,
		ICDGroup:   attacks.ICDGroupJahodaBurst, // special icd, 15s/4 hits
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.NoElement,
		Durability: 25,
		FlatDmg:    burstSkill[c.TalentLvlBurst()] * c.TotalAtk(),
	}

	heal := burstHealFlat[c.TalentLvlBurst()] + burstHealPP[c.TalentLvlBurst()]*c.TotalAtk()
	c.robotHi = info.HealInfo{
		Caller:  c.Index(),
		Target:  c.Core.Player.Active(),
		Message: "Purrsonal Coordinated Assistance Robot Healing",
		Src:     heal,
		Bonus:   c.Stat(attributes.Heal),
	}

	c.robotCount = 2

	// Apply A1 buff
	c.a1()

	// Heal ticks
	c.QueueCharTask(func() {
		for i := 0; i < burstDuration-firsHealTickDelay; i += healInterval {
			c.Core.Tasks.Add(func() {
				if src != c.burstSrc {
					return
				}

				c.Core.Player.Heal(c.robotHi)

				if c.Core.Player.ActiveChar().CurrentHPRatio() > 0.7 {
					c.a4()

					low := c.lowestHPChar()
					if low >= 0 {
						healOffField := burstAdditionalHealFlat[c.TalentLvlBurst()] + burstAdditionalHealPP[c.TalentLvlBurst()]*c.TotalAtk()

						c.Core.Player.Heal(info.HealInfo{
							Caller:  c.Index(),
							Target:  low,
							Message: "Additional Healing",
							Src:     healOffField,
							Bonus:   c.Stat(attributes.Heal),
						})
					}

				}

			}, i)

		}
	}, firsHealTickDelay)

	// Dmg ticks
	if c.Core.Player.GetMoonsignLevel() >= 2 {
		c.Core.Tasks.Add(c.absorbCheck(src, 0, burstDuration/absorptionInterval), 10+absorptionInterval) // Frames needed

	}

	c.SetCDWithDelay(action.ActionBurst, burstCD, 1)
	c.ConsumeEnergy(13)

	return action.Info{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.InvalidAction],
		CanQueueAfter:   burstFrames[action.ActionDash], // Earliest cancel, need to check
		State:           action.BurstState,
	}, nil
}

func (c *char) lowestHPChar() int {
	lowestIdx := -1
	lowestPct := 2.0 // > 1

	for i := 0; i < len(c.Core.Player.Chars()); i++ {
		ch := c.Core.Player.Chars()[i]
		if ch == nil {
			continue
		}

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

func (c *char) absorbCheck(src, count, maxcount int) func() {
	return func() {
		if src != c.burstSrc {
			return
		}

		if count == maxcount {
			return
		}

		c.robotAi.Element = c.Core.Combat.AbsorbCheck(c.Index(), c.burstAbsorbCheckLocation, attributes.Pyro, attributes.Hydro, attributes.Electro, attributes.Cryo)
		if c.robotAi.Element != attributes.NoElement {
			switch c.robotAi.Element {
			case attributes.Pyro:
				c.robotAi.ICDTag = attacks.ICDTagElementalBurstPyro
			case attributes.Hydro:
				c.robotAi.ICDTag = attacks.ICDTagElementalBurstHydro
			case attributes.Electro:
				c.robotAi.ICDTag = attacks.ICDTagElementalBurstElectro
			case attributes.Cryo:
				c.robotAi.ICDTag = attacks.ICDTagElementalBurstCryo
			}

			c.Core.Log.NewEventBuildMsg(glog.LogCharacterEvent, c.Index(),
				"jahoda burst absorbed ", c.robotAi.Element.String(),
			)

			c.c4()

			for i := 0; i < burstDuration-firstrobotHitmark; i += int(c.robotHitmarkInterval) {
				c.Core.Tasks.Add(c.robotAtkTick(src), i)
			}

			return
		}
		c.Core.Tasks.Add(c.absorbCheck(src, count+1, maxcount), absorptionInterval)
	}
}

func (c *char) robotAtkTick(src int) func() {
	return func() {
		if src != c.burstSrc {
			return
		}

		// For each robot, trigger an instance of damage on 3 closest enemies
		for i := 0; i < c.robotCount; i++ {
			c.queueOn3Closest(c.Core.Combat.Player().Pos(), c.robotAi, robotDelay*i)
		}

	}
}

// Helper to sort 3 closest enemies and attack them simultaneously
func (c *char) queueOn3Closest(origin info.Point, ai info.AttackInfo, hitDelay int) {
	enemies := c.Core.Combat.Enemies()
	type cand struct {
		t info.Target
		d float64
	}
	cands := make([]cand, 0, len(enemies))

	// Compute distance
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
		dx := p.X - origin.X
		dy := p.Y - origin.Y
		d := dx*dx + dy*dy // squared distance is enough for sorting (no sqrt)

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
		ap := combat.NewCircleHitOnTarget(t, nil, 5)
		c.Core.QueueAttack(ai, ap, hitDelay, hitDelay)
	}
}
