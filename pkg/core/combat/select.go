package combat

import (
	"sort"

	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

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
		if !e.IsWithinArea(a) {
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
		if !skipAttackPattern && !e.IsWithinArea(a) {
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
func (c *Handler) ClosestEnemy(pos geometry.Point) Enemy {
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
