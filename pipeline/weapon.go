package main

import (
	"bytes"
	"cmp"
	"fmt"
	"path"
	"reflect"
	"slices"
	"strconv"

	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/shizukayuki/excel-hk4e"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type WeaponSpec struct {
	ref *excel.Weapon `yaml:"-"`

	Name  string            `yaml:"name,omitempty"`
	Model *model.WeaponData `yaml:"model,omitempty"`
}

func (s *WeaponSpec) ClearRef() {
	if s != nil {
		s.ref = nil
	}
}

func buildWeaponSpec(cfg *Config) (*WeaponSpec, error) {
	refs := excel.Filter(excel.WeaponExcelConfigData, func(v *excel.Weapon) bool {
		if id := cfg.Override.Id; id != 0 {
			return v.Id == id
		}
		return v.StoryId != 0 && excel.SlugLower(v.Name()) == cfg.Name
	})
	if len(refs) != 1 {
		return nil, fmt.Errorf("query results in refs=%v but we expect 1", len(refs))
	}

	spec := &WeaponSpec{ref: refs[0]}
	promote := excel.Filter(excel.WeaponPromoteExcelConfigData, func(v *excel.WeaponPromote) bool { return v.WeaponPromoteId == spec.ref.WeaponPromoteId })
	if len(promote) == 0 {
		return nil, fmt.Errorf("no promote found for weapon_id=%v", spec.ref.Id)
	}
	slices.SortFunc(promote, func(a, b *excel.WeaponPromote) int { return cmp.Compare(a.PromoteLevel, b.PromoteLevel) })

	spec.Name = spec.ref.Name()
	spec.Model = &model.WeaponData{
		Id:              int32(spec.ref.Id),
		Key:             excel.SlugLower(spec.Name),
		Rarity:          int32(spec.ref.RankLevel),
		WeaponClass:     ConvertEnum[model.WeaponClass](spec.ref.WeaponType, model.WeaponClass_value, -1),
		ImageName:       spec.ref.AwakenIcon,
		BaseStats:       &model.WeaponStatsData{},
		NameTextHashMap: int64(spec.ref.NameTextMapHash),
	}
	if spec.Model.WeaponClass == -1 {
		return nil, fmt.Errorf("unknown weapon_type=%v", spec.ref.WeaponType)
	}

	for _, add := range spec.ref.WeaponProp {
		if add.InitValue == 0 || add.PropType == excel.FIGHT_PROP_NONE {
			continue
		}
		if curve := add.Type; !slices.Contains(curveTypes[KindWeapon], curve) {
			return nil, fmt.Errorf("curve not listed in known types: %v", curve)
		}
		typ := ConvertEnum[model.StatType](add.PropType, model.StatType_value, -1)
		curve := ConvertEnum[model.WeaponCurveType](add.Type, model.WeaponCurveType_value, -1)
		if typ == -1 {
			return nil, fmt.Errorf("unknown prop=%v", add.PropType)
		}
		if curve == -1 {
			return nil, fmt.Errorf("unknown curve=%v", add.Type)
		}
		spec.Model.BaseStats.BaseProps = append(spec.Model.BaseStats.BaseProps, &model.WeaponProp{
			PropType:     typ,
			InitialValue: add.InitValue,
			Curve:        curve,
		})
	}

	for _, v := range promote {
		props, err := ConvertAddProps(v.AddProps)
		if err != nil {
			return nil, err
		}
		spec.Model.BaseStats.PromoData = append(spec.Model.BaseStats.PromoData, &model.PromotionData{
			MaxLevel: int32(v.UnlockMaxLevel),
			AddProps: props,
		})
	}

	for num, id := range spec.ref.SkillAffix {
		affixes := excel.Filter(excel.EquipAffixExcelConfigData, func(v *excel.EquipAffix) bool { return v.Id == id })
		slices.SortFunc(affixes, func(a, b *excel.EquipAffix) int { return cmp.Compare(a.Level, b.Level) })
		if len(affixes) == 0 {
			continue
		}
		attr := &AttributeSpec{
			Type:   "passive",
			Name:   affixes[0].NameTextMapHash.String(),
			Desc:   affixes[0].DescTextMapHash.String(),
			Config: affixes[0].OpenConfig,
			Index:  NewTalentIndex(affixes[0].ParamList),
		}
		if num > 0 {
			attr.Type += strconv.Itoa(num)
		}
		attr.SetValues(len(affixes), func(i int) []float64 { return affixes[i].ParamList })
		cfg.Attributes = append(cfg.Attributes, attr)
	}

	return spec, nil
}

