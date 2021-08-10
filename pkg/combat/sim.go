package combat

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/genshinsim/gsim/pkg/core"
	"github.com/genshinsim/gsim/pkg/monster"
	"github.com/genshinsim/gsim/pkg/queue"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type SimStats struct {
	//these follow 4 are available in avg mode as well
	Mode                 string                   `json:"mode"`
	CharNames            []string                 `json:"char_names"`
	DamageByChar         []map[string]float64     `json:"damage_by_char"`
	CharActiveTime       []int                    `json:"char_active_time"`
	AbilUsageCountByChar []map[string]int          `json:"abil_usage_count_by_char"`
	ReactionsTriggered   map[core.ReactionType]int `json:"reactions_triggered"`
	SimDuration          int                       `json:"sim_duration"`
	//final result
	Damage float64 `json:"damage"`
	DPS    float64 `json:"dps"`
}

const (
	maxStam      = 240
	jumpFrames   = 33
	dashFrames   = 24
	swapFrames   = 20
	swapCDFrames = 60
)

type Sim struct {
	f    int
	skip int
	cfg  core.Config
	rand *rand.Rand
	log  *zap.SugaredLogger

	//action related
	swapCD       int
	lastStamUse  int
	stam         float64
	stamModifier []func(a core.ActionType) float64
	querer       core.Querer
	queue        []core.ActionItem

	//characters
	charPos            map[string]int
	chars              []core.Character
	active             int
	charActiveDuration int
	status             map[string]int

	//enemies
	targets []core.Target

	//combat
	onAttackWillLand []attackWillLandHook
	onAttackLanded   []attackLandedHook
	onAmpReaction    []onReactionDamageHook
	onTransReaction  []onReactionDamageHook
	onReaction       []onReactionHook
	onTargetDefeated []defeatHook

	//shields
	shields          []core.Shield
	DRFunc           []func() float64
	ShieldBonusFunc  []func() float64
	IncHealBonusFunc []func() float64 //% to add to amount healed

	//constructs
	constructs  []core.Construct
	consNoLimit []core.Construct

	//hurt event
	lastHurt    int
	nextHurt    int
	nextHurtAmt float64
	hurt        core.HurtEvent
	onHurt      []func(s core.Sim)

	//initializing
	initHooks []func()

	//event hooks
	eventHooks [][]eHook

	//flags
	flags core.Flags

	//result
	stats SimStats
}

