// Package combat handles all combat related functionalities including
//   - target tracking
//   - target selection
//   - hitbox collision checking
//   - attack queueing
package combat

import (
	"math/rand"

	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/core/task"
)

type CharHandler interface {
	CombatByIndex(int) Character
	ApplyHitlag(char int, factor, dur float64)
}

type Character interface {
	ApplyAttackMods(a *AttackEvent, t Target) []interface{}
}

type Handler struct {
	Opt
	enemies     []Target
	gadgets     []Gadget
	player      Target
	TotalDamage float64
	gccount     int
	keycount    targets.TargetKey
}

type Opt struct {
	Events        event.Eventter
	Tasks         task.Tasker
	Team          CharHandler
	Rand          *rand.Rand
	Debug         bool
	Log           glog.Logger
	DamageMode    bool
	DefHalt       bool
	EnableHitlag  bool
	DefaultTarget targets.TargetKey // index for default target
}

func New(opt Opt) *Handler {
	h := &Handler{
		Opt:      opt,
		keycount: 1,
	}
	h.enemies = make([]Target, 0, 5)
	h.gadgets = make([]Gadget, 0, 10)

	return h
}

func (h *Handler) nextkey() targets.TargetKey {
	h.keycount++
	return h.keycount - 1
}

func (h *Handler) Tick() {
	// collision check happens before each object ticks (as collision may remove the object)
	// enemy and player does not check for collision
	// gadgets check against player and enemy
	for i := 0; i < len(h.gadgets); i++ {
		if h.gadgets[i] != nil && h.gadgets[i].CollidableWith(targets.TargettablePlayer) {
			if h.gadgets[i].WillCollide(h.player.Shape()) {
				h.gadgets[i].CollidedWith(h.player)
			}
		}
		// sanity check in case gadget is gone
		if h.gadgets[i] != nil && h.gadgets[i].CollidableWith(targets.TargettableEnemy) {
			for j := 0; j < len(h.enemies) && h.gadgets[i] != nil; j++ {
				if h.gadgets[i].WillCollide(h.enemies[j].Shape()) {
					h.gadgets[i].CollidedWith(h.enemies[j])
				}
			}
		}
	}
	h.player.Tick()
	for _, v := range h.enemies {
		v.Tick()
	}
	for _, v := range h.gadgets {
		if v != nil {
			v.Tick()
		}
	}
	//TODO: clean up every 100 tick reasonable?
	h.gccount++
	if h.gccount > 100 {
		n := 0
		for i, v := range h.gadgets {
			if v != nil {
				h.gadgets[n] = h.gadgets[i]
				n++
			}
		}
		h.gadgets = h.gadgets[:n]
		h.gccount = 0
	}
}
