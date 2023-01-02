package collision_test

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/gadget"
	"github.com/genshinsim/gcsim/pkg/testhelper"
)

func init() {
	core.RegisterCharFunc(keys.TestCharDoNotUse, testhelper.NewChar)
	core.RegisterWeaponFunc(keys.DullBlade, testhelper.NewFakeWeapon)
}

func makeCore(trgCount int) (*core.Core, []*enemy.Enemy) {
	c, _ := core.New(core.CoreOpt{
		Seed:  time.Now().Unix(),
		Debug: true,
	})
	a := avatar.New(c, combat.Point{X: 0, Y: 0}, 1)
	c.Combat.SetPlayer(a)
	var trgs []*enemy.Enemy

	for i := 0; i < trgCount; i++ {
		e := enemy.New(c, enemy.EnemyProfile{
			Level:  100,
			Resist: make(map[attributes.Element]float64),
			Pos: core.Coord{
				X: 0,
				Y: 0,
				R: 1,
			},
		})
		trgs = append(trgs, e)
		c.Combat.AddEnemy(e)
	}

	for i := 0; i < 4; i++ {
		p := profile.CharacterProfile{}
		p.Base.Key = keys.TestCharDoNotUse
		p.Stats = make([]float64, attributes.EndStatType)
		p.StatsByLabel = make(map[string][]float64)
		p.Params = make(map[string]int)
		p.Sets = make(map[keys.Set]int)
		p.SetParams = make(map[keys.Set]map[string]int)
		p.Weapon.Params = make(map[string]int)
		p.Base.StartHP = -1
		p.Base.Element = attributes.Geo
		p.Weapon.Key = keys.DullBlade

		p.Stats[attributes.EM] = 100
		p.Base.Level = 90
		p.Base.MaxLevel = 90
		p.Talents = profile.TalentProfile{Attack: 1, Skill: 1, Burst: 1}

		_, err := c.AddChar(p)
		if err != nil {
			panic(err)
		}
	}
	c.Player.SetActive(0)

	return c, trgs
}

func advanceCoreFrame(c *core.Core) {
	c.F++
	c.Tick()
}

func TestMain(m *testing.M) {
	rand.New(rand.NewSource(time.Now().Unix()))
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestSingleTarget(t *testing.T) {
	c, trgs := makeCore(rand.Intn(10))
	count := 0
	c.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		count++
		return false
	}, "dmg-count")
	c.Events.Subscribe(event.OnPlayerHit, func(args ...interface{}) bool {
		count++
		return false
	}, "dmg-count")
	c.Events.Subscribe(event.OnGadgetHit, func(args ...interface{}) bool {
		count++
		return false
	}, "dmg-count")

	for i := 0; i < len(trgs); i++ {
		count = 0

		c.QueueAttackEvent(&combat.AttackEvent{
			Pattern: combat.NewSingleTargetHit(trgs[i].Key()),
		}, 0)
		advanceCoreFrame(c)

		if count != 1 {
			t.Errorf("expecting 1 damage count, got %v", count)
		}
	}
}

func TestMultipleEnemies(t *testing.T) {

	c, trgs := makeCore(rand.Intn(10))
	count := 0
	c.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		count++
		return false
	}, "dmg-count")
	c.Events.Subscribe(event.OnPlayerHit, func(args ...interface{}) bool {
		count++
		return false
	}, "dmg-count")
	c.Events.Subscribe(event.OnGadgetHit, func(args ...interface{}) bool {
		count++
		return false
	}, "dmg-count")

	//last one should be moved aside
	trgs[len(trgs)-1].SetPos(combat.Point{X: 2, Y: 0})

	for i := 0; i < len(trgs)-1; i++ {
		count = 0

		c.QueueAttackEvent(&combat.AttackEvent{
			Pattern: combat.NewCircleHitOnTarget(trgs[i], nil, 0.5),
		}, 0)
		advanceCoreFrame(c)

		if count != len(trgs)-1 {
			t.Errorf("expecting %v damage count, got %v", len(trgs)-1, count)
		}
	}
}

type testGadget struct {
	*gadget.Gadget
}

func (t *testGadget) HandleAttack(atk *combat.AttackEvent) float64 {
	t.Core.Events.Emit(event.OnGadgetHit, t, atk)
	return 0
}

