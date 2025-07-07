package canqueueafter

import (
	"errors"
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
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	_ "github.com/genshinsim/gcsim/pkg/simulation"
)

// purpose of this test is to check that characters abilities have correct can queue after
func TestCanQueue(t *testing.T) {
	baseActions := [][]action.Action{
		{action.ActionSkill},
		{action.ActionSkill},
		{action.ActionBurst},
		{action.ActionAttack},
		{action.ActionDash},
		{action.ActionJump},
		{action.ActionAttack, action.ActionCharge},
		{action.ActionAim},
		{action.ActionAim},
		{action.ActionAim},
		{action.ActionAim},
		{action.ActionSkill, action.ActionSkill},
		{action.ActionSkill, action.ActionBurst},
		{action.ActionSkill, action.ActionAttack},
		{action.ActionSkill, action.ActionDash},
		{action.ActionSkill, action.ActionJump},
		{action.ActionSkill, action.ActionAttack, action.ActionCharge},
		{action.ActionSkill, action.ActionAim},
		{action.ActionBurst, action.ActionSkill},
		{action.ActionBurst, action.ActionBurst},
		{action.ActionBurst, action.ActionAttack},
		{action.ActionBurst, action.ActionDash},
		{action.ActionBurst, action.ActionJump},
		{action.ActionBurst, action.ActionAttack, action.ActionCharge},
		{action.ActionBurst, action.ActionAim},
	}

	emptyParams := make(map[string]int)

	baseParams := [][]map[string]int{
		{emptyParams},
		{map[string]int{"hold": 1}},
		{emptyParams},
		{emptyParams},
		{emptyParams},
		{emptyParams},
		{emptyParams, emptyParams},
		{emptyParams},
		{map[string]int{"hold": 0}},
		{map[string]int{"hold": 1}},
		{map[string]int{"hold": 2}},
		{emptyParams, emptyParams},
		{emptyParams, emptyParams},
		{emptyParams, emptyParams},
		{emptyParams, emptyParams},
		{emptyParams, emptyParams},
		{emptyParams, emptyParams, emptyParams},
		{emptyParams, emptyParams},
		{emptyParams, emptyParams},
		{emptyParams, emptyParams},
		{emptyParams, emptyParams},
		{emptyParams, emptyParams},
		{emptyParams, emptyParams},
		{emptyParams, emptyParams, emptyParams},
		{emptyParams, emptyParams},
	}

	for k := range core.NewCharFuncMap {
		actions := make([][]action.Action, len(baseActions))
		_ = copy(actions, baseActions)
		// insert character specific combos, params
		for i, a := range actions {
			testQueue(t, k, a, baseParams[i])
		}
	}
}

func testQueue(t *testing.T, k keys.Char, acts []action.Action, params []map[string]int) {
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
	c.Flags.IgnoreBurstEnergy = true

	advanceCoreFrame(t, c)

	for i, a := range acts {
		p := params[i]
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
			advanceCoreFrame(t, c)
		}

		evt, err := getAction(c.Player.ActiveChar(), a)(p)
		if err != nil {
			return
		}
		c.Player.AnimationHandler.SetActionUsed(c.Player.Active(), a, &evt)
		for act := range action.EndActionType {
			// Normal attacks trigger this panic when there is atkspd
			if evt.State == action.NormalAttackState && c.Player.ActiveChar().Stat(attributes.AtkSpd) > 0 {
				break
			}

			// Charge attacks for Itto/Wanderer trigger this
			if evt.State == action.ChargeAttackState && c.Player.ActiveChar().Stat(attributes.AtkSpd) > 0 {
				break
			}
			if evt.Frames(act) < evt.CanQueueAfter {
				t.Errorf("character %s action %s params: %v CanQueueAfter (%d) is larger than the output of evt.Frames[%s] (%d)", c.Player.ActiveChar().Base.Key.String(), a, p, evt.CanQueueAfter, act.String(), evt.Frames(act))
			}
		}
		for !c.Player.CanQueueNextAction() {
			advanceCoreFrame(t, c)
		}
	}
}

func getAction(char *character.CharWrapper, t action.Action) func(p map[string]int) (action.Info, error) {
	switch t {
	case action.ActionCharge:
		return char.ChargeAttack
	case action.ActionDash:
		return char.Dash
	case action.ActionJump:
		return char.Jump
	case action.ActionWalk:
		return char.Walk
	case action.ActionAim:
		return char.Aimed
	case action.ActionSkill:
		return char.Skill
	case action.ActionBurst:
		return char.Burst
	case action.ActionAttack:
		return char.Attack
	case action.ActionHighPlunge:
		return char.HighPlungeAttack
	case action.ActionLowPlunge:
		return char.LowPlungeAttack
	default:
		return char.Walk
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

func advanceCoreFrame(t *testing.T, c *core.Core) {
	c.F++
	c.Tick()

	if c.F > 10000 {
		t.Errorf("Test did not complete after 10000 frames. Check if IgnoreBurstEnergy is correctly taken into account. Otherwise check if all actions can be exec'd")
	}
}
