package amber

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
	core.RegisterCharFunc(keys.Amber, NewChar)
}

type char struct {
	*tmpl.Character
	bunnies []bunny
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Pyro
	c.EnergyMax = 40
	c.Weapon.Class = weapon.WeaponClassBow
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	c.bunnies = make([]bunny, 0, 2)

	if c.Base.Cons >= 4 {
		c.SetNumCharges(action.ActionSkill, 2)
	}

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1()
	if c.Base.Cons >= 2 {
		c.overloadExplode()
	}
	return nil
}

func initCancelFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 15)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 18)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 39)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 41)
	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4], 42)

	// aimed -> x
	aimedFrames = frames.InitAbilSlice(94)

	// skill -> x
	skillFrames = frames.InitAbilSlice(35)

	// burst -> x
	burstFrames = frames.InitAbilSlice(74)
}
