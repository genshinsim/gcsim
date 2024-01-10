package chevreuse

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
)

func init() {
	core.RegisterCharFunc(keys.Chevreuse, NewChar)
}

func (c *char) Init() error {

	// setup for a1
	chars := c.Core.Player.Chars()
	count := make(map[attributes.Element]int)
	for _, this := range chars {
		count[this.Base.Element]++
	}

	c.onlyPyroElectro = count[attributes.Pyro] > 0 && count[attributes.Electro] > 0 && count[attributes.Electro]+count[attributes.Pyro] == len(chars)

	// setup overcharged ball
	c.Core.Events.Subscribe(event.OnOverload, c.AddOverchargedBall, "chev-E")

	// start subscribing for a1/c1
	c.a1()
	c.c1()
	return nil
}

type char struct {
	*tmpl.Character
	onlyPyroElectro bool
	overChargedBall bool
	c1Icd           int
	c2Icd           int
	c4ShotsLeft     int
	c6Stack         int
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.BurstCon = 5
	c.SkillCon = 3
	c.NormalHitNum = 4
	w.Character = &c

	return nil
}
