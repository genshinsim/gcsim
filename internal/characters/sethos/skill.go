package sethos

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames []int

const skillParticleICDKey = "sethos-particle-icd"
const skillRefundICDKey = "sethos-refund-icd"

func init() {
	skillFrames = frames.InitAbilSlice(28)
}

func (c *char) skillRefundHook() {
	refundCB := func(args ...interface{}) bool {
		// TODO: Check if Sethos E filters by enemy
		// a := args[0].(combat.Target)
		// if a.Type() != targets.TargettableEnemy {
		// 	return false
		// }
		ae := args[1].(*combat.AttackEvent)
		if ae.Info.ActorIndex != c.Index {
			return false
		}
		if ae.Info.AttackTag != attacks.AttackTagElementalArt {
			return false
		}
		if c.StatusIsActive(skillRefundICDKey) {
			return false
		}
		c.AddStatus(skillRefundICDKey, 0.05*60, true)
		c.AddEnergy("sethos-skill", skillEnergyRegen[c.TalentLvlSkill()])
		c.c2AddStack()

		return false
	}

	c.Core.Events.Subscribe(event.OnOverload, refundCB, "sethos-e-refund")
	c.Core.Events.Subscribe(event.OnElectroCharged, refundCB, "sethos-e-refund")
	c.Core.Events.Subscribe(event.OnSuperconduct, refundCB, "sethos-e-refund")
	c.Core.Events.Subscribe(event.OnSwirlElectro, refundCB, "sethos-e-refund")
	c.Core.Events.Subscribe(event.OnCrystallizeElectro, refundCB, "sethos-e-refund")
	c.Core.Events.Subscribe(event.OnHyperbloom, refundCB, "sethos-e-refund")
	c.Core.Events.Subscribe(event.OnQuicken, refundCB, "sethos-e-refund")
	c.Core.Events.Subscribe(event.OnAggravate, refundCB, "sethos-e-refund")
}
func (c *char) Skill(p map[string]int) (action.Info, error) {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Ancient Rite: Thunderous Roar of Sand",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}

	snap := c.Snapshot(&ai)
	ap := combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 4.5)

	c.Core.QueueAttackWithSnap(ai, snap, ap, 14, c.particleCB)

	c.SetCDWithDelay(action.ActionSkill, 8*60, 10)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(skillParticleICDKey) {
		return
	}
	c.AddStatus(skillParticleICDKey, 0.05*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 2, attributes.Electro, c.ParticleDelay)
}
