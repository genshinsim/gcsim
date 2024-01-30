package enemy

import (
	"log"
	"math"

	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

func (e *Enemy) HandleAttack(atk *combat.AttackEvent) float64 {
	// at this point attack will land
	e.Core.Combat.Events.Emit(event.OnEnemyHit, e, atk)

	var amp string
	var cata string
	var dmg float64
	var crit bool

	evt := e.Core.Combat.Log.NewEvent(atk.Info.Abil, glog.LogDamageEvent, atk.Info.ActorIndex).
		Write("target", e.Key()).
		Write("attack-tag", atk.Info.AttackTag).
		Write("ele", atk.Info.Element.String()).
		Write("damage", &dmg).
		Write("crit", &crit).
		Write("amp", &amp).
		Write("cata", &cata).
		Write("abil", atk.Info.Abil).
		Write("source_frame", atk.SourceFrame)
	evt.WriteBuildMsg(atk.Snapshot.Logs...)

	if !atk.Info.SourceIsSim {
		if atk.Info.ActorIndex < 0 {
			log.Println(atk)
		}
		preDmgModDebug := e.Core.Combat.Team.CombatByIndex(atk.Info.ActorIndex).ApplyAttackMods(atk, e)
		evt.Write("pre_damage_mods", preDmgModDebug)
	}

	dmg, crit = e.attack(atk, evt)

	// delay damage event to end of the frame
	e.Core.Combat.Tasks.Add(func() {
		// apply the damage
		actualDmg := e.applyDamage(atk, dmg)
		e.Core.Combat.TotalDamage += actualDmg
		e.Core.Combat.Events.Emit(event.OnEnemyDamage, e, atk, actualDmg, crit)
		// callbacks
		cb := combat.AttackCB{
			Target:      e,
			AttackEvent: atk,
			Damage:      actualDmg,
			IsCrit:      crit,
		}
		for _, f := range atk.Callbacks {
			f(cb)
		}
	}, 0)

	// this works because string in golang is a slice underneath, so the &amp points to the slice info
	// that's why when the underlying string in amp changes (has to be reallocated) the pointer doesn't
	// change since it's just pointing to the slice "header"
	if atk.Info.Amped {
		amp = string(atk.Info.AmpType)
	}
	if atk.Info.Catalyzed {
		cata = string(atk.Info.CatalyzedType)
	}
	return dmg
}

func (e *Enemy) attack(atk *combat.AttackEvent, evt glog.Event) (float64, bool) {
	// if target is frozen prior to attack landing, set impulse to 0
	// let the break freeze attack to trigger actual impulse
	if e.Durability[reactable.Frozen] > reactable.ZeroDur {
		atk.Info.NoImpulse = true
	}

	// check poise dmg and then shatter first
	e.PoiseDMGCheck(atk)
	e.ShatterCheck(atk)

	checkBurningICD := func() {
		// special global ICD for Burning DMG
		if atk.Info.ICDTag != attacks.ICDTagBurningDamage {
			return
		}
		// checks for ICD on all the other characters as well
		for i := 0; i < len(e.Core.Player.Chars()); i++ {
			if i == atk.Info.ActorIndex {
				continue
			}
			// burning durability wiped out to 0 if any of the other char still on icd re burning dmg
			atk.Info.Durability *= reactions.Durability(e.WillApplyEle(atk.Info.ICDTag, atk.Info.ICDGroup, i))
		}
	}
	// check tags
	if atk.Info.Durability > 0 {
		// check for ICD first
		atk.Info.Durability *= reactions.Durability(e.WillApplyEle(atk.Info.ICDTag, atk.Info.ICDGroup, atk.Info.ActorIndex))
		checkBurningICD()
		if atk.Info.Durability > 0 && atk.Info.Element != attributes.Physical {
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
					Write("target", e.Key()).
					Write("existing", existing).
					Write("after", e.Reactable.ActiveAuraString())
			}
		}
	}

	damage, isCrit := e.calc(atk, evt)

	// check for hitlag
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
			// apply hit lag to enemy
			e.ApplyHitlag(atk.Info.HitlagFactor, dur)
			// also apply hitlag to reactable
			// e.Reactable.ApplyHitlag(atk.Info.HitlagFactor, dur)
		}
	}

	// check for particle drops
	if e.prof.ParticleDropThreshold > 0 {
		next := int(e.damageTaken / e.prof.ParticleDropThreshold)
		if next > e.lastParticleDrop {
			// check the count too
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

func (e *Enemy) applyDamage(atk *combat.AttackEvent, damage float64) float64 {
	// record dmg
	// do not let hp become negative because this function can be called multiple times in same frame
	actualDmg := min(damage, e.hp) // do not let dmg be greater than remaining enemy hp
	e.hp -= actualDmg
	e.damageTaken += actualDmg //TODO: do we actually need this?

	// check if target is dead
	if e.Core.Flags.DamageMode && e.hp <= 0 {
		e.Kill()
		e.Core.Events.Emit(event.OnTargetDied, e, atk)
		return actualDmg
	}

	// apply auras
	if atk.Info.Durability > 0 && !atk.Reacted && atk.Info.Element != attributes.Physical {
		// check for ICD first
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
				Write("target", e.Key()).
				Write("existing", existing).
				Write("after", e.Reactable.ActiveAuraString())
		}
	}
	// just return damage without considering enemy hp here for both:
	// - damage mode if target not dead (otherwise would have entered the death if statement)
	// - duration mode (no concept of killing blow)
	return damage
}
