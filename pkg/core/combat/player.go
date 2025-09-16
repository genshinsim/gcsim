package combat

import (
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (h *Handler) SetPlayer(t info.Target) {
	h.player = t
	t.SetKey(0)
}

func (h *Handler) Player() info.Target {
	return h.player
}

func (h *Handler) SetPlayerPos(p geometry.Point) {
	h.player.SetPos(p)
}
