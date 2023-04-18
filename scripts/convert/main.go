package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/tools/go/ast/astutil"
)

func main() {
	// read every file in directory

	//check args if we want to print ast out
	dump := false
	if len(os.Args) > 1 {
		dump = true
	}

	files, err := ioutil.ReadDir("./")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		err = fix(file.Name(), dump)
		if err != nil {
			panic(err)
		}
	}
}

func fix(path string, dump bool) error {
	//do nothing
	if filepath.Ext(path) != ".go" {
		return nil
	}
	fs := token.NewFileSet()
	f, err := parser.ParseFile(fs, path, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	if dump {
		spew.Dump(f)
	}

	//fix any core package names
	astutil.Apply(f, func(cr *astutil.Cursor) bool {
		found, next := fixCorePkgName(cr.Node())
		if !found {
			return true
		}
		cr.Replace(next)
		return false
	}, nil)

	astutil.Apply(f, func(cr *astutil.Cursor) bool {
		// found, next := findAndReplacePreDamageBlock(cr.Node())
		// if !found {
		// 	return true
		// }
		// cr.Replace(next)
		// return false

		rep := func(f func(ast.Node) (bool, ast.Node)) bool {
			found, next := f(cr.Node())
			if !found {
				return true
			}
			cr.Replace(next)
			return false
		}

		var notFound bool

		//if found then we have to return false so we don't keep drilling down on a replaced node
		notFound = rep(findAndReplaceStatBlock)
		if !notFound {
			fmt.Println("add mod block found")
			return false
		}
		notFound = rep(findAndReplacePreDamageBlock)
		if !notFound {
			fmt.Println("add predmgmod block found")
			return false
		}
		notFound = rep(findAndReplaceActionSignature)
		if !notFound {
			fmt.Println("abil block found")
			return false
		}

		return true
	}, nil)
	// Print result
	out, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer out.Close()
	printer.Fprint(out, fs, f)

	// printer.Fprint(os.Stdout, fs, f)
	return nil
}

func findAndReplacePreDamageBlock(n ast.Node) (bool, ast.Node) {
	if expr, ok := n.(*ast.ExprStmt); ok {
		block, ok := expr.X.(*ast.CallExpr)
		if !ok {
			return false, nil
		}

		//FUN should be a SelectorExpr
		fun, ok := block.Fun.(*ast.SelectorExpr)
		if !ok {
			return false, nil
		}

		//Sel should be AddPreDamageMod
		if fun.Sel.Name != "AddPreDamageMod" {
			return false, nil
		}

		fmt.Println("found pre damage block")

		//work through the args and find amount, expiry, and key
		//args should be len 1
		if len(block.Args) != 1 {
			fmt.Println("unexpected args length > 1")
			return false, nil
		}

		//check to make sure it's a CompositeLit
		lit, ok := block.Args[0].(*ast.CompositeLit)
		if !ok {
			fmt.Println("unexpected arg type, not a composite lit")
			return false, nil
		}

		//loop through Elts to find amount, expiry and key
		var amount, expiry, key ast.Expr
		for _, v := range lit.Elts {
			t, ok := v.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			switch t.Key.(*ast.Ident).Name {
			case "Amount":
				amount = t.Value
			case "Expiry":
				expiry = t.Value
			case "Key":
				key = t.Value
			}
		}

		caller, ok := fun.X.(*ast.Ident)
		if !ok {
			fmt.Println("unexpected fun.X type, not an ident")
		}

		next := &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun:  ast.NewIdent(fmt.Sprintf("%v.AddAttackMod", caller.Name)),
				Args: []ast.Expr{key, expiry, amount},
			},
		}

		return true, next
	}
	return false, nil
}

