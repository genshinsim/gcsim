package xiangling

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/target"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

type panda struct {
	*target.Tmpl
	pyroWindowStart int
	pyroWindowEnd   int
}

func newGuoba(c *core.Core) *panda {
	p := &panda{}
	p.Tmpl = &target.Tmpl{}
	p.Reactable = &reactable.Reactable{}
	p.Reactable.Init(p, c)
	p.Tmpl.Init(0, 0, 0.5)
	p.Core = c
	return p
}

func (p *panda) Attack(atk *core.AttackEvent) (float64, bool) {
	//don't take damage, trigger swirl reaction only on sucrose E
	if p.Core.Chars[atk.Info.ActorIndex].Key() != keys.Sucrose {
		return 0, false
	}
	if atk.Info.AttackTag != core.AttackTagElementalArt {
		return 0, false
	}
	//check pyro window
	if p.Core.F < p.pyroWindowStart || p.Core.F > p.pyroWindowEnd {
		return 0, false
	}

	//cheat a bit, set the durability just enough to match incoming sucrose E gauge
	p.Durability[core.Pyro] = 25
	p.React(atk)
	//wipe out the durability after
	p.Durability[core.Pyro] = 0

	return 0, false
}

func (p *panda) Type() core.TargettableType                 { return core.TargettableObject }
func (p *panda) MaxHP() float64                             { return 1 }
func (p *panda) HP() float64                                { return 1 }
func (p *panda) Shape() core.Shape                          { return &p.Hitbox }
func (p *panda) AddDefMod(key string, val float64, dur int) {}
func (p *panda) AddResMod(key string, val core.ResistMod)   {}
func (p *panda) RemoveResMod(key string)                    {}
func (p *panda) RemoveDefMod(key string)                    {}
func (p *panda) HasDefMod(key string) bool                  { return false }
func (p *panda) HasResMod(key string) bool                  { return false }
func (p *panda) AddReactBonusMod(mod core.ReactionBonusMod) {}
func (p *panda) ReactBonus(atk core.AttackInfo) float64     { return 0 }
func (p *panda) Kill()                                      {}
func (p *panda) SetTag(key string, val int)                 {}
func (p *panda) GetTag(key string) int                      { return 0 }
func (p *panda) RemoveTag(key string)                       {}
