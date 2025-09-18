package yaoyao

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

const (
	skillParticleICD  = "skill-particle-icd"
	skillTargetingRad = 8
	radishRad         = 2.0
	travelDelay       = 13
	c6TravelDelay     = 20
)

type yuegui struct {
	*gadget.Gadget
	// *reactable.Reactable
	c            *char
	ai           info.AttackInfo
	snap         info.Snapshot
	aoe          info.AttackPattern
	throwCounter int
}

func (c *char) newYueguiThrow() *yuegui {
	yg := &yuegui{
		ai:   c.skillRadishAI,
		snap: c.Snapshot(&c.skillRadishAI),
		c:    c,
	}
	player := c.Core.Combat.Player()
	pos := info.CalcOffsetPoint(player.Pos(), info.Point{Y: 2}, player.Direction())
	yg.Gadget = gadget.New(c.Core, pos, 0.5, info.GadgetTypYueguiThrowing)

	yg.Duration = 600
	yg.OnThinkInterval = yg.throw

	// they start throwing 29f after being spawned
	yg.ThinkInterval = 29

	yg.OnKill = func() {
		yg.Core.Log.NewEvent("Yuegui (Throwing) removed", glog.LogCharacterEvent, yg.c.Index())
	}
	yg.Core.Log.NewEvent("Yuegui (Throwing) summoned", glog.LogCharacterEvent, yg.c.Index())
	// yg.Reactable = &reactable.Reactable{}
	// yg.Reactable.Init(yg, c.Core)
	yg.aoe = combat.NewCircleHitOnTarget(pos, nil, skillTargetingRad)

	return yg
}

func (c *char) newYueguiJump() {
	if !c.StatusIsActive(burstKey) || c.numYueguiJumping >= 3 {
		return
	}
	yg := &yuegui{
		ai:   c.burstRadishAI,
		snap: c.Snapshot(&c.burstRadishAI),
		c:    c,
	}
	player := c.Core.Combat.Player()
	pos := info.CalcOffsetPoint(player.Pos(), info.Point{Y: -2}, player.Direction())
	yg.Gadget = gadget.New(c.Core, pos, 0.5, info.GadgetTypYueguiJumping)
	yg.Duration = -1 // They last until they get deleted by the burst
	yg.OnThinkInterval = yg.throw

	// they start throwing 29f after being spawned
	yg.ThinkInterval = 29

	yg.OnKill = func() {
		yg.Core.Log.NewEvent("Yuegui (Jumping) removed", glog.LogCharacterEvent, yg.c.Index())
	}
	yg.Core.Log.NewEvent("Yuegui (Jumping) summoned", glog.LogCharacterEvent, yg.c.Index())
	// yg.Reactable = &reactable.Reactable{}
	// yg.Reactable.Init(yg, c.Core)
	yg.aoe = combat.NewCircleHitOnTarget(pos, nil, skillTargetingRad)

	c.Core.Combat.AddGadget(yg)
	c.yueguiJumping[c.numYueguiJumping] = yg
	c.numYueguiJumping += 1
}

func (c *char) heal(area info.AttackPattern, hi info.HealInfo) func() {
	return func() {
		if !c.Core.Combat.Player().IsWithinArea(area) {
			return
		}
		if hi.Target != -1 {
			hi.Target = c.Core.Player.Active()
		}
		c.radishHeal(hi)
	}
}

func (yg *yuegui) Tick() {
	// this is needed since both reactable and gadget tick
	// yg.Reactable.Tick()
	yg.Gadget.Tick()
}

func (yg *yuegui) makeParticleCB() info.AttackCBFunc {
	if yg.GadgetTyp() != info.GadgetTypYueguiThrowing {
		return nil
	}
	return func(a info.AttackCB) {
		if a.Target.Type() != info.TargettableEnemy {
			return
		}

		if yg.c.StatusIsActive(skillParticleICD) {
			return
		}
		yg.c.AddStatus(skillParticleICD, 1.5*60, true)
		yg.Core.QueueParticle(yg.c.Base.Key.String(), 1, attributes.Dendro, yg.c.ParticleDelay)
	}
}

func (yg *yuegui) throw() {
	yg.ThinkInterval = 60
	currHPPerc := yg.Core.Player.ActiveChar().CurrentHPRatio()
	enemy := yg.Core.Combat.RandomEnemyWithinArea(yg.aoe, nil)

	var target info.Point
	if currHPPerc > 0.7 && enemy != nil {
		target = enemy.Pos()
	} else {
		// really it should be random if no targets are in range and the character's HP is full but we aren't really simming that
		target = yg.Core.Combat.Player().Pos()
	}
	radishExplodeAoE := combat.NewCircleHitOnTarget(target, nil, radishRad)
	yg.c.QueueCharTask(func() {
		ai, hi := yg.getInfos()

		delay := 1
		yg.Core.Tasks.Add(yg.c.heal(radishExplodeAoE, hi), delay)
		yg.Core.QueueAttackWithSnap(
			ai,
			yg.snap,
			radishExplodeAoE,
			delay,
			yg.makeParticleCB(),
			yg.c.makeC2CB(),
		)
	}, travelDelay-1)
	if yg.GadgetTyp() == info.GadgetTypYueguiThrowing && yg.c.Base.Cons >= 6 && (yg.throwCounter == 2 || yg.throwCounter == 5) {
		yg.c6(target)
	}
	yg.throwCounter += 1
}

func (yg *yuegui) getInfos() (info.AttackInfo, info.HealInfo) {
	var ai info.AttackInfo
	var hi info.HealInfo

	if yg.c.StatusIsActive(burstKey) {
		ai = yg.c.burstRadishAI
		hi = yg.c.getBurstHealInfo(&yg.snap)
	} else {
		ai = yg.ai
		hi = yg.c.getSkillHealInfo(&yg.snap)
	}
	return ai, hi
}

// TODO: Confirm if yueguis can infuse cryo
func (yg *yuegui) HandleAttack(atk *info.AttackEvent) float64 {
	// yg.Core.Events.Emit(event.OnGadgetHit, yg, atk)
	// yg.Attack(atk, nil)
	return 0
}

func (yg *yuegui) Attack(*info.AttackEvent, glog.Event) (float64, bool) { return 0, false }
func (yg *yuegui) SetDirection(trg info.Point)                          {}
func (yg *yuegui) SetDirectionToClosestEnemy()                          {}
func (yg *yuegui) CalcTempDirection(trg info.Point) info.Point {
	return info.DefaultDirection()
}
