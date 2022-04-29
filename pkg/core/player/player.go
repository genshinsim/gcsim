//Package player contains player related tracking and functionalities:
// - tracking characters on the team
// - handling animations state
// - handling normal attack state
// - handling character stats and attributes
// - handling shielding
package player

import (
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/mods"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

type Handler struct {
	team    []character.CharWrapper
	Shields shield.Handler
	mods.Handler
}

func (h *Handler) ByIndex(i int) character.CharWrapper {
	return h.team[i]
}
