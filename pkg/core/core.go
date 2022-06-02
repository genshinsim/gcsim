// Package core provides core functionality for a simulation:
//	- combat
//	- tasks
//	- event handling
//	- logging
// 	- constructs (really should be just generic objects?)
//	- status
package core

import (
	"fmt"
	"math/rand"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/construct"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/core/status"
	"github.com/genshinsim/gcsim/pkg/core/task"
)

type Core struct {
	F     int
	Flags Flags
	Rand  *rand.Rand
	//various functionalities of core
	Log        glog.Logger    //we use an interface here so that we can pass in a nil logger for all except 1 run
	Events     *event.Handler //track events: subscribe/unsubscribe/emit
	Status     *status.Handler
	Tasks      *task.Handler
	Combat     *combat.Handler
	Constructs *construct.Handler
	Player     *player.Handler
}

type Flags struct {
	LogDebug bool // Used to determine logging level
	Custom   map[string]int
}
type Coord struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	R float64 `json:"r"`
}

type Reactable interface {
	React(a *combat.AttackEvent)
	AuraContains(e ...attributes.Element) bool
	AuraType() attributes.Element
	Tick()
}

type Enemy interface {
	AddResistMod(key string, dur int, ele attributes.Element, val float64)
	DeleteResistMod(key string)
	ResistModIsActive(key string) bool
	AddDefMod(key string, dur int, val float64)
	DeleteDefMod(key string)
	DefModIsActive(key string) bool
}

const MaxTeamSize = 4

func New(seed int64, debug bool) (*Core, error) {
	c := &Core{}
	c.Rand = rand.New(rand.NewSource(seed))
	c.Flags.Custom = make(map[string]int)
	if debug {
		c.Log = glog.New(&c.F, 500)
	} else {
		c.Log = &glog.NilLogger{}
	}

	c.Events = event.New()
	c.Status = status.New(&c.F, c.Log)
	c.Tasks = task.New(&c.F)
	c.Constructs = construct.New(&c.F, c.Log)
	c.Player = player.New(&c.F, c.Log, c.Events, c.Tasks, debug)
	c.Combat = combat.New(c.Log, c.Events, c.Player, false)

	return c, nil
}

func (c *Core) Init() error {
	var err error
	//setup list
	//	- resonance
	//	- on hit energy
	//	- base stats
	//	- char inits
	//	- init call backs
	c.SetupResonance()
	c.SetupOnNormalHitEnergy()
	err = c.Player.InitializeTeam()
	if err != nil {
		return err
	}

	c.Events.Emit(event.OnInitialize)
	return nil
}

func (c *Core) Tick() {
	// things to tick:
	//	- targets
	//	- constructs
	//	- player (stamina, swap, animation, etc...)
	//		- character
	//		- shields
	//		- animation
	//		- stamina
	//		- swap
	//	- tasks
	c.Combat.Tick()
	c.Constructs.Tick()
	c.Player.Tick()

}

func (c *Core) AddChar(p character.CharacterProfile) (int, error) {
	var err error

	// initialize character
	char := character.New(p, &c.F, c.Flags.LogDebug, c.Log, c.Events, c.Tasks)

	f, ok := charMap[p.Base.Key]
	if !ok {
		return -1, fmt.Errorf("invalid character: %v", p.Base.Key.String())
	}
	err = f(c, char, p)
	if err != nil {
		return -1, err
	}
	index := c.Player.AddChar(char)

	// initialize weapon
	wf, ok := weaponMap[p.Weapon.Key]
	if !ok {
		return -1, fmt.Errorf("unrecognized weapon %v for character %v", p.Weapon.Key, p.Base.Key.String())
	}
	weap, err := wf(c, char, p.Weapon)
	if err != nil {
		return -1, err
	}
	char.SetWeapon(weap)

	//set bonus
	total := 0
	for key, count := range p.Sets {
		total += count
		af, ok := setMap[key]
		if ok {
			s, err := af(c, char, count, p.SetParams[key])
			if err != nil {
				return -1, err
			}
			char.SetArtifactSet(key, s)
		} else {
			return -1, fmt.Errorf("character %v has unrecognized artifact: %v", p.Base.Name, key)
		}
	}
	//TODO: this should be handled by parser
	if total > 5 {
		return -1, fmt.Errorf("total set count cannot exceed 5, got %v", total)
	}

	return index, nil
}
