package target

import (
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

const MaxTeamSize = 4

type Target struct {
	Core            *core.Core
	key             combat.TargetKey
	Hitbox          combat.Circle
	Tags            map[string]int
	CollidableTypes [combat.TargettableTypeCount]bool
	OnCollision     func(combat.Target)

	Alive bool

	//icd related
	icdTagOnTimer       [MaxTeamSize][combat.ICDTagLength]bool
	icdTagCounter       [MaxTeamSize][combat.ICDTagLength]int
	icdDamageTagOnTimer [MaxTeamSize][combat.ICDTagLength]bool
	icdDamageTagCounter [MaxTeamSize][combat.ICDTagLength]int

	direction combat.Point
}

func New(core *core.Core, p combat.Point, r float64) *Target {
	t := &Target{
		Core: core,
	}
	t.Hitbox = *combat.NewCircle(p, r, combat.DefaultDirection(), 360)
	t.direction = combat.DefaultDirection()
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
func (t *Target) Shape() combat.Shape       { return &t.Hitbox }
func (t *Target) SetPos(p combat.Point)     { t.Hitbox.SetPos(p) }
func (t *Target) Pos() combat.Point         { return t.Hitbox.Pos() }
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

func (t *Target) AttackWillLand(a combat.AttackPattern) (bool, string) {
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
	// swirl aoe shouldn't hit the src of the aoe
	for _, v := range a.IgnoredKeys {
		if t.Key() == v {
			return false, "no self harm"
		}
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

func (t *Target) Direction() combat.Point { return t.direction }
func (t *Target) SetDirection(trg combat.Point) {
	src := t.Pos()
	t.direction = combat.CalcDirection(src, trg)
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
	enemies := t.Core.Combat.EnemyByDistance(src, combat.InvalidTargetKey)
	if len(enemies) == 0 {
		t.direction = combat.DefaultDirection()
		return
	}

	enemy := t.Core.Combat.Enemy(enemies[0])
	t.SetDirection(enemy.Pos())
	t.Core.Combat.Log.NewEvent("set target direction to closest enemy", glog.LogDebugEvent, -1).
		Write("enemy index", enemies[0]).
		Write("enemy key", enemy.Key()).
		Write("direction", t.direction)
}
func (t *Target) CalcTempDirection(trg combat.Point) combat.Point {
	src := t.Pos()
	direction := combat.CalcDirection(src, trg)
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
