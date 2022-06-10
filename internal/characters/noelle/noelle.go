package noelle

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

const normalHitNum = 4

func init() {
	initCancelFrames()
	core.RegisterCharFunc(keys.Noelle, NewChar)
}

type char struct {
	*tmpl.Character
	shieldTimer int
	a4Counter   int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	t := tmpl.New(s)
	t.CharWrapper = w
	c.Character = t

	c.Base.Element = attributes.Geo

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 60
	}
	c.Energy = float64(e)
	c.EnergyMax = 60
	c.Weapon.Class = weapon.WeaponClassClaymore
	c.NormalHitNum = normalHitNum

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a1()
	return nil
}

func initCancelFrames() {
	// NA cancels
	attackFrames = make([][]int, normalHitNum)

	attackFrames[0] = frames.InitNormalCancelSlice(attackHitmarks[0], 38)
	attackFrames[1] = frames.InitNormalCancelSlice(attackHitmarks[1], 46)
	attackFrames[2] = frames.InitNormalCancelSlice(attackHitmarks[2], 31)
	attackFrames[3] = frames.InitNormalCancelSlice(attackHitmarks[3], 107)

	// skill -> x
	skillFrames = frames.InitAbilSlice(78)
	skillFrames[action.ActionAttack] = 12
	skillFrames[action.ActionSkill] = 14 // uses burst frames
	skillFrames[action.ActionBurst] = 14
	skillFrames[action.ActionDash] = 11
	skillFrames[action.ActionJump] = 11
	skillFrames[action.ActionWalk] = 43

	// burst -> x
	burstFrames = frames.InitAbilSlice(121)
	burstFrames[action.ActionAttack] = 83
	burstFrames[action.ActionSkill] = 82
	burstFrames[action.ActionDash] = 81
	burstFrames[action.ActionJump] = 81
	burstFrames[action.ActionWalk] = 90
}

// Noelle Geo infusion can't be overridden, so it must be a snapshot modification rather than a weapon infuse
func (c *char) Snapshot(ai *combat.AttackInfo) combat.Snapshot {
	ds := c.Character.Snapshot(ai)

	if c.Core.Status.Duration("noelleq") > 0 {
		//infusion to attacks only
		switch ai.AttackTag {
		case combat.AttackTagNormal:
		case combat.AttackTagPlunge:
		case combat.AttackTagExtra:
		default:
			return ds
		}
		ai.Element = attributes.Geo
	}

	return ds
}
