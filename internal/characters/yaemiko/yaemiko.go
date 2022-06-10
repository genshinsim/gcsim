package yaemiko

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

const (
	yaeTotemCount  = "totems"
	yaeTotemStatus = "yae_oldest_totem_expiry"
)

func init() {
	core.RegisterCharFunc(keys.YaeMiko, NewChar)
}

type char struct {
	*tmpl.Character
	kitsunes         []*kitsune
	totemParticleICD int
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.Base.Element = attributes.Electro
	c.EnergyMax = 90
	c.Weapon.Class = weapon.WeaponClassCatalyst
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3

	c.SetNumCharges(action.ActionSkill, 3)

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()
	return nil
}
