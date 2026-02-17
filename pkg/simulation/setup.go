package simulation

import (
	"errors"
	"fmt"

	"github.com/genshinsim/gcsim/pkg/avatar"
	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attacks"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/info"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character"
	"github.com/genshinsim/gcsim/pkg/enemy"
	"github.com/genshinsim/gcsim/pkg/modifier"
)

const nightsoulBurstICDStatus = "nightsoul-burst-icd"

// first is 0 because you can't proc it without the natlan character
var nightsoulBurstICD = []int{0, 18 * 60, 12 * 60, 9 * 60, 9 * 60}

func SetupTargetsInCore(core *core.Core, p info.Point, r float64, targets []info.EnemyProfile) error {
	// s.stats.ElementUptime = make([]map[core.EleType]int, len(s.C.Targets))
	// s.stats.ElementUptime[0] = make(map[core.EleType]int)

	if r == 0 {
		return errors.New("player cannot have 0 radius")
	}
	player := avatar.New(core, p, r)
	core.Combat.SetPlayer(player)

	// add targets
	for i := range targets {
		v := &targets[i]
		if v.Pos.R == 0 {
			return fmt.Errorf("target cannot have 0 radius (index %v): %v", i, v)
		}
		e := enemy.New(core, *v)
		core.Combat.AddEnemy(e)
		// s.stats.ElementUptime[i+1] = make(map[core.EleType]int)
	}

	// default target is closest to player?
	defaultEnemy := core.Combat.ClosestEnemy(player.Pos())
	if defaultEnemy == nil {
		return errors.New("cannot set default target, got nil")
	}
	core.Combat.DefaultTarget = defaultEnemy.Key()

	// initialize player direction
	core.Combat.Player().SetDirection(defaultEnemy.Pos())

	return nil
}

func SetupCharactersInCore(core *core.Core, chars []info.CharacterProfile, initial keys.Char) error {
	if len(chars) > info.MaxChars {
		return fmt.Errorf("cannot have more than %v characters per team", info.MaxChars)
	}
	dup := make(map[keys.Char]bool)

	active := -1
	for i := range chars {
		// if using random stats, ignore all stats except main
		if chars[i].RandomSubstats != nil {
			stats, err := generateRandSubs(chars[i].RandomSubstats, core.Rand)
			if err != nil {
				return err
			}
			chars[i].Stats = stats
			clear(chars[i].StatsByLabel)
		}
		i, err := core.AddChar(chars[i])
		if err != nil {
			return err
		}

		if chars[i].Base.Key == initial {
			core.Player.SetActive(i)
			active = i
		}

		if _, ok := dup[chars[i].Base.Key]; ok {
			return fmt.Errorf("duplicated character %v", chars[i].Base.Key)
		}
		dup[chars[i].Base.Key] = true
	}

	if active == -1 {
		return errors.New("no active character set")
	}

	return nil
}

