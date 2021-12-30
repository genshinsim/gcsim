package main

import (
	"encoding/json"
	"log"
	"os"
	"text/template"
)

type stats struct {
	Max         int     `json:"maxlevel"`
	HP          float64 `json:"hp"`
	Atk         float64 `json:"attack"`
	Def         float64 `json:"defense"`
	Specialized float64 `json:"specialized"`
}

type curve struct {
	Atk         string `json:"attack"`
	Specialized string `json:"specialized"`
}

type data struct {
	Base          stats   `json:"base"`
	Curve         curve   `json:"curve"`
	Specialized   string  `json:"specialized"`
	PromotionData []stats `json:"promotion"`
}

func main() {

	f, err := os.ReadFile("./weapons.json")
	if err != nil {
		log.Panic(err)
	}
	var d map[string]data
	err = json.Unmarshal(f, &d)
	if err != nil {
		log.Panic(err)
	}
	//fix the specialized key
	for k, v := range d {
		v.Specialized = SpecKeyToStat[v.Specialized]
		if v.Specialized == "" {
			v.Specialized = "core.NoStat"
		}
		d[k] = v
		// fmt.Println(k)
	}
	// fmt.Println(d)
	t, err := template.New("out").Parse(tmpl)
	if err != nil {
		log.Panic(err)
	}
	os.Remove("./weap.txt")
	of, err := os.Create("./weap.txt")
	if err != nil {
		log.Panic(err)
	}
	err = t.Execute(of, d)
	if err != nil {
		log.Panic(err)
	}
}

var tmpl = `package curves

import (
	"github.com/genshinsim/gcsim/pkg/core"
)


var WeaponBaseMap = map[string]WeaponBase{
	{{- range $key, $value := . }}
	"{{$key}}": {
		AtkCurve: {{$value.Curve.Atk}},
		SpecializedCurve: {{$value.Curve.Specialized}},
		BaseAtk: {{$value.Base.Atk}},
		BaseSpecialized: {{$value.Base.Specialized}},
		Specialized: {{$value.Specialized}},
		PromotionBonus: []PromoData{
			{{- range $e := $value.PromotionData}}
			{
				MaxLevel: {{$e.Max}},
				Atk:      {{$e.Atk}},
			},
			{{- end }}
		},
	},
	{{- end }}
}

`

var SpecKeyToStat = map[string]string{
	"FIGHT_PROP_CRITICAL_HURT":     "core.CD",
	"FIGHT_PROP_HEAL_ADD":          "core.Heal",
	"FIGHT_PROP_ATTACK_PERCENT":    "core.ATKP",
	"FIGHT_PROP_ELEMENT_MASTERY":   "core.EM",
	"FIGHT_PROP_HP_PERCENT":        "core.HPP",
	"FIGHT_PROP_CHARGE_EFFICIENCY": "core.ER",
	"FIGHT_PROP_CRITICAL":          "core.CR",
	"FIGHT_PROP_PHYSICAL_ADD_HURT": "core.PhyP",
	"FIGHT_PROP_ELEC_ADD_HURT":     "core.EleP",
	"FIGHT_PROP_ROCK_ADD_HURT":     "core.GeoP",
	"FIGHT_PROP_FIRE_ADD_HURT":     "core.PyroP",
	"FIGHT_PROP_WATER_ADD_HURT":    "core.HydroP",
	"FIGHT_PROP_DEFENSE_PERCENT":   "core.DEFP",
	"FIGHT_PROP_ICE_ADD_HURT":      "core.CryoP",
	"FIGHT_PROP_WIND_ADD_HURT":     "core.AnemoP",
}
