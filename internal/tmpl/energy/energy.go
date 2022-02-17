package energy

import "github.com/genshinsim/gcsim/pkg/core"

type Ctrl struct {
	core *core.Core
}

func NewCtrl(c *core.Core) *Ctrl {
	return &Ctrl{
		core: c,
	}
}

func (e *Ctrl) DistributeParticle(p core.Particle) {
	for i, char := range e.core.Chars {
		char.ReceiveParticle(p, e.core.ActiveChar == i, len(e.core.Chars))
	}
	e.core.Events.Emit(core.OnParticleReceived, p)
}
