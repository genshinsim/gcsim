package target

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

type Target struct {
	Core            *core.Core
	TargetIndex     int
	key             combat.TargetKey
	Hitbox          combat.Circle
	Tags            map[string]int
	CollidableTypes [combat.TargettableTypeCount]bool
	OnCollision     func(combat.Target)

	Alive bool
}

func New(core *core.Core, x, y, z float64) *Target {
	t := &Target{
		Core: core,
	}
	t.Hitbox = *combat.NewCircle(x, y, z)
	t.Tags = make(map[string]int)
	t.Alive = true

	return t
}

func (t *Target) Collidable() bool                             { return t.OnCollision != nil }
func (t *Target) CollidableWith(x combat.TargettableType) bool { return t.CollidableTypes[x] }
func (t *Target) CollidedWith(x combat.Target) {
	if t.OnCollision != nil {
		t.OnCollision(x)
	}
}

func (t *Target) Key() combat.TargetKey     { return t.key }
func (t *Target) SetKey(x combat.TargetKey) { t.key = x }
func (t *Target) Index() int                { return t.TargetIndex }
func (t *Target) SetIndex(ind int)          { t.TargetIndex = ind }
func (t *Target) Shape() combat.Shape       { return &t.Hitbox }
func (t *Target) SetPos(x, y float64)       { t.Hitbox.SetPos(x, y) }
func (t *Target) Pos() (float64, float64)   { return t.Hitbox.Pos() }
func (t *Target) Kill()                     { t.Alive = false }
func (t *Target) IsAlive() bool             { return t.Alive }
func (t *Target) SetTag(key string, val int) {
	t.Tags[key] = val
}

func (t *Target) GetTag(key string) int {
	return t.Tags[key]
}

func (t *Target) RemoveTag(key string) {
	delete(t.Tags, key)
}

func (t *Target) WillCollide(s combat.Shape) bool {
	if !t.Alive {
		return false
	}
	switch v := s.(type) {
	case *combat.Circle:
		return t.Shape().IntersectCircle(*v)
	case *combat.Rectangle:
		return t.Shape().IntersectRectangle(*v)
	default:
		return false
	}
}

func (t *Target) AttackWillLand(a combat.AttackPattern, src combat.TargetKey) (bool, string) {
	//shape shouldn't be nil; panic here
	if a.Shape == nil {
		panic("unexpected nil shape")
	}
	if !t.Alive {
		return false, "target dead"
	}
	//shape can't be nil now, check if type matches
	// if !a.Targets[t.typ] {
	// 	return false, "wrong type"
	// }
	//skip if self harm is false and dmg src == i
	if !a.SelfHarm && src == t.key {
		return false, "no self harm"
	}

	//check if shape matches
	switch v := a.Shape.(type) {
	case *combat.Circle:
		return t.Shape().IntersectCircle(*v), "intersect circle"
	case *combat.Rectangle:
		return t.Shape().IntersectRectangle(*v), "intersect rectangle"
	case *combat.SingleTarget:
		//only true if
		return v.Target == t.key, "target"
	default:
		return false, "unknown shape"
	}
}
