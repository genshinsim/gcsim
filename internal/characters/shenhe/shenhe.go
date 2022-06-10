package shenhe

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

const (
	normalHitNum = 5
	quillKey     = "shenhequill"
)

func init() {
	initCancelFrames()
	core.RegisterCharFunc(keys.Shenhe, NewChar)
}

type char struct {
	*tmpl.Character
	quillcount []int
	c4count    int
	c4expiry   int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)
	c.EnergyMax = 80
	c.Weapon.Class = weapon.WeaponClassSpear
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3
	c.CharZone = character.ZoneLiyue
	c.Base.Element = attributes.Cryo

	c.c4count = 0
	c.c4expiry = 0

	if c.Base.Cons >= 1 {
		c.SetNumCharges(action.ActionSkill, 2)
	}

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.quillcount = make([]int, len(c.Core.Player.Chars()))
	c.quillDamageMod()
	if c.Base.Cons >= 4 {
		c.c4()
	}
	return nil
}

func initCancelFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 23)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 19)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 42)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 30)
	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 81)
	attackFrames[4][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it

	// charge -> x
	chargeFrames = frames.InitAbilSlice(49)
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark

	// skill (press) -> x
	skillPressFrames = frames.InitAbilSlice(31)

	// skill (hold) -> x
	skillHoldFrames = frames.InitAbilSlice(44)

	// burst -> x
	burstFrames = frames.InitAbilSlice(99)
}
