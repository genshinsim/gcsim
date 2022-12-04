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

	direction float64
}

func New(core *core.Core, x, y, r float64) *Target {
	t := &Target{
		Core: core,
	}
	t.Hitbox = *combat.NewSimpleCircle(x, y, r)
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

func (t *Target) Direction() float64 { return t.direction }
func (t *Target) SetDirection(trgX, trgY float64) {
	srcX, srcY := t.Pos()
	// setting direction to self resets direction
	if srcX == trgX && srcY == trgY {
		t.Core.Combat.Log.NewEvent("reset target direction to 0", glog.LogDebugEvent, -1)
		t.direction = 0
		return
	}
	t.direction = combat.CalcDirection(srcX, srcY, trgX, trgY)
	t.Core.Combat.Log.NewEvent("set target direction", glog.LogDebugEvent, -1).
		Write("src target key", t.key).
		Write("srcX", srcX).
		Write("srcY", srcY).
		Write("trgX", trgX).
		Write("trgY", trgY).
		Write("direction (in degrees)", combat.DirectionToDegrees(t.direction))
}
func (t *Target) SetDirectionToClosestEnemy() {
	srcX, srcY := t.Pos()
	// calculate direction towards closest enemy
	enemyIndex := t.Core.Combat.EnemyByDistance(srcX, srcY, combat.InvalidTargetKey)[0]
	enemy := t.Core.Combat.Enemy(enemyIndex)
	if enemy == nil {
		panic("there should be an enemy to calculate direction")
	}
	t.SetDirection(enemy.Pos())
	t.Core.Combat.Log.NewEvent("set target direction to closest enemy", glog.LogDebugEvent, -1).
		Write("enemy index", enemyIndex).
		Write("enemy key", enemy.Key()).
		Write("direction (in degrees)", combat.DirectionToDegrees(t.direction))
}
func (t *Target) CalcTempDirection(trgX, trgY float64) float64 {
	srcX, srcY := t.Pos()
	direction := combat.CalcDirection(srcX, srcY, trgX, trgY)
	t.Core.Combat.Log.NewEvent("using temporary target direction", glog.LogDebugEvent, -1).
		Write("src target key", t.key).
		Write("srcX", srcX).
		Write("srcY", srcY).
		Write("trgX", trgX).
		Write("trgY", trgY).
		Write("existing direction (in degrees)", combat.DirectionToDegrees(t.direction)).
		Write("temporary direction (in degrees)", combat.DirectionToDegrees(direction))
	return direction
}
