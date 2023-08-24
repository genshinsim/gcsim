package ayaka

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Ayaka, NewChar)
}

type char struct {
	*tmpl.Character
	c6CDTimerAvail bool // Flag that controls whether the 0.5 C6 CD timer is available to be started
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 80
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = normalHitNum

	c.c6CDTimerAvail = false

	// Start with C6 ability active
	if c.Base.Cons >= 6 {
		c.c6CDTimerAvail = true
	}

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	// Start with C6 ability active
	if c.Base.Cons >= 6 {
		c.c6AddBuff()
	}
	return nil
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	switch a {
	case action.ActionDash:
		f, ok := p["f"]
		if !ok {
			return 10 //tap = 36 frames, so under 1 second
		}
		//for every 1 second passed, consume extra 15
		extra := f / 60
		return float64(10 + 15*extra)
	}
	return c.Character.ActionStam(a, p)
}