func TestDefaultHitGadget(t *testing.T) {

	c, trgs := makeCore(rand.Intn(10))
	count := 0
	c.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		count++
		return false
	}, "dmg-count")
	c.Events.Subscribe(event.OnPlayerHit, func(args ...interface{}) bool {
		count++
		return false
	}, "dmg-count")
	c.Events.Subscribe(event.OnGadgetHit, func(args ...interface{}) bool {
		count++
		return false
	}, "dmg-count")

	g := &testGadget{
		Gadget: gadget.New(c, combat.Point{X: 0, Y: 0}, 0, combat.GadgetTypTest),
	}

	c.Combat.AddGadget(g)

	for i := 0; i < len(trgs); i++ {
		count = 0

		c.QueueAttackEvent(&combat.AttackEvent{
			Pattern: combat.NewCircleHitOnTarget(trgs[i], nil, 0.5),
		}, 0)
		advanceCoreFrame(c)

		if count != len(trgs)+1 {
			t.Errorf("expecting %v damage count, got %v", len(trgs)+1, count)
		}
	}
}

func TestSkipTargets(t *testing.T) {

	c, trgs := makeCore(rand.Intn(10))
	count := 0
	c.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		count++
		return false
	}, "dmg-count")
	c.Events.Subscribe(event.OnPlayerHit, func(args ...interface{}) bool {
		count++
		return false
	}, "dmg-count")
	c.Events.Subscribe(event.OnGadgetHit, func(args ...interface{}) bool {
		count++
		return false
	}, "dmg-count")

	g := &testGadget{
		Gadget: gadget.New(c, combat.Point{X: 0, Y: 0}, 0, combat.GadgetTypTest),
	}

	c.Combat.AddGadget(g)

	for i := 0; i < len(trgs); i++ {
		count = 0
		ae := &combat.AttackEvent{
			Pattern: combat.NewCircleHitOnTarget(trgs[i], nil, 0.5),
		}
		ae.Pattern.SkipTargets[combat.TargettableEnemy] = true

		c.QueueAttackEvent(ae, 0)
		advanceCoreFrame(c)

		if count != 1 {
			t.Errorf("expecting %v damage count, got %v", 1, count)
		}
	}

	for i := 0; i < len(trgs); i++ {
		count = 0
		ae := &combat.AttackEvent{
			Pattern: combat.NewCircleHitOnTarget(trgs[i], nil, 0.5),
		}
		ae.Pattern.SkipTargets[combat.TargettablePlayer] = false
		ae.Pattern.SkipTargets[combat.TargettableEnemy] = true

		c.QueueAttackEvent(ae, 0)
		advanceCoreFrame(c)

		if count != 2 {
			t.Errorf("expecting %v damage count, got %v", 2, count)
		}
	}
}

