package combat

import (
	"math"
	"sort"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

type Target interface {
	Index() int              //should correspond to index
	SetIndex(index int)      //update the current index
	Type() TargettableType   //type of target
	Shape() Shape            // shape of target
	Pos() (float64, float64) // center of target
	SetPos(x, y float64)     // move target
	IsAlive() bool
	AttackWillLand(a AttackPattern, src int) (bool, string)
	Attack(*AttackEvent, glog.Event) (float64, bool)
	Tick() //called every tick
	Kill()
}

type TargetWithAura interface {
	Target
	AuraContains(e ...attributes.Element) bool
}

type TargettableType int

const (
	TargettableEnemy TargettableType = iota
	TargettablePlayer
	TargettableObject
	TargettableTypeCount
)

// EnemyByDistance returns an array of indices of the enemies sorted by distance
func (c *Handler) EnemyByDistance(x, y float64, excl int) []int {
	//we dont actually need to know the exact distance. just find the lowest
	//of x^2 + y^2 to avoid sqrt

	var tuples []struct {
		ind  int
		dist float64
	}

	for i, v := range c.targets {
		if v.Type() != TargettableEnemy {
			continue
		}
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
	result := make([]int, 0, len(c.targets))
	for i, v := range c.targets {
		if v.Type() != TargettableEnemy {
			continue
		}
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
	result := make([]int, 0, len(c.targets))

	for i, v := range c.targets {
		if v.Type() != TargettableEnemy {
			continue
		}
		if i == self {
			continue
		}

		result = append(result, i)
	}

	return result
}

func (c *Handler) RandomEnemyTarget() int {

	count := 0
	for _, v := range c.targets {
		if v.Type() == TargettableEnemy {
			count++
		}
	}
	if count == 0 {
		//this will basically cause that attack to hit nothing
		return -1
	}
	n := c.Rand.Intn(count)
	count = 0
	for i, v := range c.targets {
		if v.Type() == TargettableEnemy {
			if n == count {
				return i
			}
			count++
		}
	}
	panic("no random target found?? should not happen")
}

func (c *Handler) RandomTargetIndex(typ TargettableType) int {
	count := 0
	for _, v := range c.targets {
		if v.Type() == typ {
			count++
		}
	}
	if count == 0 {
		//this will basically cause that attack to hit nothing
		return -1
	}
	n := c.Rand.Intn(count)
	count = 0
	for i, v := range c.targets {
		if v.Type() == typ {
			if n == count {
				return i
			}
			count++
		}
	}
	panic("no random target found?? should not happen")
}
