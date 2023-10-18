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

var (
	skillPressFrames           [][]int
	skillShortHoldFrames       [][]int
	skillShortHold0TicksFrames [][]int

	skillPressHitmarks = []int{25, 26}
)

const (
	skillPressCdStart            = 24
	skillPressSpiritThornHitmark = 70

	skillShortHold0TicksCdStart                  = 11
	skillShortHold0TicksTorrentSurgeHitmark      = 13
	skillShortHold0TicksSpiritbreathThornHitmark = 57

	skillShortHoldCdStart                  = 56
	skillShortHoldFirstDewdropRelease      = 54
	skillShortHoldTorrentSurgeHitmark      = 57
	skillShortHoldSpiritbreathThornHitmark = 103

	particleICDKey          = "travelerhydro-particle-icd"
	spiritbreathThornICDKey = "travelerhydro-spiritbreath-icd"
	skillLosingHPICDKey     = "travelerhydro-losing-hp-icd"
)

func init() {
	// Tap E
	skillPressFrames = make([][]int, 2)

	// Male
	skillPressFrames[0] = frames.InitAbilSlice(45) // Tap E -> E
	skillPressFrames[0][action.ActionAttack] = 44  // Tap E -> N1
	skillPressFrames[0][action.ActionBurst] = 44   // Tap E -> Q
	skillPressFrames[0][action.ActionDash] = 40    // Tap E -> D
	skillPressFrames[0][action.ActionJump] = 41    // Tap E -> J
	skillPressFrames[0][action.ActionWalk] = 44    // Tap E -> Walk
	skillPressFrames[0][action.ActionSwap] = 43    // Tap E -> Swap

	// Female
	skillPressFrames[1] = frames.InitAbilSlice(44) // Tap E -> E/Q/Walk
	skillPressFrames[1][action.ActionAttack] = 43  // Tap E -> N1
	skillPressFrames[1][action.ActionDash] = 41    // Tap E -> D
	skillPressFrames[1][action.ActionJump] = 40    // Tap E -> J
	skillPressFrames[1][action.ActionSwap] = 43    // Tap E -> Swap

	// Short Hold E (0 ticks)
	skillShortHold0TicksFrames = make([][]int, 2)

	// Male
	skillShortHold0TicksFrames[0] = frames.InitAbilSlice(29) // Short Hold E (0 ticks) -> D/J
	skillShortHold0TicksFrames[0][action.ActionAttack] = 36  // Short Hold E (0 ticks) -> N1
	skillShortHold0TicksFrames[0][action.ActionSkill] = 36   // Short Hold E (0 ticks) -> E
	skillShortHold0TicksFrames[0][action.ActionBurst] = 36   // Short Hold E (0 ticks) -> Q
	skillShortHold0TicksFrames[0][action.ActionWalk] = 35    // Short Hold E (0 ticks) -> Walk
	skillShortHold0TicksFrames[0][action.ActionSwap] = 44    // Short Hold E (0 ticks) -> Swap

	// Female
	skillShortHold0TicksFrames[1] = frames.InitAbilSlice(29) // Short Hold E (0 ticks) -> D
	skillShortHold0TicksFrames[1][action.ActionAttack] = 36  // Short Hold E (0 ticks) -> N1
	skillShortHold0TicksFrames[1][action.ActionSkill] = 37   // Short Hold E (0 ticks) -> E
	skillShortHold0TicksFrames[1][action.ActionBurst] = 35   // Short Hold E (0 ticks) -> Q
	skillShortHold0TicksFrames[1][action.ActionJump] = 30    // Short Hold E (0 ticks) -> J
	skillShortHold0TicksFrames[1][action.ActionWalk] = 36    // Short Hold E (0 ticks) -> Walk
	skillShortHold0TicksFrames[1][action.ActionSwap] = 43    // Short Hold E (0 ticks) -> Swap

	// Short Hold E
	skillShortHoldFrames = make([][]int, 2)

	// Male
	skillShortHoldFrames[0] = frames.InitAbilSlice(90) // Short Hold E -> Swap
	skillShortHoldFrames[0][action.ActionAttack] = 81  // Short Hold E -> N1
	skillShortHoldFrames[0][action.ActionSkill] = 81   // Short Hold E -> E
	skillShortHoldFrames[0][action.ActionBurst] = 81   // Short Hold E -> Q
	skillShortHoldFrames[0][action.ActionDash] = 74    // Short Hold E -> D
	skillShortHoldFrames[0][action.ActionJump] = 74    // Short Hold E -> J
	skillShortHoldFrames[0][action.ActionWalk] = 81    // Short Hold E -> Walk

	// Female
	skillShortHoldFrames[1] = frames.InitAbilSlice(89) // Short Hold E -> Swap
	skillShortHoldFrames[1][action.ActionAttack] = 81  // Short Hold E -> N1
	skillShortHoldFrames[1][action.ActionSkill] = 81   // Short Hold E -> E
	skillShortHoldFrames[1][action.ActionBurst] = 81   // Short Hold E -> Q
	skillShortHoldFrames[1][action.ActionDash] = 74    // Short Hold E -> D
	skillShortHoldFrames[1][action.ActionJump] = 73    // Short Hold E -> J
	skillShortHoldFrames[1][action.ActionWalk] = 80    // Short Hold E -> Walk
}

func (c *char) skillPress(hitmark, spiritHitmark, cdStart int, skillFrames [][]int) (action.Info, error) {
	c.torrentSurge(hitmark, spiritHitmark)
	c.SetCDWithDelay(action.ActionSkill, 10*60, cdStart)

	if c.Base.Cons >= 4 {
		c.c4()
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames[c.gender]),
		AnimationLength: skillFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   skillFrames[c.gender][action.ActionDash], // earliest cancel
		State:           action.SkillState,
		OnRemoved:       func(next action.AnimationState) { c.c4Remove() },
	}, nil
}

