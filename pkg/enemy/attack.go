package enemy

import (
	"math"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func (e *Enemy) Attack(atk *combat.AttackEvent, evt glog.Event) (float64, bool) {
	//if target is frozen prior to attack landing, set impulse to 0
	//let the break freeze attack to trigger actual impulse
	if e.AuraType() == attributes.Frozen {
		atk.Info.NoImpulse = true
	}

	//check shatter first
	e.ShatterCheck(atk)

	//check tags
	if atk.Info.Durability > 0 {
		//check for ICD first
		if e.WillApplyEle(atk.Info.ICDTag, atk.Info.ICDGroup, atk.Info.ActorIndex) && atk.Info.Element != attributes.Physical {
			existing := e.Reactable.ActiveAuraString()
			applied := atk.Info.Durability
			e.React(atk)
			if e.Core.Flags.LogDebug {
				e.Core.Log.NewEvent(
					"application",
					glog.LogElementEvent,
					atk.Info.ActorIndex,
					"attack_tag", atk.Info.AttackTag,
					"applied_ele", atk.Info.Element.String(),
					"dur", applied,
					"abil", atk.Info.Abil,
					"target", e.TargetIndex,
					"existing", existing,
					"after", e.Reactable.ActiveAuraString(),
				)
			}
		}
	}

	damage, isCrit := e.calc(atk, evt)

	//record dmg
	e.HPCurrent -= damage
	e.damageTaken += damage //TODO: do we actually need this?

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
			e.Reactable.ApplyHitlag(atk.Info.HitlagFactor, dur)
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
