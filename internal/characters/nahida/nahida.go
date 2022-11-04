package nahida

import (
	"fmt"

	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

func init() {
	core.RegisterCharFunc(keys.Nahida, NewChar)
}

type char struct {
	*tmpl.Character
	pyroCount     int
	electroCount  int
	hydroCount    int
	pyroBurstBuff []float64
	a4Buff        []float64
	c4Buff        []float64
	c6count       int
	markCount     int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 50
	c.NormalHitNum = normalHitNum

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	//skill hooks
	c.Core.Events.Subscribe(event.OnEnemyDamage, c.triKarmaOnBloomDamage, "nahida-tri-karma")
	for i := event.ReactionEventStartDelim; i < event.ReactionEventEndDelim; i++ {
		c.Core.Events.Subscribe(i, c.triKarmaOnReaction(i), fmt.Sprintf("nahida-tri-karma-on-%v", i))
	}
	//burst ele counts
	for _, char := range c.Core.Player.Chars() {
		switch char.Base.Element {
		case attributes.Pyro:
			c.pyroCount++
		case attributes.Hydro:
			c.hydroCount++
		case attributes.Electro:
			c.electroCount++
		}
	}

	if c.Base.Cons > 0 {
		c.c1()
	}

	if c.Base.Cons > 4 {
		c.c4Buff = make([]float64, attributes.EndStatType)
		c.c4()
	}

	c.pyroBurstBuff = make([]float64, attributes.EndStatType)
	if c.pyroCount > 0 {
		c.pyroBurstBuff[attributes.DmgP] = burstTriKarmaDmgBonus[c.pyroCount-1][c.TalentLvlBurst()]
	}

	c.a4Buff = make([]float64, attributes.EndStatType)
	c.a4()
	c.a4tick()

	if c.Base.Cons > 1 {
		c.c2()
	}

	if c.Base.Cons > 5 {
		c.c6()
	}

	//sanity check
	if c.pyroCount > 2 {
		c.pyroCount = 2
	}
	if c.hydroCount > 2 {
		c.hydroCount = 2
	}
	if c.electroCount > 2 {
		c.electroCount = 2
	}

	return nil
}
