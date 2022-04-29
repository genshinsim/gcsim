// Package core provides core functionality for a simulation:
//	- combat
//	- tasks
//	- event handling
//	- logging
// 	- constructs (really should be just generic objects?)
//	- status
package core

import (
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/construct"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/status"
	"github.com/genshinsim/gcsim/pkg/core/task"
)

type Core struct {
	F     int
	Flags Flags
	//various functionalities of core
	Log        glog.Logger   //we use an interface here so that we can pass in a nil logger for all except 1 run
	Events     event.Handler //track events: subscribe/unsubscribe/emit
	Status     status.Handler
	Tasks      task.Handler
	Combat     combat.Handler
	Constructs construct.Handler
	Player     player.Handler
}

type Flags struct {
	LogDebug     bool // Used to determine logging level
	ChildeActive bool // Used for Childe +1 NA talent passive
	Custom       map[string]int
}
