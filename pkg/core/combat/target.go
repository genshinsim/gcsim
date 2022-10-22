package combat

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

type TargetKey int

const InvalidTargetKey TargetKey = -1

type Target interface {
	Index() int              //should correspond to index
	SetIndex(index int)      //update the current index
	Key() TargetKey          //unique key for the target
	SetKey(k TargetKey)      //update key
	Type() TargettableType   //type of target
	Shape() Shape            // shape of target
	Pos() (float64, float64) // center of target
	SetPos(x, y float64)     // move target
	IsAlive() bool
	HandleAttack(*AttackEvent) float64
	AttackWillLand(a AttackPattern) (bool, string)
	Attack(*AttackEvent, glog.Event) (float64, bool)
	ApplyDamage(*AttackEvent, float64)
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
