// package modifer provides a Handler that is used to handle
// modifiers on a per target basis, with duration/durability/stacking
// handling mirror'd to game (as much as possible)
package modifier

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core/info"
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
		if mod.Durability == 0 {
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

// Add returns true if this is a new modifier or error on unsupported stacking type
func (h *Handler) Add(m *info.Modifier) (bool, error) {
	if m.Durability <= 0 {
		return false, fmt.Errorf("modifier %v has non-positive durability: %v", m.Name, m.Durability)
	}
	// force decay rate, but only if duration is positive
	// negative duration means infinite duration
	m.DecayRate = 0
	if m.Duration > 0 {
		m.DecayRate = m.Durability / info.Durability(m.Duration)
	}
	var added bool
	switch m.Stacking {
	case info.Refresh:
		added = h.refresh(m)
	case info.Unique:
		added = h.unique(m)
	case info.Overlap:
		added = h.overlap(m)
	case info.OverlapRefreshDuration:
		added = h.overlapRefreshDuration(m)
	default:
		return false, fmt.Errorf("unsupported stacking type: %v", m.Stacking)
	}
	if added && m.OnAdd != nil {
		m.OnAdd(m)
	}
	return added, nil
}

// single instance; can't be re-applied unless expired
func (h *Handler) unique(m *info.Modifier) bool {
	for _, mod := range h.modifiers {
		if mod.Name == m.Name {
			return false
		}
	}
	h.modifiers = append(h.modifiers, m)
	return true
}

// single instance. re-apply resets durability and doesn't trigger onAdded/onRemoved
// or reset onThinkInterval
func (h *Handler) refresh(m *info.Modifier) bool {
	for _, mod := range h.modifiers {
		if mod.Name == m.Name {
			// TODO: this will overwrite everything in the dest mod, including the src
			// consider if this actually matters or not?
			*mod = *m
			return false
		}
	}
	h.modifiers = append(h.modifiers, m)
	return true
}

// multiple instances can co-exist at the same time
func (h *Handler) overlap(m *info.Modifier) bool {
	h.modifiers = append(h.modifiers, m)
	return true
}

// refresh any existing with lower durability; update decay rate and duration
func (h *Handler) overlapRefreshDuration(m *info.Modifier) bool {
	for _, mod := range h.modifiers {
		if mod.Name == m.Name && mod.Durability < m.Durability {
			mod.Durability = m.Durability
			mod.Duration = m.Duration
			mod.DecayRate = m.DecayRate
		}
	}
	h.modifiers = append(h.modifiers, m)
	return true
}