// common:
// - attack centers on player
// - player at 0,0 + t0 at 2,2 + t1 at 4,7
// - player direction is towards t1
// - everyone has radius 1
func TestCircleAttackCollision(t *testing.T) {
	tests := map[string]struct {
		attackRadius float64
		attackOffset combat.Point
		fanAngle     float64
		want         int
	}{
		// no offset
		"no offset, hit nothing": {attackRadius: 1, attackOffset: combat.Point{}, fanAngle: 360, want: 0},
		"no offset, hit t0":      {attackRadius: 7.06, attackOffset: combat.Point{}, fanAngle: 360, want: 1},
		"no offset, hit t0 & t1": {attackRadius: 7.07, attackOffset: combat.Point{}, fanAngle: 360, want: 2},
		// offset
		"offset, hit nothing": {attackRadius: 1.7, attackOffset: combat.Point{X: -1, Y: 5.5}, fanAngle: 360, want: 0},
		"offset, hit t0":      {attackRadius: 1, attackOffset: combat.Point{Y: 3}, fanAngle: 360, want: 1},
		"offset, hit t1":      {attackRadius: 1.78, attackOffset: combat.Point{X: -1, Y: 5.5}, fanAngle: 360, want: 1},
		"offset, hit t0 & t1": {attackRadius: 1.79, attackOffset: combat.Point{X: -1, Y: 5.5}, fanAngle: 360, want: 2},
		// no offset, fanAngle
		"no offset, fanAngle, hit nothing": {attackRadius: 1, attackOffset: combat.Point{}, fanAngle: 30, want: 0},
		"no offset, fanAngle, hit t0":      {attackRadius: 7.06, attackOffset: combat.Point{}, fanAngle: 30, want: 1},
		"no offset, fanAngle, hit t0 & t1": {attackRadius: 7.07, attackOffset: combat.Point{}, fanAngle: 30, want: 2},
		// offset, fanAngle
		"offset, fanAngle, hit nothing": {attackRadius: 1, attackOffset: combat.Point{X: -2, Y: 2}, fanAngle: 30, want: 0},
		"offset, fanAngle, hit t0":      {attackRadius: 1, attackOffset: combat.Point{X: -2, Y: 2}, fanAngle: 35, want: 1},
		"offset, fanAngle, hit t1":      {attackRadius: 1.79, attackOffset: combat.Point{X: -1, Y: 5.5}, fanAngle: 30, want: 1},
		"offset, fanAngle, hit t0 & t1": {attackRadius: 1.79, attackOffset: combat.Point{X: -1, Y: 5.5}, fanAngle: 350, want: 2},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := circleAttackCollision(tc.attackRadius, tc.attackOffset, tc.fanAngle)
			if got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func circleAttackCollision(attackRadius float64, attackOffset combat.Point, fanAngle float64) int {
	c, trgs := makeCore(2)
	player := c.Combat.Player()
	count := 0

	c.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		count++
		return false
	}, "dmg-count")
	trgs[0].SetPos(combat.Point{X: 2, Y: 2})
	trgs[1].SetPos(combat.Point{X: 7, Y: 4})
	player.SetDirection(trgs[1].Pos())

	c.QueueAttackEvent(&combat.AttackEvent{
		Pattern: combat.NewCircleHitOnTargetFanAngle(player, attackOffset, attackRadius, fanAngle),
	}, 0)
	advanceCoreFrame(c)
	return count
}

// common:
// - attack centers on player
// - player at 0,0 + t0 at 2,2 + t1 at 4,7
// - player direction is towards t1
// - everyone has radius 1
func TestRectangleAttackCollision(t *testing.T) {
	tests := map[string]struct {
		attackWidth  float64
		attackHeight float64
		attackOffset combat.Point
		want         int
	}{
		// no offset
		"no offset, hit nothing": {attackWidth: 3, attackHeight: 3, attackOffset: combat.Point{}, want: 0},
		"no offset, hit t0":      {attackWidth: 3, attackHeight: 4, attackOffset: combat.Point{}, want: 1},
		"no offset, hit t0 & t1": {attackWidth: 3, attackHeight: 15, attackOffset: combat.Point{}, want: 2},
		// offset
		"offset, hit nothing": {attackWidth: 2.5, attackHeight: 1, attackOffset: combat.Point{X: -3, Y: 2}, want: 0},
		"offset, hit t0":      {attackWidth: 2.6, attackHeight: 1, attackOffset: combat.Point{X: -3, Y: 2}, want: 1},
		"offset, hit t1":      {attackWidth: 1, attackHeight: 1, attackOffset: combat.Point{Y: 9}, want: 1},
		"offset, hit t0 & t1": {attackWidth: 0.1, attackHeight: 15, attackOffset: combat.Point{Y: 2}, want: 2},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := rectangleAttackCollision(tc.attackWidth, tc.attackHeight, tc.attackOffset)
			if got != tc.want {
				t.Fatalf("expected: %v, got: %v", tc.want, got)
			}
		})
	}
}

func rectangleAttackCollision(attackWidth, attackHeight float64, attackOffset combat.Point) int {
	c, trgs := makeCore(2)
	player := c.Combat.Player()
	count := 0

	c.Events.Subscribe(event.OnEnemyHit, func(args ...interface{}) bool {
		count++
		return false
	}, "dmg-count")
	trgs[0].SetPos(combat.Point{X: 2, Y: 2})
	trgs[1].SetPos(combat.Point{X: 7, Y: 4})
	player.SetDirection(trgs[1].Pos())

	c.QueueAttackEvent(&combat.AttackEvent{
		Pattern: combat.NewBoxHitOnTarget(player, attackOffset, attackWidth, attackHeight),
	}, 0)
	advanceCoreFrame(c)
	return count
}
