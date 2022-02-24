package player

import (
	"github.com/genshinsim/gcsim/internal/reactable"
	"github.com/genshinsim/gcsim/internal/tmpl/target"
	"github.com/genshinsim/gcsim/pkg/core"
)

type Player struct {
	*target.Tmpl
}

func New(index int, c *core.Core) *Player {
	p := &Player{}
	p.Tmpl = &target.Tmpl{}
	p.Reactable = &reactable.Reactable{}
	p.TargetIndex = index
	p.Reactable.Init(p, c)
	p.Tmpl.Init(0, 0, 0.5)
	p.Core = c
	return p
}

func (p *Player) Attack(atk *core.AttackEvent, evt core.LogEvent) (float64, bool) {
	//ignore attacks we don't get hit
	return 0, false
}

func (p *Player) Type() core.TargettableType                 { return core.TargettablePlayer }
func (p *Player) MaxHP() float64                             { return 1 }
func (p *Player) HP() float64                                { return 1 }
func (p *Player) Shape() core.Shape                          { return &p.Hitbox }
func (p *Player) AddDefMod(key string, val float64, dur int) {}
func (p *Player) AddResMod(key string, val core.ResistMod)   {}
func (p *Player) RemoveResMod(key string)                    {}
func (p *Player) RemoveDefMod(key string)                    {}
func (p *Player) HasDefMod(key string) bool                  { return false }
func (p *Player) HasResMod(key string) bool                  { return false }
func (p *Player) AddReactBonusMod(mod core.ReactionBonusMod) {}
func (p *Player) ReactBonus(atk core.AttackInfo) float64     { return 0 }
func (p *Player) Kill()                                      {}

func (p *Player) ApplySelfInfusion(ele core.EleType, dur core.Durability, f int) {

	p.Core.Log.NewEventBuildMsg(core.LogSimEvent, -1, "self infusion applied: "+ele.String()).Write("durability", dur, "duration", f)
	//we're assuming self infusion isn't subject to 0.8x multiplier
	//also no real sanity check
	if ele == core.Frozen {
		return
	}

	//we're assuming refill maintains the same decay rate?
	if p.Durability[ele] > reactable.ZeroDur {
		//make sure we're not adding more than incoming
		if p.Durability[ele] < dur {
			if p.Durability[ele]+dur > dur {
				dur = dur - p.Durability[ele]
			}
			p.Durability[ele] += dur
		}
		return
	}
	//otherwise calculate decay based on specified f (in frames)
	p.Durability[ele] = dur
	p.DecayRate[ele] = dur / core.Durability(f)
}
