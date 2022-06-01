package noelle

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
)

func (c *char) a1() {
	icd := 0
	c.Core.Events.Subscribe(event.OnCharacterHurt, func(args ...interface{}) bool {
		if c.Core.F < icd {
			return false
		}
		active := c.Core.Player.ActiveChar()
		if active.HPCurrent/active.MaxHP() >= 0.3 {
			return false
		}
		icd = c.Core.F + 3600
		ai := combat.AttackInfo{
			ActorIndex: c.Index,
			Abil:       "A1 Shield",
			AttackTag:  combat.AttackTagNone,
		}
		snap := c.Snapshot(&ai)

		//add shield
		x := snap.BaseDef*(1+snap.Stats[attributes.DEFP]) + snap.Stats[attributes.DEF]
		c.Core.Player.Shields.Add(&shield.Tmpl{
			Src:        c.Core.F,
			ShieldType: shield.ShieldNoelleA1,
			Name:       "Noelle A1",
			HP:         4 * x,
			Ele:        attributes.Cryo,
			Expires:    c.Core.F + 1200, //20 sec
		})
		return false
	}, "noelle-a1")
}
