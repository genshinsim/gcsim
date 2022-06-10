package bennett

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
	core.RegisterCharFunc(keys.Bennett, NewChar)
}

type char struct {
	*tmpl.Character
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)
	c.EnergyMax = 60
	c.Weapon.Class = weapon.WeaponClassSword
	c.NormalHitNum = normalHitNum
	c.Base.Element = attributes.Pyro

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	if c.Base.Cons >= 2 {
		c.c2()
	}
	return nil
}

func initCancelFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 33)
	attackFrames[0][action.ActionAttack] = 20

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 27)
	attackFrames[1][action.ActionAttack] = 17

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 46)
	attackFrames[2][action.ActionAttack] = 37

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 48)
	attackFrames[3][action.ActionAttack] = 44

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4], 60)
	attackFrames[4][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it

	// charge -> x
	chargeFrames = frames.InitAbilSlice(55)
	chargeFrames[action.ActionSkill] = 41
	chargeFrames[action.ActionBurst] = 41
	chargeFrames[action.ActionDash] = chargeHitmarks[len(chargeHitmarks)-1]
	chargeFrames[action.ActionJump] = chargeHitmarks[len(chargeHitmarks)-1]
	chargeFrames[action.ActionSwap] = 44

	// skill -> x
	skillFrames = make([][]int, 5)

	// skill (press) -> x
	skillFrames[0] = frames.InitAbilSlice(42)
	skillFrames[0][action.ActionDash] = 22
	skillFrames[0][action.ActionJump] = 23
	skillFrames[0][action.ActionSwap] = 41

	// skill (hold=1) -> x
	skillFrames[1] = frames.InitAbilSlice(98)
	skillFrames[1][action.ActionBurst] = 97
	skillFrames[1][action.ActionDash] = 65
	skillFrames[1][action.ActionJump] = 66
	skillFrames[1][action.ActionSwap] = 96

	// skill (hold=1,c4) -> x
	skillFrames[2] = frames.InitAbilSlice(107)
	skillFrames[2][action.ActionDash] = 95
	skillFrames[2][action.ActionJump] = 95
	skillFrames[2][action.ActionSwap] = 106

	// skill (hold=2) -> x
	skillFrames[3] = frames.InitAbilSlice(343)
	skillFrames[3][action.ActionSkill] = 339 // uses burst frames
	skillFrames[3][action.ActionBurst] = 339
	skillFrames[3][action.ActionDash] = 231
	skillFrames[3][action.ActionJump] = 340
	skillFrames[3][action.ActionSwap] = 337

	// skill (hold=2,a4) -> x
	skillFrames[4] = frames.InitAbilSlice(175)
	skillFrames[4][action.ActionDash] = 171
	skillFrames[4][action.ActionJump] = 174
	skillFrames[4][action.ActionSwap] = 175

	// burst -> x
	burstFrames = frames.InitAbilSlice(53)
	burstFrames[action.ActionDash] = 49
	burstFrames[action.ActionJump] = 50
	burstFrames[action.ActionSwap] = 51
}
