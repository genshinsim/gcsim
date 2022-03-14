package simulation

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/coretype"
)

func (s *Simulation) handleEnergy() {
	if s.cfg.Energy.Active && s.C.Frame-s.lastEnergyDrop >= s.cfg.Energy.Start {
		f := s.C.Rand.Intn(s.cfg.Energy.End - s.cfg.Energy.Start)
		s.lastEnergyDrop = s.C.Frame + f
		s.C.Tasks.Add(func() {
			s.C.Energy.DistributeParticle(core.Particle{
				Source: "drop",
				Num:    s.cfg.Energy.Particles,
				Ele:    core.NoElement,
			})
		}, f)
		s.C.Log.NewEvent("energy queued", coretype.LogSimEvent, -1, "last", s.lastEnergyDrop, "cfg", s.cfg.Energy, "amt", s.cfg.Energy.Particles, "energy_frame", s.C.Frame+f)
	}
}

func (s *Simulation) randomOnHitEnergy() {
	/**
	WeaponClassSword
	WeaponClassClaymore
	WeaponClassSpear
	WeaponClassBow
	WeaponClassCatalyst
	**/
	current := make([]float64, core.EndWeaponClass)
	inc := []float64{
		0.05,
		0.05,
		0.04,
		0.01,
		0.01,
	}

	//TODO not sure if there's like a 0.2s icd on this. for now let's add it in to be safe
	icd := 0
	s.C.Subscribe(coretype.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*coretype.AttackEvent)
		if atk.Info.AttackTag != coretype.AttackTagNormal && atk.Info.AttackTag != coretype.AttackTagExtra {
			return false
		}
		//check icd
		if icd > s.C.Frame {
			return false
		}
		//check chance
		char := s.C.Chars[atk.Info.ActorIndex]
		w := char.WeaponClass()
		if s.C.Rand.Float64() > current[w] {
			//increment chance
			current[w] += inc[w]
			return false
		}
		//add energy
		char.AddEnergy("na-ca-on-hit", 1)
		// Add this log in sim if necessary to see as AddEnergy already generates a log
		s.C.Log.NewEvent("random energy on normal", coretype.LogSimEvent, char.Index(), "char", atk.Info.ActorIndex, "chance", current[w])
		//set icd
		icd = s.C.Frame + 12
		current[w] = 0
		return false
	}, "random-energy-restore-on-hit")
	s.C.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		//TODO: assuming we clear the probability on swap
		for i := range current {
			current[i] = 0
		}
		return false
	}, "random-energy-restore-on-hit-swap")
}
