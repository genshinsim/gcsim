package main

import (
	"bytes"
	"cmp"
	"fmt"
	"path"
	"slices"

	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/shizukayuki/excel-hk4e"
)

var curveTypes = map[Kind][]excel.GrowCurveType{
	KindCharacter: {
		excel.GROW_CURVE_ATTACK_S4,
		excel.GROW_CURVE_ATTACK_S5,
		excel.GROW_CURVE_HP_S4,
		excel.GROW_CURVE_HP_S5,
	},
	KindWeapon: {
		excel.GROW_CURVE_ATTACK_101,
		excel.GROW_CURVE_ATTACK_102,
		excel.GROW_CURVE_ATTACK_104,
		excel.GROW_CURVE_ATTACK_201,
		excel.GROW_CURVE_ATTACK_202,
		excel.GROW_CURVE_ATTACK_203,
		excel.GROW_CURVE_ATTACK_204,
		excel.GROW_CURVE_ATTACK_301,
		excel.GROW_CURVE_ATTACK_302,
		excel.GROW_CURVE_ATTACK_303,
		excel.GROW_CURVE_ATTACK_304,
		excel.GROW_CURVE_CRITICAL_101,
		excel.GROW_CURVE_CRITICAL_201,
		excel.GROW_CURVE_CRITICAL_301,
	},
	KindMonster: {
		excel.GROW_CURVE_HP,
		excel.GROW_CURVE_HP_2,
		excel.GROW_CURVE_HP_ENVIRONMENT,
	},
}

func (c *Compiled) buildCurveData() error {
	for kind, types := range curveTypes {
		var ref []*excel.Curve
		switch kind {
		case KindCharacter:
			ref = excel.AvatarCurveExcelConfigData
		case KindWeapon:
			ref = excel.WeaponCurveExcelConfigData
		case KindMonster:
			ref = excel.MonsterCurveExcelConfigData
		default:
			panic("unreachable")
		}
		ref = excel.Filter(ref, func(v *excel.Curve) bool { return v.Level >= 1 })
		slices.SortFunc(ref, func(a, b *excel.Curve) int { return cmp.Compare(a.Level, b.Level) })

		for lv, curve := range ref {
			lv := uint32(lv + 1)
			if curve.Level != lv {
				return fmt.Errorf("curve level mismatch: level=%v, expected=%v", curve.Level, lv)
			}
			for _, typ := range types {
				arith, value := curve.Type(typ)
				if arith != "ARITH_MULTI" {
					return fmt.Errorf("lv=%v,curve=%v has unexpected arith type: %v", curve.Level, typ, arith)
				}
				c.GrowCurveData[typ] = append(c.GrowCurveData[typ], value)
			}
		}
	}
	return nil
}

func (c *Compiled) GenerateCurve() error {
	var (
		avatar  []map[model.AvatarCurveType]float64
		monster []map[model.MonsterCurveType]float64
		weapon  []map[model.WeaponCurveType]float64
	)
	for _, curve := range curveTypes[KindCharacter] {
		typ := ConvertEnum[model.AvatarCurveType](curve, model.AvatarCurveType_value, -1)
		if typ == -1 {
			return fmt.Errorf("unknown curve=%v", curve)
		}
		data := c.GrowCurveData[curve]
		if len(avatar) == 0 {
			avatar = make([]map[model.AvatarCurveType]float64, len(data))
		}
		for lv, mul := range data {
			if avatar[lv] == nil {
				avatar[lv] = make(map[model.AvatarCurveType]float64)
			}
			avatar[lv][typ] = mul
		}
	}
	for _, curve := range curveTypes[KindWeapon] {
		typ := ConvertEnum[model.WeaponCurveType](curve, model.WeaponCurveType_value, -1)
		if typ == -1 {
			return fmt.Errorf("unknown curve=%v", curve)
		}
		data := c.GrowCurveData[curve]
		if len(weapon) == 0 {
			weapon = make([]map[model.WeaponCurveType]float64, len(data))
		}
		for lv, mul := range data {
			if weapon[lv] == nil {
				weapon[lv] = make(map[model.WeaponCurveType]float64)
			}
			weapon[lv][typ] = mul
		}
	}
	for _, curve := range curveTypes[KindMonster] {
		typ := ConvertEnum[model.MonsterCurveType](curve, model.MonsterCurveType_value, -1)
		if typ == -1 {
			return fmt.Errorf("unknown curve=%v", curve)
		}
		data := c.GrowCurveData[curve]
		if len(monster) == 0 {
			monster = make([]map[model.MonsterCurveType]float64, len(data))
		}
		for lv, mul := range data {
			if monster[lv] == nil {
				monster[lv] = make(map[model.MonsterCurveType]float64)
			}
			monster[lv][typ] = mul
		}
	}

	b := bytes.NewBuffer(nil)
	b.WriteString("package catalog\n")
	b.WriteString("import (\n")
	fmt.Fprintf(b, "\t\"%s\"\n", path.Join(baseModule, "pkg/model"))
	b.WriteString(")\n")

	b.WriteString("var AvatarGrowCurveByLvl = ")
	b.WriteString(dumpGo(avatar, false))
	b.WriteString("\n")

	b.WriteString("var WeaponGrowCurveByLvl = ")
	b.WriteString(dumpGo(weapon, false))
	b.WriteString("\n")

	b.WriteString("var EnemyStatGrowthMult = ")
	b.WriteString(dumpGo(monster, false))
	b.WriteString("\n")

	writeFile("pkg/catalog/curve.dm.go", b.Bytes())
	return nil
}

func (c *Compiled) buildElementCoeff() error {
	ref := excel.Filter(excel.ElementCoeffExcelConfigData, func(v *excel.ElementCoeff) bool { return v.Level >= 1 })
	slices.SortFunc(ref, func(a, b *excel.ElementCoeff) int { return cmp.Compare(a.Level, b.Level) })
	for lv, coeff := range ref {
		lv := uint32(lv + 1)
		if coeff.Level != lv {
			return fmt.Errorf("element coeff level mismatch: level=%v, expected=%v", coeff.Level, lv)
		}
		c.ElementCoeff = append(c.ElementCoeff, coeff.PlayerElementLevelCo)
	}
	return nil
}

func (c *Compiled) GenerateElementCoeff() error {
	b := bytes.NewBuffer(nil)
	b.WriteString("package combat\n")
	b.WriteString("var reactionLvlBase = ")
	b.WriteString(dumpGo(c.ElementCoeff, false))
	b.WriteString("\n")
	writeFile("pkg/core/combat/reaction.dm.go", b.Bytes())
	return nil
}

func (c *Compiled) GenerateLocalization() error {
	data, err := dumpJSON(c.Localization)
	if err != nil {
		return fmt.Errorf("failed to marshal localization json: %w", err)
	}
	writeFile("ui/packages/localization/src/locales/names.dm.json", data)
	return nil
}

func (c *Compiled) GenerateEditorJS() error {
	var chars []string
	for _, config := range c.Configuration {
		if config.Kind != KindCharacter {
			continue
		}
		chars = append(chars, config.Character.Model.Key)
		chars = append(chars, config.Shortcuts...)
	}
	input := struct{ Characters []string }{
		Characters: chars,
	}
	data := useTemplate("ui_editor.js.templ", input)
	writeFile("ui/packages/components/src/Editor/mode-gcsim.dm.js", data)
	writeFile("ui/packages/ui/src/util/mode-gcsim.dm.js", data)
	return nil
}
