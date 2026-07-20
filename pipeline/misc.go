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
		ref = excel.Filter(ref, func(v *excel.Curve) bool { return v.Level > 0 })
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
				k := ConvertEnum[model.GrowCurveType](typ, model.GrowCurveType_value, -1)
				if k == -1 {
					return fmt.Errorf("unknown curve=%v at lv=%v", typ, curve.Level)
				}
				c.GrowCurveData[k] = append(c.GrowCurveData[k], value)
			}
		}
	}
	return nil
}

func (c *Compiled) GenerateCurve() error {
	b := bytes.NewBuffer(nil)
	b.WriteString("package catalog\n")
	b.WriteString("import (\n")
	fmt.Fprintf(b, "\t\"%s\"\n", path.Join(baseModule, "pkg/model"))
	b.WriteString(")\n")

	b.WriteString("var GrowCurve = ")
	b.WriteString(dumpGo(c.GrowCurveData, false))
	b.WriteString("\n")

	writeFile("pkg/catalog/curve.dm.go", b.Bytes())
	return nil
}

func (c *Compiled) buildElementCoeff() error {
	ref := excel.Filter(excel.ElementCoeffExcelConfigData, func(v *excel.ElementCoeff) bool { return v.Level > 0 })
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