func SetupResonance(s *core.Core) {
	chars := s.Player.Chars()
	if len(chars) < 4 {
		return // no resonance if less than 4 chars
	}
	// count number of ele first
	count := make(map[attributes.Element]int)
	for _, c := range chars {
		count[c.Base.Element]++
	}

	for k, v := range count {
		if v < 2 {
			continue
		}
		switch k {
		case attributes.Pyro:
			m := make([]float64, attributes.EndStatType)
			m[attributes.ATKP] = 0.25
			for _, c := range chars {
				c.AddStatMod(character.StatMod{
					Base:         modifier.NewBase("pyro-res", -1),
					AffectedStat: attributes.NoStat,
					Amount: func() []float64 {
						return m
					},
				})
			}
		case attributes.Hydro:
			// TODO: reduce pyro duration not implemented; may affect bennett Q?
			m := make([]float64, attributes.EndStatType)
			m[attributes.HPP] = 0.25
			for _, c := range chars {
				c.AddStatMod(character.StatMod{
					Base:         modifier.NewBase("hydro-res-hpp", -1),
					AffectedStat: attributes.HPP,
					Amount: func() []float64 {
						return m
					},
				})
			}
		case attributes.Cryo:
			m := make([]float64, attributes.EndStatType)
			m[attributes.CR] = .15
			for _, c := range chars {
				c.AddAttackMod(character.AttackMod{
					Base: modifier.NewBase("cryo-res", -1),
					Amount: func(ae *info.AttackEvent, t info.Target) []float64 {
						r, ok := t.(*enemy.Enemy)
						if !ok {
							return nil
						}
						if r.AuraContains(attributes.Cryo) || r.AuraContains(attributes.Frozen) {
							return m
						}
						return nil
					},
				})
			}
		case attributes.Electro:
			last := 0

			recoverParticle := func(_ ...any) {
				if s.F-last < 300 && last != 0 { // every 5 seconds
					return
				}
				s.Player.DistributeParticle(character.Particle{
					Source: "electro-res",
					Num:    1,
					Ele:    attributes.Electro,
				})
				last = s.F
			}

			recoverNoGadget := func(args ...any) {
				if _, ok := args[0].(*enemy.Enemy); ok {
					recoverParticle(args...)
				}
			}
			s.Events.Subscribe(event.OnOverload, recoverNoGadget, "electro-res")
			s.Events.Subscribe(event.OnSuperconduct, recoverNoGadget, "electro-res")
			s.Events.Subscribe(event.OnElectroCharged, recoverNoGadget, "electro-res")
			s.Events.Subscribe(event.OnLunarCharged, recoverNoGadget, "electro-res")
			s.Events.Subscribe(event.OnQuicken, recoverNoGadget, "electro-res")
			s.Events.Subscribe(event.OnAggravate, recoverNoGadget, "electro-res")
			s.Events.Subscribe(event.OnHyperbloom, recoverParticle, "electro-res")
		case attributes.Geo:
			// Increases shield strength by 15%. Additionally, characters protected by a shield will have the
			// following special characteristics:

			//	DMG dealt increased by 15%, dealing DMG to enemies will decrease their Geo RES by 20% for 15s.
			f := func() (float64, bool) { return 0.15, true }
			s.Player.Shields.AddShieldBonusMod("geo-res", -1, f)

			// shred geo res of target
			s.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
				t, ok := args[0].(*enemy.Enemy)
				if !ok {
					return
				}
				atk := args[1].(*info.AttackEvent)
				if s.Player.Shields.CharacterIsShielded(atk.Info.ActorIndex, s.Player.Active()) {
					t.AddResistMod(info.ResistMod{
						Base:  modifier.NewBaseWithHitlag("geo-res", 15*60),
						Ele:   attributes.Geo,
						Value: -0.2,
					})
				}
			}, "geo res")

			m := make([]float64, attributes.EndStatType)
			m[attributes.DmgP] = .15
			for _, c := range chars {
				c.AddAttackMod(character.AttackMod{
					Base: modifier.NewBase("geo-res", -1),
					Amount: func(ae *info.AttackEvent, t info.Target) []float64 {
						if s.Player.Shields.CharacterIsShielded(ae.Info.ActorIndex, s.Player.Active()) {
							return m
						}
						return nil
					},
				})
			}

		case attributes.Anemo:
			s.Player.AddStamPercentMod("anemo-res-stam", -1, func(a action.Action) (float64, bool) {
				return -0.15, false
			})
			// TODO: movement spd increase?
			for _, c := range chars {
				c.AddCooldownMod(character.CooldownMod{
					Base:   modifier.NewBase("anemo-res-cd", -1),
					Amount: func(a action.Action) float64 { return -0.05 },
				})
			}
		case attributes.Dendro:
			m := make([]float64, attributes.EndStatType)
			m[attributes.EM] = 50
			for _, c := range chars {
				c.AddStatMod(character.StatMod{
					Base:         modifier.NewBase("dendro-res-50", -1),
					AffectedStat: attributes.EM,
					Amount: func() []float64 {
						return m
					},
				})
			}

			twoBuff := make([]float64, attributes.EndStatType)
			twoBuff[attributes.EM] = 30
			twoEl := func(args ...any) {
				if _, ok := args[0].(*enemy.Enemy); !ok {
					return
				}
				for _, c := range chars {
					c.AddStatMod(character.StatMod{
						Base:         modifier.NewBaseWithHitlag("dendro-res-30", 6*60),
						AffectedStat: attributes.EM,
						Amount: func() []float64 {
							return twoBuff
						},
					})
				}
			}
			s.Events.Subscribe(event.OnBurning, twoEl, "dendro-res")
			s.Events.Subscribe(event.OnBloom, twoEl, "dendro-res")
			s.Events.Subscribe(event.OnLunarBloom, twoEl, "dendro-res")
			s.Events.Subscribe(event.OnQuicken, twoEl, "dendro-res")

			threeBuff := make([]float64, attributes.EndStatType)
			threeBuff[attributes.EM] = 20
			threeEl := func(_ ...any) {
				for _, c := range chars {
					c.AddStatMod(character.StatMod{
						Base:         modifier.NewBaseWithHitlag("dendro-res-20", 6*60),
						AffectedStat: attributes.EM,
						Amount: func() []float64 {
							return threeBuff
						},
					})
				}
			}
			threeElNoGadget := func(args ...any) {
				if _, ok := args[0].(*enemy.Enemy); ok {
					threeEl(nil)
				}
			}
			s.Events.Subscribe(event.OnAggravate, threeElNoGadget, "dendro-res")
			s.Events.Subscribe(event.OnSpread, threeElNoGadget, "dendro-res")
			s.Events.Subscribe(event.OnHyperbloom, threeEl, "dendro-res")
			s.Events.Subscribe(event.OnBurgeon, threeEl, "dendro-res")
		}
	}
}

