package combat

import "github.com/genshinsim/gcsim/pkg/core/glog"

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
	Attack(*AttackEvent, glog.Event) (float64, bool)

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
