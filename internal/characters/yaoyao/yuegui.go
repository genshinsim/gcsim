package yaoyao

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/targets"
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
	ai           combat.AttackInfo
	snap         combat.Snapshot
	aoe          combat.AttackPattern
	throwCounter int
}

func (c *char) newYueguiThrow() *yuegui {
	yg := &yuegui{
		ai:   c.skillRadishAI,
		snap: c.Snapshot(&c.skillRadishAI),
		c:    c,
	}
	player := c.Core.Combat.Player()
	pos := geometry.CalcOffsetPoint(player.Pos(), geometry.Point{Y: 2}, player.Direction())
	yg.Gadget = gadget.New(c.Core, pos, 0.5, combat.GadgetTypYueguiThrowing)

	yg.Gadget.Duration = 600
	yg.Gadget.OnThinkInterval = yg.throw

	// they start throwing 29f after being spawned
	yg.Gadget.ThinkInterval = 29

	yg.Gadget.OnKill = func() {
		yg.Core.Log.NewEvent("Yuegui (Throwing) removed", glog.LogCharacterEvent, yg.c.Index)
	}
	yg.Core.Log.NewEvent("Yuegui (Throwing) summoned", glog.LogCharacterEvent, yg.c.Index)
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
	pos := geometry.CalcOffsetPoint(player.Pos(), geometry.Point{Y: -2}, player.Direction())
	yg.Gadget = gadget.New(c.Core, pos, 0.5, combat.GadgetTypYueguiJumping)
	yg.Gadget.Duration = -1 // They last until they get deleted by the burst
	yg.Gadget.OnThinkInterval = yg.throw

	// they start throwing 29f after being spawned
	yg.Gadget.ThinkInterval = 29

	yg.Gadget.OnKill = func() {
		yg.Core.Log.NewEvent("Yuegui (Jumping) removed", glog.LogCharacterEvent, yg.c.Index)
	}
	yg.Core.Log.NewEvent("Yuegui (Jumping) summoned", glog.LogCharacterEvent, yg.c.Index)
	// yg.Reactable = &reactable.Reactable{}
	// yg.Reactable.Init(yg, c.Core)
	yg.aoe = combat.NewCircleHitOnTarget(pos, nil, skillTargetingRad)

	c.Core.Combat.AddGadget(yg)
	c.yueguiJumping[c.numYueguiJumping] = yg
	c.numYueguiJumping += 1
}

func (c *char) makeHealCB(area combat.AttackPattern, hi player.HealInfo) func(combat.AttackCB) {
	done := false
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy && a.Target.Type() != targets.TargettablePlayer {
			return
		}

		if done {
			return
		}
		if c.Core.Combat.Player().IsWithinArea(area) {
			if hi.Target != -1 {
				hi.Target = c.Core.Player.Active()
			}
			c.radishHeal(hi)
			done = true
		}
	}
}

func (yg *yuegui) Tick() {
	//this is needed since both reactable and gadget tick
	// yg.Reactable.Tick()
	yg.Gadget.Tick()
}

func (yg *yuegui) makeParticleCB() combat.AttackCBFunc {
	if yg.GadgetTyp() != combat.GadgetTypYueguiThrowing {
		return nil
	}
	return func(a combat.AttackCB) {
		if a.Target.Type() != targets.TargettableEnemy {
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
	yg.Gadget.ThinkInterval = 60
	currHPPerc := yg.Core.Player.ActiveChar().CurrentHPRatio()
	enemy := yg.Core.Combat.RandomEnemyWithinArea(yg.aoe, nil)

	var target geometry.Point
	if currHPPerc > 0.7 && enemy != nil {
		target = enemy.Pos()
	} else {
		// really it should be random if no targets are in range and the character's HP is full but we aren't really simming that
		target = yg.Core.Combat.Player().Pos()
	}
	radishExplodeAoE := combat.NewCircleHitOnTarget(target, nil, radishRad)
	radishExplodeAoE.SkipTargets[targets.TargettablePlayer] = false
	yg.c.QueueCharTask(func() {
		ai, hi := yg.getInfos()

		yg.Core.QueueAttackWithSnap(
			ai,
			yg.snap,
			radishExplodeAoE,
			1,
			yg.c.makeHealCB(radishExplodeAoE, hi),
			yg.makeParticleCB(),
			yg.c.makeC2CB(),
		)
	}, travelDelay-1)
	if yg.GadgetTyp() == combat.GadgetTypYueguiThrowing && yg.c.Base.Cons >= 6 && (yg.throwCounter == 2 || yg.throwCounter == 5) {
		yg.c6(target)
	}
	yg.throwCounter += 1
}

func (yg *yuegui) getInfos() (combat.AttackInfo, player.HealInfo) {
	var ai combat.AttackInfo
	var hi player.HealInfo

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
func (yg *yuegui) HandleAttack(atk *combat.AttackEvent) float64 {
	// yg.Core.Events.Emit(event.OnGadgetHit, yg, atk)
	// yg.Attack(atk, nil)
	return 0
}

func (yg *yuegui) Attack(*combat.AttackEvent, glog.Event) (float64, bool) { return 0, false }
func (yg *yuegui) SetDirection(trg geometry.Point)                        {}
func (yg *yuegui) SetDirectionToClosestEnemy()                            {}
func (yg *yuegui) CalcTempDirection(trg geometry.Point) geometry.Point {
	return geometry.DefaultDirection()
}
