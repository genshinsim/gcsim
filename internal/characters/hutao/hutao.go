package hutao

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Hutao, NewChar)
}

type char struct {
	*tmpl.Character
	a1buff  []float64
	a4buff  []float64
	ppbuff  []float64
	c4buff  []float64
	c6buff  []float64
	applyA1 bool

	burstHealCount  int
	burstHealAmount player.HealInfo
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.onExitField()

	c.a1buff = make([]float64, attributes.EndStatType)
	c.a1buff[attributes.CR] = 0.12

	c.ppbuff = make([]float64, attributes.EndStatType)

	c.a4()

	if c.Base.Cons > 4 {
		c.c4()
	}

	if c.Base.Cons >= 6 {
		c.c6()
	}
	return nil
}

func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(_ ...interface{}) bool {
		if c.StatModIsActive(paramitaBuff) {
			c.a1()
			c.DeleteStatMod(paramitaBuff)
		}
		return false
	}, "hutao-exit")
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	switch a {
	case action.ActionCharge:
		if c.StatModIsActive(paramitaBuff) && c.Base.Cons >= 1 {
			return 0
		}
		return 25
	}
	return c.Character.ActionStam(a, p)
}

func (c *char) Snapshot(ai *combat.AttackInfo) combat.Snapshot {
	ds := c.Character.Snapshot(ai)

	if c.StatModIsActive(paramitaBuff) {
		switch ai.AttackTag {
		case attacks.AttackTagNormal:
		case attacks.AttackTagExtra:
		default:
			return ds
		}
		ai.Element = attributes.Pyro
	}
	return ds
}
