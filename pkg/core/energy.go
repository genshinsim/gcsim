package core

type EnergyHandler interface {
	DistributeParticle(p Particle)
}

type EnergyCtrl struct {
	core *Core
}

func NewEnergyCtrl(c *Core) *EnergyCtrl {
	return &EnergyCtrl{
		core: c,
	}
}

func (e *EnergyCtrl) DistributeParticle(p Particle) {
	for i, char := range e.core.Chars {
		char.ReceiveParticle(p, e.core.ActiveChar == i, len(e.core.Chars))
	}
	e.core.Events.Emit(OnParticleReceived, p)
}
