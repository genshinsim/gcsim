package main

import (
	"fmt"
	"log/slog"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/shizukayuki/excel-hk4e"
	"google.golang.org/protobuf/proto"
)

type MonsterSpec struct {
	ref *excel.Monster `yaml:"-"`

	Name  string             `yaml:"name,omitempty"`
	Model *model.MonsterData `yaml:"model,omitempty"`
}

func (s *MonsterSpec) ClearRef() {
	if s != nil {
		s.ref = nil
	}
}

func getMonsters() []*excel.Monster {
	visited := make(map[string]struct{})
	return excel.Filter(excel.MonsterExcelConfigData, func(v *excel.Monster) bool {
		desc := v.Describe()
		if desc == nil {
			return false
		}
		slug := excel.SlugLower(desc.Name())
		if slug == "" {
			return false
		}

		curve := excel.Filter(v.PropGrowCurves, func(v *excel.FightPropGrow) bool { return v.Type == excel.FIGHT_PROP_BASE_HP })
		if len(curve) != 1 {
			Log(slog.LevelWarn, "query for id=%v,monster=%v curve results in refs=%v but we expect 1", v.Id, excel.Slug(desc.Name()), len(curve))
			return false
		}

		if _, ok := visited[slug]; ok {
			return false
		}
		visited[slug] = struct{}{}
		return true
	})
}

func buildMonsterSpec(ref *excel.Monster) (*MonsterSpec, error) {
	spec := &MonsterSpec{ref: ref}
	desc := spec.ref.Describe()

	spec.Name = desc.Name()
	spec.Model = &model.MonsterData{
		Id:  spec.ref.Id,
		Key: excel.SlugLower(spec.Name),
		BaseStats: &model.MonsterStatsData{
			BaseHp: spec.ref.HpBase,
			Resist: &model.MonsterResistData{
				FireResist:     spec.ref.FireSubHurt,
				GrassResist:    spec.ref.GrassSubHurt,
				WaterResist:    spec.ref.WaterSubHurt,
				ElectricResist: spec.ref.ElecSubHurt,
				WindResist:     spec.ref.WindSubHurt,
				IceResist:      spec.ref.IceSubHurt,
				RockResist:     spec.ref.RockSubHurt,
				PhysicalResist: spec.ref.PhysicalSubHurt,
			},
			HpDrop: []*model.MonsterHPDrop{},
		},
	}
	if spec.ref.Type == "MONSTER_BOSS" {
		spec.Model.BaseStats.FreezeResist = 1
	}

	curve := excel.Find(spec.ref.PropGrowCurves, func(v *excel.FightPropGrow) bool { return v.Type == excel.FIGHT_PROP_BASE_HP })
	if curve := curve.GrowCurve; !slices.Contains(curveTypes[KindMonster], curve) {
		return nil, fmt.Errorf("curve not listed in known types: %v", curve)
	}
	if typ := ConvertEnum[model.GrowCurveType](curve.GrowCurve, model.GrowCurveType_value, -1); typ != -1 {
		spec.Model.BaseStats.HpCurve = typ
	} else {
		return nil, fmt.Errorf("unknown curve=%v", curve.GrowCurve)
	}

	add := func(drop *model.MonsterHPDrop) {
		spec.Model.BaseStats.HpDrop = append(spec.Model.BaseStats.HpDrop, drop)
	}
	for _, d := range excel.Filter(spec.ref.HpDrops, func(v *excel.MonsterDrop) bool { return v.DropId != 0 && v.HpPercent != 0 }) {
		add(&model.MonsterHPDrop{
			DropId:    d.DropId,
			HpPercent: d.HpPercent / 100,
		})
	}
	if id := spec.ref.KillDropId; id != 0 {
		add(&model.MonsterHPDrop{DropId: id})
	}

	return spec, nil
}

func (c *Compiled) GenerateMonsters() error {
	kind := KindMonster
	inputs := excel.Filter(c.Configuration, func(v *Config) bool { return v.Kind == kind })

	shortcut := ShortcutTmpl{Kind: kind, Variable: "MonsterNameToID", Type: "int"}
	catalog := CatalogTmpl{Kind: kind, Variable: "MonsterMap", Type: "int", ModelName: reflect.TypeFor[*model.MonsterData]().String()}
	doc := DocTmpl{Kind: kind}
	defer shortcut.Write()
	defer catalog.Write()
	defer doc.Write()

	type HPData struct {
		Level int     `json:"level"`
		HP    float64 `json:"hp"`
	}
	type ResistData struct {
		Pyro     float64 `json:"pyro"`
		Dendro   float64 `json:"dendro"`
		Hydro    float64 `json:"hydro"`
		Electro  float64 `json:"electro"`
		Anemo    float64 `json:"anemo"`
		Cryo     float64 `json:"cryo"`
		Geo      float64 `json:"geo"`
		Physical float64 `json:"physical"`
		Freeze   float64 `json:"freeze"`
	}
	type ParticleData struct {
		DropId    uint32  `json:"drop_id"`
		HpPercent float64 `json:"hp_percent"`
	}

	hpData := make(map[string][]HPData)
	particleData := make(map[string][]ParticleData)
	resistData := make(map[string]ResistData)
	for _, config := range inputs {
		spec := config.Monster

		doc.Key = append(doc.Key, spec.Model.Key)
		doc.Name = append(doc.Name, spec.Name)

		shortcut.Slug = append(shortcut.Slug, strconv.Itoa(int(spec.Model.Id)))
		shortcut.Names = append(shortcut.Names, []string{spec.Model.Key})

		stats := spec.Model.BaseStats
		for _, lv := range []int{1, 10, 20, 30, 40, 50, 60, 70, 80, 90, 95, 98, 100} {
			hpData[spec.Model.Key] = append(hpData[spec.Model.Key], HPData{
				Level: lv,
				HP:    stats.BaseHp * c.GrowCurveData[stats.HpCurve][lv-1],
			})
		}
		for _, drop := range stats.HpDrop {
			particleData[spec.Model.Key] = append(particleData[spec.Model.Key], ParticleData{
				DropId:    drop.DropId,
				HpPercent: drop.HpPercent,
			})
		}
		resistData[spec.Model.Key] = ResistData{
			Pyro:     stats.Resist.FireResist,
			Dendro:   stats.Resist.GrassResist,
			Hydro:    stats.Resist.WaterResist,
			Electro:  stats.Resist.ElectricResist,
			Anemo:    stats.Resist.WindResist,
			Cryo:     stats.Resist.IceResist,
			Geo:      stats.Resist.RockResist,
			Physical: stats.Resist.PhysicalResist,
			Freeze:   stats.FreezeResist,
		}

		catalog.Slug = append(catalog.Slug, strconv.Itoa(int(spec.Model.Id)))
		catalog.Model = append(catalog.Model, proto.Clone(spec.Model))
	}

	for typ, input := range map[string]any{
		"HP":       hpData,
		"Particle": particleData,
		"Resist":   resistData,
	} {
		data, err := dumpJSON(input)
		if err != nil {
			return fmt.Errorf("failed to marshal monster %s data: %w", strings.ToLower(typ), err)
		}
		writeFile(fmt.Sprintf("ui/packages/docs/src/components/%s/%s.dm.json", typ, kind), data)
	}

	return nil
}