func findAndReplaceStatBlock(n ast.Node) (bool, ast.Node) {
	if expr, ok := n.(*ast.ExprStmt); ok {
		block, ok := expr.X.(*ast.CallExpr)
		if !ok {
			return false, nil
		}

		//FUN should be a SelectorExpr
		fun, ok := block.Fun.(*ast.SelectorExpr)
		if !ok {
			return false, nil
		}

		//Sel should be AddPreDamageMod
		if fun.Sel.Name != "AddMod" {
			return false, nil
		}

		fmt.Println("found stat block")

		//work through the args and find amount, expiry, and key
		//args should be len 1
		if len(block.Args) != 1 {
			fmt.Println("unexpected args length > 1")
			return false, nil
		}

		//check to make sure it's a CompositeLit
		lit, ok := block.Args[0].(*ast.CompositeLit)
		if !ok {
			fmt.Println("unexpected arg type, not a composite lit")
			return false, nil
		}

		//loop through Elts to find amount, expiry and key
		var amount, expiry, key ast.Expr
		var stat ast.Expr = &ast.SelectorExpr{
			X:   ast.NewIdent("attributes"),
			Sel: ast.NewIdent("NoStat"),
		}
		for _, v := range lit.Elts {
			t, ok := v.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			switch t.Key.(*ast.Ident).Name {
			case "Amount":
				amount = t.Value
			case "Expiry":
				expiry = t.Value
			case "Key":
				key = t.Value
			case "AffectedStat":
				stat = t.Value
			}
		}

		caller, ok := fun.X.(*ast.Ident)
		if !ok {
			fmt.Println("unexpected fun.X type, not an ident")
		}

		next := &ast.ExprStmt{
			X: &ast.CallExpr{
				Fun:  ast.NewIdent(fmt.Sprintf("%v.AddStatMod", caller.Name)),
				Args: []ast.Expr{key, expiry, stat, amount},
			},
		}

		return true, next
	}
	return false, nil
}

var abilFunctions = map[string]bool{
	"Attack":           true,
	"Aimed":            true,
	"ChargeAttack":     true,
	"HighPlungeAttack": true,
	"LowPlungeAttack":  true,
	"Skill":            true,
	"Burst":            true,
	"Dash":             true,
	"Walk":             true,
	"Jump":             true,
}

func findAndReplaceActionSignature(n ast.Node) (bool, ast.Node) {
	//look for function declaration that matches and return (int, int)
	fd, ok := n.(*ast.FuncDecl)
	if !ok {
		return false, nil
	}
	//check name is in map
	if _, ok := abilFunctions[fd.Name.Name]; !ok {
		return false, nil
	}
	//check result is len 2 and both int
	if len(fd.Type.Results.List) != 2 {
		return false, nil
	}
	bad := false
	for _, v := range fd.Type.Results.List {
		//should be an ident with name int
		id, ok := v.Type.(*ast.Ident)
		if !ok {
			bad = true
			break
		}
		if id.Name != "int" {
			bad = true
			break
		}
	}
	if bad {
		return false, nil
	}

	//now we want to over write the result list with action.ActionInfo
	field := &ast.Field{
		Type: &ast.SelectorExpr{
			X:   ast.NewIdent("action"),
			Sel: ast.NewIdent("ActionInfo"),
		},
	}
	//replace the list
	fd.Type.Results.List = []*ast.Field{field}
	return true, fd
}

func fixCorePkgName(n ast.Node) (bool, *ast.SelectorExpr) {
	//check if selector
	sel, ok := n.(*ast.SelectorExpr)
	if !ok {
		return false, nil
	}

	//check if starts with core
	x, ok := sel.X.(*ast.Ident)
	if !ok {
		return false, nil
	}

	if x.Name != "core" {
		return false, nil
	}

	//check if Ident matches one of the ones we're replacing
	s, ok := pkgNameReplace[fmt.Sprintf("%s.%s", x.Name, sel.Sel.Name)]
	if !ok {
		return false, nil
	}

	return true, &ast.SelectorExpr{
		X:   ast.NewIdent(s[0]),
		Sel: ast.NewIdent(s[1]),
	}
}

