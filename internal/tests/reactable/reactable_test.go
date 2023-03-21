package reactable_test

import (
	"time"

	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
	"github.com/genshinsim/gcsim/pkg/enemy"
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
	a := avatar.New(c, geometry.Point{X: 0, Y: 0}, 1)
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
