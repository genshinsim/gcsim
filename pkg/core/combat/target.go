package combat

import "github.com/genshinsim/gcsim/pkg/core/glog"

type Target interface {
	Index() int         //should correspond to index
	SetIndex(index int) //update the current index
	MaxHP() float64
	HP() float64

	Type() TargettableType   //type of target
	Shape() Shape            // shape of target
	Pos() (float64, float64) // center of target
	SetPos(x, y float64)     // move target

	//apply attack to target
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
