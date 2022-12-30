package yaoyao

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/gadget"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

const skillParticleICD = 90
const skillTargetingRad = 8 // TODO: replace with the actual range
const radishRad = 1         // TODO: replace with the actual radish AoE
const travelDelay = 10      // TODO: replace with the actual travel delay

type yuegui struct {
	*gadget.Gadget
	*reactable.Reactable
	c    *char
	ai   combat.AttackInfo
	snap combat.Snapshot
}

func (c *char) newYuegui(procAI combat.AttackInfo) *yuegui {
	yg := &yuegui{
		ai:   procAI,
		snap: c.Snapshot(&procAI),
		c:    c,
	}
	x, y := c.Core.Combat.Player().Pos()
	//TODO: yuegui placement??
	yg.Gadget = gadget.New(c.Core, core.Coord{X: x, Y: y, R: 0.2}, combat.GadgetTypYueguiThrowing)
	yg.Gadget.Duration = 600
	yg.Reactable = &reactable.Reactable{}
	yg.Reactable.Init(yg, c.Core)

	return yg
}

func (yg *yuegui) Tick() {
	//this is needed since both reactable and gadget tick
	yg.Reactable.Tick()
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
	if currHPPerc > 0.7 {
		x, y := yg.Gadget.Pos()
		enemies := yg.Core.Combat.EnemiesWithinRadius(x, y, skillTargetingRad)
		if len(enemies) > 0 {
			idx := yg.Core.Rand.Intn(len(enemies))

			yg.Core.QueueAttackWithSnap(
				yg.ai,
				yg.snap,
				combat.NewCircleHit(yg.Core.Combat.Enemy(enemies[idx]), radishRad),
				travelDelay,
				particleCB,
			)
		}
	} else {
		yg.Core.QueueAttackWithSnap(
			yg.ai,
			yg.snap,
			combat.NewCircleHit(yg.Core.Combat.Player(), radishRad),
			travelDelay,
			particleCB,
		)
	}

}

func (yg *yuegui) Type() combat.TargettableType { return combat.TargettableGadget }

// TODO: Confirm if yueguis can infuse cryo
func (yg *yuegui) HandleAttack(atk *combat.AttackEvent) float64 {
	yg.Core.Events.Emit(event.OnGadgetHit, yg, atk)
	yg.Attack(atk, nil)
	return 0
}

func (yg *yuegui) Attack(atk *combat.AttackEvent, evt glog.Event) (float64, bool) {
	return 0, false
}

func (yg *yuegui) ApplyDamage(*combat.AttackEvent, float64) {}
