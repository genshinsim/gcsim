//Package player contains player related tracking and functionalities:
// - tracking characters on the team
// - handling animations state
// - handling normal attack state
// - handling character stats and attributes
// - handling shielding
package player

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/animation"
	"github.com/genshinsim/gcsim/pkg/core/player/artifact"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/player/infusion"
	"github.com/genshinsim/gcsim/pkg/core/player/shield"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
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
	weaps   []weapon.Weapon
	weapPos map[keys.Weapon]int
	sets    []artifact.Set
	setPos  map[keys.Set]int

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
		weapPos:         make(map[keys.Weapon]int),
		setPos:          make(map[keys.Set]int),
		stamPercentMods: make([]stamPercentMod, 0, 5),
		log:             log,
		events:          events,
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

func (h *Handler) AddWeapon(key keys.Weapon, w weapon.Weapon) int {
	h.weaps = append(h.weaps, w)
	index := len(h.weaps) - 1
	w.SetIndex(index)
	h.weapPos[key] = index
	return index
}

func (h *Handler) AddSet(key keys.Set, set artifact.Set) int {
	h.sets = append(h.sets, set)
	index := len(h.weaps) - 1
	set.SetIndex(index)
	h.setPos[key] = index
	return index
}

func (h *Handler) ByIndex(i int) *character.CharWrapper {
	return h.chars[i]
}

func (h *Handler) Chars() []*character.CharWrapper {
	return h.chars
}

func (h *Handler) Active() int {
	return h.active
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
		err = h.weaps[i].Init()
		if err != nil {
			return err
		}
		err = h.sets[i].Init()
		if err != nil {
			return err
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
