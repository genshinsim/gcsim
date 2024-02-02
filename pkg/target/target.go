package target

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

const MaxTeamSize = 4

type Target struct {
	Core            *core.Core
	key             targets.TargetKey
	Hitbox          geometry.Circle
	Tags            map[string]int
	CollidableTypes [targets.TargettableTypeCount]bool
	OnCollision     func(combat.Target)

	Alive bool

	// icd related
	icdTagOnTimer       [MaxTeamSize][attacks.ICDTagLength]bool
	icdTagCounter       [MaxTeamSize][attacks.ICDTagLength]int
	icdDamageTagOnTimer [MaxTeamSize][attacks.ICDTagLength]bool
	icdDamageTagCounter [MaxTeamSize][attacks.ICDTagLength]int

	direction geometry.Point
}

func New(core *core.Core, p geometry.Point, r float64) *Target {
	t := &Target{
		Core: core,
	}
	t.Hitbox = *geometry.NewCircle(p, r, geometry.DefaultDirection(), 360)
	t.direction = geometry.DefaultDirection()
	t.Tags = make(map[string]int)
	t.Alive = true

	return t
}

func (t *Target) Collidable() bool                              { return t.OnCollision != nil }
func (t *Target) CollidableWith(x targets.TargettableType) bool { return t.CollidableTypes[x] }
func (t *Target) CollidedWith(x combat.Target) {
	if t.OnCollision != nil {
		t.OnCollision(x)
	}
}

func (t *Target) Key() targets.TargetKey     { return t.key }
func (t *Target) SetKey(x targets.TargetKey) { t.key = x }
func (t *Target) Shape() geometry.Shape      { return &t.Hitbox }
func (t *Target) SetPos(p geometry.Point)    { t.Hitbox.SetPos(p) }
func (t *Target) Pos() geometry.Point        { return t.Hitbox.Pos() }
func (t *Target) Kill()                      { t.Alive = false }
func (t *Target) IsAlive() bool              { return t.Alive }
func (t *Target) SetTag(key string, val int) {
	t.Tags[key] = val
}

func (t *Target) GetTag(key string) int {
	return t.Tags[key]
}

func (t *Target) RemoveTag(key string) {
	delete(t.Tags, key)
}

func (t *Target) WillCollide(s geometry.Shape) bool {
	if !t.Alive {
		return false
	}
	switch v := s.(type) {
	case *geometry.Circle:
		return t.Shape().IntersectCircle(*v)
	case *geometry.Rectangle:
		return t.Shape().IntersectRectangle(*v)
	default:
		return false
	}
}

func (t *Target) AttackWillLand(a combat.AttackPattern) (bool, string) {
	// shape shouldn't be nil; panic here
	if a.Shape == nil {
		panic("unexpected nil shape")
	}
	if !t.Alive {
		return false, "target dead"
	}
	// shape can't be nil now, check if type matches
	// if !a.Targets[t.typ] {
	// 	return false, "wrong type"
	// }
	// swirl aoe shouldn't hit the src of the aoe
	for _, v := range a.IgnoredKeys {
		if t.Key() == v {
			return false, "no self harm"
		}
	}

	// check if shape matches
	switch v := a.Shape.(type) {
	case *geometry.Circle:
		return t.Shape().IntersectCircle(*v), "intersect circle"
	case *geometry.Rectangle:
		return t.Shape().IntersectRectangle(*v), "intersect rectangle"
	case *geometry.SingleTarget:
		// only true if
		return v.Target == t.key, "target"
	default:
		return false, "unknown shape"
	}
}

func (t *Target) IsWithinArea(a combat.AttackPattern) bool {
	return a.Shape.PointInShape(t.Pos())
}

func (t *Target) Direction() geometry.Point { return t.direction }
func (t *Target) SetDirection(trg geometry.Point) {
	src := t.Pos()
	t.direction = geometry.CalcDirection(src, trg)
	t.Core.Combat.Log.NewEvent("set target direction", glog.LogDebugEvent, -1).
		Write("src target key", t.key).
		Write("srcX", src.X).
		Write("srcY", src.Y).
		Write("trgX", trg.X).
		Write("trgY", trg.Y).
		Write("direction", t.direction)
}
func (t *Target) SetDirectionToClosestEnemy() {
	src := t.Pos()
	// calculate direction towards closest enemy, or forward direction if none
	enemy := t.Core.Combat.ClosestEnemy(src)
	if enemy == nil {
		t.direction = geometry.DefaultDirection()
		return
	}
	t.SetDirection(enemy.Pos())
	t.Core.Combat.Log.NewEvent("set target direction to closest enemy", glog.LogDebugEvent, -1).
		Write("enemy key", enemy.Key()).
		Write("direction", t.direction)
}
func (t *Target) CalcTempDirection(trg geometry.Point) geometry.Point {
	src := t.Pos()
	direction := geometry.CalcDirection(src, trg)
	t.Core.Combat.Log.NewEvent("using temporary target direction", glog.LogDebugEvent, -1).
		Write("src target key", t.key).
		Write("srcX", src.X).
		Write("srcY", src.Y).
		Write("trgX", trg.X).
		Write("trgY", trg.Y).
		Write("direction", t.direction).
		Write("temporary direction", direction)
	return direction
}
