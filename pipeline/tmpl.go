package main

import (
	"bytes"
	"fmt"
	"maps"
	"path"
	"slices"
	"strconv"

	"github.com/shizukayuki/excel-hk4e"
	"google.golang.org/protobuf/proto"
)

const baseModule = "github.com/genshinsim/gcsim"

//nolint:goconst // using a constant is not that useful here
var abilities = []string{
	"default",
	"attack", "aim", "charge",
	"plunge", "low_plunge", "high_plunge",
	"skill", "burst",
	"asc", "a1", "a4",
	"cons",
	"dash", "jump", "walk", "swap",
	"c1", "c2", "c3", "c4", "c5", "c6",
	"set1", "set2", "set4",
	"passive",
}

type ImportTmpl struct {
	Kind        Kind
	PackagePath []string
}

func (t *ImportTmpl) Write() {
	b := bytes.NewBuffer(nil)
	b.WriteString("package simulation\n")
	b.WriteString("import (\n")
	for _, pkg := range t.PackagePath {
		b.WriteString("\t_ \"")
		b.WriteString(path.Join(baseModule, pkg))
		b.WriteString("\"\n")
	}
	b.WriteString(")\n")
	writeFile(fmt.Sprintf("pkg/simulation/imports.%s.dm.go", t.Kind), b.Bytes())
}

type KeysTmpl struct {
	Kind Kind
	Name []string
}

func (t *KeysTmpl) Write() {
	typ := map[Kind]string{
		KindArtifact:  "Set",
		KindCharacter: "Char",
		KindWeapon:    "Weapon",
	}[t.Kind]
	if typ == "" {
		panic("unreachable")
	}
	end := "Invalid" + typ

	b := bytes.NewBuffer(nil)
	b.WriteString("package keys\n")
	fmt.Fprintf(b, "//go:generate go tool github.com/dmarkham/enumer -text -json -linecomment -type=%s -output %[2]s.enumer.dm.go -- %[2]s.dm.go\n", typ, t.Kind)
	fmt.Fprintf(b, "type %s int\n", typ)
	b.WriteString("const (\n")
	fmt.Fprintf(b, "\tNo%[1]s %[1]s = iota //\n", typ)
	for _, name := range t.Name {
		fmt.Fprintf(b, "\t%s // %s\n", excel.Slug(name), excel.SlugLower(name))
	}
	fmt.Fprintf(b, "\t%s // %s\n", excel.Slug(end), excel.SlugLower(end))
	b.WriteString(")\n")
	writeFile(fmt.Sprintf("pkg/core/keys/%s.dm.go", t.Kind), b.Bytes())
}

type ShortcutTmpl struct {
	Kind     Kind
	Variable string
	Type     string
	Slug     []string
	Names    [][]string
}

func (t *ShortcutTmpl) Write() {
	b := bytes.NewBuffer(nil)
	b.WriteString("package shortcut\n")
	if t.Kind != KindMonster {
		b.WriteString("import (\n")
		fmt.Fprintf(b, "\t\"%s\"\n", path.Join(baseModule, "pkg/core/keys"))
		b.WriteString(")\n")
	}
	fmt.Fprintf(b, "var %s = map[string]%s{\n", t.Variable, t.Type)
	for i := range t.Slug {
		add := func(shortcut string) {
			b.WriteString("\t")
			fmt.Fprintf(b, "\"%s\":", shortcut)
			if t.Kind != KindMonster {
				b.WriteString("keys.")
			}
			b.WriteString(t.Slug[i])
			b.WriteString(",\n")
		}
		for _, s := range t.Names[i] {
			add(s)
		}
	}
	b.WriteString("}\n")
	writeFile(fmt.Sprintf("pkg/shortcut/%s.dm.go", t.Kind), b.Bytes())

	input := make(map[string][]string)
	for _, alt := range t.Names {
		if len(alt) > 1 {
			input[alt[0]] = alt[1:]
		}
	}
	data, err := dumpJSON(input)
	assert(err)
	writeFile(fmt.Sprintf("ui/packages/docs/src/components/Names/%s.dm.json", t.Kind), data)
}

type AssetsTmpl struct {
	Kind     Kind
	Variable string
	Key      []string
	Image    []string
}

func (t *AssetsTmpl) Write() {
	b := bytes.NewBuffer(nil)
	b.WriteString("package assets\n")
	fmt.Fprintf(b, "var %s = map[string]string{\n", t.Variable)
	for i := range t.Key {
		b.WriteString("\t")
		fmt.Fprintf(b, "\"%s\":", t.Key[i])
		fmt.Fprintf(b, "\"%s\",", t.Image[i])
		b.WriteString("\n")
	}
	b.WriteString("}\n")
	writeFile(fmt.Sprintf("internal/services/assets/%s.dm.go", t.Kind), b.Bytes())
}

