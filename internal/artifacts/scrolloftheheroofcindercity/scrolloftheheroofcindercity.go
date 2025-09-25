package scrolloftheheroofcindercity

import (
	"fmt"
	"slices"

	"github.com/genshinsim/gcsim/internal/template/nightsoul"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var reactToElements = map[info.ReactionType][]attributes.Element{
	info.ReactionTypeOverload:           {attributes.Electro, attributes.Pyro},
	info.ReactionTypeSuperconduct:       {attributes.Electro, attributes.Cryo},
	info.ReactionTypeMelt:               {attributes.Pyro, attributes.Cryo},
	info.ReactionTypeVaporize:           {attributes.Pyro, attributes.Hydro},
	info.ReactionTypeFreeze:             {attributes.Cryo, attributes.Hydro},
	info.ReactionTypeElectroCharged:     {attributes.Electro, attributes.Hydro},
	info.ReactionTypeSwirlHydro:         {attributes.Anemo, attributes.Hydro},
	info.ReactionTypeSwirlCryo:          {attributes.Anemo, attributes.Cryo},
	info.ReactionTypeSwirlElectro:       {attributes.Anemo, attributes.Electro},
	info.ReactionTypeSwirlPyro:          {attributes.Anemo, attributes.Pyro},
	info.ReactionTypeCrystallizeHydro:   {attributes.Geo, attributes.Hydro},
	info.ReactionTypeCrystallizeCryo:    {attributes.Geo, attributes.Cryo},
	info.ReactionTypeCrystallizeElectro: {attributes.Geo, attributes.Electro},
	info.ReactionTypeCrystallizePyro:    {attributes.Geo, attributes.Pyro},
	info.ReactionTypeAggravate:          {attributes.Dendro, attributes.Electro},
	info.ReactionTypeSpread:             {attributes.Dendro},
	info.ReactionTypeQuicken:            {attributes.Dendro, attributes.Electro},
	info.ReactionTypeBloom:              {attributes.Dendro, attributes.Hydro},
	info.ReactionTypeHyperbloom:         {attributes.Dendro, attributes.Electro},
	info.ReactionTypeBurgeon:            {attributes.Dendro, attributes.Pyro},
	info.ReactionTypeBurning:            {attributes.Dendro, attributes.Pyro},
}

func init() {
	core.RegisterSetFunc(keys.ScrollOfTheHeroOfCinderCity, NewSet)
}

type Set struct {
	Index int
	Count int

	c             *core.Core
	char          *character.CharWrapper
	buff          []float64
	nightsoulBuff []float64
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error      { return nil }

func (s *Set) buffCB(react info.ReactionType, gadgetEmit bool) func(args ...any) bool {
	return func(args ...any) bool {
		trg := args[0].(info.Target)
		if gadgetEmit && trg.Type() != info.TargettableGadget {
			return false
		}
		if !gadgetEmit && trg.Type() != info.TargettableEnemy {
			return false
		}

		ae := args[1].(*info.AttackEvent)
		if ae.Info.ActorIndex != s.char.Index() {
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
						clear(s.nightsoulBuff)
						s.nightsoulBuff[stat] = 0.28
						return s.nightsoulBuff, true
					},
				})
			}
		}
		return false
	}
}

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		Count:         count,
		c:             c,
		char:          char,
		buff:          make([]float64, attributes.EndStatType),
		nightsoulBuff: make([]float64, attributes.EndStatType),
	}
	// 2 Piece: When a nearby party member triggers a Nightsoul Burst, the equipping
	// character regenerates 6 Elemental Energy.
	if count >= 2 {
		c.Combat.Events.Subscribe(event.OnNightsoulBurst, func(args ...any) bool {
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
		for evt, react := range map[event.Event]info.ReactionType{
			event.OnOverload:           info.ReactionTypeOverload,
			event.OnSuperconduct:       info.ReactionTypeSuperconduct,
			event.OnMelt:               info.ReactionTypeMelt,
			event.OnVaporize:           info.ReactionTypeVaporize,
			event.OnFrozen:             info.ReactionTypeFreeze,
			event.OnElectroCharged:     info.ReactionTypeElectroCharged,
			event.OnSwirlHydro:         info.ReactionTypeSwirlHydro,
			event.OnSwirlCryo:          info.ReactionTypeSwirlCryo,
			event.OnSwirlElectro:       info.ReactionTypeSwirlElectro,
			event.OnSwirlPyro:          info.ReactionTypeSwirlPyro,
			event.OnCrystallizeHydro:   info.ReactionTypeCrystallizeHydro,
			event.OnCrystallizeCryo:    info.ReactionTypeCrystallizeCryo,
			event.OnCrystallizeElectro: info.ReactionTypeCrystallizeElectro,
			event.OnCrystallizePyro:    info.ReactionTypeCrystallizePyro,
			event.OnAggravate:          info.ReactionTypeAggravate,
			event.OnSpread:             info.ReactionTypeSpread,
			event.OnQuicken:            info.ReactionTypeQuicken,
			event.OnBloom:              info.ReactionTypeBloom,
			event.OnHyperbloom:         info.ReactionTypeHyperbloom,
			event.OnBurgeon:            info.ReactionTypeBurgeon,
			event.OnBurning:            info.ReactionTypeBurning,
		} {
			elements := reactToElements[react]
			if !slices.Contains(elements, char.Base.Element) {
				continue
			}
			gadgetEmit := false
			switch react {
			case info.ReactionTypeBurgeon, info.ReactionTypeHyperbloom:
				gadgetEmit = true
			}
			c.Combat.Events.Subscribe(evt, s.buffCB(react, gadgetEmit), fmt.Sprintf("scroll-4pc-%v-%v", react, char.Base.Key.String()))
		}
	}

	return &s, nil
}
