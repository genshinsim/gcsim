package barbara

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

// barbara skill - copied from bennett burst
const skillDuration = 15*60 + 1
const barbSkillKey = "barbara-e"
const skillCDStart = 3

var (
	skillHitmarks = []int{42, 78}
	skillFrames   []int
)

func init() {
	skillFrames = frames.InitAbilSlice(55)
	skillFrames[action.ActionWalk] = 54
	skillFrames[action.ActionDash] = 4
	skillFrames[action.ActionJump] = 5
	skillFrames[action.ActionSwap] = 53
	skillFrames[action.ActionSkill] = 54
	skillFrames[action.ActionAttack] = 54
	skillFrames[action.ActionCharge] = 54
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	// restart a4 counter
	c.a4extendCount = 0

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Let the Show Begin♪ (Droplet)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}

	// 2 Droplets
	for _, hitmark := range skillHitmarks {
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 3),
			5,
			hitmark,
		) // need to confirm snapshot timing
	}

	c.skillInitF = c.Core.F // needed for ticks

	// setup heal and wet ticks (first tick at skillCDStart, once every 5s)
	stats, _ := c.Stats()
	hpplus := stats[attributes.Heal]
	heal := skillhp[c.TalentLvlSkill()] + skillhpp[c.TalentLvlSkill()]*c.MaxHP()

	// setup Melody Loop ticks (first tick at skillCDStart, once every 1.5s)
	ai.Abil = "Let the Show Begin♪ (Melody Loop)"
	ai.AttackTag = attacks.AttackTagNone
	ai.Mult = 0
	ai.HitlagFactor = 0.05
	ai.HitlagHaltFrames = 0.05 * 60
	ai.CanBeDefenseHalted = true
	ai.IsDeployable = true

	// add skill status and queue up ticks
	c.Core.Tasks.Add(func() {
		c.Core.Status.Add(barbSkillKey, skillDuration)
		c.a1()
		c.barbaraSelfTick(heal, hpplus, c.skillInitF)()
		c.barbaraMelodyTick(ai, c.skillInitF)()
	}, skillCDStart)

	if c.Base.Cons >= 2 {
		c.c2() // c2 hydro buff
		c.SetCDWithDelay(action.ActionSkill, 32*60*0.85, skillCDStart)
	} else {
		c.SetCDWithDelay(action.ActionSkill, 32*60, skillCDStart)
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash],
		State:           action.SkillState,
	}
}

func (c *char) barbaraSelfTick(healAmt float64, hpplus float64, skillInitF int) func() {
	return func() {
		// make sure it's not overwritten
		if c.skillInitF != skillInitF {
			return
		}
		// do nothing if buff expired
		if c.Core.Status.Duration(barbSkillKey) == 0 {
			return
		}

		c.Core.Log.NewEvent("barbara heal and wet ticking", glog.LogCharacterEvent, c.Index)

		// heal
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  c.Core.Player.Active(),
			Message: "Melody Loop (Tick)",
			Src:     healAmt,
			Bonus:   hpplus,
		})

		// wet: apply self infusion for 0.3s
		p, ok := c.Core.Combat.Player().(*avatar.Player)
		if !ok {
			panic("target 0 should be Player but is not!!")
		}
		p.ApplySelfInfusion(attributes.Hydro, 25, 0.3*60)

		// tick every 5s
		c.Core.Tasks.Add(c.barbaraSelfTick(healAmt, hpplus, skillInitF), 5*60)
	}
}

func (c *char) barbaraMelodyTick(ai combat.AttackInfo, skillInitF int) func() {
	return func() {
		// make sure it's not overwritten
		if c.skillInitF != skillInitF {
			return
		}
		// do nothing if buff expired
		if c.Core.Status.Duration(barbSkillKey) == 0 {
			return
		}

		c.Core.Log.NewEvent("barbara melody loop ticking", glog.LogCharacterEvent, c.Index)

		// 0 DMG attack that causes hitlag on enemy only
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 1), -1, 0)

		// tick every 1.5s
		c.Core.Tasks.Add(c.barbaraMelodyTick(ai, skillInitF), 1.5*60)
	}
}
