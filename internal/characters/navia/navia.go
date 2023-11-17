package navia

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Navia, NewChar)
}

type char struct {
	*tmpl.Character
	shrapnel          int
	artillerySnapshot combat.AttackEvent
	naviaburst        bool
	c2ready           bool
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)
	c.EnergyMax = 60
	c.NormalHitNum = 4
	c.SkillCon = 3
	c.BurstCon = 5
	c.SetNumCharges(action.ActionSkill, 2)
	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()
	c.shrapnel = 0
	c.ShrapnelGain()
	return nil
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "shrapnel":
		return c.shrapnel, nil
	case "navia-burst":
		return c.naviaburst, nil
	case "artillery":
		return c.naviaburst, nil
	default:
		return c.Character.Condition(fields)
	}
}
