package eula

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
	core.RegisterCharFunc(keys.Eula, NewChar)
}

type char struct {
	*tmpl.Character
	grimheartReset  int
	burstCounter    int
	burstCounterICD int
	grimheartICD    int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Cryo
	c.EnergyMax = 80
	c.Weapon.Class = weapon.WeaponClassClaymore
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()
	c.burstStacks()
	c.onExitField()
	if c.Base.Cons >= 4 {
		c.c4()
	}
	return nil
}

func initCancelFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 34)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 36)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][1], 56)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][0], 44)
	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][1], 105)

	// skill (press) -> x
	skillPressFrames = frames.InitAbilSlice(48)
	skillPressFrames[action.ActionAttack] = 31
	skillPressFrames[action.ActionBurst] = 31
	skillPressFrames[action.ActionDash] = 29
	skillPressFrames[action.ActionJump] = 30
	skillPressFrames[action.ActionSwap] = 29

	// skill (hold) -> x
	skillHoldFrames = frames.InitAbilSlice(77)
	skillHoldFrames[action.ActionDash] = 75
	skillHoldFrames[action.ActionJump] = 75
	skillHoldFrames[action.ActionSwap] = 100
	skillHoldFrames[action.ActionWalk] = 75

	// burst -> x
	burstFrames = frames.InitAbilSlice(122)
	burstFrames[action.ActionDash] = 121
	burstFrames[action.ActionJump] = 121
	burstFrames[action.ActionSwap] = 121
	burstFrames[action.ActionWalk] = 117
}

func (c *char) Tick() {
	c.Character.Tick()
	c.grimheartReset--
	if c.grimheartReset == 0 {
		c.Tags["grimheart"] = 0
	}
}
