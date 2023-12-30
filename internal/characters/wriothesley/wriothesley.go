package wriothesley

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Wriothesley, NewChar)
}

type char struct {
	*tmpl.Character
	caHeal               float64
	a4Stack              int
	c1N5Proc             bool
	c1SkillExtensionProc bool
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.NormalCon = 3
	c.BurstCon = 5

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.onExit()

	c.a4()
	c.c4()

	return nil
}

func (c *char) Snapshot(ai *combat.AttackInfo) combat.Snapshot {
	ds := c.Character.Snapshot(ai)

	// apply skill multiplier
	if c.StatusIsActive(skillKey) && ai.AttackTag == attacks.AttackTagNormal {
		ai.Mult = skill[c.TalentLvlSkill()] * ai.Mult
	}

	return ds
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	if a == action.ActionCharge && (c.a1Ready() || c.c1Ready()) {
		return 0
	}
	return c.Character.ActionStam(a, p)
}
