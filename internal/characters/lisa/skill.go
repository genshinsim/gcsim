package lisa

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var (
	skillPressFrames []int
	skillHoldFrames  []int
)

const (
	skillPressHitmark = 22
	skillHoldHitmark  = 117
	particleICDKey    = "lisa-particle-icd"
)

func init() {
	// skill (press) -> x
	skillPressFrames = frames.InitAbilSlice(40)
	skillPressFrames[action.ActionAttack] = 37
	skillPressFrames[action.ActionCharge] = 38
	skillPressFrames[action.ActionDash] = 35
	skillPressFrames[action.ActionJump] = 20
	skillPressFrames[action.ActionSwap] = 20

	// skill (hold) -> x
	skillHoldFrames = frames.InitAbilSlice(143)
	skillHoldFrames[action.ActionCharge] = 138
	skillHoldFrames[action.ActionBurst] = 138
	skillHoldFrames[action.ActionDash] = 116
	skillHoldFrames[action.ActionJump] = 117
	skillHoldFrames[action.ActionSwap] = 117
}

// p = 0 for no hold, p = 1 for hold
func (c *char) Skill(p map[string]int) action.Info {
	hold := p["hold"]
	if hold == 1 {
		return c.skillHold()
	}
	return c.skillPress()
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	if c.StatusIsActive(particleICDKey) {
		return
	}
	c.AddStatus(particleICDKey, 0.2*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 5, attributes.Electro, c.ParticleDelay)
}

// TODO: how long do stacks last?
func (c *char) skillPress() action.Info {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Violet Arc",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagLisaElectro,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 25,
		Mult:       skillPress[c.TalentLvlSkill()],
	}

	cb := func(a combat.AttackCB) {
		// doesn't stack off-field
		if c.Core.Player.Active() != c.Index {
			return
		}
		t, ok := a.Target.(*enemy.Enemy)
		if !ok {
			return
		}
		count := t.GetTag(conductiveTag)
		if count < 3 {
			t.SetTag(conductiveTag, count+1)
		}
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(c.Core.Combat.Player(), c.Core.Combat.PrimaryTarget(), nil, 1),
		0,
		skillPressHitmark,
		cb,
	)

	c.SetCDWithDelay(action.ActionSkill, 60, 17)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillPressFrames),
		AnimationLength: skillPressFrames[action.InvalidAction],
		CanQueueAfter:   skillPressFrames[action.ActionSwap], // earliest cancel
		State:           action.SkillState,
	}
}

// After an extended casting time, calls down lightning from the heavens, dealing massive Electro DMG to all nearby opponents.
// Deals great amounts of extra damage to opponents based on the number of Conductive stacks applied to them, and clears their Conductive status.
func (c *char) skillHold() action.Info {
	// no multiplier as that's target dependent
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Violet Arc (Hold)",
		AttackTag:  attacks.AttackTagElementalArt,
		ICDTag:     attacks.ICDTagNone,
		ICDGroup:   attacks.ICDGroupDefault,
		StrikeType: attacks.StrikeTypeDefault,
		Element:    attributes.Electro,
		Durability: 50,
	}

	// c2 add defense? no interruptions either way
	if c.Base.Cons >= 2 {
		// increase def for the duration of this abil in however many frames
		m := make([]float64, attributes.EndStatType)
		m[attributes.DEFP] = 0.25
		c.AddStatMod(character.StatMod{
			Base:         modifier.NewBase("lisa-c2", 126),
			AffectedStat: attributes.DEFP,
			Amount: func() ([]float64, bool) {
				return m, true
			},
		})
	}

	count := 0
	var c1cb func(a combat.AttackCB)
	if c.Base.Cons > 0 {
		c1cb = func(a combat.AttackCB) {
			if a.Target.Type() != targets.TargettableEnemy {
				return
			}
			if count == 5 {
				return
			}
			count++
			c.AddEnergy("lisa-c1", 2)
		}
	}

	// [8:31 PM] ArchedNosi | Lisa Unleashed: yeah 4-5 50/50 with Hold
	// [9:13 PM] ArchedNosi | Lisa Unleashed: @gimmeabreak actually wait, xd i noticed i misread my sheet, Lisa Hold E always gens 5 orbs
	enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 10), nil)
	for _, e := range enemies {
		c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(e, nil, 0.2), 0, skillHoldHitmark, c1cb, c.particleCB)
	}

	c.SetCDWithDelay(action.ActionSkill, 960, 114)

	return action.Info{
		Frames:          frames.NewAbilFunc(skillHoldFrames),
		AnimationLength: skillHoldFrames[action.InvalidAction],
		CanQueueAfter:   skillHoldFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}

func (c *char) skillHoldMult() {
	c.Core.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		atk := args[1].(*combat.AttackEvent)
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		if atk.Info.Abil != "Violet Arc (Hold)" {
			return false
		}
		stacks := t.GetTag(conductiveTag)

		atk.Info.Mult = skillHold[stacks][c.TalentLvlSkill()]

		// consume the stacks
		t.SetTag(conductiveTag, 0)

		return false
	}, "lisa-skill-hold-mul")
}