func SetupMisc(c *core.Core) {
	c.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		// dmg tag is superconduct, target is enemy
		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return
		}
		atk := args[1].(*info.AttackEvent)
		if atk.Info.AttackTag != attacks.AttackTagSuperconductDamage {
			return
		}
		// add shred
		t.AddResistMod(info.ResistMod{
			Base:  modifier.NewBaseWithHitlag("superconduct-phys-shred", 12*60),
			Ele:   attributes.Physical,
			Value: -0.4,
		})
	}, "superconduct")
}

func setupNightsoulBurst(core *core.Core) {
	chars := core.Player.Chars()
	count := 0
	for _, this := range chars {
		if this.CharZone == info.ZoneNatlan {
			count++
		}
	}
	if count == 0 {
		return
	}

	triggerCD := nightsoulBurstICD[count]
	core.Events.Subscribe(event.OnEnemyDamage, func(args ...any) {
		if core.Status.Duration(nightsoulBurstICDStatus) > 0 {
			return
		}

		t, ok := args[0].(*enemy.Enemy)
		if !ok {
			return
		}
		atk := args[1].(*info.AttackEvent)
		switch atk.Info.Element {
		case attributes.Electro:
		case attributes.Pyro:
		case attributes.Cryo:
		case attributes.Hydro:
		case attributes.Dendro:
		case attributes.Anemo:
		case attributes.Geo:
		default:
			return
		}

		core.Events.Emit(event.OnNightsoulBurst, t, atk)
		core.Status.Add(nightsoulBurstICDStatus, triggerCD)
	}, "nightsoul-burst")
}

