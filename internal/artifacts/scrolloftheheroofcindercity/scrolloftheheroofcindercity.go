package scrolloftheheroofcindercity

import (
	"fmt"
	"slices"

	"github.com/genshinsim/gcsim/internal/template/nightsoul"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var reactToElements = map[reactions.ReactionType][]attributes.Element{
	reactions.Overload:           {attributes.Electro, attributes.Pyro},
	reactions.Superconduct:       {attributes.Electro, attributes.Cryo},
	reactions.Melt:               {attributes.Pyro, attributes.Cryo},
	reactions.Vaporize:           {attributes.Pyro, attributes.Hydro},
	reactions.Freeze:             {attributes.Cryo, attributes.Hydro},
	reactions.ElectroCharged:     {attributes.Electro, attributes.Hydro},
	reactions.SwirlHydro:         {attributes.Anemo, attributes.Hydro},
	reactions.SwirlCryo:          {attributes.Anemo, attributes.Cryo},
	reactions.SwirlElectro:       {attributes.Anemo, attributes.Electro},
	reactions.SwirlPyro:          {attributes.Anemo, attributes.Pyro},
	reactions.CrystallizeHydro:   {attributes.Geo, attributes.Hydro},
	reactions.CrystallizeCryo:    {attributes.Geo, attributes.Cryo},
	reactions.CrystallizeElectro: {attributes.Geo, attributes.Electro},
	reactions.CrystallizePyro:    {attributes.Geo, attributes.Pyro},
	reactions.Aggravate:          {attributes.Dendro, attributes.Electro},
	reactions.Spread:             {attributes.Dendro},
	reactions.Quicken:            {attributes.Dendro, attributes.Electro},
	reactions.Bloom:              {attributes.Dendro, attributes.Hydro},
	reactions.Hyperbloom:         {attributes.Dendro, attributes.Electro},
	reactions.Burgeon:            {attributes.Dendro, attributes.Pyro},
	reactions.Burning:            {attributes.Dendro, attributes.Pyro},
}

func init() {
	core.RegisterSetFunc(keys.ScrollOfTheHeroOfCinderCity, NewSet)
}

type Set struct {
	Index int
	Count int

	c    *core.Core
	char *character.CharWrapper
	buff []float64
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

func (s *Set) buffCB(react reactions.ReactionType, gadgetEmit bool) func(args ...interface{}) bool {
	return func(args ...interface{}) bool {
		trg := args[0].(combat.Target)
		if gadgetEmit && trg.Type() != targets.TargettableGadget {
			return false
		}
		if !gadgetEmit && trg.Type() != targets.TargettableEnemy {
			return false
		}

		ae := args[1].(*combat.AttackEvent)
		if ae.Info.ActorIndex != s.char.Index {
			return false
		}

		hasNightsoul := s.char.StatusIsActive(nightsoul.NightsoulBlessingStatus)
		for _, other := range s.c.Player.Chars() {
			elements := reactToElements[react]
			for _, ele := range elements {
				stat := attributes.EleToDmgP(ele)
				other.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag(fmt.Sprintf("scroll-4pc-%v", ele), 15*60),
					AffectedStat: stat,
					Amount: func() ([]float64, bool) {
						clear(s.buff)
						s.buff[stat] = 0.12
						return s.buff, true
					},
				})

				if !hasNightsoul {
					continue
				}
				other.AddStatMod(character.StatMod{
					Base:         modifier.NewBaseWithHitlag(fmt.Sprintf("scroll-4pc-nightsoul-%v", ele), 20*60),
					AffectedStat: stat,
					Amount: func() ([]float64, bool) {
						clear(s.buff)
						s.buff[stat] = 0.28
						return s.buff, true
					},
				})
			}
		}
		return false
	}
}

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		Count: count,
		c:     c,
		char:  char,
		buff:  make([]float64, attributes.EndStatType),
	}
	// 2 Piece: When a nearby party member triggers a Nightsoul Burst, the equipping
	// character regenerates 6 Elemental Energy.
	if count >= 2 {
		c.Combat.Events.Subscribe(event.OnNightsoulBurst, func(args ...interface{}) bool {
			char.AddEnergy("scroll-2pc", 6)
			return false
		}, fmt.Sprintf("scroll-2pc-%v", char.Base.Key.String()))
	}
	// 4 Piece: After the equipping character triggers a reaction related to their
	// Elemental Type, all nearby party members gain a 12% Elemental DMG Bonus for
	// the Elemental Types involved in the elemental reaction for 15s. If the
	// equipping character is in the Nightsoul's Blessing state when triggering this
	// effect, all nearby party members gain an additional 28% Elemental DMG Bonus
	// for the Elemental Types involved in the elemental reaction for 20s. The
	// equipping character can trigger this effect while off-field, and the DMG bonus
	// from Artifact Sets with the same name do not stack.
	if count >= 4 {
		for evt, react := range map[event.Event]reactions.ReactionType{
			event.OnOverload:           reactions.Overload,
			event.OnSuperconduct:       reactions.Superconduct,
			event.OnMelt:               reactions.Melt,
			event.OnVaporize:           reactions.Vaporize,
			event.OnFrozen:             reactions.Freeze,
			event.OnElectroCharged:     reactions.ElectroCharged,
			event.OnSwirlHydro:         reactions.SwirlHydro,
			event.OnSwirlCryo:          reactions.SwirlCryo,
			event.OnSwirlElectro:       reactions.SwirlElectro,
			event.OnSwirlPyro:          reactions.SwirlPyro,
			event.OnCrystallizeHydro:   reactions.CrystallizeHydro,
			event.OnCrystallizeCryo:    reactions.CrystallizeCryo,
			event.OnCrystallizeElectro: reactions.CrystallizeElectro,
			event.OnCrystallizePyro:    reactions.CrystallizePyro,
			event.OnAggravate:          reactions.Aggravate,
			event.OnSpread:             reactions.Spread,
			event.OnQuicken:            reactions.Quicken,
			event.OnBloom:              reactions.Bloom,
			event.OnHyperbloom:         reactions.Hyperbloom,
			event.OnBurgeon:            reactions.Burgeon,
			event.OnBurning:            reactions.Burning,
		} {
			elements := reactToElements[react]
			if !slices.Contains(elements, char.Base.Element) {
				continue
			}
			gadgetEmit := false
			switch react {
			case reactions.Burgeon, reactions.Hyperbloom:
				gadgetEmit = true
			}
			c.Combat.Events.Subscribe(evt, s.buffCB(react, gadgetEmit), fmt.Sprintf("scroll-4pc-%v-%v", react, char.Base.Key.String()))
		}
	}

	return &s, nil
}
