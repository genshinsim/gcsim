package travelergeo

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/construct"
)

var skillFrames []int

// isn't exactly hitmark
const skillHitmark = 34

func init() {
	skillFrames = frames.InitAbilSlice(24)
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Starfell Sword",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagElementalArt,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeBlunt,
		Element:    attributes.Geo,
		Durability: 50,
		Mult:       skill[c.TalentLvlSkill()],
	}
	// TODO: check snapshot timing
	c.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), 24, skillHitmark)

	var count float64 = 3
	if c.Core.Rand.Float64() < 0.33 {
		count = 4
	}
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Geo, skillHitmark+c.Core.Flags.ParticleDelay)

	c.Core.Tasks.Add(func() {
		dur := 30 * 60
		if c.Base.Cons >= 6 {
			dur += 600
		}
		c.Core.Constructs.New(c.newStone(dur), false)
	}, skillHitmark)

	c.SetCD(action.ActionSkill, 360)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

type stone struct {
	src    int
	expiry int
	char   *char
}

func (c *char) newStone(dur int) *stone {
	return &stone{
		src:    c.Core.F,
		expiry: c.Core.F + dur,
		char:   c,
	}
}

func (s *stone) OnDestruct() {
	if s.char.Base.Cons >= 2 {
		ai := combat.AttackInfo{
			ActorIndex: s.char.Index,
			Abil:       "Rockcore Meltdown",
			AttackTag:  combat.AttackTagElementalArt,
			ICDTag:     combat.ICDTagElementalArt,
			ICDGroup:   combat.ICDGroupDefault,
			StrikeType: combat.StrikeTypeBlunt,
			Element:    attributes.Geo,
			Durability: 50,
			Mult:       skill[s.char.TalentLvlSkill()],
		}
		s.char.Core.QueueAttack(ai, combat.NewDefCircHit(2, false, combat.TargettableEnemy), 0, 0)
	}
}

func (s *stone) Key() int                         { return s.src }
func (s *stone) Type() construct.GeoConstructType { return construct.GeoConstructTravellerSkill }
func (s *stone) Expiry() int                      { return s.expiry }
func (s *stone) IsLimited() bool                  { return true }
func (s *stone) Count() int                       { return 1 }
