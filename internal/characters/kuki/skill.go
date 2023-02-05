package kuki

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
)

var skillFrames []int

const (
	skillHitmark     = 11 // Initial Hit
	hpDrainThreshold = 0.2
	ringKey          = "kuki-e"
	particleICDKey   = "kuki-particle-icd"
)

func init() {
	skillFrames = frames.InitAbilSlice(52) // E -> Q
	skillFrames[action.ActionAttack] = 50  // E -> N1
	skillFrames[action.ActionDash] = 12    // E -> D
	skillFrames[action.ActionJump] = 11    // E -> J
	skillFrames[action.ActionSwap] = 50    // E -> Swap
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	// only drain HP when above 20% HP
	if c.HPCurrent/c.MaxHP() > hpDrainThreshold {
		hpdrain := 0.3 * c.HPCurrent
		// The HP consumption from using this skill can only bring her to 20% HP.
		if (c.HPCurrent-hpdrain)/c.MaxHP() <= hpDrainThreshold {
			hpdrain = c.HPCurrent - hpDrainThreshold*c.MaxHP()
		}
		c.Core.Player.Drain(player.DrainInfo{
			ActorIndex: c.Index,
			Abil:       "Sanctifying Ring",
			Amount:     hpdrain,
		})
	}

	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Sanctifying Ring",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
		FlatDmg:    c.Stat(attributes.EM) * 0.25,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 4), skillHitmark, skillHitmark)

	// C2: Grass Ring of Sanctification's duration is increased by 3s.
	skilldur := 720
	if c.Base.Cons >= 2 {
		skilldur = 900 //12+3s
	}

	// this gets executed before kuki can experience hitlag so no need for char queue
	// ring duration starts after hitmark
	c.Core.Tasks.Add(func() {
		// E duration and ticks are not affected by hitlag
		c.Core.Status.Add(ringKey, skilldur)
		c.ringSrc = c.Core.F
		c.Core.Tasks.Add(c.bellTick(c.Core.F), 90) // Assuming this executes every 90 frames = 1.5s
		c.Core.Log.NewEvent("Bell activated", glog.LogCharacterEvent, c.Index).
			Write("next expected tick", c.Core.F+90).
			Write("expected end", c.Core.F+skilldur)
	}, 23)

	c.SetCDWithDelay(action.ActionSkill, 15*60, 7)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionJump], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != combat.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.2*60, false)
	if c.Core.Rand.Float64() < .45 {
		c.Core.QueueParticle(c.Base.Key.String(), 1, attributes.Electro, c.ParticleDelay)
	}
}

func (c *char) bellTick(src int) func() {
	return func() {
		if src != c.ringSrc {
			return
		}
		c.Core.Log.NewEvent("Bell ticking", glog.LogCharacterEvent, c.Index)

		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Grass Ring of Sanctification",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagElementalArt,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeDefault,
			Element:    attributes.Electro,
			Durability: 25,
			Mult:       skilldot[c.TalentLvlSkill()],
			FlatDmg:    c.a4Damage(),
		}
		//trigger damage
		//TODO: Check for snapshots
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 4), 2, 2, c.particleCB)

		//A4 is considered here
		c.Core.Player.Heal(player.HealInfo{
			Caller:  c.Index,
			Target:  c.Core.Player.Active(),
			Message: "Grass Ring of Sanctification Healing",
			Src:     (skillhealpp[c.TalentLvlSkill()]*c.MaxHP() + skillhealflat[c.TalentLvlSkill()] + c.a4Healing()),
			Bonus:   c.Stat(attributes.Heal),
		})

		if c.Core.Status.Duration(ringKey) == 0 {
			return
		}
		c.Core.Tasks.Add(c.bellTick(src), 90)
	}
}
