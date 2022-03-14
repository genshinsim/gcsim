package player

import "github.com/genshinsim/gcsim/pkg/coretype"

func (p *Player) DistributeParticle(particle coretype.Particle) {
	for i, char := range p.Chars {
		char.ReceiveParticle(particle, p.ActiveChar == i, len(p.Chars))
	}
	p.core.Emit(coretype.OnParticleReceived, particle)
}
