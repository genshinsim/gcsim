package modifier

import (
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
)

type Manager struct {
	handlers []*handler
}

func NewManager(size int) *Manager {
	h := make([]*handler, size)
	for i := range h {
		h[i] = &handler{}
	}
	return &Manager{
		handlers: h,
	}
}

func (m *Manager) Tick(id keys.TargetID) {
	// the assumption here is that id is index of handlers
	m.handlers[id].tick()
}

func (m *Manager) AddModifier(mod *info.Modifier, s info.StackingType) (bool, error) {
	return m.handlers[mod.Source].add(mod, s)
}
