package travelerelectro

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
	core.RegisterCharFunc(keys.TravelerElectroMale, NewChar(0))
	core.RegisterCharFunc(keys.TravelerElectroFemale, NewChar(1))
}

type char struct {
	*tmpl.Character
	abundanceAmulets      int
	burstC6Hits           int
	burstC6WillGiveEnergy bool
	burstSnap             combat.Snapshot
	burstAtk              *combat.AttackEvent
	burstSrc              int
	female                int
}

func NewChar(isFemale int) core.NewCharacterFunc {
	return func(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
		c := char{
			female: isFemale,
		}
		c.Character = tmpl.NewWithWrapper(s, w)

		c.Base.Element = attributes.Electro
		c.EnergyMax = 80
		c.BurstCon = 3
		c.SkillCon = 5
		c.NormalHitNum = normalHitNum

		w.Character = &c

		return nil
	}
}

func (c *char) Init() error {
	c.burstProc()
	return nil
}