var pkgNameReplace = map[string][2]string{
	//stats
	"core.EleType":     {"attributes", "Element"},
	"core.Electro":     {"attributes", "Electro"},
	"core.Pyro":        {"attributes", "Pyro"},
	"core.Cryo":        {"attributes", "Cryo"},
	"core.Hydro":       {"attributes", "Hydro"},
	"core.Frozen":      {"attributes", "Frozen"},
	"core.Anemo":       {"attributes", "Anemo"},
	"core.Dendro":      {"attributes", "Dendro"},
	"core.Geo":         {"attributes", "Geo"},
	"core.NoElement":   {"attributes", "NoElement"},
	"core.Physical":    {"attributes", "Physical"},
	"core.DEFP":        {"attributes", "DEFP"},
	"core.DEF":         {"attributes", "DEF"},
	"core.HP":          {"attributes", "HP"},
	"core.HPP":         {"attributes", "HPP"},
	"core.ATK":         {"attributes", "ATK"},
	"core.ATKP":        {"attributes", "ATKP"},
	"core.ER":          {"attributes", "ER"},
	"core.EM":          {"attributes", "EM"},
	"core.CR":          {"attributes", "CR"},
	"core.CD":          {"attributes", "CD"},
	"core.Heal":        {"attributes", "Heal"},
	"core.PyroP":       {"attributes", "PyroP"},
	"core.HydroP":      {"attributes", "HydroP"},
	"core.CryoP":       {"attributes", "CryoP"},
	"core.ElectroP":    {"attributes", "ElectroP"},
	"core.AnemoP":      {"attributes", "AnemoP"},
	"core.GeoP":        {"attributes", "GeoP"},
	"core.PhyP":        {"attributes", "PhyP"},
	"core.DendroP":     {"attributes", "DendroP"},
	"core.AtkSpd":      {"attributes", "AtkSpd"},
	"core.DmgP":        {"attributes", "DmgP"},
	"core.EndStatType": {"attributes", "EndStatType"},
	//events
	"core.OnAttackWillLand":         {"event", "OnAttackWillLand"},
	"core.OnDamage":                 {"event", "OnDamage"},
	"core.OnAuraDurabilityAdded":    {"event", "OnAuraDurabilityAdded"},
	"core.OnAuraDurabilityDepleted": {"event", "OnAuraDurabilityDepleted"},
	"core.ReactionEventStartDelim":  {"event", "ReactionEventStartDelim"},
	"core.OnOverload":               {"event", "OnOverload"},
	"core.OnSuperconduct":           {"event", "OnSuperconduct"},
	"core.OnMelt":                   {"event", "OnMelt"},
	"core.OnVaporize":               {"event", "OnVaporize"},
	"core.OnFrozen":                 {"event", "OnFrozen"},
	"core.OnElectroCharged":         {"event", "OnElectroCharged"},
	"core.OnSwirlHydro":             {"event", "OnSwirlHydro"},
	"core.OnSwirlCryo":              {"event", "OnSwirlCryo"},
	"core.OnSwirlElectro":           {"event", "OnSwirlElectro"},
	"core.OnSwirlPyro":              {"event", "OnSwirlPyro"},
	"core.OnCrystallizeHydro":       {"event", "OnCrystallizeHydro"},
	"core.OnCrystallizeCryo":        {"event", "OnCrystallizeCryo"},
	"core.OnCrystallizeElectro":     {"event", "OnCrystallizeElectro"},
	"core.OnCrystallizePyro":        {"event", "OnCrystallizePyro"},
	"core.ReactionEventEndDelim":    {"event", "ReactionEventEndDelim"},
	"core.OnStamUse":                {"event", "OnStamUse"},
	"core.OnShielded":               {"event", "OnShielded"},
	"core.OnCharacterSwap":          {"event", "OnCharacterSwap"},
	"core.OnParticleReceived":       {"event", "OnParticleReceived"},
	"core.OnEnergyChange":           {"event", "OnEnergyChange"},
	"core.OnTargetDied":             {"event", "OnTargetDied"},
	"core.OnCharacterHurt":          {"event", "OnCharacterHurt"},
	"core.OnHeal":                   {"event", "OnHeal"},
	"core.OnActionExec":             {"event", "OnActionExec"},
	"core.PreSkill":                 {"event", "PreSkill"},
	"core.PostSkill":                {"event", "PostSkill"},
	"core.PreBurst":                 {"event", "PreBurst"},
	"core.PostBurst":                {"event", "PostBurst"},
	"core.PreAttack":                {"event", "PreAttack"},
	"core.PostAttack":               {"event", "PostAttack"},
	"core.PreChargeAttack":          {"event", "PreChargeAttack"},
	"core.PostChargeAttack":         {"event", "PostChargeAttack"},
	"core.PrePlunge":                {"event", "PrePlunge"},
	"core.PostPlunge":               {"event", "PostPlunge"},
	"core.PreAimShoot":              {"event", "PreAimShoot"},
	"core.PostAimShoot":             {"event", "PostAimShoot"},
	"core.PreDash":                  {"event", "PreDash"},
	"core.PostDash":                 {"event", "PostDash"},
	"core.OnInitialize":             {"event", "OnInitialize"},
	"core.OnStateChange":            {"event", "OnStateChange"},
	"core.OnTargetAdded":            {"event", "OnTargetAdded"},
	"core.EndEventTypes":            {"event", "EndEventTypes"},
	//cobmat related
	"core.AttackEvent":                 {"combat", "AttackEvent"},
	"core.Target":                      {"combat", "Target"},
	"core.AttackInfo":                  {"combat", "AttackInfo"},
	"core.AttackCB":                    {"combat", "AttackCB"},
	"core.AttackTag":                   {"combat", "AttackTag"},
	"core.NewDefSingleTarget":          {"combat", "NewDefSingleTarget"},
	"core.NewDefCircHit":               {"combat", "NewDefCircHit"},
	"core.AttackTagNone":               {"combat", "AttackTagNone"},
	"core.AttackTagNormal":             {"combat", "AttackTagNormal"},
	"core.AttackTagExtra":              {"combat", "AttackTagExtra"},
	"core.AttackTagPlunge":             {"combat", "AttackTagPlunge"},
	"core.AttackTagElementalArt":       {"combat", "AttackTagElementalArt"},
	"core.AttackTagElementalArtHold":   {"combat", "AttackTagElementalArtHold"},
	"core.AttackTagElementalBurst":     {"combat", "AttackTagElementalBurst"},
	"core.AttackTagWeaponSkill":        {"combat", "AttackTagWeaponSkill"},
	"core.AttackTagMonaBubbleBreak":    {"combat", "AttackTagMonaBubbleBreak"},
	"core.AttackTagNoneStat":           {"combat", "AttackTagNoneStat"},
	"core.ReactionAttackDelim":         {"combat", "ReactionAttackDelim"},
	"core.AttackTagOverloadDamage":     {"combat", "AttackTagOverloadDamage"},
	"core.AttackTagSuperconductDamage": {"combat", "AttackTagSuperconductDamage"},
	"core.AttackTagECDamage":           {"combat", "AttackTagECDamage"},
	"core.AttackTagShatter":            {"combat", "AttackTagShatter"},
	"core.AttackTagSwirlPyro":          {"combat", "AttackTagSwirlPyro"},
	"core.AttackTagSwirlHydro":         {"combat", "AttackTagSwirlHydro"},
	"core.AttackTagSwirlCryo":          {"combat", "AttackTagSwirlCryo"},
	"core.AttackTagSwirlElectro":       {"combat", "AttackTagSwirlElectro"},
	"core.AttackTagLength":             {"combat", "AttackTagLength"},
	"core.ICDTagNone":                  {"combat", "ICDTagNone"},
	"core.ICDTagNormalAttack":          {"combat", "ICDTagNormalAttack"},
	"core.ICDTagExtraAttack":           {"combat", "ICDTagExtraAttack"},
	"core.ICDTagElementalArt":          {"combat", "ICDTagElementalArt"},
	"core.ICDTagElementalBurst":        {"combat", "ICDTagElementalBurst"},
	"core.ICDTagDash":                  {"combat", "ICDTagDash"},
	"core.ICDTagLisaElectro":           {"combat", "ICDTagLisaElectro"},
	"core.ICDTagYanfeiFire":            {"combat", "ICDTagYanfeiFire"},
	"core.ICDTagVentiBurstAnemo":       {"combat", "ICDTagVentiBurstAnemo"},
	"core.ICDTagVentiBurstPyro":        {"combat", "ICDTagVentiBurstPyro"},
	"core.ICDTagVentiBurstHydro":       {"combat", "ICDTagVentiBurstHydro"},
	"core.ICDTagVentiBurstCryo":        {"combat", "ICDTagVentiBurstCryo"},
	"core.ICDTagVentiBurstElectro":     {"combat", "ICDTagVentiBurstElectro"},
	"core.ICDTagMonaWaterDamage":       {"combat", "ICDTagMonaWaterDamage"},
	"core.ICDTagTravelerWakeOfEarth":   {"combat", "ICDTagTravelerWakeOfEarth"},
	"core.ICDTagKleeFireDamage":        {"combat", "ICDTagKleeFireDamage"},
	"core.ICDTagTartagliaRiptideFlash": {"combat", "ICDTagTartagliaRiptideFlash"},
	"core.ICDReactionDamageDelim":      {"combat", "ICDReactionDamageDelim"},
	"core.ICDTagOverloadDamage":        {"combat", "ICDTagOverloadDamage"},
	"core.ICDTagSuperconductDamage":    {"combat", "ICDTagSuperconductDamage"},
	"core.ICDTagECDamage":              {"combat", "ICDTagECDamage"},
	"core.ICDTagShatter":               {"combat", "ICDTagShatter"},
	"core.ICDTagSwirlPyro":             {"combat", "ICDTagSwirlPyro"},
	"core.ICDTagSwirlHydro":            {"combat", "ICDTagSwirlHydro"},
	"core.ICDTagSwirlCryo":             {"combat", "ICDTagSwirlCryo"},
	"core.ICDTagSwirlElectro":          {"combat", "ICDTagSwirlElectro"},
	"core.ICDTagLength":                {"combat", "ICDTagLength"},
	"core.ICDGroupDefault":             {"combat", "ICDGroupDefault"},
	"core.ICDGroupAmber":               {"combat", "ICDGroupAmber"},
	"core.ICDGroupVenti":               {"combat", "ICDGroupVenti"},
	"core.ICDGroupFischl":              {"combat", "ICDGroupFischl"},
	"core.ICDGroupDiluc":               {"combat", "ICDGroupDiluc"},
	"core.ICDGroupPoleExtraAttack":     {"combat", "ICDGroupPoleExtraAttack"},
	"core.ICDGroupXiaoDash":            {"combat", "ICDGroupXiaoDash"},
	"core.ICDGroupReactionA":           {"combat", "ICDGroupReactionA"},
	"core.ICDGroupReactionB":           {"combat", "ICDGroupReactionB"},
	"core.ICDGroupLength":              {"combat", "ICDGroupLength"},
	"core.TargettableEnemy":            {"combat", "TargettableEnemy"},
	"core.TargettablePlayer":           {"combat", "TargettablePlayer"},
	"core.TargettableObject":           {"combat", "TargettableObject"},
	"core.TargettableTypeCount":        {"combat", "TargettableTypeCount"},
	"core.StrikeTypeDefault":           {"combat", "StrikeTypeDefault"},
	"core.StrikeTypePierce":            {"combat", "StrikeTypePierce"},
	"core.StrikeTypeBlunt":             {"combat", "StrikeTypeBlunt"},
	"core.StrikeTypeSlash":             {"combat", "StrikeTypeSlash"},
	"core.StrikeTypeSpear":             {"combat", "StrikeTypeSpear"},

	//actions
	"core.ActionType":                {"action", "Action"},
	"core.InvalidAction":             {"action", "InvalidAction"},
	"core.ActionSkill":               {"action", "ActionSkill"},
	"core.ActionBurst":               {"action", "ActionBurst"},
	"core.ActionAttack":              {"action", "ActionAttack"},
	"core.ActionCharge":              {"action", "ActionCharge"},
	"core.ActionHighPlunge":          {"action", "ActionHighPlunge"},
	"core.ActionLowPlunge":           {"action", "ActionLowPlunge"},
	"core.ActionAim":                 {"action", "ActionAim"},
	"core.ActionDash":                {"action", "ActionDash"},
	"core.ActionJump":                {"action", "ActionJump"},
	"core.ActionSwap":                {"action", "ActionSwap"},
	"core.ActionWalk":                {"action", "ActionWalk"},
	"core.ActionWait":                {"action", "ActionWait"},
	"core.ActionDelay":               {"action", "ActionDelay"},
	"core.EndActionType":             {"action", "EndActionType"},
	"core.ActionSkillHoldFramesOnly": {"action", "ActionSkillHoldFramesOnly"},
	//logs
	"core.LogProcs":             {"glog", "LogProcs"},
	"core.LogDamageEvent":       {"glog", "LogDamageEvent"},
	"core.LogPreDamageMod":      {"glog", "LogPreDamageMod"},
	"core.LogHurtEvent":         {"glog", "LogHurtEvent"},
	"core.LogHealEvent":         {"glog", "LogHealEvent"},
	"core.LogCalc":              {"glog", "LogCalc"},
	"core.LogReactionEvent":     {"glog", "LogReactionEvent"},
	"core.LogElementEvent":      {"glog", "LogElementEvent"},
	"core.LogSnapshotEvent":     {"glog", "LogSnapshotEvent"},
	"core.LogSnapshotModsEvent": {"glog", "LogSnapshotModsEvent"},
	"core.LogStatusEvent":       {"glog", "LogStatusEvent"},
	"core.LogActionEvent":       {"glog", "LogActionEvent"},
	"core.LogQueueEvent":        {"glog", "LogQueueEvent"},
	"core.LogEnergyEvent":       {"glog", "LogEnergyEvent"},
	"core.LogCharacterEvent":    {"glog", "LogCharacterEvent"},
	"core.LogEnemyEvent":        {"glog", "LogEnemyEvent"},
	"core.LogHookEvent":         {"glog", "LogHookEvent"},
	"core.LogSimEvent":          {"glog", "LogSimEvent"},
	"core.LogTaskEvent":         {"glog", "LogTaskEvent"},
	"core.LogArtifactEvent":     {"glog", "LogArtifactEvent"},
	"core.LogWeaponEvent":       {"glog", "LogWeaponEvent"},
	"core.LogShieldEvent":       {"glog", "LogShieldEvent"},
	"core.LogConstructEvent":    {"glog", "LogConstructEvent"},
	"core.LogICDEvent":          {"glog", "LogICDEvent"},
	//character
	"core.Character":     {"*character", "CharWrapper"},
	"core.ZoneMondstadt": {"*character", "ZoneMondstadt"},
	"core.ZoneLiyue":     {"*character", "ZoneLiyue"},
	"core.ZoneInazuma":   {"*character", "ZoneInazuma"},
	"core.ZoneSumeru":    {"*character", "ZoneSumeru"},
	//player
	"core.HealInfo":         {"player", "HealInfo"},
	"core.HealTypeAbsolute": {"player", "HealTypeAbsolute"},
	"core.HealTypePercent":  {"player", "HealTypePercent"},
	//shield
	"core.ShieldCrystallize":         {"shield", "ShieldCrystallize"},
	"core.ShieldNoelleSkill":         {"shield", "ShieldNoelleSkill"},
	"core.ShieldNoelleA1":            {"shield", "ShieldNoelleA1"},
	"core.ShieldZhongliJadeShield":   {"shield", "ShieldZhongliJadeShield"},
	"core.ShieldDionaSkill":          {"shield", "ShieldDionaSkill"},
	"core.ShieldBeidouThunderShield": {"shield", "ShieldBeidouThunderShield"},
	"core.ShieldXinyanSkill":         {"shield", "ShieldXinyanSkill"},
	"core.ShieldXinyanC2":            {"shield", "ShieldXinyanC2"},
	"core.ShieldKaeyaC4":             {"shield", "ShieldKaeyaC4"},
	"core.ShieldYanfeiC4":            {"shield", "ShieldYanfeiC4"},
	"core.ShieldBell":                {"shield", "ShieldBell"},
	"core.ShieldYunjinSkill":         {"shield", "ShieldYunjinSkill"},
	"core.EndShieldType":             {"shield", "EndShieldType"},

	//weapon
	"core.WeaponClassSword":    {"weapon", "WeaponClassSword"},
	"core.WeaponClassClaymore": {"weapon", "WeaponClassClaymore"},
	"core.WeaponClassSpear":    {"weapon", "WeaponClassSpear"},
	"core.WeaponClassBow":      {"weapon", "WeaponClassBow"},
	"core.WeaponClassCatalyst": {"weapon", "WeaponClassCatalyst"},
	"core.EndWeaponClass":      {"weapon", "EndWeaponClass"},
}
