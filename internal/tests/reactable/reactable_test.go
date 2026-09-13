package reactable_test

import (
	"time"

	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/testhelper"
)

func init() {
	testhelper.RegisterTestCharacter()
	testhelper.RegisterTestWeapon()
}

func makeCore(trgCount int) (*core.Core, []*enemy.Enemy) {
	c, _ := core.New(core.Opt{
		Seed:  time.Now().Unix(),
		Debug: true,
	})
	a := avatar.New(c, info.Point{X: 0, Y: 0}, 1)
	c.Combat.SetPlayer(a)
	var trgs []*enemy.Enemy

	for range trgCount {
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

	for range 4 {
		p := testhelper.DefaultProfile(testhelper.TestCharKey, testhelper.TestWeaponKey)
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