func setupAscendantGleam(core *core.Core) {
	chars := core.Player.Chars()
	count := 0
	for _, char := range chars {
		count += char.Moonsign
	}
	if count < 2 {
		return
	}
	buff := 0.0
	src := -1

	var gleamBuffUpdateGen func(*character.CharWrapper, int) func()
	gleamBuffUpdateGen = func(char *character.CharWrapper, s int) func() {
		return func() {
			if s != src {
				return
			}

			if !char.ReactBonusModIsActive("ascendant-gleam") {
				return
			}

			// Ascendant gleam uses nonExtraStats
			switch char.Base.Element {
			case attributes.Electro, attributes.Pyro, attributes.Cryo:
				stats := char.SelectStat(true, attributes.BaseATK, attributes.ATKP, attributes.ATK)
				atk := stats.TotalATK()
				buff = min(atk*0.009/100, 0.36)
			case attributes.Hydro:
				stats := char.SelectStat(true, attributes.BaseHP, attributes.HPP, attributes.HP)
				hp := stats.MaxHP()
				buff = min(hp*0.006/1000, 0.36)
			case attributes.Dendro, attributes.Anemo:
				buff = min(char.NonExtraStat(attributes.EM)*0.0225/100, 0.36)
			case attributes.Geo:
				buff = min(char.TotalDef(true)*0.01/100, 0.36)
			default:
				return
			}
			if core.Flags.LogDebug {
				core.Log.NewEvent("Updating ascendant gleam react bonus", glog.LogDebugEvent, char.Index()).Write("amt", buff)
			}
			core.Tasks.Add(gleamBuffUpdateGen(char, s), 60)
		}
	}

	hook := func(args ...any) {
		src = core.F
		char := core.Player.ActiveChar()
		if char.Moonsign != 0 {
			return
		}

		gleamBuffUpdateGen(char, src)()

		for _, c := range core.Player.Chars() {
			c.AddReactBonusMod(character.ReactBonusMod{
				Base: modifier.NewBase("ascendant-gleam", 20*60),
				Amount: func(ai info.AttackInfo) float64 {
					if !attacks.AttackTagIsLunar(ai.AttackTag) {
						return 0
					}
					if core.Flags.LogDebug {
						core.Log.NewEvent("Adding ascendant gleam react bonus", glog.LogPreDamageMod, char.Index()).Write("amt", buff)
					}
					return buff
				},
			})
		}
	}
	core.Events.Subscribe(event.OnSkill, hook, "ascendant-gleam-on-skill")
	core.Events.Subscribe(event.OnBurst, hook, "ascendant-gleam-on-burst")
}

func (s *Simulation) handleEnergy() {
	// energy once interval=300 amount=1 #once at frame 300
	if s.cfg.EnergySettings.Active && s.cfg.EnergySettings.Once {
		f := s.cfg.EnergySettings.Start
		s.cfg.EnergySettings.Active = false
		s.C.Tasks.Add(func() {
			s.C.Player.DistributeParticle(character.Particle{
				Source: "enemy",
				Num:    float64(s.cfg.EnergySettings.Amount),
				Ele:    attributes.NoElement,
			})
		}, f)
		s.C.Log.NewEventBuildMsg(glog.LogEnergyEvent, -1, "energy queued (once)").
			Write("last", s.cfg.EnergySettings.LastEnergyDrop).
			Write("cfg", s.cfg.EnergySettings).
			Write("amt", s.cfg.EnergySettings.Amount).
			Write("energy_frame", s.C.F+f)
	}
	// energy every interval=300,600 amount=1 #randomly every 300 to 600 frames
	if s.cfg.EnergySettings.Active && s.C.F-s.cfg.EnergySettings.LastEnergyDrop >= s.cfg.EnergySettings.Start {
		f := s.C.Rand.Intn(s.cfg.EnergySettings.End - s.cfg.EnergySettings.Start)
		s.cfg.EnergySettings.LastEnergyDrop = s.C.F + f
		s.C.Tasks.Add(func() {
			s.C.Player.DistributeParticle(character.Particle{
				Source: "drop",
				Num:    float64(s.cfg.EnergySettings.Amount),
				Ele:    attributes.NoElement,
			})
		}, f)
		s.C.Log.NewEventBuildMsg(glog.LogEnergyEvent, -1, "energy queued").
			Write("last", s.cfg.EnergySettings.LastEnergyDrop).
			Write("cfg", s.cfg.EnergySettings).
			Write("amt", s.cfg.EnergySettings.Amount).
			Write("energy_frame", s.C.F+f)
	}
}

