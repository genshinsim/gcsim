package simulation

import "github.com/genshinsim/gcsim/pkg/core"

func (s *Simulation) handleEnergy() {
	if s.cfg.Energy.Active && s.C.F-s.lastEnergyDrop >= s.cfg.Energy.Start {
		f := s.C.Rand.Intn(s.cfg.Energy.End - s.cfg.Energy.Start)
		s.lastEnergyDrop = s.C.F + f
		s.C.Tasks.Add(func() {
			s.C.Energy.DistributeParticle(core.Particle{
				Source: "drop",
				Num:    s.cfg.Energy.Particles,
				Ele:    core.NoElement,
			})
		}, f)
		s.C.Log.Debugw("energy queued", "frame", s.C.F, "event", core.LogSimEvent, "last", s.lastEnergyDrop, "cfg", s.cfg.Energy, "amt", s.cfg.Energy.Particles, "energy_frame", s.C.F+f)
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
	s.C.Events.Subscribe(core.OnDamage, func(args ...interface{}) bool {
		atk := args[1].(*core.AttackEvent)
		if atk.Info.AttackTag != core.AttackTagNormal && atk.Info.AttackTag != core.AttackTagExtra {
			return false
		}
		//check icd
		if icd > s.C.F {
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
		char.AddEnergy(1)
		s.C.Log.Debugw("random energy on normal", "frame", s.C.F, "event", core.LogEnergyEvent, "char", atk.Info.ActorIndex, "chance", current[w])
		//set icd
		icd = s.C.F + 12
		current[w] = 0
		return false
	}, "random-energy-restore-on-hit")
	s.C.Events.Subscribe(core.OnCharacterSwap, func(args ...interface{}) bool {
		//TODO: assuming we clear the probability on swap
		for i := range current {
			current[i] = 0
		}
		return false
	}, "random-energy-restore-on-hit-swap")
}
