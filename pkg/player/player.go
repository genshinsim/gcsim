// Package player handles all player related functionalities including:
//	- handle all characters
//	- resonance
// 	- character actions
//	- animation state
package player

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/reactable"
)

type Player struct {
	Core *core.Core
	*reactable.Reactable

	//Team tracking
	Team       []*MasterChar
	ActiveChar int
}

const (
	MaxStam            = 240
	StamCDFrames       = 90
	JumpFrames         = 33
	DashFrames         = 24
	WalkFrames         = 1
	SwapCDFrames       = 60
	MaxTeamPlayerCount = 4
	DefaultTargetIndex = 1
)

func (p *Player) Active() *MasterChar {
	return p.Team[p.ActiveChar]
}

func (p *Player) Chars() []*MasterChar {
	return p.Team
}

func (p *Player) ActiveIndex() int {
	return p.ActiveChar
}

func (p *Player) ByIndex(idx int) *MasterChar {
	return p.Team[idx]
}
