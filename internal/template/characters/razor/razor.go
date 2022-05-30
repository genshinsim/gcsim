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

func (c *char) Init() error {
	c.InitCancelFrames()

	// skill
	c.EnergySigil()

	// burst
	c.SpeedBurst()
	c.WolfBurst()
	c.onSwapClearBurst()

	c.a4()

	c.c1()

	return nil
}

func (c *char) InitCancelFrames() {
	// TODO: need to update frames
	c.initNormalCancels()

	skillPressFrames = frames.InitAbilSlice(74)
	skillHoldFrames = frames.InitAbilSlice(92)
	burstFrames = frames.InitAbilSlice(62)
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

// When Razor's Energy is below 50%, increases Energy Recharge by 30%.
func (c *char) a4() {
	val := make([]float64, attributes.EndStatType)
	val[attributes.ER] = 0.3
	c.AddStatMod("er-sigil", -1, attributes.ER, func() ([]float64, bool) {
		if c.Energy/c.EnergyMax < 0.5 {
			return nil, false
		}

		return val, true
	})
}
