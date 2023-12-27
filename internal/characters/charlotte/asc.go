package charlotte

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}

	count := 0
	c.Core.Events.Subscribe(event.OnTargetDied, func(args ...interface{}) bool {
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return false
		}
		if !t.StatusIsActive(skillHoldMarkKey) {
			return false
		}
		if count == 4 {
			return false
		}
		if count == 0 {
			c.QueueCharTask(func() {
				count = 0
			}, 720)
		}
		count++
		c.ReduceActionCooldown(action.ActionSkill, 120)
		return false
	}, "charlotte-a1")
}

func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	heal := 0
	cryop := 0
	for _, this := range c.Core.Player.Chars() {
		if c.Index == this.Index {
			continue
		}
		if this.CharZone == info.ZoneFontaine {
			heal++
		} else {
			cryop++
		}
	}
	m := make([]float64, attributes.EndStatType)
	m[attributes.Heal] = 0.05 * float64(heal)
	m[attributes.CryoP] = 0.05 * float64(cryop)
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("charlotte-a4", -1),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return m, true
		},
	})
}
