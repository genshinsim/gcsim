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
var skillHitmarks = [2]int{23, 60}
var skillDropletOffsets = [][][]float64{{{-1, 3}, {0, 3.8}, {1, 3}}, {{-2, 7}, {0, 8}, {2, 7}}, {{-3, 10}, {0, 11}, {3, 10}}}
var skillDropletRandomRanges = [][][]float64{{{0.5, 1.5}, {0.5, 1.5}, {0.5, 1.5}}, {{1, 2.5}, {3.5, 4}, {1, 2.5}}, {{2, 3}, {2, 4}, {2, 3}}}

const (
	skillAlignedICD    = 10 * 60
	skillAlignedICDKey = "neuvillette-aligned-icd"

	particleCount  = 4
	particleICD    = 0.3 * 60
	particleICDKey = "neuvillette-particle-icd"
)

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
	// TODO: if target is out of range then pos should be player pos + Y: 8 offset
	skillPos := c.Core.Combat.PrimaryTarget().Pos()
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(skillPos, nil, 6),
		skillHitmarks[0], //TODO: snapshot delay?
		skillHitmarks[0],
		c.makeDropletCB(),
		c.particleCB,
	)

	aiThorn := combat.AttackInfo{
		// TODO: Apply Pneuma
		ActorIndex:         c.Index,
		Abil:               "Spiritbreath Thorn (" + c.Base.Key.Pretty() + ")",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagNone,
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypeSpear,
		Element:            attributes.Hydro,
		Durability:         0,
		Mult:               thorn[c.TalentLvlSkill()],
		CanBeDefenseHalted: true,
	}
	c.QueueCharTask(func() {
		if c.StatusIsActive(skillAlignedICDKey) {
			return
		}
		c.AddStatus(skillAlignedICDKey, skillAlignedICD, true)

		c.Core.QueueAttack(
			aiThorn,
			combat.NewCircleHitOnTarget(skillPos, nil, 4.5),
			skillHitmarks[1]-skillHitmarks[0], // TODO: snapshot delay?
			skillHitmarks[1]-skillHitmarks[0],
		)
	}, skillHitmarks[1])

	c.SetCDWithDelay(action.ActionSkill, 12*60, 20)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionCharge], // earliest cancel
		State:           action.SkillState,
	}, nil
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, particleICD, true)

	c.Core.QueueParticle(c.Base.Key.String(), particleCount, attributes.Hydro, c.ParticleDelay)
}

func (c *char) makeDropletCB() combat.AttackCBFunc {
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
			return
		}
		if done {
			return
		}
		done = true

		// determine which droplet offset and random ranges to use based on distance to first target hit
		player := c.Core.Combat.Player()
		i := 2
		if a.Target.IsWithinArea(combat.NewCircleHitOnTarget(player.Pos(), nil, 5)) {
			i = 0
		} else if a.Target.IsWithinArea(combat.NewCircleHitOnTarget(player.Pos(), nil, 10)) {
			i = 1
		}

		for j := 0; j < 3; j++ {
			common.NewSourcewaterDroplet(
				c.Core,
				geometry.CalcRandomPointFromCenter(
					geometry.CalcOffsetPoint(
						player.Pos(),
						geometry.Point{X: skillDropletOffsets[i][j][0], Y: skillDropletOffsets[i][j][1]},
						player.Direction(),
					),
					skillDropletRandomRanges[i][j][0],
					skillDropletRandomRanges[i][j][1],
					c.Core.Rand,
				),
				combat.GadgetTypSourcewaterDropletNeuv,
			)
		}
		c.Core.Combat.Log.NewEvent("Spawned 3 droplets", glog.LogCharacterEvent, c.Index).
			Write("src_action", "skill")
	}
}
