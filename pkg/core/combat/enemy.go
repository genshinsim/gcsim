package combat

import (
	"sort"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
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

type enemyTuple struct {
	enemy Enemy
	dist  float64
}

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

// area check

// checks whether the given target is within the given area
func TargetIsWithinArea(target Target, a AttackPattern) bool {
	collision, _ := target.AttackWillLand(a)
	return collision
}

// all enemies

// returns enemies within the given area, no sorting, pass nil for no filter
func (c *Handler) EnemiesWithinArea(a AttackPattern, filter func(t Enemy) bool) []Enemy {
	var enemies []Enemy

	hasFilter := filter != nil
	for _, v := range c.enemies {
		e, ok := v.(Enemy)
		if !ok {
			panic("c.enemies should contain targets that implement the Enemy interface")
		}
		if hasFilter && !filter(e) {
			continue
		}
		if !v.IsAlive() {
			continue
		}
		if !TargetIsWithinArea(e, a) {
			continue
		}
		enemies = append(enemies, e)
	}

	if len(enemies) == 0 {
		return nil
	}

	return enemies
}

// random enemies

// returns a random enemy within the given area, pass nil for no filter
func (c *Handler) RandomEnemyWithinArea(a AttackPattern, filter func(t Enemy) bool) Enemy {
	enemies := c.EnemiesWithinArea(a, filter)
	if enemies == nil {
		return nil
	}
	return enemies[c.Rand.Intn(len(enemies))]
}

// returns a list of random enemies within the given area, pass nil for no filter
func (c *Handler) RandomEnemiesWithinArea(a AttackPattern, filter func(t Enemy) bool, maxCount int) []Enemy {
	enemies := c.EnemiesWithinArea(a, filter)
	if enemies == nil {
		return nil
	}
	enemyCount := len(enemies)

	// generate random indexes to take from enemies (no duplicates!)
	indexes := c.Rand.Perm(enemyCount)

	// determine length of slice to return
	count := maxCount
	if enemyCount < maxCount {
		count = enemyCount
	}

	// add enemies given by indexes to the result
	result := make([]Enemy, 0, count)
	for i := 0; i < count; i++ {
		result = append(result, enemies[indexes[i]])
	}
	return result
}

// closest enemies

func (c *Handler) getEnemiesWithinAreaSorted(a AttackPattern, filter func(t Enemy) bool, skipAttackPattern bool) []enemyTuple {
	var enemies []enemyTuple

	hasFilter := filter != nil
	for _, v := range c.enemies {
		e, ok := v.(Enemy)
		if !ok {
			panic("c.enemies should contain targets that implement the Enemy interface")
		}
		if hasFilter && !filter(e) {
			continue
		}
		if !e.IsAlive() {
			continue
		}
		if !skipAttackPattern && !TargetIsWithinArea(e, a) {
			continue
		}
		enemies = append(enemies, enemyTuple{enemy: e, dist: a.Shape.Pos().Sub(e.Pos()).MagnitudeSquared()})
	}

	if len(enemies) == 0 {
		return nil
	}

	sort.Slice(enemies, func(i, j int) bool {
		return enemies[i].dist < enemies[j].dist
	})

	return enemies
}

// returns the closest enemy to the given position without any range restrictions; SHOULD NOT be used outside of pkg
func (c *Handler) ClosestEnemy(pos Point) Enemy {
	enemies := c.getEnemiesWithinAreaSorted(NewCircleHitOnTarget(pos, nil, 1), nil, true)
	if enemies == nil {
		return nil
	}
	return enemies[0].enemy
}

// returns the closest enemy within the given area, pass nil for no filter
func (c *Handler) ClosestEnemyWithinArea(a AttackPattern, filter func(t Enemy) bool) Enemy {
	enemies := c.getEnemiesWithinAreaSorted(a, filter, false)
	if enemies == nil {
		return nil
	}
	return enemies[0].enemy
}

// returns enemies within the given area, sorted from closest to furthest, pass nil for no filter
func (c *Handler) ClosestEnemiesWithinArea(a AttackPattern, filter func(t Enemy) bool) []Enemy {
	enemies := c.getEnemiesWithinAreaSorted(a, filter, false)
	if enemies == nil {
		return nil
	}

	result := make([]Enemy, 0, len(enemies))
	for _, v := range enemies {
		result = append(result, v.enemy)
	}
	return result
}
