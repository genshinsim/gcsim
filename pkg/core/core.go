package core

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/genshinsim/gcsim/pkg/core/keys"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	MaxStam            = 240
	StamCDFrames       = 90
	JumpFrames         = 33
	DashFrames         = 24
	SwapFrames         = 1
	SwapCDFrames       = 60
	MaxTeamPlayerCount = 4
	DefaultTargetIndex = 1
)

type Flags struct {
	DamageMode     bool
	EnergyCalcMode bool // Allows Burst Action when not at full Energy, logs current Energy when using Burst
	LogDebug       bool // Used to determine logging level
	ChildeActive   bool // Used for Childe +1 NA talent passive
	// AmpReactionDidOccur bool
	// AmpReactionType     ReactionType
	// NextAttackMVMult    float64 // melt vape multiplier
	// ReactionDamageTriggered bool
	Custom map[string]int
}

type Core struct {
	//control
	F     int   // current frame
	Flags Flags // global flags
	Rand  *rand.Rand
	Log   *zap.SugaredLogger

	//core data
	Stam   float64
	SwapCD int

	//core stuff
	queue        []ActionItem
	stamModifier []func(a ActionType) (float64, bool)
	lastStamUse  int

	//track characters
	ActiveChar     int               // index of currently active char
	ActiveDuration int               // duration in frames that the current char has been on field for
	Chars          []Character       // array holding all the characters on the team
	CharPos        map[keys.Char]int // map of character string name to their index (for quick lookup by name)

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
	Action     ActionHandler
	Queue      QueueHandler
	Combat     CombatHandler
	Tasks      TaskHandler
	Constructs ConstructHandler
	Shields    ShieldHandler
	Health     HealthHandler
	Events     EventHandler
}

func New(cfg ...func(*Core) error) (*Core, error) {
	var err error
	c := &Core{}

	c.CharPos = make(map[keys.Char]int)
	c.Flags.Custom = make(map[string]int)
	c.Stam = MaxStam
	c.stamModifier = make([]func(a ActionType) (float64, bool), 0, 10)
	c.queue = make([]ActionItem, 0, 20)

	for _, f := range cfg {
		err := f(c)
		if err != nil {
			return nil, err
		}
	}

	if c.Log == nil {
		c.Log, err = NewDefaultLogger(c, false, false, nil)
		if err != nil {
			return nil, err
		}
	}

	if c.Rand == nil {
		c.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	if c.Status == nil {
		c.Status = NewStatusCtrl(c)
	}
	if c.Energy == nil {
		c.Energy = NewEnergyCtrl(c)
	}
	if c.Action == nil {
		c.Action = NewActionCtrl(c)
	}
	if c.Queue == nil {
		c.Queue = NewQueueCtr(c)
	}
	if c.Combat == nil {
		c.Combat = NewCombatCtrl(c)
	}
	if c.Tasks == nil {
		c.Tasks = NewTaskCtrl(&c.F)
	}
	if c.Constructs == nil {
		c.Constructs = NewConstructCtrl(c)
	}
	if c.Shields == nil {
		c.Shields = NewShieldCtrl(c)
	}
	if c.Health == nil {
		c.Health = NewHealthCtrl(c)
	}
	if c.Events == nil {
		c.Events = NewEventCtrl(c)
	}

	//check handlers
	return c, nil
}

func (c *Core) Init() {

	for i, char := range c.Chars {
		char.Init(i)
	}

	c.Events.Emit(OnInitialize)
}

func (c *Core) AddChar(v CharacterProfile) error {
	f, ok := charMap[v.Base.Key]
	if !ok {
		return fmt.Errorf("invalid character: %v", v.Base.Key.String())
	}
	char, err := f(c, v)
	if err != nil {
		return err
	}
	c.Chars = append(c.Chars, char)
	c.CharPos[v.Base.Key] = len(c.Chars) - 1

	wf, ok := weaponMap[v.Weapon.Name]
	if !ok {
		return fmt.Errorf("unrecognized weapon %v for character %v", v.Weapon.Name, v.Base.Key.String())
	}
	wf(char, c, v.Weapon.Refine, v.Weapon.Param)

	//add set bonus
	for key, count := range v.Sets {
		f, ok := setMap[key]
		if ok {
			f(char, c, count)
		} else {
			c.Log.Warnf("character %v has unrecognized set %v", v.Base.Key.String(), key)
		}
	}

	return nil
}

func (c *Core) CharByName(key keys.Char) (Character, bool) {
	pos, ok := c.CharPos[key]
	if !ok {
		return nil, false
	}
	return c.Chars[pos], true
}

func (c *Core) Swap(next keys.Char) int {
	prev := c.ActiveChar
	c.ActiveChar = c.CharPos[next]
	c.SwapCD = SwapCDFrames
	c.ResetAllNormalCounter()
	c.Events.Emit(OnCharacterSwap, prev, c.ActiveChar)
	//this duration reset needs to be after the hook for spine to behave properly
	c.ActiveDuration = 0
	return SwapFrames
}

func (c *Core) AnimationCancelDelay(next ActionType) int {
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
	return c.Chars[c.ActiveChar].ActionInterruptableDelay(next)
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

func (c *Core) AddStamMod(f func(a ActionType) (float64, bool)) {
	c.stamModifier = append(c.stamModifier, f)
}

func (c *Core) StamPercentMod(a ActionType) float64 {
	var m float64 = 1
	n := 0
	for _, f := range c.stamModifier {
		v, done := f(a)
		if !done {
			c.stamModifier[n] = f
			n++
		}
		m += v
	}
	c.stamModifier = c.stamModifier[:n]
	return m
}

func (c *Core) RestoreStam(v float64) {
	c.Stam += v
	if c.Stam > MaxStam {
		c.Stam = MaxStam
	}
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

func NewDefaultLogger(c *Core, debug bool, json bool, paths []string) (*zap.SugaredLogger, error) {
	config := zap.NewDevelopmentConfig()
	if json {
		config.Encoding = "json"
	}
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	if debug {
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
		config.OutputPaths = paths
		c.Flags.LogDebug = true
	} else {
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
		config.OutputPaths = []string{}
	}

	config.EncoderConfig.TimeKey = ""
	config.EncoderConfig.StacktraceKey = ""
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.CallerKey = ""

	// config.OutputPaths = []string{"stdout"}

	zaplog, err := config.Build()
	if err != nil {
		return nil, err
	}
	return zaplog.Sugar(), nil
}
