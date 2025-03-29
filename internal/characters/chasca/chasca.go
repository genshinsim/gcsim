package chasca

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/internal/template/nightsoul"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/model"
)

func init() {
	core.RegisterCharFunc(keys.Chasca, NewChar)
}

type char struct {
	*tmpl.Character
	nightsoulState       *nightsoul.State
	nightsoulSrc         int
	partyPHECTypes       []attributes.Element
	partyPHECTypesUnique []attributes.Element
	bulletsToFire        []attributes.Element
	bulletsNext          []attributes.Element
	bulletPool           []attributes.Element
	bulletsCharged       int
	aimSrc               int
	skillParticleICD     bool
	c2Src                int
	c4Src                int
	c6Used               bool
}

func NewChar(s *core.Core, w *character.CharWrapper, _ info.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)
	c.nightsoulState = nightsoul.New(c.Core, c.CharWrapper)
	c.EnergyMax = 60
	c.NormalHitNum = normalHitNum
	c.SkillCon = 3
	c.BurstCon = 5

	w.Character = &c

	c.partyPHECTypesUnique = make([]attributes.Element, 0)

	return nil
}

func (c *char) Init() error {
	types := map[attributes.Element]bool{}
	for _, other := range c.Core.Player.Chars() {
		switch ele := other.Base.Element; ele {
		case attributes.Pyro, attributes.Hydro, attributes.Cryo, attributes.Electro:
			if !types[ele] {
				c.partyPHECTypesUnique = append(c.partyPHECTypesUnique, ele)
			}
			types[ele] = true
			c.partyPHECTypes = append(c.partyPHECTypes, ele)
		}
	}
	c.bulletsNext = make([]attributes.Element, 6)
	c.bulletsToFire = make([]attributes.Element, 6)
	c.loadSkillHoldBullets()
	c.a1DMGBuff()
	c.a4()
	return nil
}

func (c *char) Condition(fields []string) (any, error) {
	switch fields[0] {
	case "bullet1":
		return int(c.bulletsNext[0]), nil
	case "bullet2":
		return int(c.bulletsNext[1]), nil
	case "bullet3":
		return int(c.bulletsNext[2]), nil
	case "bullet4":
		return int(c.bulletsNext[3]), nil
	case "bullet5":
		return int(c.bulletsNext[4]), nil
	case "bullet6":
		return int(c.bulletsNext[5]), nil
	case "nightsoul":
		return c.nightsoulState.Condition(fields)
	default:
		return c.Character.Condition(fields)
	}
}

func (c *char) AnimationStartDelay(k model.AnimationDelayKey) int {
	if c.nightsoulState.HasBlessing() {
		switch k {
		case model.AnimationXingqiuN0StartDelay:
			return 5
		case model.AnimationYelanN0StartDelay:
			return 0
		default:
			return c.Character.AnimationStartDelay(k)
		}
	}
	switch k {
	case model.AnimationXingqiuN0StartDelay:
		return 12
	case model.AnimationYelanN0StartDelay:
		return 5
	default:
		return c.Character.AnimationStartDelay(k)
	}
}
