package combat

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

type TargetKey int

const InvalidTargetKey TargetKey = -1

type Target interface {
	Key() TargetKey        // unique key for the target
	SetKey(k TargetKey)    // update key
	Type() TargettableType // type of target
	Shape() Shape          // shape of target
	Pos() Point            // center of target
	SetPos(p Point)        // move target
	IsAlive() bool
	SetTag(key string, val int)
	GetTag(key string) int
	RemoveTag(key string)
	HandleAttack(*AttackEvent) float64
	AttackWillLand(a AttackPattern) (bool, string) // hurtbox collides with AttackPattern
	IsWithinArea(a AttackPattern) bool             // center is in AttackPattern
	Tick()                                         // called every tick
	Kill()
	// for collision check
	CollidableWith(TargettableType) bool
	CollidedWith(t Target)
	WillCollide(Shape) bool
	// direction related
	Direction() Point                  // returns viewing direction as a Point
	SetDirection(trg Point)            // calculates viewing direction relative to default direction (0, 1)
	SetDirectionToClosestEnemy()       // looks for closest enemy
	CalcTempDirection(trg Point) Point // used for stuff like Bow CA
}

type TargetWithAura interface {
	Target
	AuraContains(e ...attributes.Element) bool
}

type TargettableType int

const (
	TargettableEnemy TargettableType = iota
	TargettablePlayer
	TargettableGadget
	TargettableTypeCount
)
