package reactable

import (
	"os"
	"testing"
	"time"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
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
	trg.Target = target.New(c, 0, 0, 1)
	trg.Reactable = &Reactable{}
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
	p.Base.StartHP = -1
	p.Base.Element = attributes.Geo
	p.Weapon.Key = keys.DullBlade

	p.Stats[attributes.EM] = 100
	p.Base.Level = 90
	p.Base.MaxLevel = 90
	p.Talents = profile.TalentProfile{1, 1, 1}

	i, err := c.AddChar(p)
	if err != nil {
		panic(err)
	}
	c.Player.SetActive(i)

	return c
}

type testTarget struct {
	*Reactable
	*target.Target
	src           int
	typ           combat.TargettableType
	onDmgCallBack func(*combat.AttackEvent) (float64, bool)
}

func (t *testTarget) Type() combat.TargettableType {
	return t.typ
}

func (t *testTarget) Attack(atk *combat.AttackEvent, evt glog.Event) (float64, bool) {
	if t.onDmgCallBack != nil {
		return t.onDmgCallBack(atk)
	}
	return 0, false
}

func addTargetToCore(c *core.Core) *testTarget {
	trg := &testTarget{}
	trg.Target = target.New(c, 0, 0, 1)
	trg.Reactable = &Reactable{}
	trg.Reactable.Init(trg, c)
	c.Combat.AddEnemy(trg)
	trg.SetIndex(c.Combat.EnemyCount() - 1)
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
	r.Durability = make([]combat.Durability, attributes.ElementDelimAttachable)
	r.DecayRate = make([]combat.Durability, attributes.ElementDelimAttachable)
	r.Durability[attributes.Electro] = 20
	r.reduce(attributes.Electro, 20, 1)
	if r.Durability[attributes.Electro] != 0 {
		t.Errorf("expecting nil electro balance, got %v", r.Durability[attributes.Electro])
	}
}

func TestTick(t *testing.T) {
	c := testCore()

	trg := addTargetToCore(c)
	trg.src = 1

	//test electro
	trg.React(&combat.AttackEvent{
		Info: combat.AttackInfo{
			Element:    attributes.Electro,
			Durability: 25,
		},
	})

	if trg.Durability[attributes.Electro] != 0.8*25 {
		t.Errorf("expecting 20 electro, got %v", trg.Durability[attributes.Electro])
	}
	if trg.DecayRate[attributes.Electro] != 20.0/(6*25+420) {
		t.Errorf("expecting %v decay rate, got %v", 1.0/(6*25+420), trg.DecayRate[attributes.Electro])
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
	trg.attach(attributes.Electro, 50, 0.8)
	trg.attach(attributes.Hydro, 50, 0.8)
	trg.attach(attributes.Cryo, 50, 0.8)
	trg.attach(attributes.Pyro, 50, 0.8)
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
	decay := combat.Durability(20.0 / (6*25 + 420))
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
	if trg.Durability[attributes.Electro] < 0 {
		t.Errorf("expecting electro not to be 0 yet, got %v", trg.Durability[attributes.Electro])
	}
	//1 more tick and should be gone
	trg.Tick()
	ok = trg.allNil(t)
	if !ok {
		t.FailNow()
	}

	//test frozen
	//50 frozen should last just over 208 frames (i.e. 0 by 209)
	trg.Durability[attributes.Frozen] = 50
	for i := 0; i < 208; i++ {
		trg.Tick()
		// log.Println(trg.Durability)
		// log.Println(trg.DecayRate)
		// log.Println("------------------------")
	}
	//should be > 0 still
	if trg.Durability[attributes.Frozen] < 0 {
		t.Errorf("expecting frozen not to be 0 yet, got %v", trg.Durability[attributes.Frozen])
	}
	//1 more tick and should be gone
	trg.Tick()
	if trg.Durability[attributes.Frozen] > 0 {
		t.Errorf("expecting frozen to be gone, got %v", trg.Durability[attributes.Frozen])
	}
	//105 more frames to full recover
	for i := 0; i < 104; i++ {
		trg.Tick()
		// log.Println(trg.Durability)
		// log.Println(trg.DecayRate)
		// log.Println("------------------------")
	}
	//decay should be > 0 still
	if trg.DecayRate[attributes.Frozen] < frzDecayCap {
		t.Errorf("expecting frozen decay to > cap, got %v", trg.Durability[attributes.Frozen])
	}
	//1 more tick to reset decay
	trg.Tick()
	if trg.DecayRate[attributes.Frozen] > frzDecayCap {
		t.Errorf("expecting frozen decay to reset, got %v", trg.Durability[attributes.Frozen])
	}

}

func (target *testTarget) allNil(t *testing.T) bool {
	ok := true
	for i, v := range target.Durability {
		ele := attributes.Element(i)
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

func durApproxEqual(expect, result, tol combat.Durability) bool {
	if expect > result {
		return expect-result < tol
	}
	return result-expect < tol
}

// func floatApproxEqual(expect, result, tol float64) bool {
// 	if expect > result {
// 		return expect-result < tol
// 	}
// 	return result-expect < tol
// }
