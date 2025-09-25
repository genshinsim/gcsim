package reactable_test

import (
	"log"
	"testing"

	"github.com/genshinsim/gcsim/internal/template/dendrocore"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/gadget"
)

// test modifying dendro core to something else
func TestModifyDendroCore(t *testing.T) {
	c, _ := makeCore(1)
	err := c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}
	count := 0
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) bool {
		trg := args[0].(info.Target)
		ae := args[1].(*info.AttackEvent)
		if trg.Type() == info.TargettableEnemy && ae.Info.Abil == "bloom" {
			count++
		}
		return false
	}, "bloom")
	c.Events.Subscribe(event.OnDendroCore, func(args ...any) bool {
		if g, ok := args[0].(*dendrocore.Gadget); ok {
			log.Println("replacing gadget on dendro core")
			c.Combat.ReplaceGadget(g.Key(), &fakeCore{
				Gadget: gadget.New(c, info.Point{X: 0, Y: 0}, 0.2, info.GadgetTypDendroCore),
			})
			// prevent blowing up
			g.OnKill = nil
			g.OnExpiry = nil
			g.OnCollision = nil
		}
		return false
	}, "modify-core")

	c.QueueAttackEvent(&info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Dendro,
			Durability: 25,
		},
		Pattern: combat.NewCircleHitOnTarget(info.Point{}, nil, 100),
	}, 0)
	advanceCoreFrame(c)

	c.QueueAttackEvent(&info.AttackEvent{
		Info: info.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 50,
		},
		Pattern: combat.NewCircleHitOnTarget(info.Point{}, nil, 100),
	}, 0)

	// should create a seed, explodes after 5s
	for range dendrocore.Delay + 1 {
		advanceCoreFrame(c)
	}

	if c.Combat.GadgetCount() != 1 {
		t.Errorf("expected created a gadget (bloom), got %v", c.Combat.GadgetCount())
	}

	if _, ok := c.Combat.Gadget(0).(*fakeCore); !ok {
		t.Errorf("gadget not a fake core??")
	}

	// make sure no blow up
	for range 600 {
		advanceCoreFrame(c)
	}

	if count != 0 {
		t.Errorf("expecting 0 dmg count, got %v", count)
	}
}

type fakeCore struct {
	*gadget.Gadget
}

func (f *fakeCore) Tick()                                                {}
func (f *fakeCore) HandleAttack(*info.AttackEvent) float64               { return 0 }
func (f *fakeCore) Attack(*info.AttackEvent, glog.Event) (float64, bool) { return 0, false }
func (f *fakeCore) SetDirection(trg info.Point)                          {}
func (f *fakeCore) SetDirectionToClosestEnemy()                          {}
func (f *fakeCore) CalcTempDirection(trg info.Point) info.Point {
	return info.DefaultDirection()
}
