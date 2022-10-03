package combat

import (
	"fmt"
	"log"

	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

// attack returns true if the attack lands
func (h *Handler) attack(t Target, a *AttackEvent) (float64, bool) {
	willHit, reason := t.AttackWillLand(a.Pattern, a.Info.DamageSrc)
	if !willHit {
		// Move target logs into the "Sim" event log to avoid cluttering main display for stuff like Guoba
		// And obvious things like "Fischl A4 is single target so it didn't hit targets 2-4"
		// TODO: Maybe want to add a separate set of log events for this?
		if t.Type() != TargettablePlayer && h.Debug {
			h.Log.NewEventBuildMsg(glog.LogDebugEvent, a.Info.ActorIndex, "skipped ", a.Info.Abil, " ", reason).
				Write("attack_tag", a.Info.AttackTag).
				Write("applied_ele", a.Info.Element).
				Write("dur", a.Info.Durability).
				Write("target", t.Index()).
				Write("shape", a.Pattern.Shape.String())

		}
		return 0, false
	}

	//make a copy first
	cpy := *a

	//at this point attack will land
	h.Events.Emit(event.OnAttackWillLand, t, &cpy)

	//check to make sure it's not cancelled for w/e reason
	if a.Cancelled {
		return 0, false
	}

	var amp string
	var cata string
	var dmg float64
	var crit bool

	evt := h.Log.NewEvent(cpy.Info.Abil, glog.LogDamageEvent, cpy.Info.ActorIndex).
		Write("target", t.Index()).
		Write("attack-tag", cpy.Info.AttackTag).
		Write("ele", cpy.Info.Element.String()).
		Write("damage", &dmg).
		Write("crit", &crit).
		Write("amp", &amp).
		Write("cata", &cata).
		Write("abil", cpy.Info.Abil).
		Write("source_frame", cpy.SourceFrame)
	evt.WriteBuildMsg(cpy.Snapshot.Logs...)

	if !cpy.Info.SourceIsSim {
		if cpy.Info.ActorIndex < 0 {
			log.Println(cpy)
		}
		preDmgModDebug := h.Team.CombatByIndex(cpy.Info.ActorIndex).ApplyAttackMods(&cpy, t)
		evt.Write("pre_damage_mods", preDmgModDebug)
	}

	dmg, crit = t.Attack(&cpy, evt)

	//delay damage event to end of the frame
	h.Tasks.Add(func() {
		//apply the damage
		t.ApplyDamage(&cpy, dmg)
		h.Events.Emit(event.OnDamage, t, &cpy, dmg, crit)
		//callbacks
		cb := AttackCB{
			Target:      t,
			AttackEvent: &cpy,
			Damage:      dmg,
			IsCrit:      crit,
		}
		for _, f := range cpy.Callbacks {
			f(cb)
		}
	}, 0)

	// this works because string in golang is a slice underneath, so the &amp points to the slice info
	// that's why when the underlying string in amp changes (has to be reallocated) the pointer doesn't
	// change since it's just pointing to the slice "header"
	if cpy.Info.Amped {
		amp = string(cpy.Info.AmpType)
	}
	if cpy.Info.Catalyzed {
		cata = string(cpy.Info.CatalyzedType)
	}
	return dmg, true
}

func (h *Handler) ApplyAttack(a *AttackEvent) float64 {
	// died := false
	var total float64
	var landed bool
	if a.Pattern.Targets[TargettablePlayer] {
		//TODO: we don't check for landed here since attack that hit player should never generate hitlag?
		h.attack(h.player, a)
	}
	if a.Pattern.Targets[TargettableEnemy] {
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
	if a.Pattern.Targets[TargettableGadget] {
		for i := 0; i < len(h.gadgets); i++ {
			//sanity check here; possible gadgets died and have not been cleaned up yet
			if h.gadgets[i] == nil {
				continue
			}
			h.attack(h.gadgets[i], a)
		}
	}
	//add hitlag to actor but ignore if this is deployable
	if h.EnableHitlag && landed && !a.Info.IsDeployable {
		dur := a.Info.HitlagHaltFrames
		if h.DefHalt && a.Info.CanBeDefenseHalted {
			dur += 3.6 //0.06
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
	h.TotalDamage += total
	return total
}
