package combat

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
)

type (
	testchar struct{}
	testteam struct{}
	testtarg struct {
		typ         TargettableType
		idx         int
		key         int
		shp         Shape
		alive       bool
		collideWith [TargettableTypeCount]bool
		onCollision func(Target)
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
func (t *testtarg) Key() int                                        { return t.key }
func (t *testtarg) SetKey(i int)                                    { t.key = i }
func (t *testtarg) Type() TargettableType                           { return t.typ }
func (t *testtarg) Shape() Shape                                    { return t.shp }
func (t *testtarg) Pos() (float64, float64)                         { return t.shp.Pos() }
func (t *testtarg) SetPos(x, y float64)                             {} //??
func (t *testtarg) IsAlive() bool                                   { return t.alive }
func (t *testtarg) Attack(*AttackEvent, glog.Event) (float64, bool) { return 0, false }
func (t *testtarg) Tick()                                           {}
func (t *testtarg) Kill()                                           { t.alive = false }
func (t *testtarg) CollidableWith(x TargettableType) bool           { return t.collideWith[x] }
func (t *testtarg) CollidedWith(x Target) {
	if t.onCollision != nil {
		t.onCollision(x)
	}
}
func (t *testtarg) WillCollide(s Shape) bool {
	if !t.alive {
		return false
	}
	switch v := s.(type) {
	case *Circle:
		return t.Shape().IntersectCircle(*v)
	case *Rectangle:
		return t.Shape().IntersectRectangle(*v)
	default:
		return false
	}
}
func (t *testtarg) AttackWillLand(a AttackPattern, src int) (bool, string) {
	//shape shouldn't be nil; panic here
	if a.Shape == nil {
		panic("unexpected nil shape")
	}
	if !t.alive {
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
	case *Circle:
		return t.Shape().IntersectCircle(*v), "intersect circle"
	case *Rectangle:
		return t.Shape().IntersectRectangle(*v), "intersect rectangle"
	case *SingleTarget:
		//only true if
		return v.Target == t.idx, "target"
	default:
		return false, "unknown shape"
	}
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
