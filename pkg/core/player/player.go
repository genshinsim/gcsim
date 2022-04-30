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
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/artifact"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/infusion"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

type Handler struct {
	log    glog.Logger
	events event.Eventter

	chars   []*character.CharWrapper
	active  int
	charPos map[keys.Char]int

	weaps   []weapon.Weapon
	weapPos map[keys.Weapon]int

	sets   []artifact.Set
	setPos map[keys.Set]int

	Shields shield.Handler
	infusion.InfusionHandler
}

func (h *Handler) AddChar(char *character.CharWrapper) int {
	h.chars = append(h.chars, char)
	index := len(h.chars) - 1
	char.SetIndex(index)
	h.charPos[char.Base.Key] = index
	return index
}

func (h *Handler) AddWeapon(key keys.Weapon, w weapon.Weapon) int {
	h.weaps = append(h.weaps, w)
	index := len(h.weaps) - 1
	w.SetIndex(index)
	h.weapPos[key] = index
	return index
}

func (h *Handler) AddSet(key keys.Set, set artifact.Set) int {
	h.sets = append(h.sets, set)
	index := len(h.weaps) - 1
	set.SetIndex(index)
	h.setPos[key] = index
	return index
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
