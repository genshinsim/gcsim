package lyney

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

func init() {
	core.RegisterCharFunc(keys.Lyney, NewChar)
}

type char struct {
	*tmpl.Character
	hats              []*GrinMalkinHat
	maxHatCount       int
	propSurplusStacks int
	pyrotechnicTravel int
	c2Src             int
	c2Stacks          int
}

func NewChar(s *core.Core, w *character.CharWrapper, p profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.NormalCon = 3

	c2Stacks, ok := p.Params["c2_stacks"]
	if !ok {
		c.c2Stacks = 0
	}
	if c2Stacks < 0 {
		c.c2Stacks = 0
	}
	if c2Stacks > 3 {
		c.c2Stacks = 3
	}

	pyrotechnicTravel, ok := p.Params["pyrotechnic_travel"]
	if ok {
		c.pyrotechnicTravel = pyrotechnicTravel
	} else {
		c.pyrotechnicTravel = 35 // TODO: proper frames
	}

	c.maxHatCount = 1
	c.c1()

	c.hats = make([]*GrinMalkinHat, 0, c.maxHatCount)

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()

	c.c2Setup()

	c.onExitField()

	return nil
}
