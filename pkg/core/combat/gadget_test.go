package combat

import (
	"log"
	"testing"

	"github.com/genshinsim/gcsim/pkg/core/info"
)

//nolint:unparam // all calls currently have y = 0 but that can change
func newSimpleCircle(x, y, r float64) *info.Circle {
	return info.NewCircle(info.Point{X: x, Y: y}, r, info.DefaultDirection(), 360)
}

func TestGadgetCollision(t *testing.T) {
	c := newCombatCtrl()
	const ecount = 2
	const gcount = 4
	// 1 player
	player := &testtarg{
		typ:   info.TargettablePlayer,
		shp:   newSimpleCircle(0, 0, 0.2),
		alive: true,
		onCollision: func(info.Target) {
			log.Printf("collision shouldn't happen with player!!")
			t.FailNow()
		},
	}
	c.SetPlayer(player)
	// 2 enemies
	for i := range ecount {
		v := &testtarg{
			typ:   info.TargettableEnemy,
			shp:   newSimpleCircle(float64(i)*0.5, 0, 0.2),
			alive: true,
			onCollision: func(info.Target) {
				log.Printf("collision shouldn't happen with enemy!!")
				t.FailNow()
			},
		}
		c.AddEnemy(v)
	}
	// gadget should overlap player and first enemy
	var cw [info.TargettableTypeCount]bool
	cw[info.TargettableEnemy] = true
	cw[info.TargettablePlayer] = true
	count := 0
	// make multiple gadgets in the same spot, so we should get gcount * 2 collision total
	for range gcount {
		v := &testtarg{
			hdlr:        c,
			typ:         info.TargettableGadget,
			shp:         newSimpleCircle(0, 0, 0.1),
			alive:       true,
			collideWith: cw,
			onCollision: func(t info.Target) {
				log.Printf("Collided with %v, type: %v!\n", t, t.Type())
				count++
			},
		}
		c.AddGadget(v)
	}

	c.Tick()

	if count < 2*gcount {
		log.Printf("Expecting %v collisions, got %v\n", gcount*2, count)
		t.Fail()
	}
}

func TestGadgetLimits(t *testing.T) {
	c := newCombatCtrl()
	const ecount = 2
	const gcount = 20
	// 1 player
	player := &testtarg{
		typ:   info.TargettablePlayer,
		shp:   newSimpleCircle(0, 0, 0.2),
		alive: true,
		onCollision: func(info.Target) {
			log.Printf("collision shouldn't happen with player!!")
			t.FailNow()
		},
	}
	c.SetPlayer(player)
	// 2 enemies
	for i := range ecount {
		v := &testtarg{
			typ:   info.TargettableEnemy,
			shp:   newSimpleCircle(float64(i)*0.5, 0, 0.2),
			alive: true,
			onCollision: func(info.Target) {
				log.Printf("collision shouldn't happen with enemy!!")
				t.FailNow()
			},
		}
		c.AddEnemy(v)
	}
	// gadget should overlap player and first enemy
	var cw [info.TargettableTypeCount]bool
	cw[info.TargettableEnemy] = true
	cw[info.TargettablePlayer] = true
	count := 0
	// make multiple gadgets; gadgets should not exceed 2
	for range gcount {
		v := &testtarg{
			hdlr:        c,
			typ:         info.TargettableGadget,
			gadgetTyp:   info.GadgetTypTest,
			shp:         newSimpleCircle(0, 0, 0.1),
			alive:       true,
			collideWith: cw,
		}
		c.AddGadget(v)
	}

	c.Tick()

	// check how many we got
	for _, v := range c.gadgets {
		if v != nil && v.GadgetTyp() == info.GadgetTypTest {
			count++
		}
	}

	if count > 2 {
		t.Errorf("Expecting max 2 gadgets, got %v", count)
	}
}

func BenchmarkCollisionCheck(b *testing.B) {
	c := newCombatCtrl()
	const ecount = 2
	const gcount = 20
	// 1 player
	player := &testtarg{
		typ:   info.TargettablePlayer,
		shp:   newSimpleCircle(0, 0, 0.2),
		alive: true,
	}
	c.SetPlayer(player)
	// 2 enemies
	for i := range ecount {
		v := &testtarg{
			typ:   info.TargettableEnemy,
			shp:   newSimpleCircle(float64(i)*0.5, 0, 0.2),
			alive: true,
		}
		c.AddEnemy(v)
	}
	// gadget should overlap player and first enemy
	var cw [info.TargettableTypeCount]bool
	cw[info.TargettableEnemy] = true
	cw[info.TargettablePlayer] = true
	// make multiple gadgets in the same spot, so we should get gcount * 2 collision total
	for range gcount {
		v := &testtarg{
			typ:         info.TargettableGadget,
			shp:         newSimpleCircle(0, 0, 0.1),
			alive:       true,
			collideWith: cw,
			onCollision: func(t info.Target) {
			},
		}
		c.AddGadget(v)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		c.Tick()
	}
}

