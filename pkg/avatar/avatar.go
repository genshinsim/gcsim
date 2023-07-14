package avatar

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/reactable"
	"github.com/genshinsim/gcsim/pkg/target"
)

type Player struct {
	*target.Target
	*reactable.Reactable
}

func New(core *core.Core, pos geometry.Point, r float64) *Player {
	p := &Player{}
	p.Target = target.New(core, pos, r)
	p.Reactable = &reactable.Reactable{}
	p.Reactable.Init(p, core)
	return p
}

func (p *Player) Type() targets.TargettableType { return targets.TargettablePlayer }

func (p *Player) HandleAttack(atk *combat.AttackEvent) float64 {
	p.Core.Combat.Events.Emit(event.OnPlayerHit, p, atk)

	//TODO: implement player taking damage here
	return 0
}

func (p *Player) ApplySelfInfusion(ele attributes.Element, dur reactions.Durability, f int) {

	p.Core.Log.NewEventBuildMsg(glog.LogPlayerEvent, -1, "self infusion applied: "+ele.String()).
		Write("durability", dur).
		Write("duration", f)
	//we're assuming self infusion isn't subject to 0.8x multiplier
	//also no real sanity check
	if ele == attributes.Frozen {
		return
	}
	var mod reactable.ReactableModifier
	switch ele {
	case attributes.Electro:
		mod = reactable.ModifierElectro
	case attributes.Hydro:
		mod = reactable.ModifierHydro
	case attributes.Pyro:
		mod = reactable.ModifierPyro
	case attributes.Cryo:
		mod = reactable.ModifierCryo
	case attributes.Dendro:
		mod = reactable.ModifierDendro
	}

	//we're assuming refill maintains the same decay rate?
	if p.Durability[mod] > reactable.ZeroDur {
		//make sure we're not adding more than incoming
		if p.Durability[mod] < dur {
			p.Durability[mod] = dur
		}
		return
	}
	//otherwise calculate decay based on specified f (in frames)
	p.Durability[mod] = dur
	p.DecayRate[mod] = dur / reactions.Durability(f)

	p.Core.Combat.Events.Emit(event.OnSelfInfusion, ele, dur, f)
}

func (p *Player) ReactWithSelf(atk *combat.AttackEvent) {
	//check if have an element
	if p.AuraCount() == 0 {
		return
	}
	//otherwise react
	existing := p.Reactable.ActiveAuraString()
	applied := atk.Info.Durability
	p.React(atk)
	p.Core.Log.NewEvent("self reaction occured", glog.LogElementEvent, atk.Info.ActorIndex).
		Write("attack_tag", atk.Info.AttackTag).
		Write("applied_ele", atk.Info.Element.String()).
		Write("dur", applied).
		Write("abil", atk.Info.Abil).
		Write("target", 0).
		Write("existing", existing).
		Write("after", p.Reactable.ActiveAuraString())

}
