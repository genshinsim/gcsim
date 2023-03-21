package combat

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/targets"
)

type (
	testchar struct{}
	testteam struct{}
	testtarg struct {
		typ         targets.TargettableType
		gadgetTyp   GadgetTyp
		hdlr        *Handler
		src         int //source of gadget
		idx         int
		key         targets.TargetKey
		shp         geometry.Shape
		alive       bool
		collideWith [targets.TargettableTypeCount]bool
		onCollision func(Target)
		direction   geometry.Point
	}
)

// char
func (t *testchar) ApplyAttackMods(a *AttackEvent, x Target) []interface{} {
	return nil
}

// team
func (t *testteam) CombatByIndex(i int) Character             { return &testchar{} }
func (t *testteam) ApplyHitlag(char int, factor, dur float64) {}

// target
func (t *testtarg) Index() int                                      { return t.idx }
func (t *testtarg) SetIndex(i int)                                  { t.idx = i }
func (t *testtarg) Key() targets.TargetKey                          { return t.key }
func (t *testtarg) SetKey(i targets.TargetKey)                      { t.key = i }
func (t *testtarg) Type() targets.TargettableType                   { return t.typ }
func (t *testtarg) Shape() geometry.Shape                           { return t.shp }
func (t *testtarg) Pos() geometry.Point                             { return t.shp.Pos() }
func (t *testtarg) SetPos(p geometry.Point)                         {} //??
func (t *testtarg) IsAlive() bool                                   { return t.alive }
func (t *testtarg) SetTag(key string, val int)                      {}
func (t *testtarg) GetTag(key string) int                           { return -1 }
func (t *testtarg) RemoveTag(key string)                            {}
func (t *testtarg) Attack(*AttackEvent, glog.Event) (float64, bool) { return 0, false }
func (t *testtarg) Tick()                                           {}
func (t *testtarg) Kill()                                           { t.hdlr.RemoveGadget(t.Key()) }
func (t *testtarg) CollidableWith(x targets.TargettableType) bool   { return t.collideWith[x] }
func (t *testtarg) GadgetTyp() GadgetTyp                            { return t.gadgetTyp }
func (t *testtarg) Src() int                                        { return t.src }
func (t *testtarg) CollidedWith(x Target) {
	if t.onCollision != nil {
		t.onCollision(x)
	}
}
func (t *testtarg) WillCollide(s geometry.Shape) bool {
	if !t.alive {
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

func (t *testtarg) HandleAttack(*AttackEvent) float64 { return 0 }

func (t *testtarg) AttackWillLand(a AttackPattern) (bool, string) {
	//geometry.Shape shouldn't be nil; panic here
	if a.Shape == nil {
		panic("unexpected nil geometry.Shape")
	}
	if !t.alive {
		return false, "target dead"
	}
	//geometry.Shape can't be nil now, check if type matches
	// if !a.Targets[t.typ] {
	// 	return false, "wrong type"
	// }
	// swirl aoe shouldn't hit the src of the aoe
	for _, v := range a.IgnoredKeys {
		if t.Key() == v {
			return false, "no self harm"
		}
	}

	//check if geometry.Shape matches
	switch v := a.Shape.(type) {
	case *geometry.Circle:
		return t.Shape().IntersectCircle(*v), "intersect circle"
	case *geometry.Rectangle:
		return t.Shape().IntersectRectangle(*v), "intersect rectangle"
	case *geometry.SingleTarget:
		//only true if
		return v.Target == t.key, "target"
	default:
		return false, "unknown geometry.Shape"
	}
}

func (t *testtarg) IsWithinArea(a AttackPattern) bool {
	return a.Shape.PointInShape(t.Pos())
}

func (t *testtarg) Direction() geometry.Point { return t.direction }
func (t *testtarg) SetDirection(trg geometry.Point) {
	src := t.Pos()
	t.direction = geometry.CalcDirection(src, trg)
}
func (t *testtarg) SetDirectionToClosestEnemy() {} // ???
func (t *testtarg) CalcTempDirection(trg geometry.Point) geometry.Point {
	return geometry.DefaultDirection()
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func newCombatCtrl() *Handler {
	return New(Opt{
		Events:       event.New(),
		Team:         &testteam{},
		Rand:         rand.New(rand.NewSource(time.Now().Unix())),
		Debug:        false,
		Log:          &glog.NilLogger{},
		DefHalt:      true,
		EnableHitlag: true,
	})
}
