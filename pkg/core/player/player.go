//Package player contains player related tracking and functionalities:
// - tracking characters on the team
// - handling animations state
// - handling normal attack state
// - handling character stats and attributes
// - handling shielding
package player

import (
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/infusion"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

type Handler struct {
	chars  []*character.CharWrapper
	active int
	log    glog.Logger
	events event.Eventter

	Shields shield.Handler
	infusion.InfusionHandler
}

func (h *Handler) ByIndex(i int) *character.CharWrapper {
	return h.chars[i]
}

func (h *Handler) Chars() []*character.CharWrapper {
	return h.chars
}

func (h *Handler) Active() int {
	return h.active
}

func (h *Handler) SetActive(i int) {
	h.active = i
}

func (e *Handler) Adjust(src string, char int, amt float64) {
	e.chars[char].AddEnergy(src, amt)
}

func (e *Handler) DistributeParticle(p character.Particle) {
	for i, char := range e.chars {
		char.ReceiveParticle(p, e.active == i, len(e.chars))
	}
	e.events.Emit(event.OnParticleReceived, p)
}
