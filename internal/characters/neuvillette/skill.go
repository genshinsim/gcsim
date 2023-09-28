package neuvillette

import (
	"github.com/genshinsim/gcsim/internal/common"
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

var skillFrames []int

func init() {
	skillFrames = frames.InitAbilSlice(46)
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
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
	skillGeneratedDroplets := false
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget().Pos(), geometry.Point{}, 3),
		35, //TODO: snapshot delay?
		35,
		func(ac combat.AttackCB) {
			if !skillGeneratedDroplets {
				skillGeneratedDroplets = true

				circ, ok := ac.Target.Shape().(*geometry.Circle)
				if !ok {
					panic("rectangle target hurtbox is not supported for on target Sourcewater droplet spawning")
				}
				for i := 0; i < 3; i++ {
					// TODO: find the actual sourcewater droplet spawn radius for Neuv E
					pos := geometry.CalcRandomPointFromCenter(ac.Target.Pos(), circ.Radius()+1.0, circ.Radius()+2.0, c.Core.Rand)
					common.NewSourcewaterDroplet(c.Core, pos)
				}
			}

			// 10s Spiritbreath Thorn Interval
			if c.Core.F-c.lastThorn > 600 {
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
					FlatDmg:    thorn[c.TalentLvlSkill()] * c.MaxHP(),
				}
				c.Core.QueueAttack(
					aiThorn,
					combat.NewCircleHitOnTarget(ac.Target.Pos(), geometry.Point{}, 2),
					30, //TODO: snapshot delay?
					30,
				)
				c.lastThorn = c.Core.F
			}
			c.Core.QueueParticle(c.Base.Key.String(), 4, attributes.Hydro, c.ParticleDelay)
		},
	)

	c.SetCDWithDelay(action.ActionSkill, 12*60, 10)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}, nil
}
