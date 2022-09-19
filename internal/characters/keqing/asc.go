package keqing

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func (c *char) a1() {
	// barely cover 5N1C + N1 combo hitlagless
	dur := 300 + skillRecastHitmark + 12
	c.Core.Player.AddWeaponInfuse(
		c.CharWrapper,
		"keqing-a1",
		attributes.Electro,
		dur,
		true,
		combat.AttackTagNormal, combat.AttackTagExtra, combat.AttackTagPlunge,
	)
}

func (c *char) a4() {
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("keqing-a4", 480),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return c.a4buff, true
		},
	})
}
