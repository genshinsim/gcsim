//Package player contains player related tracking and functionalities:
// - tracking characters on the team
// - handling animations state
// - handling normal attack state
// - handling character stats and attributes
// - handling shielding
package player

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/animation"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/infusion"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/core/task"
)

type Handler struct {
	log    glog.Logger
	events event.Eventter
	tasks  task.Tasker
	f      *int
	//handlers
	*animation.AnimationHandler
	Shields *shield.Handler
	infusion.InfusionHandler

	//tracking
	chars   []*character.CharWrapper
	active  int
	charPos map[keys.Char]int

	//stam
	Stam            float64
	LastStamUse     int
	stamPercentMods []stamPercentMod

	//last action
	LastAction struct {
		UsedAt int
		Type   action.Action
		Param  map[string]int
		Char   int
	}
}

func New(f *int, log glog.Logger, events event.Eventter, tasks task.Tasker, debug bool) *Handler {
	h := &Handler{
		chars:           make([]*character.CharWrapper, 0, 4),
		charPos:         make(map[keys.Char]int),
		stamPercentMods: make([]stamPercentMod, 0, 5),
		log:             log,
		events:          events,
		tasks:           tasks,
		f:               f,
	}
	h.Shields = shield.New(f, log, events)
	h.InfusionHandler = infusion.New(f, log, debug)
	h.AnimationHandler = animation.New(f, log, events, tasks)
	return h
}

func (h *Handler) AddChar(char *character.CharWrapper) int {
	h.chars = append(h.chars, char)
	index := len(h.chars) - 1
	char.SetIndex(index)
	h.charPos[char.Base.Key] = index

	return index
}

func (h *Handler) ByIndex(i int) *character.CharWrapper {
	return h.chars[i]
}

func (h *Handler) CombatByIndex(i int) combat.Character {
	return h.chars[i]
}

func (h *Handler) ByKey(k keys.Char) *character.CharWrapper {
	return h.chars[h.charPos[k]]
}

func (h *Handler) Chars() []*character.CharWrapper {
	return h.chars
}

func (h *Handler) Active() int {
	return h.active
}

func (h *Handler) ActiveChar() *character.CharWrapper {
	return h.chars[h.active]
}

func (h *Handler) SetActive(i int) {
	h.active = i
}

func (h *Handler) Adjust(src string, char int, amt float64) {
	h.chars[char].AddEnergy(src, amt)
}

func (h *Handler) ResetAllNormalCounter() {
	for _, char := range h.chars {
		char.ResetNormalCounter()
	}
}

func (h *Handler) DistributeParticle(p character.Particle) {
	for i, char := range h.chars {
		char.ReceiveParticle(p, h.active == i, len(h.chars))
	}
	h.events.Emit(event.OnParticleReceived, p)
}

func (h *Handler) AbilStamCost(i int, a action.Action, p map[string]int) float64 {
	return h.StamPercentMod(action.ActionDash) * h.chars[i].ActionStam(action.ActionDash, p)
}

//InitializeTeam will set up resonance event hooks and calculate
//all character base stats
func (h *Handler) InitializeTeam() error {
	var err error
	for _, c := range h.chars {
		err = c.UpdateBaseStats()
		if err != nil {
			return err
		}
	}
	//loop again to initialize
	for i := range h.chars {
		err = h.chars[i].Init()
		if err != nil {
			return err
		}
		h.chars[i].Equip.Weapon.Init()
		for k := range h.chars[i].Equip.Sets {
			h.chars[i].Equip.Sets[k].Init()
		}
	}
	return nil
}

func (h *Handler) Tick() {
	//	- player (stamina, swap, animation, etc...)
	//		- character
	//		- shields
	//		- animation
	//		- stamina
	//		- swap
	h.Shields.Tick()
	for _, c := range h.chars {
		c.Tick()
	}
	h.AnimationHandler.Tick()

}
