package combat

import (
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/model"
)

func (h *Handler) SetPlayer(t model.Target) {
	h.player = t
	t.SetKey(0)
}

func (h *Handler) Player() model.Target {
	return h.player
}

func (h *Handler) SetPlayerPos(p geometry.Point) {
	h.player.SetPos(p)
}
