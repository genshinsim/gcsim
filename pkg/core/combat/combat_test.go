package combat

import (
	"math/rand"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

type (
	testchar struct{}
	testteam struct{}
	testtarg struct {
		typ         info.TargettableType
		gadgetTyp   info.GadgetTyp
		hdlr        *Handler
		src         int // source of gadget
		idx         int
		key         info.TargetKey
		shp         info.Shape
		alive       bool
		collideWith [info.TargettableTypeCount]bool
		onCollision func(info.Target)
		direction   info.Point
	}
)

// char
func (t *testchar) ApplyAttackMods(a *info.AttackEvent, x info.Target) []any {
	return nil
}

// team
func (t *testteam) CombatByIndex(i int) Character             { return &testchar{} }
func (t *testteam) ApplyHitlag(char int, factor, dur float64) {}

// target
func (t *testtarg) Index() int                                           { return t.idx }
func (t *testtarg) SetIndex(i int)                                       { t.idx = i }
func (t *testtarg) Key() info.TargetKey                                  { return t.key }
func (t *testtarg) SetKey(i info.TargetKey)                              { t.key = i }
func (t *testtarg) Type() info.TargettableType                           { return t.typ }
func (t *testtarg) Shape() info.Shape                                    { return t.shp }
func (t *testtarg) Pos() info.Point                                      { return t.shp.Pos() }
func (t *testtarg) SetPos(p info.Point)                                  {} // ??
func (t *testtarg) IsAlive() bool                                        { return t.alive }
func (t *testtarg) SetTag(key string, val int)                           {}
func (t *testtarg) GetTag(key string) int                                { return -1 }
func (t *testtarg) RemoveTag(key string)                                 {}
func (t *testtarg) Attack(*info.AttackEvent, glog.Event) (float64, bool) { return 0, false }
func (t *testtarg) Tick()                                                {}
func (t *testtarg) Kill()                                                { t.hdlr.RemoveGadget(t.Key()) }
func (t *testtarg) CollidableWith(x info.TargettableType) bool           { return t.collideWith[x] }
func (t *testtarg) GadgetTyp() info.GadgetTyp                            { return t.gadgetTyp }
func (t *testtarg) Src() int                                             { return t.src }
func (t *testtarg) CollidedWith(x info.Target) {
	if t.onCollision != nil {
		t.onCollision(x)
	}
}

func (t *testtarg) WillCollide(s info.Shape) bool {
	if !t.alive {
		return false
	}
	switch v := s.(type) {
	case *info.Circle:
		return t.Shape().IntersectCircle(*v)
	case *info.Rectangle:
		return t.Shape().IntersectRectangle(*v)
	default:
		return false
	}
}

func (t *testtarg) HandleAttack(*info.AttackEvent) float64 { return 0 }

func (t *testtarg) AttackWillLand(a info.AttackPattern) (bool, string) {
	// info.Shape shouldn't be nil; panic here
	if a.Shape == nil {
		panic("unexpected nil info.Shape")
	}
	if !t.alive {
		return false, "target dead"
	}
	// info.Shape can't be nil now, check if type matches
	// if !a.Targets[t.typ] {
	// 	return false, "wrong type"
	// }
	// swirl aoe shouldn't hit the src of the aoe
	if slices.Contains(a.IgnoredKeys, t.Key()) {
		return false, "no self harm"
	}

	// check if info.Shape matches
	switch v := a.Shape.(type) {
	case *info.Circle:
		return t.Shape().IntersectCircle(*v), "intersect circle"
	case *info.Rectangle:
		return t.Shape().IntersectRectangle(*v), "intersect rectangle"
	case *info.SingleTarget:
		// only true if
		return v.Target == t.key, "target"
	default:
		return false, "unknown info.Shape"
	}
}

func (t *testtarg) IsWithinArea(a info.AttackPattern) bool {
	return a.Shape.PointInShape(t.Pos())
}

func (t *testtarg) Direction() info.Point { return t.direction }
func (t *testtarg) SetDirection(trg info.Point) {
	src := t.Pos()
	t.direction = info.CalcDirection(src, trg)
}
func (t *testtarg) SetDirectionToClosestEnemy() {} // ???
func (t *testtarg) CalcTempDirection(trg info.Point) info.Point {
	return info.DefaultDirection()
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
