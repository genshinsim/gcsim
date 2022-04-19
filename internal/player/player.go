// Package player provide all player related functionality including tracking the characters,
// player animation state, player aura, etc..
package player

import "github.com/genshinsim/gcsim/pkg/core"

type Player struct {
	State State
	core  *core.Core

	FramesSettings FramesSettings
	Stam           float64
	LastStamUse    int

	LastAction struct {
		Type  core.ActionType
		Param map[string]int
		Char  int
	}
}

//State is used to keep track of the current state of player
type State struct {
	Animation         core.AnimationState
	FrameStarted      int                            // track the frame the current AnimationState started
	OnStateEnd        func()                         // call back to be executed when this tate ends
	AnimationDuration func(next core.ActionType) int // duration of the current animation, given next animation
}

type FramesSettings struct {
	Jump int
	Dash int
	Swap int
}
