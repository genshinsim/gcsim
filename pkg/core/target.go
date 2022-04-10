package core

import (
	"math"
	"sort"
)

type Target interface {
	//basic info
	Type() TargettableType //type of target
	Index() int            //should correspond to index
	SetIndex(ind int)      //update the current index
	MaxHP() float64
	HP() float64

	//collision detection
	Shape() Shape

	//attacks
	Attack(*AttackEvent, LogEvent) (float64, bool)

	//reaction/aura stuff
	Tick()
	AuraContains(...EleType) bool
	AuraType() EleType

	//tags
	SetTag(key string, val int)
	GetTag(key string) int
	RemoveTag(key string)

	//target mods
	AddDefMod(key string, val float64, dur int)
	AddResMod(key string, val ResistMod)
	RemoveResMod(key string)
	RemoveDefMod(key string)
	HasDefMod(key string) bool
	HasResMod(key string) bool

	//getting rid of
	Kill()
}

type TargettableType int

const (
	TargettableEnemy TargettableType = iota
	TargettablePlayer
	TargettableObject
	TargettableTypeCount
)

// type TargetEnemy interface {
// 	Index() int
// 	SetIndex(ind int) //update the current index
// 	MaxHP() float64
// 	HP() float64
// 	//aura/reactions
// 	AuraType() EleType
// 	AuraContains(e ...EleType) bool
// 	Tick() //this should happen first before task ticks

// 	//attacks
// 	Attack(ds *Snapshot) (float64, bool)

// 	AddDefMod(key string, val float64, dur int)
// 	AddResMod(key string, val ResistMod)
// 	RemoveResMod(key string)
// 	RemoveDefMod(key string)
// 	HasDefMod(key string) bool
// 	HasResMod(key string) bool

// 	Delete() //gracefully deference everything so that it can be gc'd
// }

type ResistMod struct {
	Key      string
	Ele      EleType
	Value    float64
	Duration int
	Expiry   int
	Event    LogEvent
}

type DefMod struct {
	Key    string
	Value  float64
	Expiry int
	Event  LogEvent
}

// func (c *Core) ReindexTargets() {
// 	//wipe out nil entries
// 	n := 0
// 	for _, v := range c.Targets {
// 		if v != nil {
// 			c.Targets[n] = v
// 			c.Targets[n].SetIndex(n)
// 			n++
// 		}
// 	}
// 	c.Targets = c.Targets[:n]
// }

func (c *Core) AddTarget(t Target) {
	c.Targets = append(c.Targets, t)
	t.SetIndex(len(c.Targets) - 1)
	c.Events.Emit(OnTargetAdded, t)
}

func (c *Core) RemoveTarget(i int) {
	//can't remove player
	if i <= 0 || i >= len(c.Targets) {
		return
	}
	c.Targets[i] = nil
	// c.Targets = c.Targets[:len(c.Targets)-1]
	// c.ReindexTargets()
}

//EnemeyByDistance returns an array of indices of the enemies sorted by distance
func (c *Core) EnemyByDistance(x, y float64, excl int) []int {
	//we dont actually need to know the exact distance. just find the lowest
	//of x^2 + y^2 to avoid sqrt

	var tuples []struct {
		ind  int
		dist float64
	}

	for i, v := range c.Targets {
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

//EnemyExcl returns array of indices of enemies, excluding self
func (c *Core) EnemyExcl(self int) []int {
	result := make([]int, 0, len(c.Targets))

	for i, v := range c.Targets {
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

func (c *Core) RandomEnemyTarget() int {

	count := 0
	for _, v := range c.Targets {
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
	for i, v := range c.Targets {
		if v.Type() == TargettableEnemy {
			if n == count {
				return i
			}
			count++
		}
	}
	panic("no random target found?? should not happen")
}

func (c *Core) RandomTargetIndex(typ TargettableType) int {
	count := 0
	for _, v := range c.Targets {
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
	for i, v := range c.Targets {
		if v.Type() == typ {
			if n == count {
				return i
			}
			count++
		}
	}
	panic("no random target found?? should not happen")
}
