package cooldown

import (
	"testing"
	"time"

	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	_ "github.com/genshinsim/gcsim/pkg/simulation"
	"github.com/genshinsim/gcsim/pkg/testhelper"
)

// the charge queue panics if charges drop below zero. casting on cooldown is handled in SetCD
func TestIgnoreSkillCooldown(t *testing.T) {
	c, char := makeCore(t, keys.Xiangling, true)
	p := make(map[string]int)

	// xiangling skill cd is 720 frames. becoming ready well inside that proves the flag works
	for i := range 5 {
		waitReady(t, c, keys.Xiangling, p, animationBudget)
		c.Player.Exec(action.ActionSkill, keys.Xiangling, p)
		// skill cd is set on a delay. let it land before checking
		for range cdDelayBudget {
			advance(c)
		}
		if got := char.Cooldown(action.ActionSkill); got != 0 {
			t.Fatalf("cast %v: cooldown = %v, want 0", i, got)
		}
	}
}

// klee has 2 skill charges. the queue worker panics if charges exceed the max
func TestIgnoreSkillCooldownMultiCharge(t *testing.T) {
	c, char := makeCore(t, keys.Klee, true)
	p := make(map[string]int)

	for range 5 {
		waitReady(t, c, keys.Klee, p, animationBudget)
		c.Player.Exec(action.ActionSkill, keys.Klee, p)
	}
	if got := char.Charges(action.ActionSkill); got != 2 {
		t.Errorf("charges = %v, want 2", got)
	}
}

func TestSkillCooldownAppliesWhenFlagOff(t *testing.T) {
	c, char := makeCore(t, keys.Xiangling, false)
	p := make(map[string]int)

	c.Player.Exec(action.ActionSkill, keys.Xiangling, p)
	for range animationBudget {
		advance(c)
	}
	if char.Cooldown(action.ActionSkill) <= 0 {
		t.Fatal("expected skill to be on cooldown")
	}
	if err := c.Player.ReadyCheck(action.ActionSkill, keys.Xiangling, p); err == nil {
		t.Error("expected skill to still be on cooldown")
	}
}

// long enough for any skill animation. well short of a real skill cooldown
const animationBudget = 200

// covers characters that set their skill cd on a delay
const cdDelayBudget = 30

func waitReady(t *testing.T, c *core.Core, k keys.Char, p map[string]int, limit int) {
	t.Helper()
	for range limit {
		if err := c.Player.ReadyCheck(action.ActionSkill, k, p); err == nil {
			return
		}
		advance(c)
	}
	t.Fatalf("skill not ready after %v frames", limit)
}

func makeCore(t *testing.T, k keys.Char, ignoreSkillCD bool) (*core.Core, *character.CharWrapper) {
	t.Helper()
	c, err := core.New(core.Opt{
		Seed:                time.Now().Unix(),
		Debug:               true,
		IgnoreSkillCooldown: ignoreSkillCD,
	})
	if err != nil {
		t.Fatalf("error creating core: %v", err)
	}
	c.Combat.SetPlayer(avatar.New(c, info.Point{X: 0, Y: 0}, 1))
	e := enemy.New(c, info.EnemyProfile{
		Level:  100,
		Resist: make(map[attributes.Element]float64),
		Pos:    info.Coord{X: 0, Y: 0, R: 1},
	})
	c.Combat.AddEnemy(e)

	idx, err := c.AddChar(testhelper.DefaultProfile(k, keys.DullBlade))
	if err != nil {
		t.Fatalf("error adding char: %v", err)
	}
	c.Player.SetActive(idx)
	if err := c.Init(); err != nil {
		t.Fatalf("error initializing core: %v", err)
	}
	c.Combat.DefaultTarget = e.Key()
	c.QueueParticle("system", 1000, attributes.NoElement, 0)
	advance(c)

	return c, c.Player.ByIndex(idx)
}

func advance(c *core.Core) {
	c.F++
	c.Tick()
}
