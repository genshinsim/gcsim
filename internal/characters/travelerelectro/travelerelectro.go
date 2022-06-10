package travelerelectro

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
	core.RegisterCharFunc(keys.TravelerElectro, NewChar)
}

type char struct {
	*tmpl.Character
	abundanceAmulets      int
	burstC6Hits           int
	burstC6WillGiveEnergy bool
	burstSnap             combat.Snapshot
	burstAtk              *combat.AttackEvent
	burstSrc              int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Electro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 80
	}
	c.Energy = float64(e)
	c.EnergyMax = 80
	c.Weapon.Class = weapon.WeaponClassSword
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = normalHitNum

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.burstProc()
	return nil
}

func initCancelFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 13)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 25)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 33)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 52)
	attackFrames[4] = frames.InitNormalCancelSlice(attackHitmarks[4], 40)
	attackFrames[4][action.ActionCharge] = 500 //TODO: this action is illegal; need better way to handle it

	// charge -> x
	// TODO: charge not implemented
	//chargeFrames = frames.InitAbilSlice(41)
	//chargeFrames[action.ActionDash] = chargeHitmarks[len(chargeHitmarks)-1]
	//chargeFrames[action.ActionJump] = chargeHitmarks[len(chargeHitmarks)-1]
	//chargeFrames[action.ActionSwap] = chargeHitmarks[len(chargeHitmarks)-1]

	// skill -> x
	skillFrames = frames.InitAbilSlice(55)

	// burst -> x
	burstFrames = frames.InitAbilSlice(60)
}
