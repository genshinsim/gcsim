package yaemiko

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
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
	kitsuneDetectionRadius float64
	kitsunes               []*kitsune
	totemParticleICD       int
	c4buff                 []float64
}

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 90
	c.NormalHitNum = normalHitNum
	c.BurstCon = 5
	c.SkillCon = 3

	c.SetNumCharges(action.ActionSkill, 3)

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	c.a4()
	if c.Base.Cons >= 2 {
		c.kitsuneDetectionRadius = 20
	} else {
		c.kitsuneDetectionRadius = 12.5
	}
	if c.Base.Cons >= 4 {
		c.c4buff = make([]float64, attributes.EndStatType)
		c.c4buff[attributes.ElectroP] = .20
	}
	return nil
}
