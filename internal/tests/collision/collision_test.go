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
	a := avatar.New(c, 0, 0, 1)
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
	trgs[len(trgs)-1].SetPos(2, 0)

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
		Gadget: gadget.New(c, core.Coord{X: 0, Y: 0}, combat.GadgetTypTest),
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
		Gadget: gadget.New(c, core.Coord{X: 0, Y: 0}, combat.GadgetTypTest),
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