func NewSim(cfg core.Config) (*Sim, error) {
	var err error
	s := &Sim{}
	if cfg.FixedRand {
		s.rand = rand.New(rand.NewSource(0))
	} else {
		s.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	s.f = -1
	s.cfg = cfg
	s.stam = maxStam

	err = s.initMaps()
	if err != nil {
		return nil, err
	}
	err = s.initLogger(cfg.LogConfig)
	if err != nil {
		return nil, err
	}
	s.flags.HPMode = cfg.Mode.HPMode
	err = s.initTargets(cfg)
	if err != nil {
		return nil, err
	}
	err = s.initChars(cfg)
	if err != nil {
		return nil, err
	}
	err = s.initQueuer(cfg)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Sim) initLogger(cfg core.LogConfig) error {
	config := zap.NewDevelopmentConfig()
	config.Encoding = "json"
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	switch cfg.LogLevel {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	}
	config.EncoderConfig.TimeKey = ""
	config.EncoderConfig.StacktraceKey = ""
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	if !cfg.LogShowCaller {
		config.EncoderConfig.CallerKey = ""
	}
	if cfg.LogFile != "" {
		config.OutputPaths = []string{cfg.LogFile}
	}

	zaplog, err := config.Build()
	if err != nil {
		return err
	}

	s.log = zaplog.Sugar()
	return nil
}

func (s *Sim) initTargets(cfg core.Config) error {
	s.targets = make([]core.Target, len(cfg.Targets))
	for i := 0; i < len(cfg.Targets); i++ {
		t := monster.New(i, s, s.log, cfg.Mode.HP, cfg.Targets[i])
		// t.AddOnReactionHook("stats", func(ds *def.Snapshot) {
		// 	s.stats.ReactionsTriggered[ds.ReactionType]++
		// })
		s.targets[i] = t
	}
	return nil
}

func (s *Sim) initChars(cfg core.Config) error {
	dup := make(map[string]bool)
	res := make(map[core.EleType]int)

	count := len(cfg.Characters.Profile)

	if count > 4 {
		return fmt.Errorf("more than 4 characters in a team detected")
	}

	s.stats.CharNames = make([]string, count)
	s.stats.DamageByChar = make([]map[string]float64, count)
	s.stats.CharActiveTime = make([]int, count)
	s.stats.AbilUsageCountByChar = make([]map[string]int, count)

	s.active = -1
	for i, v := range cfg.Characters.Profile {
		//call new char function

		f, ok := charMap[v.Base.Name]
		if !ok {
			return fmt.Errorf("invalid character: %v", v.Base.Name)
		}
		c, err := f(s, s.log, v)
		if err != nil {
			return err
		}

		s.chars = append(s.chars, c)
		s.charPos[v.Base.Name] = i

		if v.Base.Name == cfg.Characters.Initial {
			s.active = i
		}

		if _, ok := dup[v.Base.Name]; ok {
			return fmt.Errorf("duplicated character %v", v.Base.Name)
		}
		dup[v.Base.Name] = true

		//initialize weapon
		wf, ok := weaponMap[v.Weapon.Name]
		if !ok {
			return fmt.Errorf("unrecognized weapon %v for character %v", v.Weapon.Name, v.Base.Name)
		}
		wf(c, s, s.log, v.Weapon.Refine, v.Weapon.Param)

		//add set bonus
		for key, count := range v.Sets {
			f, ok := setMap[key]
			if ok {
				f(c, s, s.log, count)
			} else {
				s.log.Warnf("character %v has unrecognized set %v", v.Base.Name, key)
			}
		}

		//track resonance
		res[v.Base.Element]++

		//setup maps
		s.stats.DamageByChar[i] = make(map[string]float64)
		s.stats.AbilUsageCountByChar[i] = make(map[string]int)
		s.stats.CharNames[i] = v.Base.Name

	}

	s.initResonance(res)

	return nil
}

func (s *Sim) initResonance(count map[core.EleType]int) {
	for k, v := range count {
		if v >= 2 {
			switch k {
			case core.Pyro:
				s.log.Debugw("adding pyro resonance", "frame", s.f, "event", core.LogSimEvent)
				for _, c := range s.chars {
					val := make([]float64, core.EndStatType)
					val[core.ATKP] = 0.25
					c.AddMod(core.CharStatMod{
						Key: "pyro-res",
						Amount: func(a core.AttackTag) ([]float64, bool) {
							return val, true
						},
						Expiry: -1,
					})
				}
			case core.Hydro:
				//heal not implemented yet
				s.log.Debugw("adding hydro resonance", "frame", s.f, "event", core.LogSimEvent)
				s.log.Warnw("hydro resonance not implemented", "event", core.LogSimEvent)
			case core.Cryo:
				s.log.Debugw("adding cryo resonance", "frame", s.f, "event", core.LogSimEvent)
				s.AddOnAttackWillLand(func(t core.Target, ds *core.Snapshot) {
					if t.AuraType() == core.Cryo {
						ds.Stats[core.CR] += .15
						s.log.Debugw("cryo resonance + 15% crit pre damage (cryo)", "frame", s.f, "event", core.LogCalc, "char", ds.ActorIndex, "next", ds.Stats[core.CR])
					}
					if t.AuraType() == core.Frozen {
						ds.Stats[core.CR] += .15
						s.log.Debugw("cryo resonance + 15% crit pre damage  (frozen)", "frame", s.f, "event", core.LogCalc, "char", ds.ActorIndex, "next", ds.Stats[core.CR])
					}
				}, "cryo res")
			case core.Electro:
				s.log.Debugw("adding electro resonance", "frame", s.f, "event", core.LogSimEvent)
				last := 0
				s.AddOnReaction(func(t core.Target, ds *core.Snapshot) {
					switch ds.ReactionType {
					case core.Melt:
						return
					case core.Vaporize:
						return
					}
					if s.f-last < 300 && last != 0 { // every 5 seconds
						return
					}
					s.DistributeParticle(core.Particle{
						Source: "electro res",
						Num:    1,
						Ele:    core.Electro,
					})
					last = s.f
				}, "electro res")
			case core.Geo:
				s.log.Debugw("adding geo resonance", "frame", s.f, "event", core.LogSimEvent)
				//Increases shield strength by 15%. Additionally, characters protected by a shield will have the
				//following special characteristics:
				//	DMG dealt increased by 15%, dealing DMG to enemies will decrease their Geo RES by 20% for 15s.
				s.AddShieldBonus(func() float64 {
					return 0.15 //shield bonus always active
				})
				s.AddOnAttackLanded(func(t core.Target, ds *core.Snapshot, dmg float64, crit bool) {
					if !s.IsShielded() {
						return
					}
					t.AddResMod("geo res", core.ResistMod{
						Duration: 15 * 60,
						Ele:      core.Geo,
						Value:    -0.2,
					})
				}, "geo res")

				for _, c := range s.chars {
					val := make([]float64, core.EndStatType)
					val[core.DmgP] = 0.15
					c.AddMod(core.CharStatMod{
						Key: "geo-res",
						Amount: func(a core.AttackTag) ([]float64, bool) {
							return val, s.IsShielded()
						},
						Expiry: -1,
					})
				}

			case core.Anemo:
				s.log.Debugw("adding anemo resonance", "frame", s.f, "event", core.LogSimEvent)
				for _, c := range s.chars {
					c.AddCDAdjustFunc(core.CDAdjust{
						Key:    "anemo-res",
						Amount: func(a core.ActionType) float64 { return -0.05 },
						Expiry: -1,
					})
				}
			}
		}
	}
}

func (s *Sim) initMaps() error {
	s.eventHooks = make([][]eHook, core.EndEventHook)
	s.flags.Custom = make(map[string]int)

	s.status = make(map[string]int)
	s.chars = make([]core.Character, 0, core.MaxTeamPlayerCount)
	s.charPos = make(map[string]int)

	//combat stuff
	s.onAttackLanded = make([]attackLandedHook, 0, 10)
	s.onAttackWillLand = make([]attackWillLandHook, 0, 10)
	s.onReaction = make([]onReactionHook, 0, 10)
	s.onAmpReaction = make([]onReactionDamageHook, 0, 10)
	s.onTransReaction = make([]onReactionDamageHook, 0, 10)
	s.onTargetDefeated = make([]defeatHook, 0, 10)

	//shield stuff
	s.shields = make([]core.Shield, 0, core.EndShieldType)
	s.DRFunc = make([]func() float64, 0, 5)
	s.ShieldBonusFunc = make([]func() float64, 0, 5)

	//log stuff
	s.stats.ReactionsTriggered = make(map[core.ReactionType]int)

	//qeueu stuff
	s.queue = make([]core.ActionItem, 0, 10)

	return nil
}

func (s *Sim) initQueuer(cfg core.Config) error {
	cust := make(map[string]int)
	for i, v := range cfg.Rotation {
		if v.Name != "" {
			cust[v.Name] = i
		}
		// log.Println(v.Conditions)
	}
	for i, v := range cfg.Rotation {
		if _, ok := s.charPos[v.Target]; !ok {
			return fmt.Errorf("invalid char in rotation %v", v.Target)
		}
		cfg.Rotation[i].Last = -1
	}

	s.querer = queue.New(
		s,
		cfg.Rotation,
		s.log,
	)
	return nil
}

func (s *Sim) RestoreStam(v float64) {
	s.stam += v
	if s.stam > maxStam {
		s.stam = maxStam
	}
}

func (s *Sim) SwapCD() int      { return s.swapCD }
func (s *Sim) Stam() float64    { return s.stam }
func (s *Sim) Frame() int       { return s.f }
func (s *Sim) Rand() *rand.Rand { return s.rand }

func (s *Sim) TargetHasResMod(key string, param int) bool {
	if param >= len(s.targets) {
		return false
	}
	return s.targets[param].HasResMod(key)
}
func (s *Sim) TargetHasDefMod(key string, param int) bool {
	if param >= len(s.targets) {
		return false
	}
	return s.targets[param].HasDefMod(key)
}

func (s *Sim) TargetHasElement(ele core.EleType, param int) bool {
	if param >= len(s.targets) {
		return false
	}
	return s.targets[param].AuraContains(ele)
}
