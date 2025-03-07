package character

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/model"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

type charData struct {
	Config
	Data         *model.AvatarData
	NASlice      map[string]*naSlice
	SkillLvlData []skillLvlData
	ParamKeys    map[int][]paramData
}

type naSlice struct {
	Names []string
	Count []int
	Is3D  bool
}

type skillLvlData struct {
	Name    string
	Params  []skillParam
	Comment string
}

type skillParam struct {
	Values  []float64
	Comment string
}

func (g *Generator) GenerateCharTemplate() error {
	t, err := template.New("chartemplate").Parse(charTmpl)
	if err != nil {
		return fmt.Errorf("failed to build template: %w", err)
	}
	for i := range g.chars {
		v := g.chars[i]
		dm, ok := g.data[v.Key]
		if !ok {
			log.Printf("No data found for %v; skipping", v.Key)
			continue
		}
		// get rid of unnecessary fields
		dmCopy := proto.Clone(dm).(*model.AvatarData)
		dmCopy.SkillDetails.A1 = 0
		dmCopy.SkillDetails.A4 = 0
		dmCopy.SkillDetails.A1Scaling = nil
		dmCopy.SkillDetails.A4Scaling = nil
		err = writePBToFile(fmt.Sprintf("%v/data_gen.textproto", v.RelativePath), dmCopy)
		if err != nil {
			return fmt.Errorf("failed to write PB file: %w", err)
		}

		buff := new(bytes.Buffer)
		d := charData{
			Config: v,
			Data:   dm,
		}
		if d.CharStructName == "" {
			d.CharStructName = "char"
		}
		if d.ActionParamKeys != nil {
			err := d.buildValidation()
			if err != nil {
				return fmt.Errorf("failed to build validation map for %v: %w", v.Key, err)
			}
		}
		err := d.buildSkillData(dm)
		if err != nil {
			return fmt.Errorf("%v: %w", v.Key, err)
		}
		err = t.Execute(buff, d)
		if err != nil {
			return fmt.Errorf("failed to execute template for %v: %w", v.Key, err)
		}
		src := buff.Bytes()
		dst, err := format.Source(src)
		if err != nil {
			fmt.Println(string(src))
			return fmt.Errorf("failed to gofmt on %v: %w", v.RelativePath, err)
		}
		os.WriteFile(fmt.Sprintf("%v/%v_gen.go", v.RelativePath, v.PackageName), dst, 0o644)
	}

	return nil
}

func writePBToFile(path string, dm *model.AvatarData) error {
	// get rid of unncessary fields before saving
	msg := proto.Clone(dm).(*model.AvatarData)
	msg.SkillDetails.AttackScaling = nil
	msg.SkillDetails.SkillScaling = nil
	msg.SkillDetails.BurstScaling = nil
	// write the avatar data to []byte
	opt := prototext.MarshalOptions{
		Multiline: true,
		Indent:    "   ",
	}
	b, err := opt.Marshal(msg)
	if err != nil {
		log.Printf("error marshalling %v data to proto\n", dm.Key)
		return err
	}
	// hack to work around stupid prototext not stable (on purpose - google u suck)
	b = []byte(strings.ReplaceAll(string(b), ":  ", ": "))
	return os.WriteFile(path, b, 0o644)
}

func (c *charData) buildValidation() error {
	if c.KeyVarName == "" {
		c.KeyVarName = cases.Title(language.AmericanEnglish).String(c.Key)
	}
	c.ParamKeys = make(map[int][]paramData)
	for k, v := range c.ActionParamKeys {
		a := action.StringToAction(k)
		if a == action.InvalidAction {
			return fmt.Errorf("invalid action string: %v", k)
		}
		c.ParamKeys[int(a)] = v
	}

	return nil
}

func (c *charData) buildSkillData(dm *model.AvatarData) error {
	atk, err := c.skillDataByType("attack", dm.SkillDetails.AttackScaling)
	if err != nil {
		return err
	}
	skill, err := c.skillDataByType("skill", dm.SkillDetails.SkillScaling)
	if err != nil {
		return err
	}
	burst, err := c.skillDataByType("burst", dm.SkillDetails.BurstScaling)
	if err != nil {
		return err
	}
	a1, err := c.skillDataByType("a1", dm.SkillDetails.A1Scaling)
	if err != nil {
		return err
	}
	a4, err := c.skillDataByType("a4", dm.SkillDetails.A4Scaling)
	if err != nil {
		return err
	}
	c.SkillLvlData = append(c.SkillLvlData, atk...)
	c.SkillLvlData = append(c.SkillLvlData, skill...)
	c.SkillLvlData = append(c.SkillLvlData, burst...)
	c.SkillLvlData = append(c.SkillLvlData, a1...)
	c.SkillLvlData = append(c.SkillLvlData, a4...)

	// merge NA data into a single slice
	for _, v := range c.SkillLvlData {
		if !strings.Contains(v.Name, "_") {
			continue
		}

		base := strings.Split(v.Name, "_")[0]
		if c.NASlice == nil {
			c.NASlice = make(map[string]*naSlice)
		}
		if c.NASlice[base] == nil {
			c.NASlice[base] = &naSlice{}
		}

		slice := c.NASlice[base]
		slice.Names = append(slice.Names, v.Name)
		slice.Count = append(slice.Count, len(v.Params))
		if len(v.Params) > 1 {
			slice.Is3D = true
		}
	}

	return nil
}

