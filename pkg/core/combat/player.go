package combat

func (h *Handler) SetPlayer(t Target) {
	h.player = t
}

func (h *Handler) Player() Target {
	return h.player
}

func (h *Handler) SetPlayerPos(x, y float64) {
	h.player.SetPos(x, y)
}
