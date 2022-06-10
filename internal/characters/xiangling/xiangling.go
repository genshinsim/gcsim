package xiangling

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

const normalHitNum = 5

func init() {
	initCancelFrames()
	core.RegisterCharFunc(keys.Xiangling, NewChar)
}

type char struct {
	*tmpl.Character
	guoba *panda
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Pyro
	c.EnergyMax = 80
	c.Weapon.Class = weapon.WeaponClassSpear
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = character.ZoneLiyue

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	//add in a guoba
	c.guoba = newGuoba(c.Core)
	c.Core.Combat.AddTarget(c.guoba)
	return nil
}

func initCancelFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 20)

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 17)

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 28)
	attackFrames[2][action.ActionCharge] = 24

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][3], 37)
	attackFrames[3][action.ActionCharge] = 34

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 70)
	attackFrames[4][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it

	// charge -> x
	chargeFrames = frames.InitAbilSlice(69)
	chargeFrames[action.ActionAttack] = 67
	chargeFrames[action.ActionBurst] = 67
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionSwap] = 66

	// skill -> x
	skillFrames = frames.InitAbilSlice(39)
	skillFrames[action.ActionDash] = 14
	skillFrames[action.ActionJump] = 14
	skillFrames[action.ActionSwap] = 38

	// burst -> x
	burstFrames = frames.InitAbilSlice(80)
	burstFrames[action.ActionSwap] = 79
}
