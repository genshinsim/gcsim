package kokomi

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

const normalHitNum = 3

func init() {
	initCancelFrames()
	core.RegisterCharFunc(keys.Kokomi, NewChar)
}

type char struct {
	*tmpl.Character
	skillFlatDmg  float64
	skillLastUsed int
	swapEarlyF    int
	c4ICDExpiry   int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	t := tmpl.New(s)
	t.CharWrapper = w
	c.Character = t

	c.Base.Element = attributes.Hydro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 70
	}
	c.Energy = float64(e)
	c.EnergyMax = 70
	c.Weapon.Class = weapon.WeaponClassCatalyst
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5
	c.CharZone = character.ZoneInazuma

	c.skillFlatDmg = 0
	c.skillLastUsed = 0
	c.swapEarlyF = 0
	c.c4ICDExpiry = 0

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()
	c.passive()
	c.onExitField()
	c.burstActiveHook()
	return nil
}

func initCancelFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 30)
	attackFrames[0][action.ActionAttack] = 14
	attackFrames[0][action.ActionCharge] = 19

	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 34)
	attackFrames[1][action.ActionAttack] = 30

	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 65)
	attackFrames[2][action.ActionCharge] = 60
	attackFrames[2][action.ActionWalk] = 60

	// charge -> x
	chargeFrames = frames.InitAbilSlice(76)
	chargeFrames[action.ActionAttack] = 62
	chargeFrames[action.ActionCharge] = 62
	chargeFrames[action.ActionSkill] = 62
	chargeFrames[action.ActionBurst] = 62
	chargeFrames[action.ActionDash] = chargeHitmark
	chargeFrames[action.ActionJump] = chargeHitmark
	chargeFrames[action.ActionSwap] = 62

	// skill -> x
	skillFrames = frames.InitAbilSlice(61)
	skillFrames[action.ActionDash] = 29
	skillFrames[action.ActionJump] = 29

	// burst -> x
	burstFrames = frames.InitAbilSlice(77)
}
