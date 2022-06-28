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
}

type Character interface {
	ApplyAttackMods(a *AttackEvent, t Target) []interface{}
	ApplyHitlag(factor float64, dur int)
}

type Handler struct {
	log    glog.Logger
	events event.Eventter
	team   CharHandler
	rand   *rand.Rand
	debug  bool

	targets     []Target
	TotalDamage float64
	DamageMode  bool

	defHalt bool
}

func New(log glog.Logger, events event.Eventter, team CharHandler, rand *rand.Rand, debug bool, damageMode bool, defHalt bool) *Handler {
	h := &Handler{
		log:        log,
		events:     events,
		team:       team,
		rand:       rand,
		DamageMode: damageMode,
		debug:      debug,
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

func (h *Handler) SetTargetPos(i int, x, y float64) {
	h.targets[i].SetPos(x, y)
}

func (h *Handler) Tick() {
	for _, t := range h.targets {
		t.Tick()
	}
}
