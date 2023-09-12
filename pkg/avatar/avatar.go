package avatar

import (
	"log"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
	"github.com/genshinsim/gcsim/pkg/core/event"
	"github.com/genshinsim/gcsim/pkg/core/geometry"
	"github.com/genshinsim/gcsim/pkg/core/glog"
	"github.com/genshinsim/gcsim/pkg/core/player"
	"github.com/genshinsim/gcsim/pkg/core/reactions"
	"github.com/genshinsim/gcsim/pkg/core/targets"
	"github.com/genshinsim/gcsim/pkg/reactable"
	"github.com/genshinsim/gcsim/pkg/target"
)

type Player struct {
	*target.Target
	*reactable.Reactable
}

func New(core *core.Core, pos geometry.Point, r float64) *Player {
	p := &Player{}
	p.Target = target.New(core, pos, r)
	p.Reactable = &reactable.Reactable{}
	p.Reactable.Init(p, core)
	return p
}

func (p *Player) Type() targets.TargettableType { return targets.TargettablePlayer }

func (p *Player) HandleAttack(atk *combat.AttackEvent) float64 {
	p.Core.Combat.Events.Emit(event.OnPlayerHit, p, atk)
	var amp string
	var cata string
	var dmg float64
	var crit bool
	evt := p.Core.Combat.Log.NewEvent(atk.Info.Abil, glog.LogDamageEvent, atk.Info.ActorIndex)

	// TODO: Implement reactions on player

	dmg, crit = p.calc(atk, evt)

	active := p.Core.Player.Active()
	dmgLeft := p.Core.Player.Shields.OnDamage(active, dmg, atk.Info.Element)
	if dmgLeft > 0 {
		p.Core.Player.Drain(player.DrainInfo{
			ActorIndex: active,
			Abil:       atk.Info.Abil,
			Amount:     dmgLeft,
			External:   true,
		})
	}
	evt.Write("target", p.Key()).
		Write("attack-tag", atk.Info.AttackTag).
		Write("ele", atk.Info.Element.String()).
		Write("damage_pre_shield", &dmg).
		Write("damage", &dmgLeft).
		Write("crit", &crit).
		Write("amp", &amp).
		Write("cata", &cata).
		Write("abil", atk.Info.Abil).
		Write("source_frame", atk.SourceFrame)
	evt.WriteBuildMsg(atk.Snapshot.Logs...)

	if !atk.Info.SourceIsSim {
		if atk.Info.ActorIndex < 0 {
			log.Println(atk)
		}
		preDmgModDebug := p.Core.Combat.Team.CombatByIndex(atk.Info.ActorIndex).
			ApplyAttackMods(atk, p)
		evt.Write("pre_damage_mods", preDmgModDebug)
	}

	// returns 0 so that the damage done to player doesn't count
	// towards the sim's TotalDamage and DPS statistic
	return 0
}
func (p *Player) calc(atk *combat.AttackEvent, evt glog.Event) (float64, bool) {
	var isCrit bool

	st := attributes.EleToDmgP(atk.Info.Element)
	// if st < 0 {
	// 	log.Println(atk)
	// }
	elePer := 0.0
	if st > -1 {
		elePer = atk.Snapshot.Stats[st]
		// Generally not needed except for sim issues
		// p.Core.Log.NewEvent("ele lookup ok",
		// 	glog.LogCalc, atk.Info.ActorIndex,
		// 	"attack_tag", atk.Info.AttackTag,
		// 	"ele", atk.Info.Element,
		// 	"st", st,
		// 	"percent", atk.Snapshot.Stats[st],
		// 	"abil", atk.Info.Abil,
		// 	"stats", atk.Snapshot.Stats,
		// 	"target", p.TargetIndex,
		// )
	}
	dmgBonus := elePer + atk.Snapshot.Stats[attributes.DmgP]

	// calculate using attack or def
	var a float64
	totalhp := atk.Snapshot.BaseHP*(1+atk.Snapshot.Stats[attributes.HPP]) + atk.Snapshot.Stats[attributes.HP]
	if atk.Info.UseDef {
		a = atk.Snapshot.BaseDef*(1+atk.Snapshot.Stats[attributes.DEFP]) + atk.Snapshot.Stats[attributes.DEF]
	} else {
		a = atk.Snapshot.BaseAtk*(1+atk.Snapshot.Stats[attributes.ATKP]) + atk.Snapshot.Stats[attributes.ATK]
	}

	base := atk.Info.Mult*a + atk.Info.FlatDmg
	damage := base * (1 + dmgBonus)

	// make sure 0 <= cr <= 1
	if atk.Snapshot.Stats[attributes.CR] < 0 {
		atk.Snapshot.Stats[attributes.CR] = 0
	}
	if atk.Snapshot.Stats[attributes.CR] > 1 {
		atk.Snapshot.Stats[attributes.CR] = 1
	}

	char := p.Core.Player.ActiveChar()
	// TODO: Players don't have resistances right now
	res := 0.0

	def := char.Base.Def*(1+char.Stat(attributes.DEFP)) + char.Stat(attributes.DEF)

	def *= (1 - atk.Info.IgnoreDefPercent)
	defmod := 1 - def/(def+float64(5*atk.Snapshot.CharLvl)+500)

	// apply def mod
	damage *= defmod
	// apply resist mod

	resmod := 1 - res/2
	if res >= 0 && res < 0.75 {
		resmod = 1 - res
	} else if res > 0.75 {
		resmod = 1 / (4*res + 1)
	}
	damage *= resmod

	precritdmg := damage

	// check if crit
	if atk.Info.HitWeakPoint || p.Core.Rand.Float64() <= atk.Snapshot.Stats[attributes.CR] {
		damage *= (1 + atk.Snapshot.Stats[attributes.CD])
		isCrit = true
	}

	preampdmg := damage

	// calculate em bonus
	em := atk.Snapshot.Stats[attributes.EM]
	emBonus := (2.78 * em) / (1400 + em)
	var reactBonus float64
	// check melt/vape
	if atk.Info.Amped {
		reactBonus = p.Core.Player.ByIndex(atk.Info.ActorIndex).ReactBonus(atk.Info)
		// t.Core.Log.Debugw("debug", "frame", t.Core.F, core.LogPreDamageMod, "char", t.Index, "char_react", char.CharIndex(), "reactbonus", char.ReactBonus(atk.Info), "damage_pre", damage)
		damage *= (atk.Info.AmpMult * (1 + emBonus + reactBonus))
	}

	// reduce damage by damage group
	x := 1.0
	if !atk.Info.SourceIsSim {
		x = p.GroupTagDamageMult(atk.Info.ICDTag, atk.Info.ICDGroup, atk.Info.ActorIndex)
		damage *= x
	}

	if p.Core.Flags.LogDebug {
		p.Core.Log.NewEvent(
			atk.Info.Abil,
			glog.LogCalc,
			atk.Info.ActorIndex,
		).
			Write("src_frame", atk.SourceFrame).
			Write("damage_grp_mult", x).
			Write("damage", damage).
			Write("abil", atk.Info.Abil).
			Write("talent", atk.Info.Mult).
			Write("base_atk", atk.Snapshot.BaseAtk).
			Write("flat_atk", atk.Snapshot.Stats[attributes.ATK]).
			Write("atk_per", atk.Snapshot.Stats[attributes.ATKP]).
			Write("use_def", atk.Info.UseDef).
			Write("base_def", atk.Snapshot.BaseDef).
			Write("flat_def", atk.Snapshot.Stats[attributes.DEF]).
			Write("def_per", atk.Snapshot.Stats[attributes.DEFP]).
			Write("base_hp", atk.Snapshot.BaseHP).
			Write("flat_hp", atk.Snapshot.Stats[attributes.HP]).
			Write("hp_per", atk.Snapshot.Stats[attributes.HPP]).
			Write("total_hp", totalhp).
			Write("catalyzed", atk.Info.Catalyzed).
			Write("flat_dmg", atk.Info.FlatDmg).
			Write("total_atk_def", a).
			Write("base_dmg", base).
			Write("ele", st).
			Write("ele_per", elePer).
			Write("bonus_dmg", dmgBonus).
			Write("ignore_def", atk.Info.IgnoreDefPercent).
			Write("def_adj", 0). // Players don't have DefMods applied
			Write("target_lvl", char.Base.Level).
			Write("char_lvl", atk.Snapshot.CharLvl).
			Write("def_mod", defmod).
			Write("res", res).
			Write("res_mod", resmod).
			Write("cr", atk.Snapshot.Stats[attributes.CR]).
			Write("cd", atk.Snapshot.Stats[attributes.CD]).
			Write("pre_crit_dmg", precritdmg).
			Write("dmg_if_crit", precritdmg*(1+atk.Snapshot.Stats[attributes.CD])).
			Write("avg_crit_dmg", (1-atk.Snapshot.Stats[attributes.CR])*precritdmg+atk.Snapshot.Stats[attributes.CR]*precritdmg*(1+atk.Snapshot.Stats[attributes.CD])).
			Write("is_crit", isCrit).
			Write("pre_amp_dmg", preampdmg).
			Write("reaction_type", atk.Info.AmpType).
			Write("melt_vape", atk.Info.Amped).
			Write("react_mult", atk.Info.AmpMult).
			Write("em", em).
			Write("em_bonus", emBonus).
			Write("react_bonus", reactBonus).
			Write("amp_mult_total", (atk.Info.AmpMult*(1+emBonus+reactBonus))).
			Write("pre_crit_dmg_react", precritdmg*(atk.Info.AmpMult*(1+emBonus+reactBonus))).
			Write("dmg_if_crit_react", precritdmg*(1+atk.Snapshot.Stats[attributes.CD])*(atk.Info.AmpMult*(1+emBonus+reactBonus))).
			Write("avg_crit_dmg_react", ((1-atk.Snapshot.Stats[attributes.CR])*precritdmg+atk.Snapshot.Stats[attributes.CR]*precritdmg*(1+atk.Snapshot.Stats[attributes.CD]))*(atk.Info.AmpMult*(1+emBonus+reactBonus))).
			Write("target", p.Key())
	}

	return damage, isCrit
}

