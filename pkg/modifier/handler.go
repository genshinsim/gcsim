// package modifer provides a Handler that is used to handle
// modifiers on a per target basis, with duration/durability/stacking
// handling mirror'd to game (as much as possible)
package modifier

import (
	"slices"

	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type Handler struct {
	modifiers []*info.Modifier
}

func (h *Handler) Tick() {
	// reduce durability
	n := 0
	for i, mod := range h.modifiers {
		if mod.PreTick != nil {
			mod.PreTick(mod)
		}
		mod.Durability -= mod.DecayRate
		if mod.Durability < 0 {
			mod.Durability = 0
		}
		if mod.PostTick != nil {
			mod.PostTick(mod)
		}
		// delete mods if durability is 0
		if mod.Durability <= info.ZeroDur {
			if mod.OnRemove != nil {
				mod.OnRemove(mod)
			}
		} else {
			h.modifiers[n] = h.modifiers[i]
			n++
		}
	}
	h.modifiers = h.modifiers[:n]
}

func (h *Handler) Get(name keys.Modifier) []*info.Modifier {
	var res []*info.Modifier
	for _, m := range h.modifiers {
		if m.Name == name {
			res = append(res, m)
		}
	}
	return res
}

func (h *Handler) GetMaxDurability(name keys.Modifier) *info.Modifier {
	var res *info.Modifier
	for _, m := range h.modifiers {
		if m.Name != name {
			continue
		}
		if res == nil || res.Durability < m.Durability {
			res = m
		}
	}
	return res
}

func (h *Handler) MatchAny(names ...keys.Modifier) bool {
	for _, m := range h.modifiers {
		if slices.Contains(names, m.Name) {
			return true
		}
	}
	return false
}