func (c *charData) skillDataByType(typ string, data []*model.AvatarSkillExcelIndexData) ([]skillLvlData, error) {
	var result []skillLvlData
	for name, params := range c.Config.SkillDataMapping[typ] {
		skill := skillLvlData{
			Name:    name,
			Params:  make([]skillParam, len(params)),
			Comment: fmt.Sprintf("%s: %s = %v", typ, name, params),
		}

		for i, param := range params {
			if param == -1 {
				skill.Params[i].Values = make([]float64, 15)
				continue
			}

			var paramData *model.AvatarSkillExcelIndexData
			for _, v := range data {
				if v.Index != int32(param) {
					continue
				}
				paramData = v
				break
			}

			if paramData == nil {
				return nil, fmt.Errorf("could not find param data for: %s %s %d", typ, name, param)
			}

			// skill.Params[i].Comment = "<paramDesc>"
			for _, ld := range paramData.LevelData {
				skill.Params[i].Values = append(skill.Params[i].Values, ld.Value)
			}
		}

		result = append(result, skill)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result, nil
}

const charTmpl = `// Code generated by "pipeline"; DO NOT EDIT.
package {{.PackageName}}

import (
	_ "embed"

	"github.com/genshinsim/gcsim/pkg/model"
	"google.golang.org/protobuf/encoding/prototext"
	{{if ne (len .ParamKeys) 0 -}}
	"slices"
	"fmt"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/gcs/validation"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	{{- end }}
)

//go:embed data_gen.textproto
var pbData []byte
var base *model.AvatarData
{{if ne (len .ParamKeys) 0 -}}
var paramKeysValidation = map[action.Action][]string {
	{{- range $key, $slice := .ParamKeys}}
	{{$key}}: {
		{{- range $val := $slice -}}
		"{{$val.Param}}",
		{{- end -}}
	},
	{{- end}}
}
{{- end}}

func init() {
	base = &model.AvatarData{}
	err := prototext.Unmarshal(pbData, base)
	if err != nil {
		panic(err)
	}
	{{- if ne (len .ParamKeys) 0}}
	validation.RegisterCharParamValidationFunc(keys.{{.KeyVarName}}, ValidateParamKeys){{end}}
}

{{if ne (len .ParamKeys) 0 -}}
func ValidateParamKeys(a action.Action, keys []string) error {
	valid, ok := paramKeysValidation[a]
	if !ok {
		return nil
	}
	for _, v := range keys {
		if !slices.Contains(valid, v) {
			return fmt.Errorf("key %v is invalid for action %v", v, a.String())
		}
	}
	return nil
}
{{- end}}

func (x *{{.CharStructName}}) Data() *model.AvatarData {
	return base
}
{{- if ne (len .NASlice) 0 }}

var (
{{- range $key, $slice := .NASlice }}
	{{ $key }} = [][] {{- if $slice.Is3D -}} [] {{- end -}} float64{
{{- range $i, $name := $slice.Names }}
	{{- if or (gt (index $slice.Count $i) 1) (not $slice.Is3D) }}
		{{ $name }},
	{{- else }}
		{ {{- $name -}} },
	{{- end }}
{{- end }}
	}
{{- end }}
)

{{- end }}
{{- if .SkillLvlData }}

{{- $single := false}}

var (
{{- range $skill := .SkillLvlData }}
{{- if eq (len $skill.Params) 1 }}
	{{- $param := index $skill.Params 0 }}
	{{- if gt (len $param.Values) 1}}
	// {{ $skill.Comment }}
	{{ $skill.Name }} = []float64{
		{{- if ne $param.Comment "" }}
		// {{ $param.Comment }}
		{{- end }}
	{{- range $v := $param.Values }}
		{{ $v }},
	{{- end }}
	}
	{{- else }}{{$single = true}}
	{{- end  }}
{{- else }}
	// {{ $skill.Comment }}
	{{ $skill.Name }} = [][]float64{
	{{- range $param := $skill.Params }}
		{ {{- if ne $param.Comment "" -}} // {{ $param.Comment }} {{- end }}
		{{- range $v := $param.Values }}
			{{ $v }},
		{{- end }}
		},
	{{- end }}
	}
{{- end }}
{{- end }}
)

{{- if eq $single true}}
const (
	{{- range $skill := .SkillLvlData }}
	{{- if eq (len $skill.Params) 1 }}
		{{- $param := index $skill.Params 0 }}
		{{- if eq (len $param.Values) 1}}
		// {{ $skill.Comment }}
		{{ $skill.Name }} float64 = {{index $param.Values 0}}
		{{- end  }}
	{{- end }}
	{{- end }}
)
{{- end }}

{{- end }}
`
