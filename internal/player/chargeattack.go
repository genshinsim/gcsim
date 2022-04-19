package player

import "github.com/genshinsim/gcsim/pkg/core"

//try using charge attack, return ErrActionNotReady if not enough stam
func (p *Player) chargeattack(c core.Character, param map[string]int) error {
	req := p.core.StamPercentMod(core.ActionCharge) * c.ActionStam(core.ActionCharge, param)
	if p.Stam < req {
		p.core.Log.NewEvent("insufficient stam: charge attack", core.LogSimEvent, -1, "have", p.Stam)
		return ErrActionNotReady
	}
	p.core.Events.Emit(core.PreChargeAttack)

	//[8:09 PM] characters frame recount beggar: anyone know whether charge attack consumes energy at the end or at the beginning?
	//[8:10 PM] BowTae: should be beginning, since you can cancel catalyst CA before it comes out
	//[8:10 PM] BowTae: and stamina is still consumed
	p.Stam -= req
	p.LastStamUse = p.core.F
	p.core.Events.Emit(core.OnStamUse, core.ActionCharge)

	res := c.ChargeAttack(param)
	p.State.FrameStarted = p.core.F
	p.State.AnimationDuration = res.Frames

	//stam should be consumed at end of animation?
	p.State.OnStateEnd = p.postChargeAttack

	return nil
}

func (p *Player) postChargeAttack() {
	p.core.Events.Emit(core.PostChargeAttack)
}
