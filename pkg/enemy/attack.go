package enemy

import (
	"math"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

func (e *Enemy) Attack(atk *combat.AttackEvent, evt glog.Event) (float64, bool) {
	//if target is frozen prior to attack landing, set impulse to 0
	//let the break freeze attack to trigger actual impulse
	if e.Durability[reactable.ModifierFrozen] > reactable.ZeroDur {
		atk.Info.NoImpulse = true
	}

	//check shatter first
	e.ShatterCheck(atk)

	//check tags
	if atk.Info.Durability > 0 {
		//check for ICD first
		atk.OnICD = !e.WillApplyEle(atk.Info.ICDTag, atk.Info.ICDGroup, atk.Info.ActorIndex)
		//special global ICD for Burning DMG
		if atk.Info.ICDTag == combat.ICDTagBurningDamage {
			//checks for ICD on all the other characters as well
			for i := 0; i < len(e.Core.Player.Chars()); i++ {
				if i != atk.Info.ActorIndex {
					atk.OnICD = atk.OnICD || !e.WillApplyEle(atk.Info.ICDTag, atk.Info.ICDGroup, i)
				}
			}
		}
		if !atk.OnICD && atk.Info.Element != attributes.Physical {
			existing := e.Reactable.ActiveAuraString()
			applied := atk.Info.Durability
			e.React(atk)
			if e.Core.Flags.LogDebug && atk.Reacted {
				e.Core.Log.NewEvent(
					"application",
					glog.LogElementEvent,
					atk.Info.ActorIndex,
				).
					Write("attack_tag", atk.Info.AttackTag).
					Write("applied_ele", atk.Info.Element.String()).
					Write("dur", applied).
					Write("abil", atk.Info.Abil).
					Write("target", e.TargetIndex).
					Write("existing", existing).
					Write("after", e.Reactable.ActiveAuraString())

			}
		}
	}

	damage, isCrit := e.calc(atk, evt)

	//check for hitlag
	if e.Core.Combat.EnableHitlag {
		willapply := true
		if atk.Info.HitlagOnHeadshotOnly {
			willapply = atk.Info.HitWeakPoint
		}
		dur := atk.Info.HitlagHaltFrames
		if e.Core.Flags.DefHalt && atk.Info.CanBeDefenseHalted {
			dur += 3.6
		}
		dur = math.Ceil(dur)
		if willapply && dur > 0 {
			//apply hit lag to enemy
			e.ApplyHitlag(atk.Info.HitlagFactor, dur)
			//also apply hitlag to reactable
			// e.Reactable.ApplyHitlag(atk.Info.HitlagFactor, dur)
		}
	}

	//check for particle drops
	if e.prof.ParticleDropThreshold > 0 {
		next := int(e.damageTaken / e.prof.ParticleDropThreshold)
		if next > e.lastParticleDrop {
			//check the count too
			count := next - e.lastParticleDrop
			e.lastParticleDrop = next
			e.Core.Log.NewEvent("particle hp threshold triggered", glog.LogEnemyEvent, atk.Info.ActorIndex)
			e.Core.Tasks.Add(
				func() {
					e.Core.Player.DistributeParticle(character.Particle{
						Source: "hp_drop",
						Num:    e.prof.ParticleDropCount * float64(count),
						Ele:    e.prof.ParticleElement,
					})
				},
				100, //TODO: should be subject to global delay maybe??
			)
		}
	}

	return damage, isCrit
}

func (e *Enemy) ApplyDamage(atk *combat.AttackEvent, damage float64) {
	//record dmg
	e.hp -= damage
	e.damageTaken += damage //TODO: do we actually need this?

	//check if target is dead
	if e.Core.Flags.DamageMode && e.hp <= 0 {
		e.Kill()
		e.Core.Events.Emit(event.OnTargetDied, &combat.AttackEventPayload{
			Target:      e,
			AttackEvent: atk,
		})
		return
	}

	//apply auras
	if atk.Info.Durability > 0 && !atk.Reacted && !atk.OnICD && atk.Info.Element != attributes.Physical {
		//check for ICD first
		existing := e.Reactable.ActiveAuraString()
		applied := atk.Info.Durability
		e.AttachOrRefill(atk)
		if e.Core.Flags.LogDebug {
			e.Core.Log.NewEvent(
				"application",
				glog.LogElementEvent,
				atk.Info.ActorIndex,
			).
				Write("attack_tag", atk.Info.AttackTag).
				Write("applied_ele", atk.Info.Element.String()).
				Write("dur", applied).
				Write("abil", atk.Info.Abil).
				Write("target", e.TargetIndex).
				Write("existing", existing).
				Write("after", e.Reactable.ActiveAuraString())

		}
	}
}
