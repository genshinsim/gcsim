package engine

import "github.com/genshinsim/gcsim/pkg/core/info"

type ModifierMgr interface {
	Handler(info.EntityIndex) ModifierHandler
}

type ModifierHandler interface {
	Tick()
	Add(*info.Modifier) (bool, error)
}
