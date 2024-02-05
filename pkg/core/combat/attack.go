package combat

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

// attack returns true if the attack lands
func (h *Handler) attack(t Target, a *AttackEvent) (float64, bool) {
	willHit, reason := t.AttackWillLand(a.Pattern)
	if !willHit {
		// Move target logs into the "Sim" event log to avoid cluttering main display for stuff like Guoba
		// And obvious things like "Fischl A4 is single target so it didn't hit targets 2-4"
		// TODO: Maybe want to add a separate set of log events for this?
		if h.Debug && t.Type() != targets.TargettablePlayer {
			h.Log.NewEventBuildMsg(glog.LogDebugEvent, a.Info.ActorIndex, "skipped ", a.Info.Abil, " ", reason).
				Write("attack_tag", a.Info.AttackTag).
				Write("applied_ele", a.Info.Element).
				Write("dur", a.Info.Durability).
				Write("target", t.Key()).
				Write("geometry.Shape", a.Pattern.Shape.String())
		}
		return 0, false
	}
	// make a copy first
	cpy := *a
	dmg := t.HandleAttack(&cpy)
	return dmg, true
}

func (h *Handler) ApplyAttack(a *AttackEvent) float64 {
	h.Events.Emit(event.OnApplyAttack, a)
	// died := false
	var total float64
	var landed bool
	// check player
	if !a.Pattern.SkipTargets[targets.TargettablePlayer] {
		//TODO: we don't check for landed here since attack that hit player should never generate hitlag?
		h.attack(h.player, a)
	}
	// check enemies
	if !a.Pattern.SkipTargets[targets.TargettableEnemy] {
		for _, v := range h.enemies {
			if v == nil {
				continue
			}
			if !v.IsAlive() {
				continue
			}
			a, l := h.attack(v, a)
			total += a
			if l {
				landed = true
			}
		}
	}
	// check gadgets
	if !a.Pattern.SkipTargets[targets.TargettableGadget] {
		for i := 0; i < len(h.gadgets); i++ {
			// sanity check here; possible gadgets died and have not been cleaned up yet
			if h.gadgets[i] == nil {
				continue
			}
			h.attack(h.gadgets[i], a)
		}
	}
	// add hitlag to actor but ignore if this is deployable
	if h.EnableHitlag && landed && !a.Info.IsDeployable {
		dur := a.Info.HitlagHaltFrames
		if h.DefHalt && a.Info.CanBeDefenseHalted {
			dur += 3.6 // 0.06
		}
		if dur > 0 {
			h.Team.ApplyHitlag(a.Info.ActorIndex, a.Info.HitlagFactor, dur)
			if h.Debug {
				h.Log.NewEvent(fmt.Sprintf("%v applying hitlag: %.3f", a.Info.Abil, dur), glog.LogHitlagEvent, a.Info.ActorIndex).
					Write("duration", dur).
					Write("factor", a.Info.HitlagFactor)
			}
		}
	}
	return total
}