func (s *Simulation) handleHurt() {
	// hurt once interval=300 amount=1,300 element=physical #once at frame 300 (or nearest)
	if s.cfg.HurtSettings.Active && s.cfg.HurtSettings.Once {
		f := s.cfg.HurtSettings.Start
		amt := s.cfg.HurtSettings.Min + s.C.Rand.Float64()*(s.cfg.HurtSettings.Max-s.cfg.HurtSettings.Min)
		s.cfg.HurtSettings.Active = false

		s.C.Tasks.Add(func() {
			ai := info.AttackInfo{
				ActorIndex:       s.C.Player.Active(),
				Abil:             "Hurt",
				AttackTag:        attacks.AttackTagNone,
				ICDTag:           attacks.ICDTagNone,
				ICDGroup:         attacks.ICDGroupDefault,
				StrikeType:       attacks.StrikeTypeDefault,
				Durability:       0,
				Element:          s.cfg.HurtSettings.Element,
				FlatDmg:          amt,
				IgnoreDefPercent: 1,
			}
			ap := combat.NewSingleTargetHit(s.C.Combat.Player().Key())
			ap.SkipTargets[info.TargettablePlayer] = false
			ap.SkipTargets[info.TargettableEnemy] = true
			ap.SkipTargets[info.TargettableGadget] = true
			s.C.QueueAttack(ai, ap, -1, 0) // -1 to avoid snapshot
		}, f)

		s.C.Log.NewEventBuildMsg(glog.LogHurtEvent, -1, "hurt queued (once)").
			Write("last", s.cfg.HurtSettings.LastHurt).
			Write("cfg", s.cfg.HurtSettings).
			Write("amt", amt).
			Write("hurt_frame", s.C.F+f)
	}
	// hurt every interval=480,720 amount=1,300 element=physical #randomly 1 to 300 dmg every 480 to 720 frames
	if s.cfg.HurtSettings.Active && s.C.F-s.cfg.HurtSettings.LastHurt >= s.cfg.HurtSettings.Start {
		f := s.C.Rand.Intn(s.cfg.HurtSettings.End - s.cfg.HurtSettings.Start)
		amt := s.cfg.HurtSettings.Min + s.C.Rand.Float64()*(s.cfg.HurtSettings.Max-s.cfg.HurtSettings.Min)
		s.cfg.HurtSettings.LastHurt = s.C.F + f

		s.C.Tasks.Add(func() {
			ai := info.AttackInfo{
				ActorIndex:       s.C.Player.Active(),
				Abil:             "Hurt",
				AttackTag:        attacks.AttackTagNone,
				ICDTag:           attacks.ICDTagNone,
				ICDGroup:         attacks.ICDGroupDefault,
				StrikeType:       attacks.StrikeTypeDefault,
				Durability:       0,
				Element:          s.cfg.HurtSettings.Element,
				FlatDmg:          amt,
				IgnoreDefPercent: 1,
			}
			ap := combat.NewSingleTargetHit(s.C.Combat.Player().Key())
			ap.SkipTargets[info.TargettablePlayer] = false
			ap.SkipTargets[info.TargettableEnemy] = true
			ap.SkipTargets[info.TargettableGadget] = true
			s.C.QueueAttack(ai, ap, -1, 0) // -1 to avoid snapshot
		}, f)

		s.C.Log.NewEventBuildMsg(glog.LogHurtEvent, -1, "hurt queued").
			Write("last", s.cfg.HurtSettings.LastHurt).
			Write("cfg", s.cfg.HurtSettings).
			Write("amt", amt).
			Write("hurt_frame", s.C.F+f)
	}
}
