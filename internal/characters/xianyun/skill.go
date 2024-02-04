package xianyun

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillLeapFrames [][]int
var skillRecastFrames []int

const (
	skillPressHitmark        = 1
	skillFirstRecastHitmark  = 41
	skillSecondRecastHitmark = 18
	leapKey                  = "xianyun-leap"
	particleICDKey           = "xianyun-particleICD-key"
)

func init() {

	skillLeapFrames = make([][]int, 3)
	// skill -> x (can only use skill, plunge or wait(?))
	skillLeapFrames[0] = frames.InitAbilSlice(41)
	skillLeapFrames[0][action.ActionHighPlunge] = 28
	skillLeapFrames[0][action.ActionSkill] = skillFirstRecastHitmark

	// skill (recast) -> x (can only use skill, plunge or wait(?))
	skillLeapFrames[1] = frames.InitAbilSlice(46)
	skillLeapFrames[1][action.ActionHighPlunge] = 10
	skillLeapFrames[1][action.ActionSkill] = skillSecondRecastHitmark

	// skill (recast) -> x (can only use skill, plunge or wait(?))
	skillLeapFrames[2] = frames.InitAbilSlice(30)
	skillLeapFrames[2][action.ActionHighPlunge] = 42

}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	// check if first leap
	if !c.StatusIsActive(eWindowKey) {
		c.eCounter = 0
	}

	if c.eCounter == 3 {
		c.eCounter = 0
	}
	//C2: After using White Clouds at Dawn, Xianyun's ATK will be increased by 20% for 15s.
	if c.Base.Cons >= 2 {
		m := make([]float64, attributes.EndStatType)
		m[attributes.ATKP] = 0.20
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("xianyun-C2", 15*60),
			AffectedStat: attributes.ATKP,
			Amount: func() ([]float64, bool) {
				if c.Core.Player.Active() == c.Index {
					return m, true
				}
				return nil, false
			},
		})
	}
	if c.eCounter == 0 {
		//Adeptal Aspect Trail DMG
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "Adeptal Aspect Trail",
			AttackTag:  attacks.AttackTagElementalArt,
			ICDTag:     attacks.ICDTagNone,
			ICDGroup:   attacks.ICDGroupDefault,
			StrikeType: attacks.StrikeTypeDefault,
			Element:    attributes.Anemo,
			Durability: 0,
			Mult:       skillPress[c.TalentLvlSkill()],
		}
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 0),
			0,
			skillPressHitmark,
		)
	}
	c.AddStatus(eWindowKey, 2*60, true)
	idx := c.eCounter
	c.eCounter++
	switch c.eCounter {
	case 1:
		c.SetCD(action.ActionSkill, 12*60)
		c.skillHeight = 3.0
		c.skillRadius = 4.0
	case 2:
		c.skillHeight = 4.0
		c.skillRadius = 5.0
	case 3:
		c.skillHeight = 5.0
		c.skillRadius = 6.5
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(skillLeapFrames[idx]),
		AnimationLength: skillLeapFrames[idx][action.InvalidAction],
		CanQueueAfter:   skillLeapFrames[idx][action.ActionDash], // earliest cancel
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
	c.AddStatus(particleICDKey, 0.2*60, true)

	count := 5.0
	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Anemo, c.ParticleDelay)
}
