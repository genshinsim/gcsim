package neuvillette

import (
	"github.com/genshinsim/gcsim/internal/common"
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

var skillFrames []int
var skillHitmarks = [2]int{23, 54}
var skillHitbox = [2]float64{6.0, 4.5}

func init() {
	skillFrames = frames.InitAbilSlice(42)
	skillFrames[action.ActionCharge] = 21
	skillFrames[action.ActionBurst] = 30
	skillFrames[action.ActionDash] = 29
	skillFrames[action.ActionJump] = 32
	skillFrames[action.ActionWalk] = 41
	skillFrames[action.ActionSwap] = 29
	// skill -> skill is unknown
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	c.chargeEarlyCancelled = false
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "O Tears, I Shall Repay",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Hydro,
		Durability: 25,
		FlatDmg:    skill[c.TalentLvlSkill()] * c.MaxHP(),
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), geometry.Point{}, skillHitbox[0]),
		skillHitmarks[0], //TODO: snapshot delay?
		skillHitmarks[0],
		c.skillcb,
	)
	// 10s Spiritbreath Thorn Interval
	if c.Core.F-c.lastThorn > 600 {
		c.lastThorn = c.Core.F
		aiThorn := combat.AttackInfo{
			// TODO: Apply Pneuma
			ActorIndex: c.Index,
			Abil:       "Spiritbreath Thorn",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Hydro,
			Durability: 0,
			Mult:       thorn[c.TalentLvlSkill()],
		}
		c.Core.QueueAttack(
			aiThorn,
			combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), geometry.Point{}, skillHitbox[1]),
			skillHitmarks[1], //TODO: snapshot delay?
			skillHitmarks[1],
		)
	}
	c.SetCDWithDelay(action.ActionSkill, 12*60, 20)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) skillcb(ac combat.AttackCB) {
	if ac.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.Core.F-c.lastSkillParticle > 18 {
		c.lastSkillParticle = c.Core.F

		circ, ok := ac.Target.Shape().(*geometry.Circle)
		if !ok {
			panic("rectangle target hurtbox is not supported for on target Sourcewater droplet spawning")
		}
		for i := 0; i < 3; i++ {
			// TODO: find the actual sourcewater droplet spawn radius for Neuv E
			pos := geometry.CalcRandomPointFromCenter(ac.Target.Pos(), circ.Radius()+1.0, circ.Radius()+4.0, c.Core.Rand)
			common.NewSourcewaterDroplet(c.Core, pos, combat.GadgetTypSourcewaterDropletNeuv)
		}
		c.Core.Combat.Log.NewEvent("Spawned 3 droplets", glog.LogCharacterEvent, c.Index).
			Write("src_action", "skill")
		c.Core.QueueParticle(c.Base.Key.String(), 4, attributes.Hydro, c.ParticleDelay)
	}
}
