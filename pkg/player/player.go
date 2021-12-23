package player

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/target"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

type Player struct {
	*target.Tmpl
}

func New(index int, c *core.Core) *Player {
	p := &Player{}
	p.Tmpl = &target.Tmpl{}
	p.Reactable = &reactable.Reactable{}
	p.TargetIndex = index
	p.Reactable.Init(p, p.Core)
	p.Tmpl.Init(0, 0, 0.5)
	p.Core = c
	return p
}

func (p *Player) Attack(atk *core.AttackEvent) (float64, bool) {
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
