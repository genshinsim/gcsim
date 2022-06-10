package diluc

import (
	"github.com/genshinsim/gcsim/internal/frames"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

const normalHitNum = 4

func init() {
	initCancelFrames()
	core.RegisterCharFunc(keys.Diluc, NewChar)
}

type char struct {
	*tmpl.Character
	eCounter int
	eWindow  int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Pyro
	c.EnergyMax = 40
	c.Weapon.Class = weapon.WeaponClassClaymore
	c.NormalHitNum = normalHitNum

	c.eCounter = 0
	c.eWindow = -1

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	if c.Base.Cons >= 1 && c.Core.Combat.DamageMode {
		c.c1()
	}
	if c.Base.Cons >= 2 {
		c.c2()
	}
	if c.Base.Cons >= 4 {
		c.c4()
	}
	return nil
}

func initCancelFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 32)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 46)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 34)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 99)

	// skill -> x
	skillFrames = make([][]int, 3)

	// skill (1st) -> x
	skillFrames[0] = frames.InitAbilSlice(32)
	skillFrames[0][action.ActionSkill] = 31
	skillFrames[0][action.ActionDash] = 24
	skillFrames[0][action.ActionJump] = 24
	skillFrames[0][action.ActionSwap] = 30

	// skill (2nd) -> x
	skillFrames[1] = frames.InitAbilSlice(38)
	skillFrames[1][action.ActionSkill] = 37
	skillFrames[1][action.ActionBurst] = 37
	skillFrames[1][action.ActionDash] = 28
	skillFrames[1][action.ActionJump] = 31
	skillFrames[1][action.ActionSwap] = 36

	// skill (3rd) -> x
	// TODO: missing counts for skill -> skill
	skillFrames[2] = frames.InitAbilSlice(66)
	skillFrames[2][action.ActionAttack] = 58
	skillFrames[2][action.ActionSkill] = 57 // uses burst frames
	skillFrames[2][action.ActionBurst] = 57
	skillFrames[2][action.ActionDash] = 47
	skillFrames[2][action.ActionJump] = 48

	// burst -> x
	burstFrames = frames.InitAbilSlice(141)
	burstFrames[action.ActionAttack] = 140
	burstFrames[action.ActionSkill] = 139
	burstFrames[action.ActionDash] = 139
	burstFrames[action.ActionSwap] = 138
}

func (c *char) ActionReady(a action.Action, p map[string]int) bool {
	// check if it is possible to use next skill
	if a == action.ActionSkill && c.Core.F < c.eWindow {
		return true
	}
	return c.Character.ActionReady(a, p)
}