func TestKillGadgetOnCollision(t *testing.T) {
	c := newCombatCtrl()
	const ecount = 2
	const gcount = 4
	// 1 player
	player := &testtarg{
		typ:   info.TargettablePlayer,
		shp:   newSimpleCircle(0, 0, 0.2),
		alive: true,
		onCollision: func(info.Target) {
			log.Printf("collision shouldn't happen with player!!")
			t.FailNow()
		},
	}
	c.SetPlayer(player)
	// 2 enemies
	for i := range ecount {
		v := &testtarg{
			typ:   info.TargettableEnemy,
			shp:   newSimpleCircle(float64(i)*0.5, 0, 0.2),
			alive: true,
			onCollision: func(info.Target) {
				log.Printf("collision shouldn't happen with enemy!!")
				t.FailNow()
			},
		}
		c.AddEnemy(v)
	}
	// gadget should overlap player and first enemy
	var cw [info.TargettableTypeCount]bool
	cw[info.TargettableEnemy] = true
	cw[info.TargettablePlayer] = true
	count := 0
	// make multiple gadgets in the same spot, so we should get gcount * 2 collision total
	for range gcount {
		v := &testtarg{
			typ:         info.TargettableGadget,
			shp:         newSimpleCircle(0, 0, 0.1),
			alive:       true,
			collideWith: cw,
		}
		v.onCollision = func(t info.Target) {
			count++
			// kill self
			c.RemoveGadget(v.key)
		}
		c.AddGadget(v)
	}

	c.Tick()

	// only 1 collision per since it should kill self
	if count < gcount {
		log.Printf("Expecting %v collisions, got %v\n", gcount, count)
		t.FailNow()
	}

	count = 0

	c.Tick()

	if count != 0 {
		log.Printf("expecting 0 collision, got %v\n", count)
		t.FailNow()
	}

	if c.GadgetCount() != 0 {
		log.Printf("expecting 0 gadgets, got %v\n", c.GadgetCount())
		t.FailNow()
	}
}

type sharedTarg struct {
	*testtarg
	got *info.AttackEvent
}

func (s *sharedTarg) HandleSharedAttack(a *info.AttackEvent) (float64, bool) {
	s.got = a
	return 0, true
}

func TestHandleSharedAttackGetsEvent(t *testing.T) {
	c := newCombatCtrl()
	c.SetPlayer(&testtarg{
		typ:   info.TargettablePlayer,
		shp:   newSimpleCircle(0, 0, 0.2),
		alive: true,
	})
	g := &sharedTarg{testtarg: &testtarg{
		typ:   info.TargettableGadget,
		shp:   newSimpleCircle(0, 0, 0.1),
		alive: true,
	}}
	c.AddGadget(g)

	a := &info.AttackEvent{Pattern: NewCircleHitOnTarget(info.Point{}, nil, 1)}
	c.ApplyAttack(a)
	if g.got != a {
		t.Fatalf("expecting HandleSharedAttack to get the attack event itself, got %p want %p", g.got, a)
	}
}

func TestGadgetMadeCollidableAfterAdd(t *testing.T) {
	c := newCombatCtrl()
	c.SetPlayer(&testtarg{
		typ:   info.TargettablePlayer,
		shp:   newSimpleCircle(0, 0, 0.2),
		alive: true,
	})
	count := 0
	g := &testtarg{
		typ:   info.TargettableGadget,
		shp:   newSimpleCircle(0, 0, 0.1),
		alive: true,
		onCollision: func(info.Target) {
			count++
		},
	}
	c.AddGadget(g)

	c.Tick()
	if count != 0 {
		t.Fatalf("expecting 0 collisions before the gadget is collidable, got %v", count)
	}

	g.collideWith[info.TargettablePlayer] = true
	c.EnableCollisionChecks()
	c.Tick()
	if count != 1 {
		t.Fatalf("expecting 1 collision after the gadget is collidable, got %v", count)
	}
}
