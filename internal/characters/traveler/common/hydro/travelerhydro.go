package hydro

import (
	"github.com/genshinsim/gcsim/internal/characters/traveler/common"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

type Traveler struct {
	*tmpl.Character
	a4Bonus float64
	gender  int
}

func NewTraveler(s *core.Core, w *character.CharWrapper, p info.CharacterProfile, gender int) (*Traveler, error) {
	c := Traveler{
		gender: gender,
	}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Hydro
	c.EnergyMax = 80
	c.BurstCon = 5
	c.SkillCon = 3
	c.HasArkhe = true

	c.NormalHitNum = normalHitNum
	common.TravelerBaseAtkIncrease(c.CharWrapper, p)
	return &c, nil
}

func (c *Traveler) Init() error {
	return nil
}

func (c *Traveler) AnimationStartDelay(k model.AnimationDelayKey) int {
	switch k {
	case model.AnimationXingqiuN0StartDelay:
		if c.gender == 0 {
			return 8
		} else {
			return 7
		}
	default:
		return c.Character.AnimationStartDelay(k)
	}
}
