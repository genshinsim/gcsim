package core

import (
	"fmt"
	"log"
	"math/rand"
)

const (
	MaxStam            = 240
	StamCDFrames       = 90
	JumpFrames         = 33
	DashFrames         = 24
	WalkFrames         = 1
	SwapCDFrames       = 60
	MaxTeamPlayerCount = 4
	DefaultTargetIndex = 1
)

type Flags struct {
	DamageMode     bool
	EnergyCalcMode bool // Allows Burst Action when not at full Energy, logs current Energy when using Burst
	LogDebug       bool // Used to determine logging level
	ChildeActive   bool // Used for Childe +1 NA talent passive
	Delays         Delays
	// AmpReactionDidOccur bool
	// AmpReactionType     ReactionType
	// NextAttackMVMult    float64 // melt vape multiplier
	// ReactionDamageTriggered bool
	Custom map[string]int
}

type Delays struct {
	Skill  int
	Burst  int
	Attack int
	Charge int
	Aim    int
	Dash   int
	Jump   int
	Swap   int
}

type Core struct {
	//control
	F     int   // current frame
	Flags Flags // global flags
	Rand  *rand.Rand
	Log   LogCtrl

	//core data
	Stam   float64
	SwapCD int

	//core stuff
	// queue        []Command
	stamModifier []stamMod
	lastStamUse  int

	//track characters
	ActiveChar     int             // index of currently active char
	ActiveDuration int             // duration in frames that the current char has been on field for
	Chars          []Character     // array holding all the characters on the team
	CharPos        map[CharKey]int // map of character string name to their index (for quick lookup by name)

	//track targets
	Targets     []Target
	TotalDamage float64 // keeps tracks of total damage dealt for the purpose of final results

	//last action taken by the sim
	LastAction ActionItem

	//tracks the current animation state
	state       AnimationState
	stateExpiry int

	//handlers
	Status     StatusHandler
	Energy     EnergyHandler
	Action     CommandHandler
	Queue      QueueHandler
	Combat     CombatHandler
	Tasks      TaskHandler
	Constructs ConstructHandler
	Shields    ShieldHandler
	Health     HealthHandler
	Events     EventHandler
}

func New() *Core {
	// var err error
	c := &Core{}

	c.CharPos = make(map[CharKey]int)
	c.Flags.Custom = make(map[string]int)
	c.Stam = MaxStam
	c.stamModifier = make([]stamMod, 0, 10)
	//make a default nil writer
	c.Log = &NilLogger{}

	return c
}

func (c *Core) Init() {
	for _, char := range c.Chars {
		char.Init()
	}

	c.Events.Emit(OnInitialize)
}

func (c *Core) AddChar(v CharacterProfile) (Character, error) {
	f, ok := charMap[v.Base.Key]
	if !ok {
		return nil, fmt.Errorf("invalid character: %v", v.Base.Key.String())
	}
	char, err := f(c, v)
	if err != nil {
		return nil, err
	}
	c.Chars = append(c.Chars, char)
	i := len(c.Chars) - 1
	c.CharPos[v.Base.Key] = i
	char.SetIndex(i)

	wf, ok := weaponMap[v.Weapon.Name]
	if !ok {
		return nil, fmt.Errorf("unrecognized weapon %v for character %v", v.Weapon.Name, v.Base.Key.String())
	}
	wk := wf(char, c, v.Weapon.Refine, v.Weapon.Params)
	char.SetWeaponKey(wk)

	//add set bonus
	total := 0
	for key, count := range v.Sets {
		total += count
		f, ok := setMap[key]
		if ok {
			f(char, c, count, v.SetParams[key])
		} else {
			return nil, fmt.Errorf("character %v has unrecognized artifact: %v", v.Base.Name, key)
		}
	}
	if total > 5 {
		return nil, fmt.Errorf("total set count cannot exceed 5, got %v", total)
	}

	err = char.CalcBaseStats()

	return char, err
}

func (c *Core) CharByName(key CharKey) (Character, bool) {
	pos, ok := c.CharPos[key]
	if !ok {
		return nil, false
	}
	return c.Chars[pos], true
}

func (c *Core) Swap(next CharKey) int {
	prev := c.ActiveChar
	c.ActiveChar = c.CharPos[next]
	c.SwapCD = SwapCDFrames
	c.ResetAllNormalCounter()
	c.Events.Emit(OnCharacterSwap, prev, c.ActiveChar)
	//this duration reset needs to be after the hook for spine to behave properly
	c.ActiveDuration = 0
	return 1
}

func (c *Core) AnimationCancelDelay(next ActionType, p map[string]int) int {
	//if last action is jump, dash, swap,
	switch c.LastAction.Typ {
	case ActionSwap:
		fallthrough
	case ActionDash:
		fallthrough
	case ActionJump:
		return 0
	}
	//other wise check with the current character
	return c.Chars[c.ActiveChar].ActionInterruptableDelay(next, p)
}

func (c *Core) UserCustomDelay() int {
	d := 0
	switch c.LastAction.Typ {
	case ActionSkill:
		d = c.Flags.Delays.Skill
	case ActionBurst:
		d = c.Flags.Delays.Burst
	case ActionAttack:
		d = c.Flags.Delays.Attack
	case ActionCharge:
		d = c.Flags.Delays.Charge
	case ActionDash:
		d = c.Flags.Delays.Dash
	case ActionJump:
		d = c.Flags.Delays.Jump
	case ActionSwap:
		d = c.Flags.Delays.Swap
	case ActionAim:
		d = c.Flags.Delays.Aim
	}
	return c.LastAction.Param["delay"] + d
}

func (c *Core) ResetAllNormalCounter() {
	for _, char := range c.Chars {
		char.ResetNormalCounter()
	}
}

func (c *Core) SetCustomFlag(key string, val int) {
	c.Flags.Custom[key] = val
}

func (c *Core) GetCustomFlag(key string) (int, bool) {
	val, ok := c.Flags.Custom[key]
	return val, ok
}

func (c *Core) Skip(n int) {
	for i := 0; i < n; i++ {
		c.Tick()
	}
}

func (c *Core) Tick() {
	//increment frame count
	c.F++
	//tick auras
	for _, t := range c.Targets {
		if t == nil {
			log.Print("unexpected nil target?")
			log.Println(c.Targets)
		}
		t.Tick()
	}
	//tick shields
	c.Shields.Tick()
	//tick constructs
	c.Constructs.Tick()
	//tick characters
	for _, v := range c.Chars {
		v.Tick()
	}
	//run queued tasks
	c.Tasks.Run()
	//recover stamina
	if c.Stam < MaxStam && c.F-c.lastStamUse > StamCDFrames {
		c.Stam += 25.0 / 60
		if c.Stam > MaxStam {
			c.Stam = MaxStam
		}
	}
	//recover swap cd
	if c.SwapCD > 0 {
		c.SwapCD--
	}
	//update activeduration
	c.ActiveDuration++
}