func (c *char) skillShortHold(travel int) (action.Info, error) {
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

	c.QueueCharTask(func() {
		c.skillLosingHP(&aiHold)
		c.Core.QueueAttack(
			aiHold,
			combat.NewBoxHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), geometry.Point{Y: -0.4}, 0.3, 1.3),
			0,
			1,
			c.makeA1CB(),
			c.makeC4CB(),
		)
	}, skillShortHoldFirstDewdropRelease+travel)

	c.torrentSurge(skillShortHoldTorrentSurgeHitmark, skillShortHoldSpiritbreathThornHitmark)
	c.SetCDWithDelay(action.ActionSkill, 10*60, skillShortHoldCdStart)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillShortHoldFrames[c.gender]),
		AnimationLength: skillShortHoldFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   skillShortHoldFrames[c.gender][action.ActionJump], // earliest cancel
		State:           action.SkillState,
		OnRemoved:       func(next action.AnimationState) { c.c4Remove() },
	}, nil
}

func (c *char) skillHold(travel, holdTicks int) (action.Info, error) {
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

	extend := 15 * (holdTicks - 1)
	a1cb := c.makeA1CB()
	c4cb := c.makeC4CB()
	for i := 0; i <= extend; i += 15 {
		c.QueueCharTask(func() {
			c.skillLosingHP(&aiHold)
			c.Core.QueueAttack(
				aiHold,
				combat.NewBoxHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), geometry.Point{Y: -0.4}, 0.3, 1.3),
				0,
				1,
				a1cb,
				c4cb,
			)
			aiHold.FlatDmg = 0
		}, skillShortHoldFirstDewdropRelease+i+travel)
	}

	c.torrentSurge(skillShortHoldTorrentSurgeHitmark+extend, skillShortHoldSpiritbreathThornHitmark+extend)
	c.SetCDWithDelay(action.ActionSkill, 10*60, skillShortHoldCdStart+extend)

	return action.Info{
		Frames:          func(next action.Action) int { return skillShortHoldFrames[c.gender][next] + extend },
		AnimationLength: skillShortHoldFrames[c.gender][action.InvalidAction] + extend,
		CanQueueAfter:   skillShortHoldFrames[c.gender][action.ActionJump] + extend, // earliest cancel
		State:           action.SkillState,
		OnRemoved:       func(next action.AnimationState) { c.c4Remove() },
	}, nil
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	hold := p["hold"] == 1
	holdTicks := p["hold_ticks"]
	if hold {
		holdTicks = 22
	} else if holdTicks > 0 {
		hold = true
	}
	if holdTicks > 22 {
		holdTicks = 22
	}

	// for skill hold
	travel, ok := p["travel"]
	if !ok {
		travel = 6
	}
	switch {
	case !hold:
		// hold=0
		return c.skillPress(
			skillPressHitmarks[c.gender],
			skillPressSpiritThornHitmark,
			skillPressCdStart,
			skillPressFrames,
		)
	case holdTicks == 0:
		// hold=1, hold_ticks=0
		return c.skillPress(
			skillShortHold0TicksTorrentSurgeHitmark,
			skillShortHold0TicksSpiritbreathThornHitmark,
			skillShortHold0TicksCdStart,
			skillShortHold0TicksFrames,
		)
	case holdTicks == 1:
		// hold=1, hold_ticks=1
		return c.skillShortHold(travel)
	default:
		// hold=1, hold_ticks>1
		return c.skillHold(travel, holdTicks)
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

func (c *char) torrentSurge(hitmark, spiritHitmark int) {
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

	if c.Base.Ascension >= 4 {
		ai.FlatDmg += c.a4Bonus
		c.a4Bonus = 0
	}

	hitbox := combat.NewBoxHitOnTarget(c.Core.Combat.Player(), nil, 1.2, 15)
	c.Core.QueueAttack(ai, hitbox, hitmark, hitmark, c.skillParticleCB)

	if !c.StatusIsActive(spiritbreathThornICDKey) {
		ai = combat.AttackInfo{
			ActorIndex:         c.Index,
			Abil:               "Spiritbreath Thorn",
			AttackTag:          attacks.AttackTagElementalArt,
			ICDTag:             attacks.ICDTagNone,
			ICDGroup:           attacks.ICDGroupDefault,
			StrikeType:         attacks.StrikeTypePierce,
			Element:            attributes.Hydro,
			Durability:         0,
			Mult:               spiritbreathThorn[c.TalentLvlSkill()],
			CanBeDefenseHalted: true,
		}

		c.Core.QueueAttack(ai, hitbox, spiritHitmark, spiritHitmark)
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

	// If HP has been consumed via Suffusion while using the Hold Mode Aquacrest Saber, the Torrent Surge at the skill's end
	// will deal Bonus DMG equal to 45% of the total HP the Traveler has consumed in this skill use via Suffusion.
	// The maximum DMG Bonus that can be gained this way is 5,000.
	if c.Base.Ascension >= 4 {
		c.a4Bonus += drainHP * 0.45
		if c.a4Bonus > 5000 {
			c.a4Bonus = 5000
		}

		c.Core.Log.NewEvent("travelerhydro a4 adding dmg bonus", glog.LogCharacterEvent, c.Index).
			Write("dmg bonus", c.a4Bonus)
	}
	c.AddStatus(skillLosingHPICDKey, 0.9*60, true)
}