func (c *Compiled) GenerateWeapons() error {
	kind := KindWeapon
	inputs := excel.Filter(c.Configuration, func(v *Config) bool { return v.Kind == kind })

	imports := ImportTmpl{Kind: kind}
	keys := KeysTmpl{Kind: kind}
	assets := AssetsTmpl{Kind: kind, Variable: "weaponMap"}
	shortcut := ShortcutTmpl{Kind: kind, Variable: "WeaponNameToKey", Type: "keys.Weapon"}
	catalog := CatalogTmpl{Kind: kind, Variable: "WeaponMap", Type: "keys.Weapon", ModelName: reflect.TypeFor[*model.WeaponData]().String()}
	doc := DocTmpl{Kind: kind}
	issues := IssuesTmpl{Kind: kind}
	fields := FieldsTmpl{Kind: kind}
	hitbox := HitboxTmpl{Kind: kind}
	hitlag := HitlagTmpl{Kind: kind}
	params := ParamsTmpl{Kind: kind}
	defer imports.Write()
	defer keys.Write()
	defer assets.Write()
	defer shortcut.Write()
	defer catalog.Write()
	defer doc.Write()
	defer issues.Write()
	defer fields.Write()
	defer hitbox.Write()
	defer hitlag.Write()
	defer params.Write()

	models := &model.WeaponDataMap{Data: make(map[string]*model.WeaponData)}
	for _, config := range inputs {
		spec := config.Weapon

		imports.PackagePath = append(imports.PackagePath, config.Dir())
		keys.Name = append(keys.Name, spec.Name)
		shortcut.Slug = append(shortcut.Slug, excel.Slug(spec.Name))
		shortcut.Names = append(shortcut.Names, append([]string{spec.Model.Key}, config.Shortcuts...))

		doc.Key = append(doc.Key, spec.Model.Key)
		doc.Name = append(doc.Name, spec.Name)
		issues.Key = append(issues.Key, spec.Model.Key)
		issues.Data = append(issues.Data, config.Docs.Issues)
		fields.Key = append(fields.Key, spec.Model.Key)
		fields.Data = append(fields.Data, config.Docs.Fields)
		hitbox.Key = append(hitbox.Key, spec.Model.Key)
		hitbox.Data = append(hitbox.Data, config.Abilities)
		hitlag.Key = append(hitlag.Key, spec.Model.Key)
		hitlag.Data = append(hitlag.Data, config.Abilities)
		params.Key = append(params.Key, spec.Model.Key)
		params.Data = append(params.Data, config.Abilities)

		assets.Key = append(assets.Key, spec.Model.Key)
		assets.Image = append(assets.Image, spec.Model.ImageName)

		catalog.Slug = append(catalog.Slug, excel.Slug(spec.Name))
		catalog.Model = append(catalog.Model, proto.Clone(spec.Model))

		m := proto.Clone(spec.Model).(*model.WeaponData)
		m.BaseStats = nil
		models.Data[m.Key] = m

		b := bytes.NewBuffer(nil)
		for _, attr := range config.Attributes {
			b.WriteString(attr.EmitDesc("// "))
		}
		fmt.Fprintf(b, "package %s\n", path.Base(config.Dir()))
		b.WriteString("import (\n")
		for _, pkg := range []string{
			path.Join(baseModule, "pkg/core"),
			path.Join(baseModule, "pkg/core/keys"),
		} {
			fmt.Fprintf(b, "\t\"%s\"\n", pkg)
		}
		b.WriteString(")\n")

		b.WriteString("func init() {\n")
		fmt.Fprintf(b, "core.Register%[2]sFunc(keys.%[1]s, New%[2]s)\n", excel.Slug(spec.Name), "Weapon")
		b.WriteString("}\n")

		if err := emitTalents(b, config.Talents, config.Attributes); err != nil {
			return fmt.Errorf("%v: %w", config.Path, err)
		}
		writeFile(path.Join(config.Dir(), fmt.Sprintf("zz_%s.dm.go", spec.Model.Key)), b.Bytes())
	}

	data, err := protojson.Marshal(models)
	if err != nil {
		return fmt.Errorf("failed to marshal %v models: %w", kind, err)
	}
	writeFile(fmt.Sprintf("ui/packages/ui/src/Data/%s.dm.json", kind), data)

	return nil
}
