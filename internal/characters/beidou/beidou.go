package beidou

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
	core.RegisterCharFunc(keys.Beidou, NewChar)
}

type char struct {
	*tmpl.Character
	burstSnapshot combat.Snapshot
	burstAtk      *combat.AttackEvent
	burstSrc      int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Electro
	c.EnergyMax = 80
	c.Weapon.Class = weapon.WeaponClassClaymore
	c.NormalHitNum = normalHitNum
	c.CharZone = character.ZoneLiyue

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.burstProc()
	if c.Base.Cons >= 4 {
		c.c4()
	}
	return nil
}

func initCancelFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 31)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 36)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 54)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 36)
	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4], 96)

	// skill -> x
	skillFrames = frames.InitAbilSlice(45)
	skillFrames[action.ActionAttack] = 44
	skillFrames[action.ActionDash] = 24
	skillFrames[action.ActionJump] = 24
	skillFrames[action.ActionSwap] = 44

	// burst -> x
	burstFrames = frames.InitAbilSlice(58)
	burstFrames[action.ActionAttack] = 55
	burstFrames[action.ActionDash] = 48
	burstFrames[action.ActionJump] = 48
	burstFrames[action.ActionSwap] = 46
}
