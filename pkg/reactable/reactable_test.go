package reactable

import (
	"os"
	"testing"
	"time"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/target"
	"github.com/genshinsim/gcsim/pkg/testhelper"
)

func init() {
	core.RegisterCharFunc(keys.TestCharDoNotUse, testhelper.NewChar)
	core.RegisterWeaponFunc(keys.DullBlade, testhelper.NewFakeWeapon)
}

// make our own core because otherwise we run into problems with circular import
func testCore() *core.Core {
	c, _ := core.New(core.CoreOpt{
		Seed:  time.Now().Unix(),
		Debug: true,
	})
	//add player (first target)
	trg := &testTarget{}
	trg.Target = target.New(c, geometry.Point{X: 0, Y: 0}, 1)
	trg.Reactable = &Reactable{}
	trg.typ = targets.TargettablePlayer
	trg.Reactable.Init(trg, c)
	c.Combat.SetPlayer(trg)

	//add default character
	p := profile.CharacterProfile{}
	p.Base.Key = keys.TestCharDoNotUse
	p.Stats = make([]float64, attributes.EndStatType)
	p.StatsByLabel = make(map[string][]float64)
	p.Params = make(map[string]int)
	p.Sets = make(map[keys.Set]int)
	p.SetParams = make(map[keys.Set]map[string]int)
	p.Weapon.Params = make(map[string]int)
	p.Base.Element = attributes.Geo
	p.Weapon.Key = keys.DullBlade

	p.Stats[attributes.EM] = 100
	p.Base.Level = 90
	p.Base.MaxLevel = 90
	p.Talents = profile.TalentProfile{Attack: 1, Skill: 1, Burst: 1}

	i, err := c.AddChar(p)
	if err != nil {
		panic(err)
	}
	c.Player.SetActive(i)

	return c
}

func testCoreWithTrgs(count int) (*core.Core, []*testTarget) {
	c := testCore()
	var r []*testTarget
	for i := 0; i < count; i++ {
		r = append(r, addTargetToCore(c))
	}
	return c, r
}

func makeAOEAttack(c *core.Core, ele attributes.Element, dur reactions.Durability) *combat.AttackEvent {
	return &combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    ele,
			Durability: dur,
		},
		Pattern: combat.NewCircleHitOnTarget(geometry.Point{}, nil, 100),
	}
}

func makeSTAttack(ele attributes.Element, dur reactions.Durability, trg targets.TargetKey) *combat.AttackEvent {
	return &combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    ele,
			Durability: dur,
		},
		Pattern: combat.NewSingleTargetHit(trg),
	}

}

type testTarget struct {
	*Reactable
	*target.Target
	src  int
	typ  targets.TargettableType
	last combat.AttackEvent
}

func (t *testTarget) Type() targets.TargettableType { return t.typ }

func (t *testTarget) HandleAttack(atk *combat.AttackEvent) float64 {
	t.Attack(atk, nil)
	//delay damage event to end of the frame
	t.Core.Combat.Tasks.Add(func() {
		//apply the damage
		t.applyDamage(atk, 1)
		t.Core.Combat.Events.Emit(event.OnEnemyDamage, t, atk, 1.0, false)
	}, 0)
	return 1
}

func (t *testTarget) Attack(atk *combat.AttackEvent, evt glog.Event) (float64, bool) {
	t.last = *atk
	t.ShatterCheck(atk)
	if atk.Info.Durability > 0 {
		//don't care about icd
		t.React(atk)
	}
	return 0, false
}

func (t *testTarget) applyDamage(atk *combat.AttackEvent, amt float64) {
	if !atk.Reacted {
		t.Reactable.AttachOrRefill(atk)
	}
}

func addTargetToCore(c *core.Core) *testTarget {
	trg := &testTarget{}
	trg.Target = target.New(c, geometry.Point{X: 0, Y: 0}, 1)
	trg.Reactable = &Reactable{}
	trg.Reactable.Init(trg, c)
	c.Combat.AddEnemy(trg)
	return trg
}

