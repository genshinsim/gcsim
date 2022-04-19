package player

import "github.com/genshinsim/gcsim/pkg/core"

//try dashing, return ErrActionNotReady if not enough stam
func (p *Player) dash(c core.Character, param map[string]int) error {
	req := p.core.StamPercentMod(core.ActionDash) * c.ActionStam(core.ActionDash, param)
	if p.Stam < req {
		p.core.Log.NewEvent("insufficient stam: dash", core.LogSimEvent, -1, "have", p.Stam)
		return ErrActionNotReady
	}
	p.core.Events.Emit(core.PreDash)
	//stam should be consumed at end of animation?
	p.core.Tasks.Add(func() {
		p.Stam -= req
		p.LastStamUse = p.core.F
		p.core.Events.Emit(core.OnStamUse, core.ActionDash)
		p.core.Events.Emit(core.PostDash)
	}, p.FramesSettings.Dash-1)

	return nil
}
