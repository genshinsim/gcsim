package ganyu

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

const normalHitNum = 6

func init() {
	initCancelFrames()
	core.RegisterCharFunc(keys.Ganyu, NewChar)
}

type char struct {
	*tmpl.Character
	a1Expiry int
	c4Stacks int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	t := tmpl.New(s)
	t.CharWrapper = w
	c.Character = t

	c.Base.Element = attributes.Cryo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 60
	}
	c.Energy = float64(e)
	c.EnergyMax = 60
	c.Weapon.Class = weapon.WeaponClassBow
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = character.ZoneLiyue

	c.a1Expiry = -1
	c.c4Stacks = 0

	if c.Base.Cons >= 2 {
		c.SetNumCharges(action.ActionSkill, 2)
	}

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	if c.Base.Cons >= 1 {
		c.c1()
	}
	if c.Base.Cons >= 4 {
		c.c4()
	}
	return nil
}

func initCancelFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 19)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 27)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 38)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 37)
	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4], 28)
	attackFrames[5] = frames.InitNormalCancelSlice(attackHitmarks[5], 59)

	// aimed -> x
	aimedFrames = frames.InitAbilSlice(113) // is this 103 or 113?

	// skill -> x
	skillFrames = frames.InitAbilSlice(28)
	skillFrames[action.ActionSwap] = 27

	// burst -> x
	burstFrames = frames.InitAbilSlice(124)
	burstFrames[action.ActionSwap] = 122
}
