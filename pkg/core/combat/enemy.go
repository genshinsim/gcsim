package combat

import (
	"sort"

	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

func (h *Handler) Enemy(i int) Target {
	if i < 0 || i > len(h.enemies) {
		return nil
	}
	return h.enemies[i]
}

func (h *Handler) SetEnemyPos(i int, p Point) bool {
	if i < 0 || i > len(h.enemies)-1 {
		return false
	}
	h.enemies[i].SetPos(p)
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

// EnemyByDistance returns an array of indices of the enemies sorted by distance
func (c *Handler) EnemyByDistance(p Point, excl TargetKey) []int {
	//we dont actually need to know the exact distance. just find the lowest
	//of x^2 + y^2 to avoid sqrt

	var tuples []struct {
		ind  int
		dist float64
	}

	for i, v := range c.enemies {
		if v.Key() == excl {
			continue
		}
		if !v.IsAlive() {
			continue
		}
		vPos := v.Shape().Pos()
		dist := p.Sub(vPos).MagnitudeSquared()
		tuples = append(tuples, struct {
			ind  int
			dist float64
		}{ind: i, dist: dist})
	}

	sort.Slice(tuples, func(i, j int) bool {
		return tuples[i].dist < tuples[j].dist
	})

	result := make([]int, 0, len(tuples))

	for _, v := range tuples {
		result = append(result, v.ind)
	}

	return result
}

// EnemiesWithinRadius returns an array of indices of the enemies within radius r
func (c *Handler) EnemiesWithinRadius(p Point, r float64) []int {
	result := make([]int, 0, len(c.enemies))
	for i, v := range c.enemies {
		vPos := v.Shape().Pos()
		dist := p.Sub(vPos).MagnitudeSquared()
		if dist > r*r {
			continue
		}
		if !v.IsAlive() {
			continue
		}
		result = append(result, i)
	}

	return result
}

// EnemyExcl returns array of indices of enemies, excluding self
func (c *Handler) EnemyExcl(self TargetKey) []int {
	result := make([]int, 0, len(c.enemies))

	for i, e := range c.enemies {
		if e.Key() == self {
			continue
		}
		if !e.IsAlive() {
			continue
		}
		result = append(result, i)
	}

	return result
}

func (c *Handler) RandomEnemyTarget() int {

	count := len(c.enemies)
	if count == 0 {
		//this will basically cause that attack to hit nothing
		return -1
	}
	return c.Rand.Intn(count)
}
