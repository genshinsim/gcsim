package yelan

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillFrames []int

const (
	skillHitmark        = 35
	particleICDKey      = "yelan-particle-icd"
	skillTargetCountTag = "marked"
	skillHoldDuration   = "hold_length" // not yet implemented
	skillMarkedTag      = "yelan-skill-marked"
)

func init() {
	skillFrames = frames.InitAbilSlice(42)
	skillFrames[action.ActionBurst] = 41
	skillFrames[action.ActionDash] = 41
	skillFrames[action.ActionJump] = 41
	skillFrames[action.ActionSwap] = 40
}

/*
*
Fires off a Lifeline that tractors her in rapidly, entangling and marking opponents along its path.
When her rapid movement ends, the Lifeline will explode, dealing Hydro DMG to the marked opponents based on Yelan's Max HP.
*
*/
func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Lingering Lifeline",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		Mult:       0,
		FlatDmg:    skill[c.TalentLvlSkill()] * c.MaxHP(),
	}

	// clear all existing tags
	for _, t := range c.Core.Combat.Enemies() {
		if e, ok := t.(*enemy.Enemy); ok {
			e.SetTag(skillMarkedTag, 0)
		}
	}

	if !c.StatusIsActive("yelanc4") {
		c.c4count = 0
		c.Core.Log.NewEvent("c4 stacks set to 0", glog.LogCharacterEvent, c.Index)
	}

	// add a task to loop through targets and mark them
	marked, ok := p[skillTargetCountTag]
	// default 1
	if !ok {
		marked = 1
	}
	c.Core.Tasks.Add(func() {
		for _, t := range c.Core.Combat.Enemies() {
			if marked == 0 {
				break
			}
			e, ok := t.(*enemy.Enemy)
			if !ok {
				continue
			}
			e.SetTag(skillMarkedTag, 1)
			c.Core.Log.NewEvent("marked by Lifeline", glog.LogCharacterEvent, c.Index).
				Write("target", e.Key())
			marked--
			c.c4count++
			if c.Base.Cons >= 4 {
				c.AddStatus("yelanc4", 25*60, true)
			}
		}
	}, skillHitmark) //TODO: frames for hold e

	// hold := p["hold"]

	cb := func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		// check for breakthrough
		if c.Core.Rand.Float64() < 0.34 {
			c.breakthrough = true
			c.Core.Log.NewEvent("breakthrough state added", glog.LogCharacterEvent, c.Index)
		}
		//TODO: icd on this??
		if c.StatusIsActive(burstKey) {
			c.summonExquisiteThrow()
			c.Core.Log.NewEvent("yelan burst on skill", glog.LogCharacterEvent, c.Index)
		}
	}

	// add a task to loop through targets and deal damage if marked
	c.Core.Tasks.Add(func() {
		for _, t := range c.Core.Combat.Enemies() {
			e, ok := t.(*enemy.Enemy)
			if !ok {
				continue
			}
			if e.GetTag(skillMarkedTag) == 0 {
				continue
			}
			e.SetTag(skillMarkedTag, 0)
			c.Core.Log.NewEvent("damaging marked target", glog.LogCharacterEvent, c.Index).
				Write("target", e.Key())
			marked--
			// queueing attack one frame later
			//TODO: does hold have different attack size? don't think so?
			c.Core.QueueAttack(ai, combat.NewSingleTargetHit(e.Key()), 1, 1, c.particleCB, cb)
		}

		// activate c4 if relevant
		//TODO: check if this is accurate
		if c.Base.Cons >= 4 && c.c4count > 0 {
			m := make([]float64, attributes.EndStatType)
			m[attributes.HPP] = float64(c.c4count) * 0.1
			if m[attributes.HPP] > 0.4 {
				m[attributes.HPP] = 0.4
			}
			c.Core.Log.NewEvent("c4 activated", glog.LogCharacterEvent, c.Index).
				Write("enemies count", c.c4count)
			for _, char := range c.Core.Player.Chars() {
				char.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag("yelan-c4", 25*60),
					AffectedStat: attributes.HPP,
					Amount: func() ([]float64, bool) {
						return m, true
					},
				})
			}
		}

	}, skillHitmark) //TODO: frames for e dmg? possibly 5 second after attaching?

	c.SetCDWithDelay(action.ActionSkill, 600, skillHitmark-2)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.3*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 4, attributes.Hydro, c.ParticleDelay) // TODO: this used to be 82?
}
