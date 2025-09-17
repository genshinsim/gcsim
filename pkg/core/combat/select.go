package combat

import (
	"sort"

	"github.com/genshinsim/gcsim/pkg/core/info"
)

// all targets

func enemiesWithinAreaFiltered(a info.AttackPattern, filter func(t info.Enemy) bool, originalEnemies []info.Target) []info.Enemy {
	var enemies []info.Enemy
	hasFilter := filter != nil
	for _, v := range originalEnemies {
		e, ok := v.(info.Enemy)
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

func gadgetsWithinAreaFiltered(a info.AttackPattern, filter func(t info.Gadget) bool, originalGadgets []info.Gadget) []info.Gadget {
	var gadgets []info.Gadget
	hasFilter := filter != nil
	for _, v := range originalGadgets {
		if v == nil {
			continue
		}
		// check if info.Gadget is enemy camp, abilities don't target allied gadgets
		if v.GadgetTyp() <= info.StartGadgetTypEnemy || v.GadgetTyp() >= info.EndGadgetTypEnemy {
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
func (h *Handler) EnemiesWithinArea(a info.AttackPattern, filter func(t info.Enemy) bool) []info.Enemy {
	enemies := enemiesWithinAreaFiltered(a, filter, h.enemies)
	if len(enemies) == 0 {
		return nil
	}
	return enemies
}

// returns gadgets within the given area, no sorting, pass nil for no filter
func (h *Handler) GadgetsWithinArea(a info.AttackPattern, filter func(t info.Gadget) bool) []info.Gadget {
	gadgets := gadgetsWithinAreaFiltered(a, filter, h.gadgets)
	if len(gadgets) == 0 {
		return nil
	}
	return gadgets
}

// random targets

// returns a random enemy within the given area, pass nil for no filter
func (h *Handler) RandomEnemyWithinArea(a info.AttackPattern, filter func(t info.Enemy) bool) info.Enemy {
	enemies := h.EnemiesWithinArea(a, filter)
	if enemies == nil {
		return nil
	}
	return enemies[h.Rand.Intn(len(enemies))]
}

// returns a random info.Gadget within the given area, pass nil for no filter
func (h *Handler) RandomGadgetWithinArea(a info.AttackPattern, filter func(t info.Gadget) bool) info.Gadget {
	gadgets := h.GadgetsWithinArea(a, filter)
	if gadgets == nil {
		return nil
	}
	return gadgets[h.Rand.Intn(len(gadgets))]
}

// returns a list of random enemies within the given area, pass nil for no filter
func (h *Handler) RandomEnemiesWithinArea(a info.AttackPattern, filter func(t info.Enemy) bool, maxCount int) []info.Enemy {
	enemies := h.EnemiesWithinArea(a, filter)
	if enemies == nil {
		return nil
	}
	enemyCount := len(enemies)

	// generate random indexes to take from enemies (no duplicates!)
	indexes := h.Rand.Perm(enemyCount)

	// determine length of slice to return
	count := min(enemyCount, maxCount)

	// add enemies given by indexes to the result
	result := make([]info.Enemy, 0, count)
	for i := range count {
		result = append(result, enemies[indexes[i]])
	}
	return result
}

// returns a list of random gadgets within the given area, pass nil for no filter
func (h *Handler) RandomGadgetsWithinArea(a info.AttackPattern, filter func(t info.Gadget) bool, maxCount int) []info.Gadget {
	gadgets := h.GadgetsWithinArea(a, filter)
	if gadgets == nil {
		return nil
	}
	gadgetCount := len(gadgets)

	// generate random indexes to take from gadgets (no duplicates!)
	indexes := h.Rand.Perm(gadgetCount)

	// determine length of slice to return
	count := min(gadgetCount, maxCount)

	// add gadgets given by indexes to the result
	result := make([]info.Gadget, 0, count)
	for i := range count {
		result = append(result, gadgets[indexes[i]])
	}
	return result
}

// closest targets

type enemyTuple struct {
	enemy info.Enemy
	dist  float64
}

func enemiesWithinAreaSorted(a info.AttackPattern, filter func(t info.Enemy) bool, skipAttackPattern bool, originalEnemies []info.Target) []enemyTuple {
	var enemies []enemyTuple

	hasFilter := filter != nil
	for _, v := range originalEnemies {
		e, ok := v.(info.Enemy)
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
	Gadget info.Gadget
	dist   float64
}

func gadgetsWithinAreaSorted(a info.AttackPattern, filter func(t info.Gadget) bool, skipAttackPattern bool, originalGadgets []info.Gadget) []gadgetTuple {
	var gadgets []gadgetTuple

	hasFilter := filter != nil
	for _, v := range originalGadgets {
		if v == nil {
			continue
		}
		// check if info.Gadget is enemy camp, abilities don't target allied gadgets
		if v.GadgetTyp() <= info.StartGadgetTypEnemy || v.GadgetTyp() >= info.EndGadgetTypEnemy {
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
		gadgets = append(gadgets, gadgetTuple{Gadget: v, dist: a.Shape.Pos().Sub(v.Pos()).MagnitudeSquared()})
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
func (h *Handler) ClosestEnemy(pos info.Point) info.Enemy {
	enemies := enemiesWithinAreaSorted(NewCircleHitOnTarget(pos, nil, 1), nil, true, h.enemies)
	if enemies == nil {
		return nil
	}
	return enemies[0].enemy
}

// returns the closest info.Gadget to the given position without any range restrictions; SHOULD NOT be used outside of pkg
func (h *Handler) ClosestGadget(pos info.Point) info.Gadget {
	gadgets := gadgetsWithinAreaSorted(NewCircleHitOnTarget(pos, nil, 1), nil, true, h.gadgets)
	if gadgets == nil {
		return nil
	}
	return gadgets[0].Gadget
}

// returns the closest enemy within the given area, pass nil for no filter
func (h *Handler) ClosestEnemyWithinArea(a info.AttackPattern, filter func(t info.Enemy) bool) info.Enemy {
	enemies := enemiesWithinAreaSorted(a, filter, false, h.enemies)
	if enemies == nil {
		return nil
	}
	return enemies[0].enemy
}

// returns the closest info.Gadget within the given area, pass nil for no filter
func (h *Handler) ClosestGadgetWithinArea(a info.AttackPattern, filter func(t info.Gadget) bool) info.Gadget {
	gadgets := gadgetsWithinAreaSorted(a, filter, false, h.gadgets)
	if gadgets == nil {
		return nil
	}
	return gadgets[0].Gadget
}

// returns enemies within the given area, sorted from closest to furthest, pass nil for no filter
func (h *Handler) ClosestEnemiesWithinArea(a info.AttackPattern, filter func(t info.Enemy) bool) []info.Enemy {
	enemies := enemiesWithinAreaSorted(a, filter, false, h.enemies)
	if enemies == nil {
		return nil
	}

	result := make([]info.Enemy, 0, len(enemies))
	for _, v := range enemies {
		result = append(result, v.enemy)
	}
	return result
}

// returns enemies within the given area, sorted from closest to furthest, pass nil for no filter
func (h *Handler) ClosestGadgetsWithinArea(a info.AttackPattern, filter func(t info.Gadget) bool) []info.Gadget {
	gadgets := gadgetsWithinAreaSorted(a, filter, false, h.gadgets)
	if gadgets == nil {
		return nil
	}

	result := make([]info.Gadget, 0, len(gadgets))
	for _, v := range gadgets {
		result = append(result, v.Gadget)
	}
	return result
}
