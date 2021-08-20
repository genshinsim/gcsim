package team

import "github.com/genshinsim/gsim/pkg/core"

type Handler struct {
	status map[string]int
	core   *core.Core
}

func New(c *core.Core) *Handler {
	h := &Handler{core: c}

	return h
}

func (h *Handler) DistributeParticle(p core.Particle) {

}

func (h *Handler) Tick() {

}
