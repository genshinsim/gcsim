package combat

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/genshinsim/gsim/pkg/core"
	"github.com/genshinsim/gsim/pkg/monster"
)

type SimStats struct {
	//these follow 4 are available in avg mode as well
	Mode                 string                    `json:"mode"`
	CharNames            []string                  `json:"char_names"`
	DamageByChar         []map[string]float64      `json:"damage_by_char"`
	CharActiveTime       []int                     `json:"char_active_time"`
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
	// f    int
	skip int
	c    *core.Core
	//hurt event
	lastHurt    int
	nextHurt    int
	nextHurtAmt float64
	//result
	stats SimStats
}

func NewSim(cfg core.Config) (*Sim, error) {
	var err error
	s := &Sim{}

	c, err := core.New(
		func(c *core.Core) error {

			if cfg.FixedRand {
				c.Rand = rand.New(rand.NewSource(0))
			} else {
				c.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
			}
			c.F = -1
			c.Flags.DamageMode = cfg.RunOptions.DamageMode
			c.Log, err = core.NewDefaultLogger(cfg.RunOptions.Debug)
			if err != nil {
				return err
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}
	s.c = c

	err = s.initMaps()
	if err != nil {
		return nil, err
	}
	err = s.initTargets(cfg)
	if err != nil {
		return nil, err
	}
	err = s.initChars(cfg)
	if err != nil {
		return nil, err
	}

	c.Init()

	return s, nil
}

func (s *Sim) initTargets(cfg core.Config) error {
	s.c.Targets = make([]core.Target, len(cfg.Targets))
	for i := 0; i < len(cfg.Targets); i++ {
		s.c.Targets[i] = monster.New(i, s.c, cfg.Targets[i])
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

	s.c.ActiveChar = -1
	for i, v := range cfg.Characters.Profile {
		//call new char function

		f, ok := charMap[v.Base.Name]
		if !ok {
			return fmt.Errorf("invalid character: %v", v.Base.Name)
		}
		c, err := f(s.c, v)
		if err != nil {
			return err
		}

		s.c.Chars = append(s.c.Chars, c)

		if v.Base.Name == cfg.Characters.Initial {
			s.c.ActiveChar = i
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
		wf(c, s.c, v.Weapon.Refine, v.Weapon.Param)

		//add set bonus
		for key, count := range v.Sets {
			f, ok := setMap[key]
			if ok {
				f(c, s.c, count)
			} else {
				s.c.Log.Warnf("character %v has unrecognized set %v", v.Base.Name, key)
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
				s.c.Log.Debugw("adding pyro resonance", "frame", s.c.F, "event", core.LogSimEvent)
				for _, c := range s.c.Chars {
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
				s.c.Log.Debugw("adding hydro resonance", "frame", s.c.F, "event", core.LogSimEvent)
				s.c.Log.Warnw("hydro resonance not implemented", "event", core.LogSimEvent)
			case core.Cryo:
				s.c.Log.Debugw("adding cryo resonance", "frame", s.c.F, "event", core.LogSimEvent)
				s.c.Events.Subscribe(core.OnAttackWillLand, func(t core.Target, ds *core.Snapshot) bool {
					if t.AuraType() == core.Cryo {
						ds.Stats[core.CR] += .15
						s.c.Log.Debugw("cryo resonance + 15% crit pre damage (cryo)", "frame", s.c.F, "event", core.LogCalc, "char", ds.ActorIndex, "next", ds.Stats[core.CR])
					}
					if t.AuraType() == core.Frozen {
						ds.Stats[core.CR] += .15
						s.c.Log.Debugw("cryo resonance + 15% crit pre damage  (frozen)", "frame", s.c.F, "event", core.LogCalc, "char", ds.ActorIndex, "next", ds.Stats[core.CR])
					}
					return false
				}, "cryo res")
			case core.Electro:
				s.c.Log.Debugw("adding electro resonance", "frame", s.c.F, "event", core.LogSimEvent)
				last := 0
				s.c.Events.Subscribe(core.OnReactionOccured, func(t core.Target, ds *core.Snapshot) bool {
					switch ds.ReactionType {
					case core.Melt:
						return false
					case core.Vaporize:
						return false
					}
					if s.c.F-last < 300 && last != 0 { // every 5 seconds
						return false
					}
					s.c.Energy.DistributeParticle(core.Particle{
						Source: "electro res",
						Num:    1,
						Ele:    core.Electro,
					})
					last = s.c.F
					return false
				}, "electro res")
			case core.Geo:
				s.c.Log.Debugw("adding geo resonance", "frame", s.c.F, "event", core.LogSimEvent)
				//Increases shield strength by 15%. Additionally, characters protected by a shield will have the
				//following special characteristics:
				//	DMG dealt increased by 15%, dealing DMG to enemies will decrease their Geo RES by 20% for 15s.
				s.c.Shields.AddBonus(func() float64 {
					return 0.15 //shield bonus always active
				})
				s.c.Events.Subscribe(core.OnDamage, func(t core.Target, ds *core.Snapshot) bool {
					if s.c.Shields.IsShielded() {
						t.AddResMod("geo res", core.ResistMod{
							Duration: 15 * 60,
							Ele:      core.Geo,
							Value:    -0.2,
						})
					}
					return false
				}, "geo res")

				for _, c := range s.c.Chars {
					val := make([]float64, core.EndStatType)
					val[core.DmgP] = 0.15
					c.AddMod(core.CharStatMod{
						Key: "geo-res",
						Amount: func(a core.AttackTag) ([]float64, bool) {
							return val, s.c.Shields.IsShielded()
						},
						Expiry: -1,
					})
				}

			case core.Anemo:
				s.c.Log.Debugw("adding anemo resonance", "frame", s.c.F, "event", core.LogSimEvent)
				for _, c := range s.c.Chars {
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

	//log stuff
	s.stats.ReactionsTriggered = make(map[core.ReactionType]int)

	return nil
}
