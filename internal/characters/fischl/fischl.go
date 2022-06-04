package fischl

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

const normalHitNum = 5

func init() {
	initCancelFrames()
	core.RegisterCharFunc(keys.Fischl, NewChar)
}

type char struct {
	*tmpl.Character
	//field use for calculating oz damage
	ozSnapshot    combat.AttackEvent
	ozSource      int //keep tracks of source of oz aka resets
	ozActiveUntil int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	t := tmpl.New(s)
	t.CharWrapper = w
	c.Character = t

	c.Base.Element = attributes.Electro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 60
	}
	c.Energy = float64(e)
	c.EnergyMax = 60
	c.Weapon.Class = weapon.WeaponClassBow
	c.NormalHitNum = normalHitNum

	c.ozSource = -1
	c.ozActiveUntil = -1

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()
	if c.Base.Cons >= 6 {
		c.c6()
	}
	return nil
}

func initCancelFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	// normal cancels (missing Nx -> Aim)
	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 25)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 22)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 38)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 32)
	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4], 67)

	// aimed -> x
	// aim cancel frames are currently generic, should record specific cancels for each one at some point
	aimedFrames = frames.InitAbilSlice(96)

	// skill -> x
	skillFrames = frames.InitAbilSlice(43)
	skillFrames[action.ActionDash] = 14
	skillFrames[action.ActionJump] = 16
	skillFrames[action.ActionSwap] = 42

	// burst -> x
	burstFrames = frames.InitAbilSlice(148)
	burstFrames[action.ActionDash] = 111
	burstFrames[action.ActionJump] = 115
	burstFrames[action.ActionSwap] = 24
}
