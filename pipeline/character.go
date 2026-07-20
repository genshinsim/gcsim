package main

import (
	"bytes"
	"cmp"
	"fmt"
	"path"
	"reflect"
	"slices"
	"strings"

	"github.com/genshinsim/gcsim/pkg/model"
	"github.com/shizukayuki/excel-hk4e"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type CharacterSpec struct {
	ref   *excel.Avatar           `yaml:"-"`
	depot *excel.AvatarSkillDepot `yaml:"-"`

	Name  string            `yaml:"name,omitempty"`
	Model *model.AvatarData `yaml:"model,omitempty"`
}

func (s *CharacterSpec) ClearRef() {
	if s != nil {
		s.ref = nil
		s.depot = nil
	}
}

func buildCharacterSpec(cfg *Config) (*CharacterSpec, error) {
	refs := excel.Filter(excel.AvatarExcelConfigData, func(v *excel.Avatar) bool {
		if id := cfg.Override.Id; id != 0 {
			return v.Id == id
		}
		if v.Codex() == nil {
			return false
		}
		return excel.SlugLower(v.Name()) == cfg.Name
	})
	if len(refs) != 1 {
		return nil, fmt.Errorf("query results in refs=%v but we expect 1", len(refs))
	}

	spec := &CharacterSpec{ref: refs[0]}
	fetter := spec.ref.FetterInfo()
	if fetter == nil {
		return nil, fmt.Errorf("no fetter found for avatar_id=%v", spec.ref.Id)
	}
	promote := excel.Filter(excel.AvatarPromoteExcelConfigData, func(v *excel.AvatarPromote) bool { return v.AvatarPromoteId == spec.ref.AvatarPromoteId })
	if len(promote) == 0 {
		return nil, fmt.Errorf("no promote found for avatar_id=%v", spec.ref.Id)
	}
	slices.SortFunc(promote, func(a, b *excel.AvatarPromote) int { return cmp.Compare(a.PromoteLevel, b.PromoteLevel) })

	spec.Name = spec.ref.Name()
	spec.Model = &model.AvatarData{
		Id:          spec.ref.Id,
		SubId:       spec.ref.SkillDepotId,
		Key:         excel.SlugLower(spec.Name),
		Rarity:      ConvertEnum[model.QualityType](spec.ref.QualityType, model.QualityType_value, -1),
		Body:        ConvertEnum[model.BodyType](spec.ref.BodyType, model.BodyType_value, -1),
		Region:      ConvertEnum[model.AssocType](fetter.AvatarAssocType, model.AssocType_value, -1),
		WeaponClass: ConvertEnum[model.WeaponType](spec.ref.WeaponType, model.WeaponType_value, -1),
		IconName:    spec.ref.IconName,
		Stats: &model.AvatarStatsData{
			BaseHp:         spec.ref.HpBase,
			BaseAtk:        spec.ref.AttackBase,
			BaseDef:        spec.ref.DefenseBase,
			ElementMastery: spec.ref.ElementMastery,
		},
		SkillDetails: &model.AvatarSkillsData{},
	}
	if spec.Model.Rarity == -1 {
		return nil, fmt.Errorf("unknown quality_type=%v", spec.ref.QualityType)
	}
	if spec.Model.Body == -1 {
		return nil, fmt.Errorf("unknown body_type=%v", spec.ref.BodyType)
	}
	if spec.Model.Region == -1 {
		return nil, fmt.Errorf("unknown assoc_type=%v", fetter.AvatarAssocType)
	}
	if spec.Model.WeaponClass == -1 {
		return nil, fmt.Errorf("unknown weapon_type=%v", spec.ref.WeaponType)
	}

	if v := cfg.Override.Depot; v != 0 {
		spec.Model.SubId = v
	}
	spec.depot = excel.FindSkillDepot(spec.Model.SubId)
	if spec.depot == nil {
		return nil, fmt.Errorf("depot=%v not found", spec.Model.SubId)
	}
	attack := excel.FindSkill(spec.depot.Skills[0])
	skill := excel.FindSkill(spec.depot.Skills[1])
	burst := excel.FindSkill(spec.depot.EnergySkill)
	if attack == nil {
		return nil, fmt.Errorf("depot=%v,attack=%v not found", spec.depot.Id, spec.depot.Skills[0])
	}
	if skill == nil {
		return nil, fmt.Errorf("depot=%v,skill=%v not found", spec.depot.Id, spec.depot.Skills[1])
	}
	if burst == nil {
		return nil, fmt.Errorf("depot=%v,burst=%v not found", spec.depot.Id, spec.depot.EnergySkill)
	}
	spec.Model.SkillDetails.Attack = attack.Id
	spec.Model.SkillDetails.Skill = skill.Id
	spec.Model.SkillDetails.Burst = burst.Id

	if spec.ref.Name() == "Traveler" {
		switch spec.Model.Body {
		case model.BodyType_BODY_BOY:
			spec.Name = excel.FindManualTextMap("INFO_MALE_PRONOUN_KONG").Name()
		case model.BodyType_BODY_GIRL:
			spec.Name = excel.FindManualTextMap("INFO_MALE_PRONOUN_YING").Name()
		default:
			panic("unreachable")
		}
		spec.Model.Key = excel.SlugLower(spec.Name)
	}
	if len(spec.ref.CandSkillDepotIds) > 0 {
		spec.Name = fmt.Sprintf("%s (%s)", spec.Name, excel.FindManualTextMap(burst.CostElemType.String()).Name())
		spec.Model.Key = excel.SlugLower(spec.Name)
	}

	for curve, prop := range map[*model.GrowCurveType]excel.FightProp{
		&spec.Model.Stats.HpCurve:  excel.FIGHT_PROP_BASE_HP,
		&spec.Model.Stats.AtkCurve: excel.FIGHT_PROP_BASE_ATTACK,
		&spec.Model.Stats.DefCruve: excel.FIGHT_PROP_BASE_DEFENSE,
	} {
		v := excel.Find(spec.ref.PropGrowCurves, func(v *excel.FightPropGrow) bool { return v.Type == prop })
		if curve := v.GrowCurve; !slices.Contains(curveTypes[KindCharacter], curve) {
			return nil, fmt.Errorf("curve not listed in known types: %v", curve)
		}
		*curve = ConvertEnum[model.GrowCurveType](v.GrowCurve, model.GrowCurveType_value, -1)
		if *curve == -1 {
			return nil, fmt.Errorf("unknown curve=%v", v.GrowCurve)
		}
	}

	for _, v := range promote {
		props, err := ConvertAddProps(v.AddProps)
		if err != nil {
			return nil, err
		}
		spec.Model.Stats.PromoData = append(spec.Model.Stats.PromoData, &model.PromotionData{
			MaxLevel: v.UnlockMaxLevel,
			AddProps: props,
		})
	}

	spec.Model.SkillDetails.BurstEnergyCost = burst.CostElemVal
	if ele := ConvertEnum[model.ElementType](burst.CostElemType, model.ElementType_value, -1); ele != -1 {
		spec.Model.Element = ele
	} else {
		return nil, fmt.Errorf("unknown element=%v", burst.CostElemType)
	}

	//nolint:goconst // using a constant is not that useful here
	skills := map[string]uint32{
		"attack": attack.ProudSkillGroupId,
		"skill":  skill.ProudSkillGroupId,
		"burst":  burst.ProudSkillGroupId,
	}
	for _, v := range spec.depot.InherentProudSkillOpens {
		if v.ProudSkillGroupId == 0 {
			continue
		}
		if v.NeedAvatarPromoteLevel > 0 {
			typ := fmt.Sprintf("a%d", v.NeedAvatarPromoteLevel)
			skills[typ] = v.ProudSkillGroupId
		}
	}
	var links []uint32
	for _, typ := range abilities {
		if _, ok := skills[typ]; !ok {
			continue
		}
		pds := excel.Filter(excel.ProudSkillExcelConfigData, func(v *excel.ProudSkill) bool { return v.ProudSkillGroupId == skills[typ] })
		slices.SortFunc(pds, func(a, b *excel.ProudSkill) int { return cmp.Compare(a.Level, b.Level) })

		pd := pds[0]
		for i := range pds {
			if pds[i].Level == 6 { // reference level
				pd = pds[i]
				break
			}
		}

		skill := pd.FindSkill()
		if skill != nil {
			attr := &AttributeSpec{
				Type:   typ,
				Name:   skill.Name(),
				Desc:   skill.DescTextMapHash.String(),
				Config: pd.OpenConfig,
			}
			links = append(links, ExtractLinks(attr.Desc)...)
			attr.SetValues(0, nil)
			cfg.Attributes = append(cfg.Attributes, attr)
		}

		hasParamDesc := false
		for _, desc := range pd.ParamDescList {
			text, params, ok := strings.Cut(desc.String(), "|")
			if !ok {
				continue
			}
			hasParamDesc = true
			text, params = strings.TrimSpace(text), strings.TrimSpace(params)
			attr := &AttributeSpec{
				Type:   typ,
				Name:   skill.Name(),
				Desc:   text,
				Config: pd.OpenConfig,
			}
			attr.Index, attr.ParamDesc = IndexFromParams(params)
			attr.SetValues(len(pds), func(i int) []float64 { return pds[i].ParamList })
			cfg.Attributes = append(cfg.Attributes, attr)
		}
		if !hasParamDesc {
			attr := &AttributeSpec{
				Type:   typ,
				Name:   pd.NameTextMapHash.String(),
				Desc:   pd.DescTextMapHash.String(),
				Config: pd.OpenConfig,
				Index:  NewTalentIndex(pd.ParamList),
			}
			links = append(links, ExtractLinks(attr.Desc)...)
			attr.SetValues(len(pds), func(i int) []float64 { return pds[i].ParamList })
			cfg.Attributes = append(cfg.Attributes, attr)
		}
	}
	for con, id := range spec.depot.Talents {
		t := excel.FindTalent(id)
		attr := &AttributeSpec{
			Type:   fmt.Sprintf("c%d", con+1),
			Name:   t.NameTextMapHash.String(),
			Desc:   t.DescTextMapHash.String(),
			Config: t.OpenConfig,
			Index:  NewTalentIndex(t.ParamList),
		}
		links = append(links, ExtractLinks(attr.Desc)...)
		attr.SetValues(1, func(i int) []float64 { return t.ParamList })
		cfg.Attributes = append(cfg.Attributes, attr)
	}

	seen := make(map[uint32]struct{})
	for _, id := range links {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		link := excel.FindHyperLink(id)
		attr := &AttributeSpec{
			Type: fmt.Sprintf("effect%d", len(seen)),
			Name: link.NameTextMapHash.String(),
			Desc: link.DescTextMapHash.String(),
		}
		attr.SetValues(0, nil)
		cfg.Attributes = append(cfg.Attributes, attr)
	}

	return spec, nil
}

func (c *Compiled) GenerateCharacters() error {
	kind := KindCharacter
	inputs := excel.Filter(c.Configuration, func(v *Config) bool { return v.Kind == kind })

	imports := ImportTmpl{Kind: kind}
	keys := KeysTmpl{Kind: kind}
	assets := AssetsTmpl{Kind: kind, Variable: "avatarMap"}
	shortcut := ShortcutTmpl{Kind: kind, Variable: "CharNameToKey", Type: "keys.Char"}
	catalog := CatalogTmpl{Kind: kind, Variable: "CharacterMap", Type: "keys.Char", ModelName: reflect.TypeFor[*model.AvatarData]().String()}
	doc := DocTmpl{Kind: kind}
	issues := IssuesTmpl{Kind: kind}
	fields := FieldsTmpl{Kind: kind}
	frames := FramesTmpl{Kind: kind}
	actions := ActionsTmpl{Kind: kind}
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
	defer frames.Write()
	defer actions.Write()
	defer hitbox.Write()
	defer hitlag.Write()
	defer params.Write()

	models := &model.AvatarDataMap{Data: make(map[string]*model.AvatarData)}
	for _, config := range inputs {
		spec := config.Character

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
		frames.Key = append(frames.Key, spec.Model.Key)
		frames.Data = append(frames.Data, config.Docs.Frames)
		actions.Key = append(actions.Key, spec.Model.Key)
		actions.Data = append(actions.Data, config.Abilities)
		hitbox.Key = append(hitbox.Key, spec.Model.Key)
		hitbox.Data = append(hitbox.Data, config.Abilities)
		hitlag.Key = append(hitlag.Key, spec.Model.Key)
		hitlag.Data = append(hitlag.Data, config.Abilities)
		params.Key = append(params.Key, spec.Model.Key)
		params.Data = append(params.Data, config.Abilities)

		assets.Key = append(assets.Key, spec.Model.Key)
		assets.Image = append(assets.Image, spec.Model.IconName)

		catalog.Slug = append(catalog.Slug, excel.Slug(spec.Name))
		catalog.Model = append(catalog.Model, proto.Clone(spec.Model))

		m := proto.Clone(spec.Model).(*model.AvatarData)
		m.Stats = nil
		models.Data[m.Key] = m

		b := bytes.NewBuffer(nil)
		// seen := make(map[string]struct{})
		// for _, attr := range config.Attributes {
		// 	if _, ok := seen[attr.Type]; ok {
		// 		continue
		// 	}
		// 	seen[attr.Type] = struct{}{}
		// 	b.WriteString(attr.EmitDesc("// "))
		// }
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
		fmt.Fprintf(b, "core.Register%[2]sFunc(keys.%[1]s, New%[2]s)\n", excel.Slug(spec.Name), "Char")
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
	writeFile(fmt.Sprintf("ui/packages/db/src/Data/%s.dm.json", kind), data)
	writeFile(fmt.Sprintf("ui/packages/ui/src/Data/%s.dm.json", kind), data)

	return nil
}
