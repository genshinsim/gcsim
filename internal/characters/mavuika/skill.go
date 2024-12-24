package mavuika

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

const (
	skillHitmark                     = 16
	ringsOfSearchingRadianceInterval = 120
)

var (
	skillFrames       []int
	skillSwitchFrames []int
)

func init() {
	skillFrames = frames.InitAbilSlice(20)       // E -> Swap
	skillSwitchFrames = frames.InitAbilSlice(18) // E -> N1
}

func (c *char) Skill(p map[string]int) (action.Info, error) {
	if c.nightsoulState.HasBlessing() {
		c.flamestriderModeActive = !c.flamestriderModeActive
		if c.flamestriderModeActive {
			c.enterFlamestrider()
		} else {
			c.exitFlamestrider()
		}

		return action.Info{
			Frames:          frames.NewAbilFunc(skillSwitchFrames),
			AnimationLength: skillSwitchFrames[action.InvalidAction],
			CanQueueAfter:   skillSwitchFrames[action.ActionAttack], // change to earliest
			State:           action.SkillState,
		}, nil
	}

	c.nightsoulState.EnterBlessing(c.nightsoulState.MaxPoints)
	c.nightsoulSrc = c.Core.F
	c.QueueCharTask(c.nightsoulPointReduceFunc(c.Core.F), 12)
	hold, ok := p["hold"]
	if !ok {
		hold = 0
	}
	switch {
	case hold < 0:
		hold = 0
	case hold > 1:
		hold = 1
	}
	ai := combat.AttackInfo{
		ActorIndex:     c.Index,
		Abil:           "Skill DMG",
		AttackTag:      attacks.AttackTagElementalArt,
		AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
		ICDTag:         attacks.ICDTagNone,
		ICDGroup:       attacks.ICDGroupDefault,
		StrikeType:     attacks.StrikeTypeDefault,
		Element:        attributes.Pyro,
		Durability:     25,
		Mult:           1.2648,
	}
	c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.Player().Pos(), geometry.Point{Y: 1}, 3.5), skillHitmark, skillHitmark, c.particleCB)
	c.SetCD(action.ActionSkill, 15*60)

	c.QueueCharTask(c.ringsOfSearchingRadianceHit(c.Core.F), ringsOfSearchingRadianceInterval)
	c.QueueCharTask(c.c6FlamestriderModeHit(c.Core.F), 3*60)

	if hold == 1 {
		c.enterFlamestrider()
	} else {
		c.flamestriderModeActive = false
		c.c2AddDefMod()
	}

	return action.Info{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionBurst], // change to earliest
		State:           action.SkillState,
	}, nil
}

func (c *char) nightsoulPointReduceFunc(src int) func() {
	return func() {
		if c.nightsoulSrc != src {
			return
		}

		if !c.nightsoulState.HasBlessing() {
			return
		}
		val := 1.
		if c.flamestriderModeActive {
			val += 0.8
			if c.Core.Player.CurrentState() == action.ChargeAttackState {
				val += 0.4
			}
		}
		if !c.StatusIsActive(crucibleOfDeathAndLifeStatus) {
			c.reduceNightsoulPoints(val)
		}
		c.QueueCharTask(c.nightsoulPointReduceFunc(src), 12)
	}
}

func (c *char) reduceNightsoulPoints(val float64) {
	c.nightsoulState.ConsumePoints(val)
	// don't exit nightsoul while in NA/Plunge/Charge of Flamestride
	if c.nightsoulState.Points() <= 0.00001 {
		if !c.flamestriderModeActive {
			c.c2DeleteDefMod()
		}
		if !c.nightsoulState.HasBlessing() {
			return
		}
		c.nightsoulState.ExitBlessing()
		c.nightsoulState.ClearPoints()
		c.flamestriderModeActive = false
		c.nightsoulSrc = -1
		c.NormalHitNum = normalHitNum
		c.NormalCounter = 0
	}
}

func (c *char) ringsOfSearchingRadianceHit(src int) func() {
	return func() {
		if src != c.nightsoulSrc {
			return
		}
		if !c.nightsoulState.HasBlessing() {
			return
		}
		if !c.flamestriderModeActive {
			ai := combat.AttackInfo{
				ActorIndex:     c.Index,
				Abil:           "Rings of Searing Radiance DMG",
				AttackTag:      attacks.AttackTagElementalArt,
				AdditionalTags: []attacks.AdditionalTag{attacks.AdditionalTagNightsoul},
				ICDTag:         attacks.ICDTagNone,
				ICDGroup:       attacks.ICDGroupDefault,
				StrikeType:     attacks.StrikeTypeDefault,
				Element:        attributes.Pyro,
				Durability:     25,
				Mult:           2.176,
			}
			// TODO: change hurtbox
			c.Core.QueueAttack(ai, combat.NewCircleHitOnTarget(c.Core.Combat.PrimaryTarget(), nil, 3.5), 0, 0, c.c6RSRModeHitCB())
			// a hit of E comsumes 3 NS points
			c.nightsoulState.ConsumePoints(3)
		}
		c.QueueCharTask(c.ringsOfSearchingRadianceHit(src), ringsOfSearchingRadianceInterval)
	}
}

func (c *char) enterFlamestrider() {
	c.Core.Log.NewEvent("switching to Flamestrider state", glog.LogCharacterEvent, c.Index)
	c.flamestriderModeActive = true
	c.NormalHitNum = bikeHitNum
	c.c2DeleteDefMod()
}

func (c *char) exitFlamestrider() {
	c.Core.Log.NewEvent("switching to Rings of Searing Flames state", glog.LogCharacterEvent, c.Index)
	c.flamestriderModeActive = false
	c.NormalHitNum = normalHitNum
	c.c2AddDefMod()
}

func (c *char) particleCB(a combat.AttackCB) {
	if a.Target.Type() != targets.TargettableEnemy {
		return
	}
	c.Core.QueueParticle(c.Base.Key.String(), 5, attributes.Pyro, c.ParticleDelay)
}
