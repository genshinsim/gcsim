package jean

import (
	tmpl "github.com/genshinsim/gcsim/internal/template/character"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterCharFunc(keys.Jean, NewChar)
}

type char struct {
	*tmpl.Character
	c2buff []float64
}

func NewChar(s *core.Core, w *character.CharWrapper, _ profile.CharacterProfile) error {
	c := char{}
	c.Character = tmpl.NewWithWrapper(s, w)

	c.EnergyMax = 80
	c.NormalHitNum = normalHitNum
	c.BurstCon = 3
	c.SkillCon = 5

	w.Character = &c

	return nil
}

func (c *char) Init() error {
	// C2:
	// When Jean picks up an Elemental Orb/Particle, all party members have their Movement SPD and ATK SPD increased by 15% for 15s.
	if c.Base.Cons >= 2 {
		c.c2buff = make([]float64, attributes.EndStatType)
		c.c2buff[attributes.AtkSpd] = 0.15
		c.Core.Events.Subscribe(event.OnParticleReceived, func(args ...interface{}) bool {
			// only trigger if Jean catches the particle
			if c.Core.Player.Active() != c.Index {
				return false
			}
			// apply C2 to all characters
			for _, this := range c.Core.Player.Chars() {
				this.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag("jean-c2", 900),
					AffectedStat: attributes.AtkSpd,
					Amount: func() ([]float64, bool) {
						return c.c2buff, true
					},
				})
			}
			return false
		}, "jean-c2")
	}
	if c.Base.Cons >= 6 {
		c.Core.Log.NewEvent("jean c6 not implemented", glog.LogCharacterEvent, c.Index)
	}
	return nil
}
