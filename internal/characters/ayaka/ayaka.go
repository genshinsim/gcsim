package ayaka

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
	core.RegisterCharFunc(keys.Ayaka, NewChar)
}

type char struct {
	*tmpl.Character
	icdC1          int
	c6CDTimerAvail bool // Flag that controls whether the 0.5 C6 CD timer is available to be started
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Cryo
	c.EnergyMax = 80
	c.Weapon.Class = weapon.WeaponClassSword
	c.CharZone = character.ZoneInazuma
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = normalHitNum

	c.icdC1 = -1
	c.c6CDTimerAvail = false

	// Start with C6 ability active
	if c.Base.Cons >= 6 {
		c.c6CDTimerAvail = true
	}

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	// Start with C6 ability active
	if c.Base.Cons >= 6 {
		c.c6AddBuff()
	}
	return nil
}

func initCancelFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0][0], 22)
	attackFrames[0][action.ActionAttack] = 9

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1][0], 20)
	attackFrames[1][action.ActionAttack] = 19

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2][0], 32)
	attackFrames[2][action.ActionCharge] = 31

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3][2], 23)
	attackFrames[3][action.ActionAttack] = 22

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4][0], 66)
	attackFrames[4][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it

	// charge -> x
	chargeFrames = frames.InitAbilSlice(71)
	chargeFrames[action.ActionSkill] = 62
	chargeFrames[action.ActionBurst] = 63
	chargeFrames[action.ActionDash] = chargeHitmarks[len(chargeHitmarks)-1]
	chargeFrames[action.ActionJump] = chargeHitmarks[len(chargeHitmarks)-1]
	chargeFrames[action.ActionSwap] = chargeHitmarks[len(chargeHitmarks)-1]

	// skill -> x
	skillFrames = frames.InitAbilSlice(49)
	skillFrames[action.ActionBurst] = 48
	skillFrames[action.ActionDash] = 30
	skillFrames[action.ActionJump] = 32
	skillFrames[action.ActionSwap] = 48

	// burst -> x
	burstFrames = frames.InitAbilSlice(125)
	burstFrames[action.ActionAttack] = 124
	burstFrames[action.ActionDash] = 124
	burstFrames[action.ActionJump] = 114
	burstFrames[action.ActionSwap] = 123

	// dash -> x
	dashFrames = frames.InitAbilSlice(35)
	dashFrames[action.ActionDash] = 30
	dashFrames[action.ActionSwap] = 34
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	switch a {
	case action.ActionDash:
		f, ok := p["f"]
		if !ok {
			return 10 //tap = 36 frames, so under 1 second
		}
		//for every 1 second passed, consume extra 15
		extra := f / 60
		return float64(10 + 15*extra)
	}
	return c.Character.ActionStam(a, p)
}
