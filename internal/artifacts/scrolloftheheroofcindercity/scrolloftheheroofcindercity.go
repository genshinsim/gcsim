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
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

var reactToElements = map[model.ReactionType][]attributes.Element{
	model.ReactionTypeOverload:           {attributes.Electro, attributes.Pyro},
	model.ReactionTypeSuperconduct:       {attributes.Electro, attributes.Cryo},
	model.ReactionTypeMelt:               {attributes.Pyro, attributes.Cryo},
	model.ReactionTypeVaporize:           {attributes.Pyro, attributes.Hydro},
	model.ReactionTypeFreeze:             {attributes.Cryo, attributes.Hydro},
	model.ReactionTypeElectroCharged:     {attributes.Electro, attributes.Hydro},
	model.ReactionTypeSwirlHydro:         {attributes.Anemo, attributes.Hydro},
	model.ReactionTypeSwirlCryo:          {attributes.Anemo, attributes.Cryo},
	model.ReactionTypeSwirlElectro:       {attributes.Anemo, attributes.Electro},
	model.ReactionTypeSwirlPyro:          {attributes.Anemo, attributes.Pyro},
	model.ReactionTypeCrystallizeHydro:   {attributes.Geo, attributes.Hydro},
	model.ReactionTypeCrystallizeCryo:    {attributes.Geo, attributes.Cryo},
	model.ReactionTypeCrystallizeElectro: {attributes.Geo, attributes.Electro},
	model.ReactionTypeCrystallizePyro:    {attributes.Geo, attributes.Pyro},
	model.ReactionTypeAggravate:          {attributes.Dendro, attributes.Electro},
	model.ReactionTypeSpread:             {attributes.Dendro},
	model.ReactionTypeQuicken:            {attributes.Dendro, attributes.Electro},
	model.ReactionTypeBloom:              {attributes.Dendro, attributes.Hydro},
	model.ReactionTypeHyperbloom:         {attributes.Dendro, attributes.Electro},
	model.ReactionTypeBurgeon:            {attributes.Dendro, attributes.Pyro},
	model.ReactionTypeBurning:            {attributes.Dendro, attributes.Pyro},
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

func (s *Set) buffCB(react model.ReactionType, gadgetEmit bool) func(args ...interface{}) bool {
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
		for evt, react := range map[event.Event]model.ReactionType{
			event.OnOverload:           model.ReactionTypeOverload,
			event.OnSuperconduct:       model.ReactionTypeSuperconduct,
			event.OnMelt:               model.ReactionTypeMelt,
			event.OnVaporize:           model.ReactionTypeVaporize,
			event.OnFrozen:             model.ReactionTypeFreeze,
			event.OnElectroCharged:     model.ReactionTypeElectroCharged,
			event.OnSwirlHydro:         model.ReactionTypeSwirlHydro,
			event.OnSwirlCryo:          model.ReactionTypeSwirlCryo,
			event.OnSwirlElectro:       model.ReactionTypeSwirlElectro,
			event.OnSwirlPyro:          model.ReactionTypeSwirlPyro,
			event.OnCrystallizeHydro:   model.ReactionTypeCrystallizeHydro,
			event.OnCrystallizeCryo:    model.ReactionTypeCrystallizeCryo,
			event.OnCrystallizeElectro: model.ReactionTypeCrystallizeElectro,
			event.OnCrystallizePyro:    model.ReactionTypeCrystallizePyro,
			event.OnAggravate:          model.ReactionTypeAggravate,
			event.OnSpread:             model.ReactionTypeSpread,
			event.OnQuicken:            model.ReactionTypeQuicken,
			event.OnBloom:              model.ReactionTypeBloom,
			event.OnHyperbloom:         model.ReactionTypeHyperbloom,
			event.OnBurgeon:            model.ReactionTypeBurgeon,
			event.OnBurning:            model.ReactionTypeBurning,
		} {
			elements := reactToElements[react]
			if !slices.Contains(elements, char.Base.Element) {
				continue
			}
			gadgetEmit := false
			switch react {
			case model.ReactionTypeBurgeon, model.ReactionTypeHyperbloom:
				gadgetEmit = true
			}
			c.Combat.Events.Subscribe(evt, s.buffCB(react, gadgetEmit), fmt.Sprintf("scroll-4pc-%v-%v", react, char.Base.Key.String()))
		}
	}

	return &s, nil
}
