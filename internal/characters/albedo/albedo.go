package albedo

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
	core.RegisterCharFunc(keys.Albedo, NewChar)
}

type char struct {
	*tmpl.Character
	lastConstruct   int
	skillAttackInfo combat.AttackInfo
	skillSnapshot   combat.Snapshot
	bloomSnapshot   combat.Snapshot
	icdSkill        int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Geo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 40
	}
	c.Energy = float64(e)
	c.EnergyMax = 40
	c.Weapon.Class = weapon.WeaponClassSword
	c.NormalHitNum = normalHitNum

	c.icdSkill = 0

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.skillHook()
	if c.Base.Cons >= 4 {
		c.c4()
	}
	if c.Base.Cons == 6 {
		c.c6()
	}
	return nil
}

func initCancelFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 12)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 18)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 29)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 39)
	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[3], 54)

	// charge -> x
	chargeFrames = frames.InitAbilSlice(54)
	chargeFrames[action.ActionDash] = chargeHitmarks[len(chargeHitmarks)-1]
	chargeFrames[action.ActionJump] = chargeHitmarks[len(chargeHitmarks)-1]

	// skill -> x
	skillFrames = frames.InitAbilSlice(32)

	// burst -> x
	burstFrames = frames.InitAbilSlice(96)
}
