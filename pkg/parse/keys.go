package parse

import "github.com/genshinsim/gcsim/pkg/core"

var key = map[string]ItemType{
	".": itemDot,
	//commands
	"chain":       itemChain,
	"wait_for":    itemWaitFor,
	"wait":        itemWait,
	"restart":     itemRestart,
	"reset_limit": itemResetLimit,
	"hurt":        itemHurt,
	"target":      itemTarget,
	"energy":      itemEnergy,
	"active":      itemActive,
	"options":     itemOptions,
	//team keywords
	"add":      itemAdd,
	"char":     itemChar,
	"stats":    itemStats,
	"weapon":   itemWeapon,
	"set":      itemSet,
	"lvl":      itemLvl,
	"refine":   itemRefine,
	"cons":     itemCons,
	"talent":   itemTalent,
	"start_hp": itemStartHP,
	"count":    itemCount,
	"params":   itemParams,
	"label":    itemLabel,
	"until":    itemUntil,
	//flags
	"if":         itemIf,
	"swap_to":    itemSwap,
	"swap_lock":  itemSwapLock,
	"is_onfield": itemOnField,
	"needs":      itemNeeds,
	"limit":      itemLimit,
	"timeout":    itemTimeout,
	"try":        itemTry,
	"drop":       itemDrop,
	// ??
	"value":  itemValue,
	"max":    itemMax,
	"filler": itemFiller,
	//energy/hurt event related
	"interval": itemInterval,
	"amount":   itemAmount,
	"once":     itemOnce,
	"every":    itemEvery,
	"ele":      itemEle,
	// target related
	"resist": itemResist,
}

var queueModeKeys = map[string]core.SimulationQueueMode{
	"apl": core.ActionPriorityList,
	"sl":  core.SequentialList,
}

var statKeys = map[string]core.StatType{
	"def%":     core.DEFP,
	"def":      core.DEF,
	"hp":       core.HP,
	"hp%":      core.HPP,
	"atk":      core.ATK,
	"atk%":     core.ATKP,
	"er":       core.ER,
	"em":       core.EM,
	"cr":       core.CR,
	"cd":       core.CD,
	"heal":     core.Heal,
	"pyro%":    core.PyroP,
	"hydro%":   core.HydroP,
	"cryo%":    core.CryoP,
	"electro%": core.ElectroP,
	"anemo%":   core.AnemoP,
	"geo%":     core.GeoP,
	"phys%":    core.PhyP,
	// "ele%":     core.ElementalP,
	"dendro%": core.DendroP,
	"atkspd%": core.AtkSpd,
	"dmg%":    core.DmgP,
}

var eleKeys = map[string]core.EleType{
	"electro":  core.Electro,
	"pyro":     core.Pyro,
	"cryo":     core.Cryo,
	"hydro":    core.Hydro,
	"frozen":   core.Frozen,
	"anemo":    core.Anemo,
	"dendro":   core.Dendro,
	"geo":      core.Geo,
	"physical": core.Physical,
}

var actionKeys = map[string]core.ActionType{
	"skill":       core.ActionSkill,
	"burst":       core.ActionBurst,
	"attack":      core.ActionAttack,
	"charge":      core.ActionCharge,
	"high_plunge": core.ActionHighPlunge,
	"low_plunge":  core.ActionLowPlunge,
	"aim":         core.ActionAim,
	"dash":        core.ActionDash,
	"jump":        core.ActionJump,
	"walk":        core.ActionWalk,
	"swap":        core.ActionSwap,
}
