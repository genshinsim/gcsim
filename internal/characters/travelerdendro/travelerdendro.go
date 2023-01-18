package travelerdendro

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
)

func init() {
	core.RegisterCharFunc(keys.AetherDendro, NewChar(0))
	core.RegisterCharFunc(keys.LumineDendro, NewChar(1))
}

type char struct {
	*tmpl.Character
	burstPos                   combat.Point
	burstRadius                float64
	burstOverflowingLotuslight int
	skillC1                    bool // this variable also ensures that C1 only restores energy once per cast
	burstTransfig              attributes.Element
	gender                     int
}

func NewChar(gender int) core.NewCharacterFunc {
	return func(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
		c := char{
			gender: gender,
		}
		c.Character = tmpl.NewWithWrapper(s, w)

		c.Base.Element = attributes.Dendro
		c.EnergyMax = 80
		c.BurstCon = 5
		c.SkillCon = 3
		c.NormalHitNum = normalHitNum

		w.Character = &c

		return nil
	}
}

func (c *char) Init() error {
	c.a1Init()
	c.a4Init()

	if c.Base.Cons >= 6 {
		c.c6Init()
	}

	c.skillC1 = false
	return nil
}
