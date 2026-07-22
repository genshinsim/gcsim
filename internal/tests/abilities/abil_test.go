package abilities

import (
	"errors"
	"testing"
	"time"

	// we import simulation like this so that import.go is pulled in

	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/enemy"
	_ "github.com/genshinsim/gcsim/pkg/simulation"
	"github.com/genshinsim/gcsim/pkg/testhelper"
)

// purpose of this test is to check that characters abilities do not randomly panic
func TestAbilities(t *testing.T) {
	for k := range core.NewCharFuncMap {
		testChar(t, k)
	}
}

func testChar(t *testing.T, k keys.Char) {
	c, trg := makeCore(2)
	prof := testhelper.DefaultProfile(k, keys.DullBlade)
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
	c.QueueParticle("system", 1000, attributes.NoElement, 0)
	advanceCoreFrame(c)

	p := make(map[string]int)
	for a := action.InvalidAction + 1; a < action.ActionSwap; a++ {
		for {
			err := c.Player.ReadyCheck(a, k, p)
			if err == nil {
				break
			}
			switch {
			case errors.Is(err, player.ErrActionNotReady):
			case errors.Is(err, player.ErrPlayerNotReady):
			case errors.Is(err, player.ErrActionNoOp):
				break
			default:
				t.Errorf("unexpected error waiting for action to be ready: %v", err)
				t.FailNow()
			}
			advanceCoreFrame(c)
		}
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

	c.Player.SetActive(0)

	return c, trgs
}

func advanceCoreFrame(c *core.Core) {
	c.F++
	c.Tick()
}
