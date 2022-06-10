package keqing

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
	stilettoKey  = "keqingstiletto"
)

func init() {
	initCancelFrames()
	core.RegisterCharFunc(keys.Keqing, NewChar)
}

type char struct {
	*tmpl.Character
	c2ICD int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Electro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 40
	}
	c.Energy = float64(e)
	c.EnergyMax = 40
	c.Weapon.Class = weapon.WeaponClassSword
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = character.ZoneLiyue

	w.Character = &c

	return nil
}

func (c *char) Init() error {
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

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 21)
	attackFrames[0][action.ActionAttack] = 15

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 24)
	attackFrames[1][action.ActionAttack] = 16

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 36)
	attackFrames[2][action.ActionAttack] = 27

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][1], 58)
	attackFrames[3][action.ActionAttack] = 31

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 66)
	attackFrames[4][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it

	// charge -> x
	chargeFrames = frames.InitAbilSlice(36)
	chargeFrames[action.ActionSkill] = 35
	chargeFrames[action.ActionBurst] = 35
	chargeFrames[action.ActionDash] = chargeHitmarks[len(chargeHitmarks)-1]
	chargeFrames[action.ActionJump] = chargeHitmarks[len(chargeHitmarks)-1]
	chargeFrames[action.ActionSwap] = chargeHitmarks[len(chargeHitmarks)-1]

	// skill -> x
	skillFrames = frames.InitAbilSlice(37)
	skillFrames[action.ActionAttack] = 36
	skillFrames[action.ActionSkill] = 35
	skillFrames[action.ActionDash] = 21
	skillFrames[action.ActionJump] = 21
	skillFrames[action.ActionSwap] = 28

	// skill (recast) -> x
	skillRecastFrames = frames.InitAbilSlice(43)
	skillRecastFrames[action.ActionAttack] = 42
	skillRecastFrames[action.ActionDash] = 15
	skillRecastFrames[action.ActionJump] = 16
	skillRecastFrames[action.ActionSwap] = 42

	// burst -> x
	burstFrames = frames.InitAbilSlice(124)
	burstFrames[action.ActionDash] = 122
	burstFrames[action.ActionSwap] = 123
}

func (c *char) ActionReady(a action.Action, p map[string]int) bool {
	// check if stiletto is on-field
	if a == action.ActionSkill && c.Core.Status.Duration(stilettoKey) > 0 {
		return true
	}
	return c.Character.ActionReady(a, p)
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	switch a {
	case action.ActionCharge:
		return 25
	}
	return c.Character.ActionStam(a, p)
}
