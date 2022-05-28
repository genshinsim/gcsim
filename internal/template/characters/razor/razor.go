package razor

import (
	"github.com/genshinsim/gcsim/internal/frames"
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

type char struct {
	*tmpl.Character
	sigils         int
	sigilsDuration int
}

func init() {
	core.RegisterCharFunc(keys.Razor, NewChar)
}

func NewChar(s *core.Core, w *character.CharWrapper, p character.CharacterProfile) error {
	c := char{}
	t := tmpl.New(s)
	t.CharWrapper = w
	c.Character = t

	c.Base.Element = attributes.Electro

	e, ok := p.Params["start_energy"]
	if !ok {
		e = 80
	}
	c.Energy = float64(e)
	c.EnergyMax = 80
	c.Weapon.Class = weapon.WeaponClassClaymore
	c.BurstCon = 3
	c.SkillCon = 5
	c.NormalHitNum = 4
	c.CharZone = character.ZoneMondstadt

	w.Character = &c

	return nil
}

func (c *char) InitCancelFrames() {
	// TODO: need to update frames
	c.initNormalCancels()

	skillPressFrames = frames.InitAbilSlice(74)
	skillHoldFrames = frames.InitAbilSlice(92)
}

func (c *char) ActionStam(a action.Action, p map[string]int) float64 {
	switch a {
	case action.ActionDash:
		return 18
	case action.ActionCharge:
		// NOT IMPLEMENTED
		return 0
	default:
		c.Core.Log.NewEvent("ActionStam not implemented", glog.LogActionEvent, c.Index, "action", a.String())
		return 0
	}
}

func (c *char) EnergySigil() {
	val := make([]float64, attributes.EndStatType)
	c.AddStatMod("er-sigil", -1, attributes.ER, func() ([]float64, bool) {
		if c.Core.F > c.sigilsDuration {
			return nil, false
		}

		val[attributes.ER] = float64(c.sigils) * 0.2
		return val, true
	})
}
