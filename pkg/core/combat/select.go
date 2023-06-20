package combat

import (
	"sort"

	"github.com/genshinsim/gcsim/pkg/core/geometry"
)

// all targets

func enemiesWithinAreaFiltered(a AttackPattern, filter func(t Enemy) bool, originalEnemies []Target) []Enemy {
	var enemies []Enemy
	hasFilter := filter != nil
	for _, v := range originalEnemies {
		e, ok := v.(Enemy)
		if !ok {
			panic("enemies should contain targets that implement the Enemy interface")
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
	return enemies
}

func gadgetsWithinAreaFiltered(a AttackPattern, filter func(t Gadget) bool, originalGadgets []Gadget) []Gadget {
	var gadgets []Gadget
	hasFilter := filter != nil
	for _, v := range originalGadgets {
		if v == nil {
			continue
		}
		// check if gadget is enemy camp, abilities don't target allied gadgets
		if !(v.GadgetTyp() > StartGadgetTypEnemy && v.GadgetTyp() < EndGadgetTypEnemy) {
			continue
		}
		if hasFilter && !filter(v) {
			continue
		}
		if !v.IsAlive() {
			continue
		}
		if !v.IsWithinArea(a) {
			continue
		}
		gadgets = append(gadgets, v)
	}
	return gadgets
}

// returns enemies within the given area, no sorting, pass nil for no filter
func (c *Handler) EnemiesWithinArea(a AttackPattern, filter func(t Enemy) bool) []Enemy {
	enemies := enemiesWithinAreaFiltered(a, filter, c.enemies)
	if len(enemies) == 0 {
		return nil
	}
	return enemies
}

// returns gadgets within the given area, no sorting, pass nil for no filter
func (c *Handler) GadgetsWithinArea(a AttackPattern, filter func(t Gadget) bool) []Gadget {
	gadgets := gadgetsWithinAreaFiltered(a, filter, c.gadgets)
	if len(gadgets) == 0 {
		return nil
	}
	return gadgets
}

// random targets

// returns a random enemy within the given area, pass nil for no filter
func (c *Handler) RandomEnemyWithinArea(a AttackPattern, filter func(t Enemy) bool) Enemy {
	enemies := c.EnemiesWithinArea(a, filter)
	if enemies == nil {
		return nil
	}
	return enemies[c.Rand.Intn(len(enemies))]
}

// returns a random gadget within the given area, pass nil for no filter
func (c *Handler) RandomGadgetWithinArea(a AttackPattern, filter func(t Gadget) bool) Gadget {
	gadgets := c.GadgetsWithinArea(a, filter)
	if gadgets == nil {
		return nil
	}
	return gadgets[c.Rand.Intn(len(gadgets))]
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

// returns a list of random gadgets within the given area, pass nil for no filter
func (c *Handler) RandomGadgetsWithinArea(a AttackPattern, filter func(t Gadget) bool, maxCount int) []Gadget {
	gadgets := c.GadgetsWithinArea(a, filter)
	if gadgets == nil {
		return nil
	}
	gadgetCount := len(gadgets)

	// generate random indexes to take from gadgets (no duplicates!)
	indexes := c.Rand.Perm(gadgetCount)

	// determine length of slice to return
	count := maxCount
	if gadgetCount < maxCount {
		count = gadgetCount
	}

	// add gadgets given by indexes to the result
	result := make([]Gadget, 0, count)
	for i := 0; i < count; i++ {
		result = append(result, gadgets[indexes[i]])
	}
	return result
}

// closest targets

type enemyTuple struct {
	enemy Enemy
	dist  float64
}

func enemiesWithinAreaSorted(a AttackPattern, filter func(t Enemy) bool, skipAttackPattern bool, originalEnemies []Target) []enemyTuple {
	var enemies []enemyTuple

	hasFilter := filter != nil
	for _, v := range originalEnemies {
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

type gadgetTuple struct {
	gadget Gadget
	dist   float64
}

func gadgetsWithinAreaSorted(a AttackPattern, filter func(t Gadget) bool, skipAttackPattern bool, originalGadgets []Gadget) []gadgetTuple {
	var gadgets []gadgetTuple

	hasFilter := filter != nil
	for _, v := range originalGadgets {
		if v == nil {
			continue
		}
		// check if gadget is enemy camp, abilities don't target allied gadgets
		if !(v.GadgetTyp() > StartGadgetTypEnemy && v.GadgetTyp() < EndGadgetTypEnemy) {
			continue
		}
		if hasFilter && !filter(v) {
			continue
		}
		if !v.IsAlive() {
			continue
		}
		if !skipAttackPattern && !v.IsWithinArea(a) {
			continue
		}
		gadgets = append(gadgets, gadgetTuple{gadget: v, dist: a.Shape.Pos().Sub(v.Pos()).MagnitudeSquared()})
	}

	if len(gadgets) == 0 {
		return nil
	}

	sort.Slice(gadgets, func(i, j int) bool {
		return gadgets[i].dist < gadgets[j].dist
	})

	return gadgets
}

// returns the closest enemy to the given position without any range restrictions; SHOULD NOT be used outside of pkg
func (c *Handler) ClosestEnemy(pos geometry.Point) Enemy {
	enemies := enemiesWithinAreaSorted(NewCircleHitOnTarget(pos, nil, 1), nil, true, c.enemies)
	if enemies == nil {
		return nil
	}
	return enemies[0].enemy
}

// returns the closest gadget to the given position without any range restrictions; SHOULD NOT be used outside of pkg
func (c *Handler) ClosestGadget(pos geometry.Point) Gadget {
	gadgets := gadgetsWithinAreaSorted(NewCircleHitOnTarget(pos, nil, 1), nil, true, c.gadgets)
	if gadgets == nil {
		return nil
	}
	return gadgets[0].gadget
}

// returns the closest enemy within the given area, pass nil for no filter
func (c *Handler) ClosestEnemyWithinArea(a AttackPattern, filter func(t Enemy) bool) Enemy {
	enemies := enemiesWithinAreaSorted(a, filter, false, c.enemies)
	if enemies == nil {
		return nil
	}
	return enemies[0].enemy
}

// returns the closest gadget within the given area, pass nil for no filter
func (c *Handler) ClosestGadgetWithinArea(a AttackPattern, filter func(t Gadget) bool) Gadget {
	gadgets := gadgetsWithinAreaSorted(a, filter, false, c.gadgets)
	if gadgets == nil {
		return nil
	}
	return gadgets[0].gadget
}

// returns enemies within the given area, sorted from closest to furthest, pass nil for no filter
func (c *Handler) ClosestEnemiesWithinArea(a AttackPattern, filter func(t Enemy) bool) []Enemy {
	enemies := enemiesWithinAreaSorted(a, filter, false, c.enemies)
	if enemies == nil {
		return nil
	}

	result := make([]Enemy, 0, len(enemies))
	for _, v := range enemies {
		result = append(result, v.enemy)
	}
	return result
}

// returns enemies within the given area, sorted from closest to furthest, pass nil for no filter
func (c *Handler) ClosestGadgetsWithinArea(a AttackPattern, filter func(t Gadget) bool) []Gadget {
	gadgets := gadgetsWithinAreaSorted(a, filter, false, c.gadgets)
	if gadgets == nil {
		return nil
	}

	result := make([]Gadget, 0, len(gadgets))
	for _, v := range gadgets {
		result = append(result, v.gadget)
	}
	return result
}
