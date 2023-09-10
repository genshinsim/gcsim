package ast

import (
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
)

var key = map[string]TokenType{
	".":           itemDot,
	"let":         keywordLet,
	"while":       keywordWhile,
	"if":          keywordIf,
	"else":        keywordElse,
	"fn":          keywordFn,
	"switch":      keywordSwitch,
	"case":        keywordCase,
	"default":     keywordDefault,
	"break":       keywordBreak,
	"continue":    keywordContinue,
	"fallthrough": keywordFallthrough,
	"return":      keywordReturn,
	"for":         keywordFor,
	//genshin specific keywords
	"options":             keywordOptions,
	"add":                 keywordAdd,
	"char":                keywordChar,
	"stats":               keywordStats,
	"weapon":              keywordWeapon,
	"set":                 keywordSet,
	"lvl":                 keywordLvl,
	"refine":              keywordRefine,
	"cons":                keywordCons,
	"talent":              keywordTalent,
	"count":               keywordCount,
	"params":              keywordParams,
	"label":               keywordLabel,
	"until":               keywordUntil,
	"active":              keywordActive,
	"target":              keywordTarget,
	"particle_threshold":  keywordParticleThreshold,
	"particle_drop_count": keywordParticleDropCount,
	"resist":              keywordResist,
	"energy":              keywordEnergy,
	"hurt":                keywordHurt,
	//commands
	//team keywords
	//flags
	// ??
	//energy/hurt event related
	// target related
}

var statKeys = map[string]attributes.Stat{
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

var eleKeys = map[string]attributes.Element{
	"electro":  attributes.Electro,
	"pyro":     attributes.Pyro,
	"cryo":     attributes.Cryo,
	"hydro":    attributes.Hydro,
	"frozen":   attributes.Frozen,
	"anemo":    attributes.Anemo,
	"dendro":   attributes.Dendro,
	"geo":      attributes.Geo,
	"physical": attributes.Physical,
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
