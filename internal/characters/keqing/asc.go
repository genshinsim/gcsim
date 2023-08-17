package keqing

import (
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

// After recasting Stellar Restoration while a Lightning Stiletto is present, Keqing's weapon gains an Electro Infusion for 5s.
func (c *char) a1() {
	if c.Base.Ascension < 1 {
		return
	}
	// starts on E Recast start
	// should actually be 5.6s
	// should barely cover EE N1 Q N1 E 2N1C / 5N1C N1
	// very susceptible to break
	dur := 5*60 + 20
	c.AddStatus("keqinginfuse", dur, true)
	c.Core.Player.AddWeaponInfuse(
		c.Index,
		"keqing-a1",
		attributes.Electro,
		dur,
		true,
		attacks.AttackTagNormal, attacks.AttackTagExtra, attacks.AttackTagPlunge,
	)
}

// When casting Starward Sword, Keqing's CRIT Rate is increased by 15%, and her Energy Recharge is increased by 15%. This effect lasts for 8s.
func (c *char) a4() {
	if c.Base.Ascension < 4 {
		return
	}
	c.AddStatMod(character.StatMod{
		Base:         modifier.NewBaseWithHitlag("keqing-a4", 480),
		AffectedStat: attributes.NoStat,
		Amount: func() ([]float64, bool) {
			return c.a4buff, true
		},
	})
}
