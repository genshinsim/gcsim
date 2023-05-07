package kaveh

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

func init() {
	core.RegisterCharFunc(keys.Kaveh, NewChar)
}

type char struct {
	*tmpl.Character
	a4Stacks int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.BurstCon = 3
	c.SkillCon = 5
	c.EnergyMax = 80
	c.NormalHitNum = len(attackHitmarks)

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	if c.Base.Cons >= 4 {
		c.c4()
	}
	if c.Base.Cons >= 6 {
		c.c6()
	}
	c.a1()
	c.a4AddStacksHandler()
	c.addBurstExitHandler()
	return nil
}

func (c *char) Snapshot(ai *combat.AttackInfo) combat.Snapshot {
	ds := c.Character.Snapshot(ai)

	if c.StatModIsActive(burstKey) {
		switch ai.AttackTag {
		case attacks.AttackTagNormal,
			attacks.AttackTagPlunge,
			attacks.AttackTagExtra:
			ai.Element = attributes.Dendro
		}
	}

	return ds
}
