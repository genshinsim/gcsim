package xingqiu

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
	core.RegisterCharFunc(keys.Xingqiu, NewChar)
}

type char struct {
	*tmpl.Character
	numSwords     int
	nextRegen     bool
	burstCounter  int
	burstTickSrc  int
	orbitalActive bool
	burstSwordICD int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Hydro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 80
	}
	c.Energy = float64(e)
	c.EnergyMax = 80
	c.Weapon.Class = weapon.WeaponClassSword
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = normalHitNum
	c.CharZone = character.ZoneLiyue

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()
	c.burstStateHook()
	return nil
}

func initCancelFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 35)
	attackFrames[0][action.ActionAttack] = 18

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 29)
	attackFrames[1][action.ActionAttack] = 24

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 35)
	attackFrames[2][action.ActionCharge] = 26

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 33)
	attackFrames[3][action.ActionAttack] = 28

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][1], 66)
	attackFrames[4][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it

	// charge -> x
	chargeFrames = frames.InitAbilSlice(58)
	chargeFrames[action.ActionSkill] = 32
	chargeFrames[action.ActionBurst] = 32
	chargeFrames[action.ActionDash] = chargeHitmarks[len(chargeHitmarks)-1]
	chargeFrames[action.ActionJump] = chargeHitmarks[len(chargeHitmarks)-1]
	chargeFrames[action.ActionSwap] = 31

	// skill -> x
	skillFrames = frames.InitAbilSlice(67)
	skillFrames[action.ActionSkill] = 65
	skillFrames[action.ActionDash] = 30
	skillFrames[action.ActionJump] = 34

	// burst -> x
	burstFrames = frames.InitAbilSlice(40)
	burstFrames[action.ActionAttack] = 33
	burstFrames[action.ActionSkill] = 33
	burstFrames[action.ActionDash] = 33
	burstFrames[action.ActionJump] = 33
}
