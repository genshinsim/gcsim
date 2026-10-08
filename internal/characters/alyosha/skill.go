package alyosha

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var (
	skillFrames     []int
	skillHoldFrames []int
)

const (
	skillHitmark     = 18
	skillHoldHitmark = 113
	particleICDKey   = "alyosha-particle-icd"
	skillMarkKey     = "alyosha-hunters-mark"
	skillBuffKey     = "alyosha-hunters-precision"
	skillBuffDur     = 902
)

func init() {
	skillFrames = frames.InitAbilSlice(31)
	skillFrames[action.ActionAttack] = 30
	skillFrames[action.ActionBurst] = 30
	skillFrames[action.ActionDash] = 30
	skillFrames[action.ActionWalk] = 30
	skillFrames[action.ActionSwap] = 29

	skillHoldFrames = frames.InitAbilSlice(142)
	skillHoldFrames[action.ActionAttack] = 138
	skillHoldFrames[action.ActionBurst] = 138
	skillHoldFrames[action.ActionDash] = 138
	skillHoldFrames[action.ActionJump] = 138
	skillHoldFrames[action.ActionSwap] = 137
}

// Follows his hunter's instincts. Tap or hold to produce different effects as well as apply the
// Hunter's Mark effect to the hit opponents.
//
// Tapping
// Fires at opponents in front of him, dealing AoE Electro DMG.
func (c *char) Skill(p map[string]int) (action.Info, error) {
	hold, ok := p["hold"]
	if ok && hold > 0 {
		return c.skillHold()
	}

	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Skill",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skillTap[c.TalentLvlSkill()],
	}

	ap := combat.NewBoxHitOnTarget(c.Core.Combat.Player(), nil, 2.5, 10)
	c.Core.QueueAttack(
		ai,
		ap,
		skillHitmark,
		skillHitmark,
		c.baseParticleCB,
		c.triggerSkillMarkCB(true),
	)

	c.SetCDWithDelay(action.ActionSkill, 15*60, 16)

	return action.Info{
		Frames:          func(next action.Action) int { return skillFrames[next] },
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}, nil
}

// Holding
// Enters aiming mode, then targets opponents within a certain area in front of him. Alyosha's Interruption RES is enhanced when in this mode.
// At the end of the button hold, he deals Electro DMG to the targeted opponents.
func (c *char) skillHold() (action.Info, error) {
	ai := info.AttackInfo{
		ActorIndex: c.Index(),
		Abil:       "Skill (Hold)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skillHold[c.TalentLvlSkill()],
	}
	ap := combat.NewCircleHitOnTargetFanAngle(c.Core.Combat.Player(), nil, 15, 150)
	c.Core.QueueAttack(
		ai,
		ap,
		skillHoldHitmark,
		skillHoldHitmark,
		c.baseParticleCB,
		c.triggerSkillMarkCB(true),
	)

	c.SetCDWithDelay(action.ActionSkill, 15*60, 111)

	return action.Info{
		Frames:          func(next action.Action) int { return skillHoldFrames[next] },
		AnimationLength: skillHoldFrames[action.InvalidAction],
		CanQueueAfter:   skillHoldFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) baseParticleCB(a info.AttackCB) {
	if a.Target.Type() != info.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.5*60, true)

	c.Core.QueueParticle(c.Base.Key.String(), 5, attributes.Electro, c.ParticleDelay)
}

// Hunter's Mark
// A special effect Alyosha can apply on the opponent.
// When this effect is applied to an opponent who is already under the Hunter's Mark effect, it will
// activate Hunter's Mark instead, which removes said mark and grants all nearby party members the
// Hunter's Precision effect, which increases their ATK when they are on the field.
func (c *char) skillInit() {
	m := make([]float64, attributes.EndStatType)
	for _, char := range c.Core.Player.Chars() {
		char.AddStatMod(character.StatMod{
			Base:         modifier.NewBase(skillBuffKey+"-atkp", -1),
			AffectedStat: attributes.ATKP,
			Amount: func() []float64 {
				if c.Core.Player.Active() != char.Index() {
					return nil
				}
				if !c.StatusIsActive(skillBuffKey) {
					return nil
				}

				m[attributes.ATKP] = skillBuff[c.TalentLvlSkill()] * float64(c.skillStacks)
				return m
			},
		})
	}
}

func (c *char) triggerSkillMarkCB(canApply bool) info.AttackCBFunc {
	return func(ac info.AttackCB) {
		if ac.Target.Type() != info.TargettableEnemy {
			return
		}

		e, ok := ac.Target.(*enemy.Enemy)
		if !ok {
			return
		}

		// Hunter's Mark is only applied 6 frames after the hitmark
		c.Core.Tasks.Add(func() {
			if e.StatusIsActive(skillMarkKey) {
				e.DeleteStatus(skillMarkKey)
				if !c.StatusIsActive(skillBuffKey) {
					c.skillStacks = 0
				}

				c.skillStacks = min(c.skillStacks+1, 1+c.c6ExtraSkillStacks())
				c.AddStatus(skillBuffKey, skillBuffDur, true)
				return
			}

			if canApply {
				e.AddStatus(skillMarkKey, skillBuffDur, true)
			}
		}, 6)
	}
}
