package combat

import (
	"math"
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

func (h *Handler) SetEnemyPos(i int, x, y float64) bool {
	if i < 0 || i > len(h.enemies)-1 {
		return false
	}
	h.enemies[i].SetPos(x, y)
	h.Log.NewEvent("target position changed", glog.LogSimEvent, -1).
		Write("index", i).
		Write("x", x).
		Write("y", y)
	return true
}

func (h *Handler) KillEnemy(i int) {
	h.enemies[i].Kill()
	h.Events.Emit(event.OnTargetDied, h.enemies[i], &AttackEvent{}) // TODO: it's fine?
	h.Log.NewEvent("enemy dead", glog.LogSimEvent, -1).Write("index", i)
}

func (h *Handler) AddEnemy(t Target) {
	h.enemies = append(h.enemies, t)
	t.SetIndex(len(h.enemies) - 1)
}

func (h *Handler) Enemies() []Target {
	return h.enemies
}

func (h *Handler) EnemyCount() int {
	return len(h.enemies)
}

func (h *Handler) PrimaryTarget() Target {
	return h.enemies[h.DefaultTarget]
}

// EnemyByDistance returns an array of indices of the enemies sorted by distance
func (c *Handler) EnemyByDistance(x, y float64, excl int) []int {
	//we dont actually need to know the exact distance. just find the lowest
	//of x^2 + y^2 to avoid sqrt

	var tuples []struct {
		ind  int
		dist float64
	}

	for i, v := range c.enemies {
		if i == excl {
			continue
		}
		vx, vy := v.Shape().Pos()
		dist := math.Pow(x-vx, 2) + math.Pow(y-vy, 2)
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
func (c *Handler) EnemiesWithinRadius(x, y, r float64) []int {
	result := make([]int, 0, len(c.enemies))
	for i, v := range c.enemies {
		vx, vy := v.Shape().Pos()
		dist := math.Pow(x-vx, 2) + math.Pow(y-vy, 2)
		if dist > r {
			continue
		}
		result = append(result, i)
	}

	return result
}

// EnemyExcl returns array of indices of enemies, excluding self
func (c *Handler) EnemyExcl(self int) []int {
	result := make([]int, 0, len(c.enemies))

	for i := range c.enemies {
		if i == self {
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
