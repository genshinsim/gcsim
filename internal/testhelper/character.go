package testhelper

import (
	"github.com/genshinsim/gcsim/internal/tmpl/enemy"
	"github.com/genshinsim/gcsim/internal/tmpl/player"
	"github.com/genshinsim/gcsim/pkg/core"
)

func testUseAction(a core.ActionType, c *core.Core, x core.Character, count int) {
	p := make(map[string]int)
	repeat := func(abil func(p map[string]int) (int, int)) {
		for i := 0; i < count; i++ {
			f, _ := abil(p)
			for j := 0; j < f; j++ {
				c.Tick()
			}
		}
	}
	switch a {
	case core.ActionAttack:
		repeat(x.Attack)
	case core.ActionAim:
		repeat(x.Aimed)
	case core.ActionCharge:
		repeat(x.ChargeAttack)
	case core.ActionLowPlunge:
		repeat(x.LowPlungeAttack)
	case core.ActionHighPlunge:
		repeat(x.HighPlungeAttack)
	case core.ActionSkill:
		repeat(x.Skill)
	case core.ActionBurst:
		repeat(x.Burst)
	}
}

func SkipFrames(c *core.Core, i int) {
	for x := 0; x < i; x++ {
		c.Tick()
	}
}

func setupChar(c *core.Core, x core.Character) {
	//create a basic core with no logger
	c.Chars = append(c.Chars, x)
	c.CharPos[x.Key()] = 0
	c.Init()

	c.Targets = append(c.Targets, player.New(0, c))
	c.Targets = append(c.Targets, enemy.New(1, c, EnemyProfile()))
}

func TestCatalystCharacter(c *core.Core, x core.Character) {
	setupChar(c, x)
	testUseAction(core.ActionSkill, c, x, 1)
	testUseAction(core.ActionBurst, c, x, 1)
	testUseAction(core.ActionAttack, c, x, 10)
	testUseAction(core.ActionCharge, c, x, 1)
	SkipFrames(c, 1200)
}

func TestSwordCharacter(c *core.Core, x core.Character) {
	setupChar(c, x)
	testUseAction(core.ActionSkill, c, x, 1)
	testUseAction(core.ActionBurst, c, x, 1)
	testUseAction(core.ActionAttack, c, x, 10)
	testUseAction(core.ActionCharge, c, x, 1)
	SkipFrames(c, 1200)
}

func TestPolearmCharacter(c *core.Core, x core.Character) {
	setupChar(c, x)
	testUseAction(core.ActionSkill, c, x, 1)
	testUseAction(core.ActionBurst, c, x, 1)
	testUseAction(core.ActionAttack, c, x, 10)
	testUseAction(core.ActionCharge, c, x, 1)
	SkipFrames(c, 1200)
}

func TestClaymoreCharacter(c *core.Core, x core.Character) {
	setupChar(c, x)
	testUseAction(core.ActionSkill, c, x, 1)
	testUseAction(core.ActionBurst, c, x, 1)
	testUseAction(core.ActionAttack, c, x, 10)
	SkipFrames(c, 1200)
}

func TestBowCharacter(c *core.Core, x core.Character) {
	setupChar(c, x)
	testUseAction(core.ActionSkill, c, x, 1)
	testUseAction(core.ActionBurst, c, x, 1)
	testUseAction(core.ActionAttack, c, x, 10)
	testUseAction(core.ActionCharge, c, x, 1)
	testUseAction(core.ActionAim, c, x, 10)
	SkipFrames(c, 1200)
}
