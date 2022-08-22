// Package combat handles all combat related functionalities including
//	- target tracking
//	- target selection
//	- hitbox collision checking
//  - attack queueing
package combat

import (
	"math/rand"

	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
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
	targets     []Target
	TotalDamage float64
}

type Opt struct {
	Events        event.Eventter
	Team          CharHandler
	Rand          *rand.Rand
	Debug         bool
	Log           glog.Logger
	DamageMode    bool
	DefHalt       bool
	EnableHitlag  bool
	DefaultTarget int //index for default target
}

func New(opt Opt) *Handler {
	h := &Handler{
		Opt: opt,
	}
	h.targets = make([]Target, 0, 5)

	return h
}

func (h *Handler) AddTarget(t Target) {
	h.targets = append(h.targets, t)
	t.SetIndex(len(h.targets) - 1)
}

func (h *Handler) Target(i int) Target {
	if i < 0 || i > len(h.targets) {
		return nil
	}
	return h.targets[i]
}

func (h *Handler) Targets() []Target {
	return h.targets
}

func (h *Handler) TargetsCount() int {
	return len(h.targets)
}

func (h *Handler) PrimaryTarget() Target {
	return h.targets[h.DefaultTarget]
}

func (h *Handler) Player() Target {
	return h.targets[0] //assuming player is always position 0
}

func (h *Handler) SetTargetPos(i int, x, y float64) bool {
	if i < 0 || i > len(h.targets)-1 {
		return false
	}
	h.targets[i].SetPos(x, y)
	h.Log.NewEvent("target position changed", glog.LogSimEvent, -1).
		Write("index", i).
		Write("x", x).
		Write("y", y)
	return true
}

func (h *Handler) KillTarget(i int) bool {
	// don't kill yourself
	if i < 1 || i > len(h.targets)-1 {
		return false
	}
	h.targets[i].Kill()
	h.Events.Emit(event.OnTargetDied, h.targets[i], &AttackEvent{}) // TODO: it's fine?
	h.Log.NewEvent("target is died", glog.LogSimEvent, -1).Write("index", i)
	return true
}

func (h *Handler) Tick() {
	for _, t := range h.targets {
		t.Tick()
	}
}
