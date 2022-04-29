//Package team provides access to each character and their abilities
package team

type Handler struct {
	team []CharWrapper
}

func (h *Handler) ByIndex(i int) CharWrapper {
	return h.team[i]
}
