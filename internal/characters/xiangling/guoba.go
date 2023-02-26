package xiangling

import (
	"github.com/genshinsim/gcsim/internal/characters/faruzan"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/gadget"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

type panda struct {
	*gadget.Gadget
	*reactable.Reactable
	c     *char
	ai    combat.AttackInfo
	snap  combat.Snapshot
	timer int
}

func (c *char) newGuoba(ai combat.AttackInfo) *panda {
	p := &panda{
		ai:   ai,
		snap: c.Snapshot(&ai),
		c:    c,
	}
	player := c.Core.Combat.Player()
	pos := combat.CalcOffsetPoint(
		player.Pos(),
		combat.Point{Y: 1.3},
		player.Direction(),
	)
	p.Gadget = gadget.New(c.Core, pos, 0.2, combat.GadgetTypGuoba)
	p.Gadget.Duration = 438
	p.Reactable = &reactable.Reactable{}
	p.Reactable.Init(p, c.Core)

	return p
}

func (p *panda) Tick() {
	// this is needed since both reactable and gadget tick
	p.Reactable.Tick()
	p.Gadget.Tick()
	p.timer++
	// guoba pew pew every 100 frames
	// first pew pew is at 126, but guoba spawns 13 in; so it's really 113
	// then every 100 after that
	// TODO: kids.. don't do this
	switch p.timer {
	case 103, 203, 303, 403: // swirl window
		p.Core.Log.NewEvent("guoba self infusion applied", glog.LogElementEvent, p.c.Index).
			SetEnded(p.c.Core.F + infuseWindow + 1)
		p.Durability[reactable.ModifierPyro] = infuseDurability
		p.Core.Tasks.Add(func() {
			p.Durability[reactable.ModifierPyro] = 0
		}, infuseWindow+1) // +1 since infuse window is inclusive
		p.SetDirectionToClosestEnemy()
		// queue this in advance because that's how it is on live
		p.breath()
	}
}

func (p *panda) breath() {
	// assume A1
	p.Core.QueueAttackWithSnap(
		p.ai,
		p.snap,
		combat.NewCircleHitOnTargetFanAngle(p, nil, p.c.guobaFlameRange, 60),
		10,
		p.c.c1,
		p.c.particleCB,
	)
}

func (p *panda) Type() combat.TargettableType { return combat.TargettableGadget }

func (p *panda) HandleAttack(atk *combat.AttackEvent) float64 {
	p.Core.Events.Emit(event.OnGadgetHit, p, atk)
	p.Attack(atk, nil)
	return 0
}

func (p *panda) Attack(atk *combat.AttackEvent, evt glog.Event) (float64, bool) {
	if atk.Info.AttackTag != combat.AttackTagElementalArt {
		return 0, false
	}
	// check pyro window
	if p.Durability[reactable.ModifierPyro] <= reactable.ZeroDur {
		return 0, false
	}

	// don't take damage, trigger swirl reaction only on sucrose E or faruzan E
	switch p.Core.Player.Chars()[atk.Info.ActorIndex].Base.Key {
	case keys.Sucrose:
		p.Core.Log.NewEvent("guoba hit by sucrose E", glog.LogCharacterEvent, p.c.Index)
	case keys.Faruzan:
		if atk.Info.Abil != faruzan.VortexAbilName {
			return 0, false
		}
		p.Core.Log.NewEvent("guoba hit by faruzan pressurized collapse", glog.LogCharacterEvent, p.c.Index)
	default:
		return 0, false
	}

	// cheat a bit, set the durability just enough to match incoming sucrose/faruzan E gauge
	oldDur := p.Durability[reactable.ModifierPyro]
	p.Durability[reactable.ModifierPyro] = infuseDurability
	p.React(atk)
	// restore the durability after
	p.Durability[reactable.ModifierPyro] = oldDur

	return 0, false
}
