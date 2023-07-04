package combat

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

type Status struct {
	modifier.Base
}
type ResistMod struct {
	Ele   attributes.Element
	Value float64
	modifier.Base
}

type DefMod struct {
	Value float64
	Dur   int
	modifier.Base
}

type Enemy interface {
	Target
	// hp related
	MaxHP() float64
	HP() float64
	// hitlag related
	ApplyHitlag(factor, dur float64)
	QueueEnemyTask(f func(), delay int)
	// modifier related
	// add
	AddStatus(key string, dur int, hitlag bool)
	AddResistMod(mod ResistMod)
	AddDefMod(mod DefMod)
	// delete
	DeleteStatus(key string)
	DeleteResistMod(key string)
	DeleteDefMod(key string)
	// active
	StatusIsActive(key string) bool
	ResistModIsActive(key string) bool
	DefModIsActive(key string) bool
	StatusExpiry(key string) int
}

func (h *Handler) Enemy(i int) Target {
	if i < 0 || i > len(h.enemies) {
		return nil
	}
	return h.enemies[i]
}

func (h *Handler) SetEnemyPos(i int, p geometry.Point) bool {
	if i < 0 || i > len(h.enemies)-1 {
		return false
	}

	h.enemies[i].SetPos(p)
	h.Events.Emit(event.OnTargetMoved, h.enemies[i])

	h.Log.NewEvent("target position changed", glog.LogSimEvent, -1).
		Write("index", i).
		Write("x", p.X).
		Write("y", p.Y)
	return true
}

func (h *Handler) KillEnemy(i int) {
	h.enemies[i].Kill()
	h.Events.Emit(event.OnTargetDied, h.enemies[i], &AttackEvent{}) // TODO: it's fine?
	h.Log.NewEvent("enemy dead", glog.LogSimEvent, -1).Write("index", i)
}

func (h *Handler) AddEnemy(t Target) {
	h.enemies = append(h.enemies, t)
	t.SetKey(h.nextkey())
}

func (h *Handler) Enemies() []Target {
	return h.enemies
}

func (h *Handler) EnemyCount() int {
	return len(h.enemies)
}

func (h *Handler) PrimaryTarget() Target {
	for _, v := range h.enemies {
		if v.Key() == h.DefaultTarget {
			if !v.IsAlive() {
				h.Log.NewEvent("default target is dead", glog.LogWarnings, -1)
			}
			return v
		}
	}
	panic("default target does not exist?!")
}
