package celestialgift

import (
	"fmt"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const (
	duration          = 60 * 20
	set4Key           = "celestialgift-4pc"
	set4KeyHolder     = set4Key + "-holder"
	set4KeyActive     = set4Key + "-active"
	set4KeyActiveBuff = set4Key + "-active-buff"
)

func init() {
	core.RegisterSetFunc(keys.CelestialGift, NewSet)
}

type Set struct {
	Index int
	Count int
	buff  []float64
	core  *core.Core
	char  *character.CharWrapper
}

func (s *Set) SetIndex(idx int) { s.Index = idx }
func (s *Set) GetCount() int    { return s.Count }
func (s *Set) Init() error {
	count := s.Count
	core := s.core

	if count < 4 || !s.char.IsHexerei {
		return nil
	}

	holderBuffIsActiveOfElement := func(e attributes.Element) bool {
		for _, char := range core.Player.Chars() {
			if char.StatusIsActive(set4KeyHolder) && char.Base.Element == e {
				return true
			}
		}
		return false
	}

	// returns false if there is another party member of the same element with the 4pc buff active,
	// otherwise returns true
	shouldRefreshElemBuff := func() bool {
		selfEle := s.char.Base.Element
		selfIndex := s.char.Index()
		for _, char := range core.Player.Chars() {
			if char.StatusIsActive(set4KeyHolder) && char.Base.Element == selfEle && char.Index() != selfIndex {
				return false
			}
		}
		return true
	}

	// active char status is active if Hexerei Secret Rite is active and there is a party member with the 4pc buff
	// don't need to check Secret Rite because the active char stat mod isn't added when not Secret Rite
	activeStatusIsActive := func() bool {
		for _, char := range core.Player.Chars() {
			if char.StatusIsActive(set4KeyActive) {
				return true
			}
		}
		return false
	}

	activeStatusDuration := func() int {
		dur := 0
		for _, char := range core.Player.Chars() {
			dur = max(dur, char.StatusDuration(set4KeyActive))
		}
		return dur
	}

	deleteActiveBuff := func() {
		for _, char := range core.Player.Chars() {
			char.DeleteStatus(set4KeyActiveBuff)
		}
	}

	elems := map[attributes.Element]bool{}
	for _, char := range core.Player.Chars() {
		elems[char.Base.Element] = true
	}

	buffStrength := 0.2
	isHexereiSecretRite := core.Player.GetHexereiCount() >= 2
	if isHexereiSecretRite {
		buffStrength = 0.4
	}

	core.Events.Subscribe(event.OnSkill, func(args ...any) {
		// NOTE: Description says "The equipping character can trigger this effect while off-field."
		// but currently no characters can use skills while off field, so this is ignored
		if core.Player.Active() != s.char.Index() {
			return
		}

		// we always refresh the duration of the active char buff
		s.char.AddStatus(set4KeyActive, duration, true)

		elem := s.char.Base.Element
		if !shouldRefreshElemBuff() {
			return
		}

		s.char.AddStatus(set4KeyHolder, duration, true)
		buffedDmgP := attributes.EleToDmgP(elem)
		for _, otherChar := range core.Player.Chars() {
			otherChar.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(set4KeyHolder+"-"+elem.String(), duration),
				AffectedStat: buffedDmgP,
				Amount: func() []float64 {
					clear(s.buff)
					s.buff[buffedDmgP] = buffStrength
					return s.buff
				},
			})
		}
		// need to remove the active char buff because the holder is currently on field
		// and is now applying a buff of their own element
		deleteActiveBuff()
	}, fmt.Sprintf(set4KeyHolder+"-%v", s.char.Base.Key.String()))

	if !isHexereiSecretRite {
		return nil
	}

	core.Events.Subscribe(event.OnCharacterSwap, func(args ...any) {
		deleteActiveBuff()
		remainingDuration := activeStatusDuration()
		if remainingDuration == 0 {
			return
		}
		elem := core.Player.ActiveChar().Base.Element

		// The active char buff doesn't stack with the holder buff
		// Even if the holder buff expires, the active char buff isn't reapplied until
		// the next swap
		if holderBuffIsActiveOfElement(elem) {
			return
		}

		buffedDmgP := attributes.EleToDmgP(elem)
		for _, otherChar := range core.Player.Chars() {
			// This stat mod needs hitlag on both the buff reciever, and the holder
			// in order for the stat mod to actually be extended
			otherChar.AddStatMod(character.StatMod{
				Base:         modifier.NewBaseWithHitlag(set4KeyActiveBuff, remainingDuration),
				AffectedStat: buffedDmgP,
				Amount: func() []float64 {
					clear(s.buff)
					s.buff[buffedDmgP] = buffStrength
					if activeStatusIsActive() {
						return s.buff
					}
					return nil
				},
			})
		}
	}, set4KeyActive) // OnSwap subscription is same for every holder, so we allow it to be overwritten
	return nil
}

func NewSet(c *core.Core, char *character.CharWrapper, count int, param map[string]int) (info.Set, error) {
	s := Set{
		Count: count,
		buff:  make([]float64, attributes.EndStatType),
		core:  c,
		char:  char,
	}

	if count < 2 {
		return &s, nil
	}

	m := make([]float64, attributes.EndStatType)
	m[attributes.ER] = 0.2
	char.AddStatMod(character.StatMod{
		Base:         modifier.NewBase("celestialgift-2pc", -1),
		AffectedStat: attributes.ER,
		Amount: func() []float64 {
			return m
		},
	})

	return &s, nil
}
