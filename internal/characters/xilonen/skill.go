package xilonen

import (
	"fmt"
	"math"

	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var skillFrames []int

const (
	skillHitmarks   = 6
	samplerInterval = 0.3 * 60

	skilRecastCD     = "xilonen-e-recast-cd"
	skillMaxDurKey   = "xilonen-e-limit"
	particleICDKey   = "xilonen-particle-icd"
	samplerShredKey  = "xilonen-e-shred"
	activeSamplerKey = "xilonen-samplers-activated"
)

func init() {
	skillFrames = frames.InitAbilSlice(20)
	skillFrames[action.ActionAttack] = 19
	skillFrames[action.ActionDash] = 15
	skillFrames[action.ActionJump] = 15
	skillFrames[action.ActionSwap] = 19
}

func (c *char) reduceNightsoulPoints(val float64) {
	if c.StatusIsActive(c6key) {
		return
	}

	c.nightsoulState.ConsumePoints(val * c.nightsoulConsumptionMul())

	// don't exit nightsoul while in NA/Plunge
	switch c.Core.Player.CurrentState() {
	case action.NormalAttackState, action.PlungeAttackState:
		return
	}

	if c.nightsoulState.Points() < 0.001 {
		c.exitNightsoul()
	}
}

func (c *char) canUseNightsoul() bool {
	return c.nightsoulState.Points() >= 0.001 || c.StatusIsActive(c6key)
}

func (c *char) enterNightsoul() {
	c.nightsoulState.EnterBlessing(45)
	c.nightsoulSrc = c.Core.F
	c.nightsoulPointReduceTask(c.nightsoulSrc)
	c.NormalHitNum = rollerHitNum
	c.NormalCounter = 0

	duration := int(9 * 60 * c.nightsoulDurationMul())
	c.setNightsoulExitTimer(duration)
	c.skillLastStamF = c.Core.Player.LastStamUse
	c.Core.Player.LastStamUse = math.MaxInt
	// Don't queue the task if C2 or higher
	if c.Base.Cons < 2 {
		c.activeGeoSampler(c.nightsoulSrc)()
	}
}

func (c *char) exitNightsoul() {
	if !c.nightsoulState.HasBlessing() {
		return
	}
	c.nightsoulState.ExitBlessing()
	c.nightsoulState.ClearPoints()
	c.nightsoulSrc = -1
	c.exitStateSrc = -1
	c.SetCD(action.ActionSkill, 7*60)
	c.NormalHitNum = normalHitNum
	c.NormalCounter = 0
	c.Core.Player.LastStamUse = c.skillLastStamF
	c.DeleteStatus(c6key)
}

func (c *char) nightsoulPointReduceTask(src int) {
	const tickInterval = .1
	c.QueueCharTask(func() {
		if c.nightsoulSrc != src {
			return
		}
		// reduce 0.5 point every 6f, which is 5 per second
		c.reduceNightsoulPoints(0.5)
	}, 60*tickInterval)
}

func (c *char) applySamplerShred(ele attributes.Element, enemies []combat.Enemy) {
	for _, e := range enemies {
		e.AddResistMod(combat.ResistMod{
			Base:  modifier.NewBaseWithHitlag(fmt.Sprintf("%v-%v", samplerShredKey, ele.String()), 60),
			Ele:   ele,
			Value: -skillShred[c.TalentLvlSkill()],
		})
	}
}

func (c *char) activeGeoSampler(src int) func() {
	return func() {
		if c.Base.Cons < 2 {
			if c.nightsoulSrc != src {
				return
			}
			if !c.nightsoulState.HasBlessing() {
				return
			}
			if c.StatusIsActive(activeSamplerKey) {
				// move to activeSamplers
				return
			}
		}
		enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 10), nil)
		c.applySamplerShred(attributes.Geo, enemies)
		c.QueueCharTask(c.activeGeoSampler(src), samplerInterval)
	}
}

func (c *char) activeSamplers(src int) func() {
	return func() {
		if c.sampleSrc != src {
			return
		}
		if !c.StatusIsActive(activeSamplerKey) {
			return
		}

		enemies := c.Core.Combat.EnemiesWithinArea(combat.NewCircleHitOnTarget(c.Core.Combat.Player(), nil, 10), nil)
		for ele := range c.shredElements {
			// skip geo when C2 or above since it's always active
			if ele == attributes.Geo && c.Base.Cons >= 2 {
				continue
			}
			c.applySamplerShred(ele, enemies)
		}

		// QueueCharTask needs to be called on the active char
		active := c.Core.Player.ActiveChar()
		active.QueueCharTask(c.activeSamplers(src), samplerInterval)
	}
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() { // don't use canUseNightsoul
		c.exitNightsoul()
		return action.Info{
			Frames:          func(_ action.Action) int { return 1 },
			AnimationLength: 1,
			CanQueueAfter:   1, // earliest cancel
			State:           action.SkillState,
		}, nil
	}

	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Yohual's Scratch",
		AttackTag:          attacks.AttackTagElementalArt,
		ICDTag:             attacks.ICDTagElementalArt,
		AdditionalTags:     []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDGroup:           attacks.ICDGroupDefault,
		StrikeType:         attacks.StrikeTypePierce,
		Element:            attributes.Geo,
		Durability:         25,
		HitlagFactor:       0.01,
		Mult:               skillDMG[c.TalentLvlSkill()],
		UseDef:             true,
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}
	ap := combat.NewCircleHitOnTarget(
		c.Core.Combat.Player(),
		geometry.Point{Y: 1.0},
		0.8,
	)
	c.Core.QueueAttack(ai, ap, skillHitmarks, skillHitmarks, c.particleCB)
	c.AddStatus(skilRecastCD, 60, true)

	if c.Core.Player.Stam >= 15 {
		c.Core.Player.RestoreStam(5)
	} else {
		// align to 15
		c.Core.Player.RestoreStam(15 - c.Core.Player.Stam)
	}

	c.enterNightsoul()
	c.c4()

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
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
	c.AddStatus(particleICDKey, 0.5*60, true)
	c.Core.QueueParticle(c.Base.Key.String(), 4, attributes.Geo, c.ParticleDelay)
}

func (c *char) setNightsoulExitTimer(duration int) {
	c.exitStateSrc = c.Core.F
	src := c.exitStateSrc
	c.QueueCharTask(func() {
		if c.exitStateSrc != src {
			return
		}
		c.nightsoulState.ClearPoints()
		if !c.canUseNightsoul() {
			// don't exit nightsoul while in NA/Plunge
			switch c.Core.Player.CurrentState() {
			case action.NormalAttackState, action.PlungeAttackState:
				return
			}
			c.exitNightsoul()
		}
	}, duration)
	c.AddStatus(skillMaxDurKey, duration, true)
}
