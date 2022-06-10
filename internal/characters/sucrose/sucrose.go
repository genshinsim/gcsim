package sucrose

import (
	"github.com/genshinsim/gcsim/internal/frames"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

const normalHitNum = 4

func init() {
	initCancelFrames()
	core.RegisterCharFunc(keys.Sucrose, NewChar)
}

type char struct {
	*tmpl.Character
	qInfused            attributes.Element
	infuseCheckLocation combat.AttackPattern
	c4Count             int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Anemo
	c.EnergyMax = 80
	c.Weapon.Class = weapon.WeaponClassCatalyst
	c.NormalHitNum = normalHitNum

	c.infuseCheckLocation = combat.NewDefCircHit(0.1, false, combat.TargettableEnemy, combat.TargettablePlayer, combat.TargettableObject)

	if c.Base.Cons >= 1 {
		c.SetNumCharges(action.ActionSkill, 2)
	}

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1()
	return nil
}

func initCancelFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	// TODO: check if hitmarks for NA->CA and CA->CA lines up
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 20)
	attackFrames[0][action.ActionAttack] = 17

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 26)
	attackFrames[1][action.ActionCharge] = 18

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 33)
	attackFrames[2][action.ActionCharge] = 28

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[2], 54)
	attackFrames[3][action.ActionAttack] = 51

	// charge -> x
	chargeFrames = frames.InitAbilSlice(69)
	chargeFrames[action.ActionCharge] = 66
	chargeFrames[action.ActionSkill] = 60
	chargeFrames[action.ActionBurst] = 61
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionSwap] = chargeHitmark // idk if this is correct or not

	// skill -> x
	skillFrames = frames.InitAbilSlice(57)
	skillFrames[action.ActionCharge] = 56
	skillFrames[action.ActionSkill] = 56
	skillFrames[action.ActionDash] = 11
	skillFrames[action.ActionJump] = 11
	skillFrames[action.ActionSwap] = 56

	// burst -> x
	burstFrames = frames.InitAbilSlice(49)
	burstFrames[action.ActionCharge] = 48
	burstFrames[action.ActionSkill] = 48
	burstFrames[action.ActionDash] = 47
	burstFrames[action.ActionJump] = 47
	burstFrames[action.ActionSwap] = 47

	// dash -> x
	dashFrames = frames.InitAbilSlice(24)
	dashFrames[action.ActionSkill] = 1
	dashFrames[action.ActionBurst] = 1
}
