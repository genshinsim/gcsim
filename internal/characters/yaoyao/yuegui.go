package yaoyao

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

const skillParticleICD = 90
const skillTargetingRad = 8 // TODO: replace with the actual range
const radishRad = 1         // TODO: replace with the actual radish AoE
const travelDelay = 10      // TODO: replace with the actual travel delay

type yuegui struct {
	*gadget.Gadget
	// *reactable.Reactable
	c    *char
	ai   combat.AttackInfo
	snap combat.Snapshot
	aoe  combat.AttackPattern
}

func (c *char) newYueguiThrow(procAI combat.AttackInfo) *yuegui {

	yg := &yuegui{
		ai:   procAI,
		snap: c.Snapshot(&procAI),
		c:    c,
	}
	pos := c.Core.Combat.Player().Pos().Add(combat.Point{0, 1})
	//TODO: yuegui placement??
	yg.Gadget = gadget.New(c.Core, pos, 0.5, combat.GadgetTypYueguiThrowing)
	yg.Gadget.Duration = 600
	yg.Gadget.OnThinkInterval = yg.throw
	yg.Gadget.ThinkInterval = 60
	// yg.Reactable = &reactable.Reactable{}
	// yg.Reactable.Init(yg, c.Core)
	yg.aoe = combat.NewCircleHitOnTarget(pos, nil, 7)

	return yg
}

func (c *char) newYueguiJump() {
	if !c.StatusIsActive(burstKey) || c.numYueguiJumping >= 3 {
		return
	}
	yg := &yuegui{
		snap: c.Snapshot(&c.burstAI),
		c:    c,
	}
	pos := c.Core.Combat.Player().Pos()
	//TODO: yuegui placement??
	yg.Gadget = gadget.New(c.Core, pos, 0.5, combat.GadgetTypYueguiJumping)
	yg.Gadget.Duration = -1 // They last until they get deleted by the burst
	yg.Gadget.OnThinkInterval = yg.throw
	yg.Gadget.ThinkInterval = 60
	// yg.Reactable = &reactable.Reactable{}
	// yg.Reactable.Init(yg, c.Core)
	yg.aoe = combat.NewCircleHitOnTarget(pos, nil, 7)

	c.yueguiJumping[c.numYueguiJumping] = yg
	c.numYueguiJumping += 1
}

func (yg *yuegui) Tick() {
	//this is needed since both reactable and gadget tick
	// yg.Reactable.Tick()
	yg.Gadget.Tick()
}

func (yg *yuegui) throw() {
	particleCB := func(_ combat.AttackCB) {
		if yg.Core.F-yg.c.lastSkillParticle < skillParticleICD {
			return
		}
		yg.c.lastSkillParticle = yg.Core.F
		yg.Core.QueueParticle("yaoyao", 1, attributes.Pyro, yg.c.ParticleDelay)
	}
	currHPPerc := yg.Core.Player.ActiveChar().HPCurrent / yg.Core.Player.ActiveChar().MaxHP()
	enemy := yg.Core.Combat.RandomEnemyWithinArea(yg.aoe, nil)

	var target combat.Point
	if currHPPerc > 0.7 && enemy != nil {
		target = enemy.Pos()
	} else {
		// really it should be random if no targets are in range and the character's HP is full but we aren't really simming that
		target = yg.Core.Combat.Player().Pos()
	}
	radishExplodeAoE := combat.NewCircleHitOnTarget(target, nil, radishRad)
	var ai combat.AttackInfo
	var hi player.HealInfo
	if yg.c.StatusIsActive(burstKey) {
		ai = yg.c.burstAI
		hi = yg.c.getBurstHealInfo()
	} else {
		ai = yg.ai
		hi = yg.getSkillHealInfo()
	}
	yg.Core.QueueAttackWithSnap(
		ai,
		yg.snap,
		radishExplodeAoE,
		travelDelay,
		particleCB,
	)
	if yg.Core.Combat.Player().IsWithinArea(radishExplodeAoE) {
		yg.Core.Player.Heal(hi)
	}

}
func (yg *yuegui) getSkillHealInfo() player.HealInfo {
	return player.HealInfo{
		Caller:  yg.c.Index,
		Target:  yg.Core.Player.Active(),
		Message: "Yuegui skill",
		// Src:     skillHeal[c.TalentLvlSkill()],
		Bonus: yg.snap.Stats[attributes.Heal],
	}
}

func (yg *yuegui) Type() combat.TargettableType { return combat.TargettableGadget }

// TODO: Confirm if yueguis can infuse cryo
func (yg *yuegui) HandleAttack(atk *combat.AttackEvent) float64 {
	// yg.Core.Events.Emit(event.OnGadgetHit, yg, atk)
	// yg.Attack(atk, nil)
	return 0
}

func (yg *yuegui) Attack(atk *combat.AttackEvent, evt glog.Event) (float64, bool) {
	return 0, false
}

func (yg *yuegui) ApplyDamage(*combat.AttackEvent, float64) {}
