package combat

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

type Target interface {
	Key() targets.TargetKey        // unique key for the target
	SetKey(k targets.TargetKey)    // update key
	Type() targets.TargettableType // type of target
	Shape() Shape                  // shape of target
	Pos() Point                    // center of target
	SetPos(p Point)                // move target
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
	CollidableWith(targets.TargettableType) bool
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
