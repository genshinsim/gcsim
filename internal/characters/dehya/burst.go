package dehya

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

var burstFrames []int
var kickFrames []int

const burstDoT1Hitmark = 102
const fastPunchHitmark = 24 //10 hits max on 240 f
const slowPunchHitmark = 40 //6 hits minimum
const kickHitmark = 46      //6 hits minimum
var punchHitmarks = []int{42, 30, 30, 27, 27, 24, 24, 24, 24}

func init() {
	burstFrames = frames.InitAbilSlice(87) // Q -> E/D/J
	burstFrames[action.ActionAttack] = 86  // Q -> N1
	burstFrames[action.ActionSwap] = 86    // Q -> Swap

	kickFrames = frames.InitAbilSlice(101) // Q -> E/D/J
}

func (c *char) Burst(p map[string]int) action.ActionInfo {
	c1var := 0.0
	if c.Base.Cons >= 1 {
		c1var = 0.06
	}
	punches, ok := p["hits"]
	if !ok || punches < 6 || punches > 10 {
		punches = 10
	}
	if c.sanctumActive {
		c.sanctumActive = false
		c.sanctumExpiry += c.sanctumPickupExtension
		c.sanctumRetrieved = true
		c.sanctumICD = c.StatusDuration("dehya-skill-icd")
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Flame-Mane's Fist",
		AttackTag:  combat.AttackTagElementalBurst,
		ICDTag:     combat.ICDTagElementalBurstPyro,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Pyro,
		Durability: 25,
		Mult:       burstPunchAtk[c.TalentLvlBurst()],
		FlatDmg:    (c1var + burstPunchHP[c.TalentLvlBurst()]) * c.MaxHP(),
	}

	c.QueueCharTask(func() {
		//c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 2}, 4), 0, 0)
		ai.CanBeDefenseHalted = false // TODO:no hitlag?
		punchCounter := 0
		punchTimer := 0
		// one punch on burst hitmark
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 2}, 4), 0, 0)

		for i := 0; i < 240; {
			if punchCounter < punches && punches == 10 {
				c.Core.QueueAttack(
					ai,
					combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 2}, 4),
					i+punchHitmarks[punchCounter],
					i+punchHitmarks[punchCounter],
				)
				i += punchHitmarks[punchCounter]
				// c.Core.Log.NewEvent("Dehya burst punch ", glog.LogCharacterEvent, c.Index).
				// 	Write("punchCount", punchCounter+1).
				// 	Write("frameCount", i)
			} else {
				c.Core.QueueAttack(
					ai,
					combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 2}, 4),
					i+slowPunchHitmark,
					i+slowPunchHitmark,
				)
				i += slowPunchHitmark
			}
			punchCounter++
			punchTimer = i
		}
		ai.Abil = "Incineration Drive"
		ai.Mult = burstKickAtk[c.TalentLvlBurst()]
		ai.FlatDmg = burstKickHP[c.TalentLvlBurst()] * c.MaxHP()
		ai.ICDTag = combat.ICDTagElementalBurst
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 2}, 4),
			punchTimer+kickHitmark,
			punchTimer+kickHitmark,
		)

		if c.sanctumRetrieved {
			c.Core.Tasks.Add(func() {
				c.sanctumActive = true
				c.sanctumExpiry += burstDoT1Hitmark + punchTimer + kickHitmark

				// snapshot for ticks
				ai.Abil = "Molten Inferno (DoT)"
				ai.ICDTag = combat.ICDTagElementalArt
				ai.Mult = skillDotAtk[c.TalentLvlSkill()]
				ai.FlatDmg = skillDotHP[c.TalentLvlSkill()] * c.MaxHP()
				c.skillAttackInfo = ai
				c.skillSnapshot = c.Snapshot(&c.skillAttackInfo)
				c.Core.Tasks.Add(c.removeSanctum(c.sanctumExpiry), c.sanctumExpiry-c.Core.F)
			}, punchTimer+kickHitmark)
			c.AddStatus(skillICDKey, 1+punchTimer+kickHitmark+c.sanctumICD, false)

			c.Core.Log.NewEvent("Sanctum Expiration Info ", glog.LogCharacterEvent, c.Index).
				Write("Duration Remaining", c.sanctumExpiry-c.Core.F+burstDoT1Hitmark).
				Write("New Expiry Frame", c.sanctumExpiry+burstDoT1Hitmark+punchTimer+kickHitmark).
				Write("Field Source", c.sanctumSource).
				Write("DoT tick CD", c.StatusDuration("dehya-skill-icd")-punchTimer-kickHitmark)
		}

	}, burstDoT1Hitmark)

	c.ConsumeEnergy(5)
	c.SetCDWithDelay(action.ActionBurst, 18*60, 1)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(burstFrames),
		AnimationLength: burstFrames[action.ActionAttack] + 240 + kickFrames[action.ActionAttack],
		CanQueueAfter:   burstFrames[action.ActionAttack] + 240 + kickFrames[action.ActionAttack], // earliest cancel
		State:           action.BurstState,
	}
}