type CatalogTmpl struct {
	Kind      Kind
	Variable  string
	Type      string
	ModelName string
	Slug      []string
	Model     []proto.Message
}

func (t *CatalogTmpl) Write() {
	b := bytes.NewBuffer(nil)
	b.WriteString("package catalog\n")
	b.WriteString("import (\n")
	fmt.Fprintf(b, "\t\"%s\"\n", path.Join(baseModule, "pkg/model"))
	if t.Kind != KindMonster {
		fmt.Fprintf(b, "\t\"%s\"\n", path.Join(baseModule, "pkg/core/keys"))
	}
	b.WriteString(")\n")
	fmt.Fprintf(b, "var %s = map[%s]%s{\n", t.Variable, t.Type, t.ModelName)
	for i := range t.Slug {
		if t.Kind != KindMonster {
			b.WriteString("keys.")
		}
		b.WriteString(t.Slug[i])
		b.WriteString(":")
		b.WriteString(dumpGo(t.Model[i], false))
		b.WriteString(",\n")
	}
	b.WriteString("}\n")
	writeFile(fmt.Sprintf("pkg/catalog/%s.dm.go", t.Kind), b.Bytes())
}

type DocTmpl struct {
	Kind Kind
	Key  []string
	Name []string
}

func (t *DocTmpl) Write() {
	dir := map[Kind]string{
		KindArtifact:  "artifacts",
		KindCharacter: "characters",
		KindMonster:   "enemies",
		KindWeapon:    "weapons",
	}[t.Kind]
	if dir == "" {
		panic("unreachable")
	}
	for i := range t.Key {
		input := struct{ Key, Name, Type string }{
			Key:  t.Key[i],
			Name: strconv.Quote(t.Name[i]),
			Type: t.Kind.String(),
		}
		data := useTemplate(fmt.Sprintf("docs_%s.md.templ", t.Kind), input)
		writeFile(fmt.Sprintf("ui/packages/docs/docs/reference/%s/%s.md", dir, input.Key), data)
	}
}

type IssuesTmpl struct {
	Kind Kind
	Key  []string
	Data [][]string
}

func (t *IssuesTmpl) Write() {
	input := make(map[string][]string)
	for i, name := range t.Key {
		if d := t.Data[i]; len(d) > 0 {
			input[name] = d
		}
	}
	data, err := dumpJSON(input)
	assert(err)
	writeFile(fmt.Sprintf("ui/packages/docs/src/components/Issues/%s.dm.json", t.Kind), data)
}

type FieldsTmpl struct {
	Kind Kind
	Key  []string
	Data []map[string]string
}

func (t *FieldsTmpl) Write() {
	type Input struct {
		Field string `json:"field"`
		Desc  string `json:"desc"`
	}
	input := make(map[string][]Input)
	for i, name := range t.Key {
		for _, field := range slices.Sorted(maps.Keys(t.Data[i])) {
			input[name] = append(input[name], Input{
				Field: field,
				Desc:  t.Data[i][field],
			})
		}
	}
	data, err := dumpJSON(input)
	assert(err)
	writeFile(fmt.Sprintf("ui/packages/docs/src/components/Fields/%s.dm.json", t.Kind), data)
}

type FramesTmpl struct {
	Kind Kind
	Key  []string
	Data [][]Frames
}

func (t *FramesTmpl) Write() {
	type Input struct {
		Count       string `json:"count"`
		CountCredit string `json:"count_credit"`
		Video       string `json:"video"`
		VideoCredit string `json:"video_credit"`
	}
	input := make(map[string][]Input)
	for i, name := range t.Key {
		for _, v := range t.Data[i] {
			input[name] = append(input[name], Input{
				Count:       v.Count,
				CountCredit: v.CountCredit,
				Video:       v.Video,
				VideoCredit: v.VideoCredit,
			})
		}
	}
	data, err := dumpJSON(input)
	assert(err)
	writeFile(fmt.Sprintf("ui/packages/docs/src/components/Frames/%s.dm.json", t.Kind), data)
}

type ActionsTmpl struct {
	Kind Kind
	Key  []string
	Data [][]*Ability
}

