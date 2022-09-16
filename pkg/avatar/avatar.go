package avatar

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/reactable"
	"github.com/genshinsim/gcsim/pkg/target"
)

type Player struct {
	*target.Target
	*reactable.Reactable
}

func New(core *core.Core, x, y, r float64) *Player {
	p := &Player{}
	p.Target = target.New(core, x, y, r)
	p.Reactable = &reactable.Reactable{}
	p.Reactable.Init(p, core)
	return p
}

func (p *Player) Type() combat.TargettableType { return combat.TargettablePlayer }

func (p *Player) Attack(ae *combat.AttackEvent, evt glog.Event) (float64, bool) {
	//TODO: consider using this to implement additional self reactions
	return 0, false
}

func (p *Player) ApplySelfInfusion(ele attributes.Element, dur combat.Durability, f int) {

	p.Core.Log.NewEventBuildMsg(glog.LogPlayerEvent, -1, "self infusion applied: "+ele.String()).
		Write("durability", dur).
		Write("duration", f)
	//we're assuming self infusion isn't subject to 0.8x multiplier
	//also no real sanity check
	if ele == attributes.Frozen {
		return
	}

	//we're assuming refill maintains the same decay rate?
	if p.Durability[ele] > reactable.ZeroDur {
		//make sure we're not adding more than incoming
		if p.Durability[ele] < dur {
			p.Durability[ele] = dur
		}
		return
	}
	//otherwise calculate decay based on specified f (in frames)
	p.Durability[ele] = dur
	p.DecayRate[ele] = dur / combat.Durability(f)
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
