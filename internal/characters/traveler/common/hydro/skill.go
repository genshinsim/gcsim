package hydro

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

// TODO: dendromc/anemomc based frames
var (
	skillPressFrames     [][]int
	skillHoldDelayFrames [][]int
)

const (
	skillPressHitmark = 28
	skillPressCdStart = 25

	particleICDKey          = "travelerhydro-particle-icd"
	spiritbreathThornICDKey = "travelerhydro-spiritbreath-icd"
	skillLosingHPICDKey     = "travelerhydro-losing-hp-icd"
)

func init() {
	// Tap E
	skillPressFrames = make([][]int, 2)

	// Male
	skillPressFrames[0] = frames.InitAbilSlice(37) // E -> N1
	skillPressFrames[0][action.ActionDash] = 29    // E -> D
	skillPressFrames[0][action.ActionJump] = 29    // E -> J
	skillPressFrames[0][action.ActionSwap] = 36    // E -> Swap

	// Female
	skillPressFrames[1] = frames.InitAbilSlice(37) // E -> N1/Q
	skillPressFrames[1][action.ActionDash] = 28    // E -> D
	skillPressFrames[1][action.ActionJump] = 28    // E -> J
	skillPressFrames[1][action.ActionSwap] = 35    // E -> Swap

	// Short Hold E as base for Hold E frames
	// "2 tick duration - 2 tick last hitmark"
	skillHoldDelayFrames = make([][]int, 2)

	// Male
	skillHoldDelayFrames[0] = frames.InitAbilSlice(98 - 54) // Short Hold E -> N1/Q - Short Hold E -> D
	skillHoldDelayFrames[0][action.ActionDash] = 0          // Short Hold E -> D - Short Hold E -> D
	skillHoldDelayFrames[0][action.ActionJump] = 0          // Short Hold E -> J - Short Hold E -> D
	skillHoldDelayFrames[0][action.ActionSwap] = 89 - 54    // Short Hold E -> Swap - Short Hold E -> D

	// Female
	skillHoldDelayFrames[1] = frames.InitAbilSlice(84 - 54) // Short Hold E -> Q - Short Hold E -> D
	skillHoldDelayFrames[1][action.ActionAttack] = 83 - 54  // Short Hold E -> N1 - Short Hold E -> D
	skillHoldDelayFrames[1][action.ActionDash] = 0          // Short Hold E -> D - Short Hold E -> D
	skillHoldDelayFrames[1][action.ActionJump] = 0          // Short Hold E -> J - Short Hold E -> D
	skillHoldDelayFrames[1][action.ActionSwap] = 83 - 54    // Short Hold E -> Swap - Short Hold E -> D
}

func (c *char) SkillPress() action.ActionInfo {
	c.QueueCharTask(c.torrentSurge, skillPressHitmark-1)
	c.SetCDWithDelay(action.ActionSkill, 10*60, skillPressCdStart)

	if c.Base.Cons >= 4 {
		c.c4()
	}

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillPressFrames[c.gender]),
		AnimationLength: skillPressFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   skillPressFrames[c.gender][action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillParticleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.3*60, true)

	count := 3.0
	if c.Core.Rand.Float64() < 0.33 {
		count = 4
	}
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Hydro, c.ParticleDelay)
}

func (c *char) SkillHold(holdTicks int) action.ActionInfo {
	c.a4Bonus = 0

	aiHold := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Dewdrop (Hold)",
		AttackTag:  attacks.AttackTagElementalArtHold,
		ICDTag:     attacks.ICDTagTravelerDewdrop,
		ICDGroup:   attacks.ICDGroupTravelerDewdrop,
		StrikeType: attacks.StrikeTypeSlash,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       skillDewdrop[c.TalentLvlSkill()],
	}

	if c.Base.Cons >= 4 {
		c.c4()
	}

	firstTick := 31
	hitmark := firstTick
	for i := 0; i < holdTicks; i++ {
		c.QueueCharTask(func() {
			c.skillLosingHP(&aiHold)
			c.Core.QueueAttack(
				aiHold,
				combat.NewBoxHitOnTarget(c.Core.Combat.Player(), geometry.Point{Y: -0.4}, 0.3, 1.3), // TODO: or single target hit?
				1,
				1,
				c.a1,
				c.c4CB,
			)
			aiHold.FlatDmg = 0
		}, hitmark-1)
		hitmark += 15
	}

	c.QueueCharTask(c.torrentSurge, hitmark)

	return action.ActionInfo{
		Frames:          func(next action.Action) int { return skillHoldDelayFrames[c.gender][next] + hitmark },
		AnimationLength: skillHoldDelayFrames[c.gender][action.InvalidAction] + hitmark,
		CanQueueAfter:   skillHoldDelayFrames[c.gender][action.ActionDash] + hitmark, // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	holdTicks := 0
	if p["hold"] == 1 {
		holdTicks = 21
	}
	if p["hold_ticks"] > 0 {
		holdTicks = p["hold_ticks"]
	}
	if holdTicks > 21 {
		holdTicks = 21
	}

	if holdTicks == 0 {
		return c.SkillPress()
	} else {
		return c.SkillHold(holdTicks)
	}
}

func (c *char) torrentSurge() {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Torrent Surge",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}

	// If HP has been consumed via Suffusion while using the Hold Mode Aquacrest Saber, the Torrent Surge at the skill's end
	// will deal Bonus DMG equal to 45% of the total HP the Traveler has consumed in this skill use via Suffusion.
	// The maximum DMG Bonus that can be gained this way is 5,000.
	if c.a4Bonus > 5000 {
		c.a4Bonus = 5000
	}
	if c.Base.Ascension >= 4 {
		ai.FlatDmg += c.a4Bonus
	}

	hitbox := combat.NewBoxHitOnTarget(c.Core.Combat.Player(), nil, 1.2, 15)
	c.Core.QueueAttack(ai, hitbox, 0, 1, c.skillParticleCB)

	if !c.StatusIsActive(spiritbreathThornICDKey) {
		ai = combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Spiritbreath Thorn",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Hydro,
			Durability: 0,
			Mult:       spiritbreathThorn[c.TalentLvlSkill()],
		}

		c.Core.QueueAttack(ai, hitbox, 0.7*60, 0.7*60, c.skillParticleCB)
		c.AddStatus(spiritbreathThornICDKey, 9*60, true)
	}
}

func (c *char) skillLosingHP(ai *combat.AttackInfo) {
	if c.StatusIsActive(skillLosingHPICDKey) {
		return
	}
	if c.CurrentHPRatio() <= 0.5 {
		return
	}
	ai.FlatDmg = dewdropBonus[c.TalentLvlSkill()] * c.MaxHP()

	drainHP := 0.04 * c.CurrentHP()
	c.Core.Player.Drain(player.DrainInfo{
		ActorIndex: c.Index,
		Abil:       "Suffusion",
		Amount:     drainHP,
	})
	if c.Base.Ascension >= 4 {
		c.a4Bonus += drainHP * 0.45
		c.Core.Log.NewEvent("hmc a4 adding dmg bonus", glog.LogCharacterEvent, c.Index).
			Write("dmg bonus", c.a4Bonus)
	}
	c.AddStatus(skillLosingHPICDKey, 0.9*60, true)
}
