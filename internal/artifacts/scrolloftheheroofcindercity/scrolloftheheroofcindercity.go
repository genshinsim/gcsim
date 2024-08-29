package scrolloftheheroofcindercity

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

func init() {
	core.RegisterSetFunc(keys.ScrollOfTheHeroOfCinderCity, NewSet)
}

type Set struct {
	Index int
	Count int
}

var reactionElementsArr map[reactions.ReactionType][]attributes.Element
var elementToReactionsArr [][]reactions.ReactionType

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }

func (s *Set) Init() error {
	reactionElementsArr[reactions.Overload] = []attributes.Element{attributes.Electro, attributes.Pyro}
	reactionElementsArr[reactions.Superconduct] = []attributes.Element{attributes.Electro, attributes.Cryo}
	reactionElementsArr[reactions.Melt] = []attributes.Element{attributes.Pyro, attributes.Cryo}
	reactionElementsArr[reactions.Vaporize] = []attributes.Element{attributes.Pyro, attributes.Hydro}
	reactionElementsArr[reactions.Freeze] = []attributes.Element{attributes.Cryo, attributes.Hydro}
	reactionElementsArr[reactions.FreezeExtend] = []attributes.Element{attributes.Cryo, attributes.Hydro}
	reactionElementsArr[reactions.ElectroCharged] = []attributes.Element{attributes.Electro, attributes.Hydro}
	reactionElementsArr[reactions.SwirlHydro] = []attributes.Element{attributes.Anemo, attributes.Hydro}
	reactionElementsArr[reactions.SwirlCryo] = []attributes.Element{attributes.Anemo, attributes.Cryo}
	reactionElementsArr[reactions.SwirlElectro] = []attributes.Element{attributes.Anemo, attributes.Electro}
	reactionElementsArr[reactions.SwirlPyro] = []attributes.Element{attributes.Anemo, attributes.Pyro}
	reactionElementsArr[reactions.CrystallizeHydro] = []attributes.Element{attributes.Geo, attributes.Hydro}
	reactionElementsArr[reactions.CrystallizeCryo] = []attributes.Element{attributes.Geo, attributes.Cryo}
	reactionElementsArr[reactions.CrystallizeElectro] = []attributes.Element{attributes.Geo, attributes.Electro}
	reactionElementsArr[reactions.CrystallizePyro] = []attributes.Element{attributes.Geo, attributes.Pyro}
	reactionElementsArr[reactions.Aggravate] = []attributes.Element{attributes.Dendro, attributes.Electro}
	reactionElementsArr[reactions.Spread] = []attributes.Element{attributes.Dendro}
	reactionElementsArr[reactions.Quicken] = []attributes.Element{attributes.Dendro, attributes.Electro}
	reactionElementsArr[reactions.Bloom] = []attributes.Element{attributes.Dendro, attributes.Hydro}
	reactionElementsArr[reactions.Hyperbloom] = []attributes.Element{attributes.Dendro, attributes.Electro}
	reactionElementsArr[reactions.Burgeon] = []attributes.Element{attributes.Dendro, attributes.Pyro}
	reactionElementsArr[reactions.Burning] = []attributes.Element{attributes.Dendro, attributes.Pyro}

	elementToReactionsArr = make([][]reactions.ReactionType, attributes.EndEleType)

	elementToReactionsArr[attributes.Anemo] = []reactions.ReactionType{reactions.SwirlHydro, reactions.SwirlCryo, reactions.SwirlElectro, reactions.SwirlPyro}
	elementToReactionsArr[attributes.Cryo] = []reactions.ReactionType{reactions.Superconduct, reactions.Freeze, reactions.Melt, reactions.SwirlCryo, reactions.CrystallizeCryo, reactions.FreezeExtend}
	elementToReactionsArr[attributes.Dendro] = []reactions.ReactionType{reactions.Quicken, reactions.Aggravate, reactions.Spread, reactions.Bloom, reactions.Burgeon, reactions.Hyperbloom, reactions.Burning}
	elementToReactionsArr[attributes.Electro] = []reactions.ReactionType{reactions.Superconduct, reactions.ElectroCharged, reactions.Overload, reactions.Quicken, reactions.Aggravate, reactions.Hyperbloom, reactions.SwirlElectro, reactions.CrystallizeElectro}
	elementToReactionsArr[attributes.Geo] = []reactions.ReactionType{reactions.CrystallizeCryo, reactions.CrystallizeElectro, reactions.CrystallizeHydro, reactions.CrystallizePyro}
	elementToReactionsArr[attributes.Hydro] = []reactions.ReactionType{reactions.Freeze, reactions.Vaporize, reactions.ElectroCharged, reactions.Bloom, reactions.SwirlHydro, reactions.CrystallizeHydro, reactions.FreezeExtend}
	elementToReactionsArr[attributes.Pyro] = []reactions.ReactionType{reactions.Vaporize, reactions.Melt, reactions.Overload, reactions.Burgeon, reactions.Burning, reactions.SwirlPyro, reactions.CrystallizePyro}
	return nil
}

func reactionElements(r reactions.ReactionType) []attributes.Element {
	return reactionElementsArr[r]
}

func elementToReactions(e attributes.Element) []reactions.ReactionType {
	return elementToReactionsArr[e]
}