func (p *Player) ApplySelfInfusion(ele attributes.Element, dur reactions.Durability, f int) {
	p.Core.Log.NewEventBuildMsg(glog.LogPlayerEvent, -1, "self infusion applied: "+ele.String()).
		Write("durability", dur).
		Write("duration", f)
	// we're assuming self infusion isn't subject to 0.8x multiplier
	// also no real sanity check
	if ele == attributes.Frozen {
		return
	}
	var mod reactable.Modifier
	switch ele {
	case attributes.Electro:
		mod = reactable.Electro
	case attributes.Hydro:
		mod = reactable.Hydro
	case attributes.Pyro:
		mod = reactable.Pyro
	case attributes.Cryo:
		mod = reactable.Cryo
	case attributes.Dendro:
		mod = reactable.Dendro
	}

	// we're assuming refill maintains the same decay rate?
	if p.Durability[mod] > reactable.ZeroDur {
		// make sure we're not adding more than incoming
		if p.Durability[mod] < dur {
			p.Durability[mod] = dur
		}
		return
	}
	// otherwise calculate decay based on specified f (in frames)
	p.Durability[mod] = dur
	p.DecayRate[mod] = dur / reactions.Durability(f)
}

func (p *Player) ReactWithSelf(atk *combat.AttackEvent) {
	// check if have an element
	if p.AuraCount() == 0 {
		return
	}
	// otherwise react
	existing := p.Reactable.ActiveAuraString()
	applied := atk.Info.Durability
	p.React(atk)
	p.Core.Log.NewEvent("self reaction occured", glog.LogElementEvent, atk.Info.ActorIndex).
		Write("attack_tag", atk.Info.AttackTag).
		Write("applied_ele", atk.Info.Element.String()).
		Write("dur", applied).
		Write("abil", atk.Info.Abil).
		Write("target", 0).
		Write("existing", existing).
		Write("after", p.Reactable.ActiveAuraString())
}
