package combat

import (
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

func (h *Handler) Enemy(i int) info.Target {
	if i < 0 || i > len(h.enemies) {
		return nil
	}
	return h.enemies[i]
}

func (h *Handler) SetEnemyPos(i int, p info.Point) bool {
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
	h.Events.Emit(event.OnTargetDied, h.enemies[i], &info.AttackEvent{}) // TODO: it's fine?
	h.Log.NewEvent("enemy dead", glog.LogSimEvent, -1).Write("index", i)
}

func (h *Handler) AddEnemy(t info.Target) {
	h.enemies = append(h.enemies, t)
	t.SetKey(h.nextkey())
}

func (h *Handler) Enemies() []info.Target {
	return h.enemies
}

func (h *Handler) EnemyCount() int {
	return len(h.enemies)
}

func (h *Handler) PrimaryTarget() info.Target {
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