func advanceCoreFrame(c *core.Core) {
	c.F++
	c.Tick()
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestReduce(t *testing.T) {
	r := &Reactable{}
	r.Durability[ModifierElectro] = 20
	r.reduce(attributes.Electro, 20, 1)
	if r.Durability[ModifierElectro] != 0 {
		t.Errorf("expecting nil electro balance, got %v", r.Durability[ModifierElectro])
	}

	//straight up consumption
	r.Durability[ModifierPyro] = 20
	r.Durability[ModifierBurning] = 50
	consumed := r.reduce(attributes.Pyro, 30, 1)
	if consumed != 30 {
		t.Errorf("expecting consumed to be 30, got %v", consumed)
	}

	//2x multiplier, i.e. 1 incoming reduces 2
	r.Durability[ModifierPyro] = 50
	consumed = r.reduce(attributes.Pyro, 20, 2)
	if consumed != 20 {
		t.Errorf("expecting consumed to be 20, got %v", consumed)
	}
	if r.Durability[ModifierPyro] != 10 {
		t.Errorf("expecting 10 remaining ModifierPyro, got %v", 10)
	}

	r.Durability[ModifierPyro] = 50
	consumed = r.reduce(attributes.Pyro, 50, 0.5)
	if consumed != 50 {
		t.Errorf("expecting consumed to be 50, got %v", consumed)
	}
	if r.Durability[ModifierPyro] != 25 {
		t.Errorf("expecting 25 remaining ModifierPyro, got %v", 25)
	}
}

func TestTick(t *testing.T) {
	c := testCore()

	trg := addTargetToCore(c)
	trg.src = 1

	//test electro
	trg.AttachOrRefill(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 25,
		},
	})

	if trg.Durability[ModifierElectro] != 0.8*25 {
		t.Errorf("expecting 20 electro, got %v", trg.Durability[ModifierElectro])
	}
	if trg.DecayRate[ModifierElectro] != 20.0/(6*25+420) {
		t.Errorf("expecting %v decay rate, got %v", 1.0/(6*25+420), trg.DecayRate[ModifierElectro])
	}

	//should deplete fully in 570 ticks
	for i := 0; i < 6*50+420; i++ {
		trg.Tick()
		// log.Println(target.Durability)
	}
	//check all durability should be nil
	ok := trg.allNil(t)
	if !ok {
		t.FailNow()
	}

	//test multiple aura
	trg.Durability[ModifierElectro] = 0 //reset from previous test
	trg.DecayRate[ModifierElectro] = 0
	trg.AttachOrRefill(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 50,
		},
	})
	trg.AttachOrRefill(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Hydro,
			Durability: 50,
		},
	})
	trg.AttachOrRefill(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Cryo,
			Durability: 50,
		},
	})
	trg.AttachOrRefill(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Pyro,
			Durability: 50,
		},
	})
	for i := 0; i < 6*50+420; i++ {
		trg.Tick()
	}
	ok = trg.allNil(t)
	if !ok {
		t.FailNow()
	}

	//test refilling
	trg.React(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 25,
		},
	})
	for i := 0; i < 100; i++ {
		trg.Tick()
	}
	//calculate expected duration
	decay := reactions.Durability(20.0 / (6*25 + 420))
	left := 20 - 100*decay
	life := int((left + 40) / decay)
	// log.Println(decay, left, life)

	trg.React(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 50,
		},
	})
	for i := 0; i < life-1; i++ {
		trg.Tick()
	}
	//make sure > 0
	if trg.Durability[ModifierElectro] < 0 {
		t.Errorf("expecting electro not to be 0 yet, got %v", trg.Durability[ModifierElectro])
	}
	//1 more tick and should be gone
	trg.Tick()
	ok = trg.allNil(t)
	if !ok {
		t.FailNow()
	}

	//test frozen
	//50 frozen should last just over 208 frames (i.e. 0 by 209)
	trg.Durability[ModifierFrozen] = 50
	for i := 0; i < 208; i++ {
		trg.Tick()
		// log.Println(trg.Durability)
		// log.Println(trg.DecayRate)
		// log.Println("------------------------")
	}
	//should be > 0 still
	if trg.Durability[ModifierFrozen] < 0 {
		t.Errorf("expecting frozen not to be 0 yet, got %v", trg.Durability[ModifierFrozen])
	}
	//1 more tick and should be gone
	trg.Tick()
	if trg.Durability[ModifierFrozen] > 0 {
		t.Errorf("expecting frozen to be gone, got %v", trg.Durability[ModifierFrozen])
	}
	//105 more frames to full recover
	for i := 0; i < 104; i++ {
		trg.Tick()
		// log.Println(trg.Durability)
		// log.Println(trg.DecayRate)
		// log.Println("------------------------")
	}
	//decay should be > 0 still
	if trg.DecayRate[ModifierFrozen] < frzDecayCap {
		t.Errorf("expecting frozen decay to > cap, got %v", trg.Durability[ModifierFrozen])
	}
	//1 more tick to reset decay
	trg.Tick()
	if trg.DecayRate[ModifierFrozen] > frzDecayCap {
		t.Errorf("expecting frozen decay to reset, got %v", trg.Durability[ModifierFrozen])
	}

}

func (target *testTarget) allNil(t *testing.T) bool {
	ok := true
	for i, v := range target.Durability {
		ele := ReactableModifier(i).Element()
		if !durApproxEqual(0, v, 0.00001) {
			t.Errorf("ele %v expected 0 durability got %v", ele, v)
			ok = false
		}
		if !durApproxEqual(0, target.DecayRate[i], 0.00001) && ele != attributes.Frozen {
			t.Errorf("ele %v expected 0 decay got %v", ele, target.DecayRate[i])
			ok = false
		} else if !durApproxEqual(frzDecayCap, target.DecayRate[i], 0.00001) && ele == attributes.Frozen {
			t.Errorf("frozen decay expected %v got %v", frzDecayCap, target.DecayRate[i])
			ok = false
		}
	}
	return ok
}

func durApproxEqual(expect, result, tol reactions.Durability) bool {
	if expect > result {
		return expect-result < tol
	}
	return result-expect < tol
}
