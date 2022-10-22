package avatar

import (
	"log"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
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

func (p *Player) HandleAttack(atk *combat.AttackEvent) float64 {
	//TODO: logging outside of LogHurtEvent?
	//at this point character will be hit by the attack
	p.Core.Combat.Events.Emit(event.OnPlayerHit)

	//calculate dmg
	var dmg float64

	if !atk.Info.SourceIsSim {
		if atk.Info.ActorIndex < 0 {
			log.Println(atk)
		}
		p.Core.Combat.Team.CombatByIndex(atk.Info.ActorIndex).ApplyAttackMods(atk, p)
	}

	dmg, _ = p.Attack(atk, nil)

	//delay damage event to end of the frame
	p.Core.Combat.Tasks.Add(func() {
		//apply the damage
		p.ApplyDamage(atk, dmg)
	}, 0)

	return dmg
}

func (p *Player) Attack(atk *combat.AttackEvent, evt glog.Event) (float64, bool) {
	//TODO: use this to handle self reactions in the future
	//TODO: no more external ApplySelfInfusion/ReactWithSelf calls

	//calc damage
	damagePreShield, isCrit := p.calc(atk, evt)
	damagePostShield := p.Core.Player.Shields.OnDamage(0, damagePreShield, atk.Info.Element)

	return damagePostShield, isCrit
}

func (p *Player) calc(atk *combat.AttackEvent, evt glog.Event) (float64, bool) {
	// TODO: calculate player dmg
	return 0, false
}

func (p *Player) ApplyDamage(atk *combat.AttackEvent, dmg float64) {
	di := player.DrainInfo{
		ActorIndex: p.Core.Player.ActiveChar().Index,
		Abil:       atk.Info.Abil,
		Amount:     dmg,
		External:   true,
	}
	p.Core.Player.Drain(di)
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
	p.DecayRate[mod] = dur / combat.Durability(f)
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