func reactionToEvent(r reactions.ReactionType) event.Event {
	switch r {
	case reactions.Overload:
		return event.OnOverload
	case reactions.Superconduct:
		return event.OnSuperconduct
	case reactions.Melt:
		return event.OnMelt
	case reactions.Vaporize:
		return event.OnVaporize
	case reactions.Freeze, reactions.FreezeExtend:
		return event.OnFrozen
	case reactions.ElectroCharged:
		return event.OnElectroCharged
	case reactions.SwirlHydro:
		return event.OnSwirlHydro
	case reactions.SwirlCryo:
		return event.OnSwirlCryo
	case reactions.SwirlElectro:
		return event.OnSwirlElectro
	case reactions.SwirlPyro:
		return event.OnSwirlPyro
	case reactions.CrystallizeHydro:
		return event.OnCrystallizeHydro
	case reactions.CrystallizeCryo:
		return event.OnCrystallizeCryo
	case reactions.CrystallizeElectro:
		return event.OnCrystallizeElectro
	case reactions.CrystallizePyro:
		return event.OnCrystallizePyro
	case reactions.Aggravate:
		return event.OnAggravate
	case reactions.Spread:
		return event.OnSpread
	case reactions.Quicken:
		return event.OnQuicken
	case reactions.Bloom:
		return event.OnBloom
	case reactions.Hyperbloom:
		return event.OnHyperbloom
	case reactions.Burgeon:
		return event.OnBurgeon
	case reactions.Burning:
		return event.OnBurning
	case reactions.Shatter:
		return event.OnShatter
	default:
		return event.ReactionEventEndDelim
	}
}

func reactionEventToReaction(e event.Event) reactions.ReactionType {
	switch e {
	case event.OnOverload:
		return reactions.Overload
	case event.OnSuperconduct:
		return reactions.Superconduct
	case event.OnMelt:
		return reactions.Melt
	case event.OnVaporize:
		return reactions.Vaporize
	case event.OnFrozen:
		return reactions.Freeze
	case event.OnElectroCharged:
		return reactions.ElectroCharged
	case event.OnSwirlHydro:
		return reactions.SwirlHydro
	case event.OnSwirlCryo:
		return reactions.SwirlCryo
	case event.OnSwirlElectro:
		return reactions.SwirlElectro
	case event.OnSwirlPyro:
		return reactions.SwirlPyro
	case event.OnCrystallizeHydro:
		return reactions.CrystallizeHydro
	case event.OnCrystallizeCryo:
		return reactions.CrystallizeCryo
	case event.OnCrystallizeElectro:
		return reactions.CrystallizeElectro
	case event.OnCrystallizePyro:
		return reactions.CrystallizePyro
	case event.OnAggravate:
		return reactions.Aggravate
	case event.OnSpread:
		return reactions.Spread
	case event.OnQuicken:
		return reactions.Quicken
	case event.OnBloom:
		return reactions.Bloom
	case event.OnHyperbloom:
		return reactions.Hyperbloom
	case event.OnBurgeon:
		return reactions.Burgeon
	case event.OnBurning:
		return reactions.Burning
	case event.OnShatter:
		return reactions.Shatter
	default:
		return reactions.NoReaction
	}
}

func Map[T, V any](ts []T, fn func(T) V) []V {
	result := make([]V, len(ts))
	for i, t := range ts {
		result[i] = fn(t)
	}
	return result
}

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{Count: count}
	// 2 Piece: When a nearby party member triggers a Nightsoul Burst, the equipping
	// character regenerates 6 Elemental Energy.
	if count >= 2 {
		// replace with event.OnNightsoulBurst
		c.Combat.Events.Subscribe(event.OnEnemyDamage, func(args ...interface{}) bool {
			char.AddEnergy("scroll-2pc", 6)
			return false
		}, "scroll-2pc")
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
		buffArrs := make([][]float64, attributes.EndEleType)
		for i := range buffArrs {
			ele := attributes.Element(i)
			buffArrs[ele] = make([]float64, attributes.EndStatType)
			buffArrs[ele][attributes.EleToDmgP(ele)] = 0.12
		}

		buffArrsNightsoul := make([][]float64, attributes.EndEleType)
		for i := range buffArrs {
			ele := attributes.Element(i)
			buffArrsNightsoul[ele] = make([]float64, attributes.EndStatType)
			buffArrsNightsoul[ele][attributes.EleToDmgP(ele)] = 0.28
		}

		reactionList := elementToReactions(char.Base.Element)
		eventList := Map(reactionList, reactionToEvent)
		for _, evt := range eventList {
			react := reactionEventToReaction(evt)
			c.Combat.Events.Subscribe(evt, func(args ...interface{}) bool {
				// core.Log.NewEvent("archaic petra proc'd", glog.LogArtifactEvent, char.Index).
				// 	Write("ele", s.element)
				for _, ele := range reactionElements(react) {
					// Apply mod to all characters
					for _, c := range c.Player.Chars() {
						c.AddStatMod(character.StatMod{
							Base:         modifier.NewBaseWithHitlag(fmt.Sprintf("scroll-4pc-%s", attributes.ElementString[ele]), 15*60),
							AffectedStat: attributes.EleToDmgP(ele),
							Amount: func() ([]float64, bool) {
								return buffArrs[ele], true
							},
						})
					}
				}

				if char.CharZone != info.ZoneNatlan {
					return false
				}

				// if char.nightsoul <= 0 {
				// 	return false
				// }
				for _, ele := range reactionElements(react) {
					// Apply mod to all characters
					for _, c := range c.Player.Chars() {
						c.AddStatMod(character.StatMod{
							Base:         modifier.NewBaseWithHitlag(fmt.Sprintf("scroll-4pc-nightsoul-%s", attributes.ElementString[ele]), 20*60),
							AffectedStat: attributes.EleToDmgP(ele),
							Amount: func() ([]float64, bool) {
								return buffArrsNightsoul[ele], true
							},
						})
					}
				}
				return false
			}, fmt.Sprintf("scroll-4pc-%s", react))
		}
	}

	return &s, nil
}
