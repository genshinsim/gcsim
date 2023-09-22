package abilities

import (
	"testing"
	"time"

	// we import simulation like this so that import.go is pulled in
	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/enemy"
	_ "github.com/genshinsim/gcsim/pkg/simulation"
)

// purpose of this test is to check that characters have all abil implemented to avoid loops
func TestAbilitiesImplemented(t *testing.T) {
	for k := range core.NewCharFuncMap {
		testChar(t, k)
	}
}

func testChar(t *testing.T, k keys.Char) {
	c, trg := makeCore(2)
	prof := defProfile(k)
	prof.Base.Cons = 6
	idx, err := c.AddChar(prof)
	if err != nil {
		t.Errorf("error adding char: %v", err)
		t.FailNow()
	}
	c.Player.SetActive(idx)
	err = c.Init()
	if err != nil {
		t.Errorf("error initializing core: %v", err)
		t.FailNow()
	}
	// initialize some settings
	c.Combat.DefaultTarget = trg[0].Key()

	p := make(map[string]int)
	for a := action.InvalidAction + 1; a < action.ActionSwap; a++ {
		c.Player.Exec(a, k, p)
		for !c.Player.CanQueueNextAction() {
			advanceCoreFrame(c)
		}
	}
}

func makeCore(trgCount int) (*core.Core, []*enemy.Enemy) {
	c, _ := core.New(core.Opt{
		Seed:  time.Now().Unix(),
		Debug: true,
	})
	a := avatar.New(c, geometry.Point{X: 0, Y: 0}, 1)
	c.Combat.SetPlayer(a)
	var trgs []*enemy.Enemy

	for i := 0; i < trgCount; i++ {
		e := enemy.New(c, info.EnemyProfile{
			Level:  100,
			Resist: make(map[attributes.Element]float64),
			Pos: info.Coord{
				X: 0,
				Y: 0,
				R: 1,
			},
		})
		trgs = append(trgs, e)
		c.Combat.AddEnemy(e)
	}

	c.Player.SetActive(0)

	return c, trgs
}

func defProfile(key keys.Char) info.CharacterProfile {
	p := info.CharacterProfile{}
	p.Base.Key = key
	p.Stats = make([]float64, attributes.EndStatType)
	p.StatsByLabel = make(map[string][]float64)
	p.Params = make(map[string]int)
	p.Sets = make(map[keys.Set]int)
	p.SetParams = make(map[keys.Set]map[string]int)
	p.Weapon.Params = make(map[string]int)
	p.Base.Element = keys.CharKeyToEle[key]
	p.Weapon.Key = keys.DullBlade

	p.Stats[attributes.EM] = 100
	p.Base.Level = 90
	p.Base.MaxLevel = 90
	p.Talents = info.TalentProfile{Attack: 1, Skill: 1, Burst: 1}

	return p
}

func advanceCoreFrame(c *core.Core) {
	c.F++
	c.Tick()
}
