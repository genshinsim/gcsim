package chongyun

import (
	"github.com/genshinsim/gcsim/internal/frames"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

const normalHitNum = 4

func init() {
	initCancelFrames()
	core.RegisterCharFunc(keys.Chongyun, NewChar)
}

type char struct {
	*tmpl.Character
	fieldSrc int
	a4Snap   *combat.AttackEvent
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Cryo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 40
	}
	c.Energy = float64(e)
	c.EnergyMax = 40
	c.Weapon.Class = weapon.WeaponClassClaymore
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = character.ZoneLiyue

	c.fieldSrc = -601

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.onSwapHook()
	if c.Base.Cons >= 4 {
		c.c4()
	}
	if c.Base.Cons >= 6 && c.Core.Combat.DamageMode {
		c.c6()
	}
	return nil
}

func initCancelFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 24)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 38)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 62)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 80)

	// skill -> x
	skillFrames = frames.InitAbilSlice(57)

	// burst -> x
	burstFrames = frames.InitAbilSlice(79)
}
