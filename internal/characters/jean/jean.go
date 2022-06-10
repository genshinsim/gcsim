package jean

import (
	"github.com/genshinsim/gcsim/internal/frames"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

const normalHitNum = 5

func init() {
	initCancelFrames()
	core.RegisterCharFunc(keys.Jean, NewChar)
}

type char struct {
	*tmpl.Character
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Anemo
	c.EnergyMax = 80
	c.Weapon.Class = weapon.WeaponClassSword
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	if c.Base.Cons >= 6 {
		c.Core.Log.NewEvent("jean c6 not implemented", glog.LogCharacterEvent, c.Index)
	}
	return nil
}

func initCancelFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 25)
	attackFrames[0][action.ActionAttack] = 22

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 20)
	attackFrames[1][action.ActionAttack] = 14

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 31)
	attackFrames[2][action.ActionAttack] = 28

	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 49)
	attackFrames[3][action.ActionAttack] = 44

	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4], 68)
	attackFrames[4][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it

	// charge -> x
	chargeFrames = frames.InitAbilSlice(57)
	chargeFrames[action.ActionBurst] = 56
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionSwap] = 39

	// skill -> x
	skillFrames = frames.InitAbilSlice(46)
	skillFrames[action.ActionDash] = 28
	skillFrames[action.ActionJump] = 28
	skillFrames[action.ActionSwap] = 45

	// burst -> x
	burstFrames = frames.InitAbilSlice(84)
	burstFrames[action.ActionAttack] = 83
	burstFrames[action.ActionSkill] = 83
	burstFrames[action.ActionDash] = 70
	burstFrames[action.ActionJump] = 70
}

func (c *char) ReceiveParticle(p character.Particle, isActive bool, partyCount int) {
	c.Character.ReceiveParticle(p, isActive, partyCount)
	if c.Base.Cons >= 2 {
		//only pop this if jean is active
		if !isActive {
			return
		}
		m := make([]float64, attributes.EndStatType)
		m[attributes.AtkSpd] = 0.15
		for _, active := range c.Core.Player.Chars() {
			active.AddStatMod("jean-c2", 900, attributes.AtkSpd, func() ([]float64, bool) {
				return m, true
			})
		}
	}
}