func (t *ActionsTmpl) Write() {
	type Input struct {
		Ability string `json:"ability"`
		Note    string `json:"note,omitempty"`
		Invalid bool   `json:"invalid"`
	}
	input := make(map[string][]Input)
	actions := []string{
		"attack", "aim", "charge",
		"low_plunge", "high_plunge",
		"skill", "burst",
		"dash", "jump", "walk", "swap",
	}
	hasDefaultImpl := []string{
		"dash", "jump", "walk", "swap",
	}
	for i, name := range t.Key {
		for _, action := range actions {
			v := Input{Ability: action}
			if ind := slices.IndexFunc(t.Data[i], func(v *Ability) bool { return v.Name == action }); ind != -1 {
				v.Note = t.Data[i][ind].Note
			} else if !slices.Contains(hasDefaultImpl, action) {
				v.Invalid = true
			}
			if v.Ability == "default" {
				v.Ability = "-"
			}
			input[name] = append(input[name], v)
		}
	}
	data, err := dumpJSON(input)
	assert(err)
	writeFile(fmt.Sprintf("ui/packages/docs/src/components/Actions/%s.dm.json", t.Kind), data)
}

type HitboxTmpl struct {
	Kind Kind
	Key  []string
	Data [][]*Ability
}

func (t *HitboxTmpl) Write() {
	type Input struct {
		Ability  string  `json:"ability"`
		Name     string  `json:"name"`
		Note     string  `json:"note,omitempty"`
		Shape    string  `json:"shape"`
		Center   string  `json:"center"`
		OffsetX  float64 `json:"offset_x,omitempty"`
		OffsetY  float64 `json:"offset_y,omitempty"`
		BoxX     float64 `json:"box_x,omitempty"`
		BoxY     float64 `json:"box_y,omitempty"`
		FanAngle float64 `json:"fan_angle,omitempty"`
		Radius   float64 `json:"radius,omitempty"`
	}
	input := make(map[string][]Input)
	for i, name := range t.Key {
		for _, abil := range t.Data[i] {
			for i := range abil.Hitbox {
				h := &abil.Hitbox[i]
				v := Input{
					Ability:  abil.Name,
					Name:     h.Name,
					Note:     h.Note,
					Shape:    h.Shape,
					Center:   h.Center,
					FanAngle: h.FanAngle,
					Radius:   h.Radius,
				}
				if v.Ability == "default" {
					v.Ability = "-"
				}
				if len(h.Offset) == 2 {
					v.OffsetX = h.Offset[0]
					v.OffsetY = h.Offset[1]
				}
				if len(h.Box) == 2 {
					v.BoxX = h.Box[0]
					v.BoxY = h.Box[1]
				}
				input[name] = append(input[name], v)
			}
		}
	}
	data, err := dumpJSON(input)
	assert(err)
	writeFile(fmt.Sprintf("ui/packages/docs/src/components/AoE/%s.dm.json", t.Kind), data)
}

type HitlagTmpl struct {
	Kind Kind
	Key  []string
	Data [][]*Ability
}

func (t *HitlagTmpl) Write() {
	type Input struct {
		Ability     string  `json:"ability"`
		Name        string  `json:"name"`
		Time        float64 `json:"time"`
		Scale       float64 `json:"scale"`
		DefenseHalt bool    `json:"defense_halt"`
		Deployable  bool    `json:"deployable"`
	}
	input := make(map[string][]Input)
	for i, name := range t.Key {
		for _, abil := range t.Data[i] {
			for i := range abil.Hitbox {
				h := abil.Hitbox[i].Hitlag
				if h.IsZero() {
					continue
				}
				v := Input{
					Ability:    abil.Name,
					Name:       abil.Hitbox[i].Name,
					Deployable: h.Deployable,
				}
				if v.Ability == "default" {
					v.Ability = "-"
				}
				if h.HasHitlag() {
					v.Time = h.Time
					v.Scale = h.Scale
					v.DefenseHalt = h.DefenseHalt
				}
				input[name] = append(input[name], v)
			}
		}
	}
	data, err := dumpJSON(input)
	assert(err)
	writeFile(fmt.Sprintf("ui/packages/docs/src/components/Hitlag/%s.dm.json", t.Kind), data)
}

type ParamsTmpl struct {
	Kind Kind
	Key  []string
	Data [][]*Ability
}

func (t *ParamsTmpl) Write() {
	type Input struct {
		Ability string `json:"ability"`
		Param   string `json:"param"`
		Desc    string `json:"desc"`
	}
	input := make(map[string][]Input)
	for i, name := range t.Key {
		for _, abil := range t.Data[i] {
			for _, param := range slices.Sorted(maps.Keys(abil.Params)) {
				v := Input{
					Ability: abil.Name,
					Param:   param,
					Desc:    abil.Params[param],
				}
				if v.Ability == "default" {
					v.Ability = "-"
				}
				input[name] = append(input[name], v)
			}
		}
	}
	data, err := dumpJSON(input)
	assert(err)
	writeFile(fmt.Sprintf("ui/packages/docs/src/components/Params/%s.dm.json", t.Kind), data)
}
