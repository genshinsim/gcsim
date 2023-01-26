package alhaitham

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

func init() {
	core.RegisterCharFunc(keys.Alhaitham, NewChar)
}

type char struct {
	*tmpl.Character
	recentlyMirrorGain bool
	mirrorCount        int
	lastInfusionSrc    int
	a1ICD              int
	c1ICD              int
	c2Counter          int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 70
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.onExitField()
	if c.Base.Ascension >= 4 {
		c.a4()
	}
	return nil
}
func (c *char) onExitField() {
	c.Core.Events.Subscribe(event.OnCharacterSwap, func(args ...interface{}) bool {
		// do nothing if previous char wasn't alhaitham
		prev := args[0].(int)
		if prev != c.Index {
			return false
		}
		c.mirrorCount = 0
		if c.mirrorCount > 0 {
			c.Core.Log.NewEvent("Alhaitham left the field, mirror lost", glog.LogCharacterEvent, c.Index)
		}

		return false
	}, "alhaitham-exit")
}

func (c *char) Snapshot(ai *combat.AttackInfo) combat.Snapshot {
	ds := c.Character.Snapshot(ai)

	if c.mirrorCount > 0 { //weapon infusion can't be overriden for haitham
		switch ai.AttackTag {
		case combat.AttackTagNormal:
		case combat.AttackTagPlunge:
		case combat.AttackTagExtra:
		default:
			return ds
		}
		ai.Element = attributes.Dendro
	}
	return ds
}
func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "mirrors":
		return c.mirrorCount, nil
	default:
		return c.Character.Condition(fields)
	}
}
