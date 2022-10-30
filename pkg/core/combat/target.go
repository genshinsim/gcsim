package combat

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

type TargetKey int

const InvalidTargetKey TargetKey = -1

type Target interface {
	Key() TargetKey          //unique key for the target
	SetKey(k TargetKey)      //update key
	Type() TargettableType   //type of target
	Shape() Shape            // shape of target
	Pos() (float64, float64) // center of target
	SetPos(x, y float64)     // move target
	IsAlive() bool
	HandleAttack(*AttackEvent) float64
	AttackWillLand(a AttackPattern) (bool, string)
	Tick() //called every tick
	Kill()
	//for collision check
	CollidableWith(TargettableType) bool
	CollidedWith(t Target)
	WillCollide(Shape) bool
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
