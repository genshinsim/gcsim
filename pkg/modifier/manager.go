package modifier

import (
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/engine"
)

type Manager struct {
	handlers []*Handler
}

func NewManager(size int) *Manager {
	h := make([]*Handler, size)
	for i := range h {
		h[i] = &Handler{}
	}
	return &Manager{
		handlers: h,
	}
}

func (m *Manager) Handler(id info.EntityIndex) engine.ModifierHandler {
	return m.handlers[id]
}

func (m *Manager) Add(mod *info.Modifier) (bool, error) {
	return m.handlers[mod.Target].Add(mod)
}
