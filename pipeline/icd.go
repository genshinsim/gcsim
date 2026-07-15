package main

import (
	"bytes"
	"cmp"
	"fmt"
	"maps"
	"slices"

	"github.com/shizukayuki/excel-hk4e"
)

func (c *Compiled) buildICDGroup() error {
	for name, att := range c.ICDGroup {
		if att.Timer != 0 && len(att.Damage) > 0 && len(att.Durability) > 0 {
			continue
		}
		refs := excel.Filter(excel.AttackAttenuationExcelConfigData, func(v *excel.AttackAttenuation) bool {
			return excel.Slug(v.Group) == att.Name
		})
		if len(refs) != 1 {
			return fmt.Errorf("attack_attenuation=%v results in refs=%v but we expect 1", excel.Slug(att.Name), len(refs))
		}
		ref := refs[0]
		if att.Timer == 0 {
			att.Timer = ref.ResetCycle
		}
		if len(att.Damage) == 0 {
			att.Damage = ref.DamageSequence
		}
		if len(att.Durability) == 0 {
			att.Durability = ref.DurabilitySequence
		}
		c.ICDGroup[name] = att
	}
	return nil
}

func (c *Compiled) GenerateICDGroup() error {
	sorted := slices.SortedFunc(maps.Keys(c.ICDGroup), func(a, b string) int {
		n := cmp.Compare(a, b)
		if n == 0 {
			return n
		}
		for _, s := range []string{
			"Default",
			"PoleExtraAttack",
			"ReactionA",
			"ReactionB",
			"Burning",
		} {
			switch s {
			case a:
				return -1
			case b:
				return +1
			}
		}
		return n
	})
	typ := "ICDGroup"
	consts := make(map[string]string)

	b := bytes.NewBuffer(nil)
	b.WriteString("package attacks\n")
	fmt.Fprintf(b, "type %s int\n", typ)
	b.WriteString("const (\n")
	for i, name := range sorted {
		consts[name] = typ + excel.Slug(name)
		b.WriteString("\t")
		b.WriteString(consts[name])
		if i == 0 {
			fmt.Fprintf(b, " %s = iota", typ)
		}
		b.WriteString("\n")
	}
	b.WriteString(")\n")

	b.WriteString("\n")
	fmt.Fprintf(b, "var %s = []int{\n", "ICDGroupResetTimer")
	for _, name := range sorted {
		att := c.ICDGroup[name]
		fmt.Fprintf(b, "\t%s: %v, // %vs\n", consts[name], int(att.Timer*60), att.Timer)
	}
	b.WriteString("}\n")

	b.WriteString("\n")
	fmt.Fprintf(b, "var %s = [][]float64{\n", "ICDGroupEleApplicationSequence")
	for _, name := range sorted {
		att := c.ICDGroup[name]
		fmt.Fprintf(b, "\t%s: %v,\n", consts[name], dumpGo(att.Durability, true))
	}
	b.WriteString("}\n")

	b.WriteString("\n")
	fmt.Fprintf(b, "var %s = [][]float64{\n", "ICDGroupDamageSequence")
	for _, name := range sorted {
		att := c.ICDGroup[name]
		fmt.Fprintf(b, "\t%s: %v,\n", consts[name], dumpGo(att.Damage, true))
	}
	b.WriteString("}\n")

	writeFile("pkg/core/attacks/icd_groups.dm.go", b.Bytes())
	return nil
}

func (c *Compiled) GenerateICDTag() error {
	typ := "ICDTag"
	tags := make(map[string]struct{})
	for _, config := range c.Configuration {
		for _, name := range config.Combat.ICDTags {
			name = typ + excel.Slug(name)
			tags[name] = struct{}{}
		}
	}
	var sorted []string
	sorted = append(sorted, typ+"None")
	sorted = append(sorted, slices.Sorted(maps.Keys(tags))...)
	sorted = append(sorted, typ+"Length")

	b := bytes.NewBuffer(nil)
	b.WriteString("package attacks\n")
	fmt.Fprintf(b, "type %s int\n", typ)
	b.WriteString("const (\n")
	for i, name := range sorted {
		b.WriteString("\t")
		b.WriteString(name)
		if i == 0 {
			fmt.Fprintf(b, " %s = iota", typ)
		}
		b.WriteString("\n")
	}
	b.WriteString(")\n")

	writeFile("pkg/core/attacks/icd_tags.dm.go", b.Bytes())
	return nil
}
