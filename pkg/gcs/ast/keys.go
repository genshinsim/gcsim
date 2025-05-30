package ast

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

var key = map[string]TokenType{
	".":           ItemDot,
	"let":         KeywordLet,
	"while":       KeywordWhile,
	"if":          KeywordIf,
	"else":        KeywordElse,
	"fn":          KeywordFn,
	"switch":      KeywordSwitch,
	"case":        KeywordCase,
	"default":     KeywordDefault,
	"break":       KeywordBreak,
	"continue":    KeywordContinue,
	"fallthrough": KeywordFallthrough,
	"return":      KeywordReturn,
	"for":         KeywordFor,
	// genshin specific keywords
	"options":             KeywordOptions,
	"add":                 KeywordAdd,
	"char":                KeywordChar,
	"stats":               KeywordStats,
	"weapon":              KeywordWeapon,
	"set":                 KeywordSet,
	"lvl":                 KeywordLvl,
	"refine":              KeywordRefine,
	"cons":                KeywordCons,
	"talent":              KeywordTalent,
	"count":               KeywordCount,
	"params":              KeywordParams,
	"label":               KeywordLabel,
	"until":               KeywordUntil,
	"active":              KeywordActive,
	"target":              KeywordTarget,
	"particle_threshold":  KeywordParticleThreshold,
	"particle_drop_count": KeywordParticleDropCount,
	"particle_element":    KeywordParticleElement,
	"resist":              KeywordResist,
	"energy":              KeywordEnergy,
	"hurt":                KeywordHurt,
	// commands
	// team keywords
	// flags
	// ??
	// energy/hurt event related
	// target related
}

var StatKeys = map[string]attributes.Stat{
	"def%":     attributes.DEFP,
	"def":      attributes.DEF,
	"hp":       attributes.HP,
	"hp%":      attributes.HPP,
	"atk":      attributes.ATK,
	"atk%":     attributes.ATKP,
	"er":       attributes.ER,
	"em":       attributes.EM,
	"cr":       attributes.CR,
	"cd":       attributes.CD,
	"heal":     attributes.Heal,
	"pyro%":    attributes.PyroP,
	"hydro%":   attributes.HydroP,
	"cryo%":    attributes.CryoP,
	"electro%": attributes.ElectroP,
	"anemo%":   attributes.AnemoP,
	"geo%":     attributes.GeoP,
	"phys%":    attributes.PhyP,
	// "ele%":     attributes.ElementalP,
	"dendro%": attributes.DendroP,
	"atkspd%": attributes.AtkSpd,
	"dmg%":    attributes.DmgP,
}

var EleKeys = map[string]attributes.Element{
	"electro":  attributes.Electro,
	"pyro":     attributes.Pyro,
	"cryo":     attributes.Cryo,
	"hydro":    attributes.Hydro,
	"frozen":   attributes.Frozen,
	"anemo":    attributes.Anemo,
	"dendro":   attributes.Dendro,
	"geo":      attributes.Geo,
	"physical": attributes.Physical,
	"none":     attributes.NoElement,
}

var actionKeys = map[string]action.Action{
	"skill":       action.ActionSkill,
	"burst":       action.ActionBurst,
	"attack":      action.ActionAttack,
	"charge":      action.ActionCharge,
	"high_plunge": action.ActionHighPlunge,
	"low_plunge":  action.ActionLowPlunge,
	"aim":         action.ActionAim,
	"dash":        action.ActionDash,
	"jump":        action.ActionJump,
	"walk":        action.ActionWalk,
	"swap":        action.ActionSwap,
}
