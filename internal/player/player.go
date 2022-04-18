// Package player provide all player related functionality including tracking the characters,
// player animation state, player aura, etc..
package player

import "github.com/genshinsim/gcsim/pkg/core"

type Player struct {
	State State
	core  *core.Core
}

//State is used to keep track of the current state of player
type State struct {
	FrameStarted      int                            // track the frame the current AnimationState started
	OnStateEnd        func()                         // call back to be executed when this tate ends
	AnimationDuration func(next core.ActionType) int // duration of the current animation, given next animation
}
