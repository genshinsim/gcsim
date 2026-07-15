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

type ArtifactSpec struct {
	ref *excel.ReliquarySet `yaml:"-"`

	Name  string              `yaml:"name,omitempty"`
	Model *model.ArtifactData `yaml:"model,omitempty"`
}

func (s *ArtifactSpec) ClearRef() {
	if s != nil {
		s.ref = nil
	}
}

func buildArtifactSpec(cfg *Config) (*ArtifactSpec, error) {
	refs := excel.Filter(excel.ReliquarySetExcelConfigData, func(v *excel.ReliquarySet) bool {
		if id := cfg.Override.Id; id != 0 {
			return v.SetId == id
		}
		affix := v.Affix(0)
		if affix == nil {
			return false
		}
		return excel.SlugLower(affix.Name()) == cfg.Name
	})
	if len(refs) != 1 {
		return nil, fmt.Errorf("query results in refs=%v but we expect 1", len(refs))
	}

	spec := &ArtifactSpec{ref: refs[0]}
	affixes := excel.Filter(excel.EquipAffixExcelConfigData, func(v *excel.EquipAffix) bool { return v.Id == spec.ref.EquipAffixId })
	codex := excel.Filter(excel.ReliquaryCodexExcelConfigData, func(v *excel.ReliquaryCodex) bool { return v.SuitId == spec.ref.SetId })
	if len(affixes) == 0 {
		return nil, fmt.Errorf("no affix found for reliquary_set=%v", spec.ref.SetId)
	}
	if len(codex) == 0 {
		return nil, fmt.Errorf("no codex found for reliquary_set=%v", spec.ref.SetId)
	}
	slices.SortFunc(affixes, func(a, b *excel.EquipAffix) int { return cmp.Compare(a.Level, b.Level) })
	slices.SortFunc(codex, func(a, b *excel.ReliquaryCodex) int { return cmp.Compare(a.Level, b.Level) })

	spec.Name = affixes[0].Name()
	spec.Model = &model.ArtifactData{
		SetId:      int64(spec.ref.SetId),
		TextMapId:  int64(affixes[len(affixes)-1].NameTextMapHash), // FIXME: this is unstable
		Key:        excel.SlugLower(spec.Name),
		ImageNames: &model.ArtifactImageData{},
	}

	for image, typ := range map[*string]excel.EquipType{
		&spec.Model.ImageNames.Flower:  excel.EQUIP_BRACER,
		&spec.Model.ImageNames.Plume:   excel.EQUIP_NECKLACE,
		&spec.Model.ImageNames.Sands:   excel.EQUIP_SHOES,
		&spec.Model.ImageNames.Goblet:  excel.EQUIP_RING,
		&spec.Model.ImageNames.Circlet: excel.EQUIP_DRESS,
	} {
		if r := codex[0].Reliquary(typ); r != nil {
			*image = r.Icon
		}
	}

	for _, affix := range affixes {
		needNum := spec.ref.SetNeedNum[affix.Level]
		attr := &AttributeSpec{
			Type:   fmt.Sprintf("set%d", needNum),
			Name:   affix.NameTextMapHash.String(),
			Desc:   affix.DescTextMapHash.String(),
			Config: affix.OpenConfig,
			Index:  NewTalentIndex(affix.ParamList),
		}
		attr.SetValues(1, func(int) []float64 { return affix.ParamList })
		cfg.Attributes = append(cfg.Attributes, attr)
	}

	return spec, nil
}

func (c *Compiled) GenerateArtifacts() error {
	kind := KindArtifact
	inputs := excel.Filter(c.Configuration, func(v *Config) bool { return v.Kind == kind })

	imports := ImportTmpl{Kind: kind}
	keys := KeysTmpl{Kind: kind}
	assets := AssetsTmpl{Kind: kind, Variable: "artifactMap"}
	shortcut := ShortcutTmpl{Kind: kind, Variable: "SetNameToKey", Type: "keys.Set"}
	catalog := CatalogTmpl{Kind: kind, Variable: "ArtifactMap", Type: "keys.Set", ModelName: reflect.TypeFor[*model.ArtifactData]().String()}
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

	models := &model.ArtifactDataMap{Data: make(map[string]*model.ArtifactData)}
	for _, config := range inputs {
		spec := config.Artifact

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

		add := func(slot, image string) {
			assets.Key = append(assets.Key, strings.Join([]string{spec.Model.Key, slot}, "_"))
			assets.Image = append(assets.Image, image)
		}
		add("flower", spec.Model.ImageNames.Flower)
		add("plume", spec.Model.ImageNames.Plume)
		add("sands", spec.Model.ImageNames.Sands)
		add("goblet", spec.Model.ImageNames.Goblet)
		add("circlet", spec.Model.ImageNames.Circlet)

		catalog.Slug = append(catalog.Slug, excel.Slug(spec.Name))
		catalog.Model = append(catalog.Model, proto.Clone(spec.Model))

		m := proto.Clone(spec.Model).(*model.ArtifactData)
		m.ImageNames = nil
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
		fmt.Fprintf(b, "core.Register%[2]sFunc(keys.%[1]s, New%[2]s)\n", excel.Slug(spec.Name), "Set")
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
