package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/curves"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)

type Avatars []struct {
	// ScriptDataPathHashSuffix     int           `json:"scriptDataPathHashSuffix"`
	// ScriptDataPathHashPre        int           `json:"scriptDataPathHashPre"`
	// SideIconName                 string        `json:"sideIconName"`
	// ChargeEfficiency             float64       `json:"chargeEfficiency"`
	// CombatConfigHashSuffix       int           `json:"combatConfigHashSuffix"`
	// CombatConfigHashPre          int           `json:"combatConfigHashPre"`
	// ManekinPathHashSuffix        int           `json:"manekinPathHashSuffix"`
	// ManekinPathHashPre           int           `json:"manekinPathHashPre"`
	// Mjgngjhbagi                  int64         `json:"MJGNGJHBAGI"`
	// GachaCardNameHashPre         int           `json:"gachaCardNameHashPre"`
	// Cbfoekcenea                  int64         `json:"CBFOEKCENEA"`
	// Pagadeakhac                  int           `json:"PAGADEAKHAC"`
	// CutsceneShow                 string        `json:"cutsceneShow"`
	// SkillDepotID                 int           `json:"skillDepotId"`
	// StaminaRecoverSpeed          float64       `json:"staminaRecoverSpeed"`
	// CandSkillDepotIds            []interface{} `json:"candSkillDepotIds"`
	// ManekinJSONConfigHashSuffix  int           `json:"manekinJsonConfigHashSuffix"`
	// ManekinJSONConfigHashPre     int           `json:"manekinJsonConfigHashPre"`
	// ManekinMotionConfig          int           `json:"manekinMotionConfig"`
	// DescTextMapHash              int           `json:"descTextMapHash"`
	// AvatarIdentityType string `json:"avatarIdentityType"`
	// AvatarPromoteRewardLevelList []int         `json:"avatarPromoteRewardLevelList"`
	// AvatarPromoteRewardIDList    []int         `json:"avatarPromoteRewardIdList"`
	// FeatureTagGroupID            int           `json:"featureTagGroupID"`
	// InfoDescTextMapHash          int           `json:"infoDescTextMapHash"`
	// PrefabPathRagdollHashSuffix    int    `json:"prefabPathRagdollHashSuffix"`
	// PrefabPathRagdollHashPre       int    `json:"prefabPathRagdollHashPre"`
	// AnimatorConfigPathHashSuffix   int64  `json:"animatorConfigPathHashSuffix"`
	// Njiekklklld                    int    `json:"NJIEKKLKLLD"`
	// PrefabPathHashSuffix           int64  `json:"prefabPathHashSuffix"`
	// PrefabPathHashPre              int    `json:"prefabPathHashPre"`
	// PrefabPathRemoteHashSuffix     int64  `json:"prefabPathRemoteHashSuffix"`
	// PrefabPathRemoteHashPre        int    `json:"prefabPathRemoteHashPre"`
	// ControllerPathHashSuffix       int64  `json:"controllerPathHashSuffix"`
	// ControllerPathHashPre          int    `json:"controllerPathHashPre"`
	// ControllerPathRemoteHashSuffix int64  `json:"controllerPathRemoteHashSuffix"`
	// ControllerPathRemoteHashPre    int    `json:"controllerPathRemoteHashPre"`
	// LODPatternName                 string `json:"LODPatternName"`
	BodyType        string  `json:"bodyType"`
	IconName        string  `json:"iconName"`
	QualityType     string  `json:"qualityType"`
	InitialWeapon   int     `json:"initialWeapon"`
	WeaponType      string  `json:"weaponType"`
	ImageName       string  `json:"imageName"`
	AvatarPromoteID int     `json:"avatarPromoteId"`
	HpBase          float64 `json:"hpBase"`
	AttackBase      float64 `json:"attackBase"`
	DefenseBase     float64 `json:"defenseBase"`
	// Critical        float64 `json:"critical"`
	// CriticalHurt    float64 `json:"criticalHurt"`
	PropGrowCurves []struct {
		Type      string `json:"type"`
		GrowCurve string `json:"growCurve"`
	} `json:"propGrowCurves"`
	ID              int    `json:"id"`
	NameTextMapHash int    `json:"nameTextMapHash"`
	UseType         string `json:"useType,omitempty"`
}
type AvatarPromotes []struct {
	AvatarPromoteID int       `json:"avatarPromoteId"`
	AddProps        []AddProp `json:"addProps"`
	PromoteLevel    int       `json:"promoteLevel,omitempty"`
	// PromoteAudio    string `json:"promoteAudio"`
	// CostItems       []struct {
	// } `json:"costItems"`
	// UnlockMaxLevel int `json:"unlockMaxLevel"`
	// ScoinCost           int `json:"scoinCost,omitempty"`
	// RequiredPlayerLevel int `json:"requiredPlayerLevel,omitempty"`
}

type AddProp struct {
	PropType string  `json:"propType"`
	Value    float64 `json:"value"`
}

type FetterInfo []struct {
	AvatarAssocType string `json:"avatarAssocType"`
	AvatarId        int    `json:"avatarId"`
}

func main() {
	avatarData, err := fetchJsonFromUrl("https://raw.githubusercontent.com/Dimbreath/GenshinData/master/ExcelBinOutput/AvatarExcelConfigData.json")
	if err != nil {
		log.Fatal(err)
	}

	avatarPromoteData, err := fetchJsonFromUrl("https://raw.githubusercontent.com/Dimbreath/GenshinData/master/ExcelBinOutput/AvatarPromoteExcelConfigData.json")
	if err != nil {
		log.Fatal(err)
	}

	fetterInfoData, err := fetchJsonFromUrl("https://raw.githubusercontent.com/Dimbreath/GenshinData/master/ExcelBinOutput/FetterInfoExcelConfigData.json")
	if err != nil {
		log.Fatal(err)
	}
	var avatars Avatars
	if err := json.Unmarshal([]byte(avatarData), &avatars); err != nil {
		log.Fatal(err)
	}
	var avatarPromotes AvatarPromotes
	if err := json.Unmarshal([]byte(avatarPromoteData), &avatarPromotes); err != nil {
		log.Fatal(err)
	}
	var fetterInfo FetterInfo
	if err := json.Unmarshal([]byte(fetterInfoData), &fetterInfo); err != nil {
		log.Fatal(err)
	}
	// reshape fetterInfo to map of avatarId to AvatarAssocType
	locationMap := make(map[int]string)
	for _, v := range fetterInfo {
		locationMap[v.AvatarId] = v.AvatarAssocType
	}
	//reshape avatarPromotes and remove duplicates to map of avatarPromoteId to AddProps
	specializedStatMap := make(map[int]string)
	for _, v := range avatarPromotes {
		specializedStatMap[v.AvatarPromoteID] = v.AddProps[3].PropType
	}

	// fmt.Printf("%+v\n", specializedStatMap)

	// reshape avatarPromotes to map of avatarPromoteId to PromoData array of 7 elements
	promoDataMap := make(map[int][]curves.PromoData)
	for _, v := range avatarPromotes {
		partialPromoData := addPropArraytoPromoData(v.AddProps)

		partialPromoData.MaxLevel = promoteLevelToMaxLevel(v.PromoteLevel)
		promoDataMap[v.AvatarPromoteID] = append(promoDataMap[v.AvatarPromoteID], partialPromoData)
	}

	// fmt.Printf("%+v\n", promoDataMap)

	characterArray := make([]curves.CharBase, len(avatars))

	for _, avatar := range avatars {
		if avatar.UseType != "AVATAR_FORMAL" {
			continue
		}
		char := curves.CharBase{}

		switch avatar.QualityType {
		case "QUALITY_PURPLE":
			char.Rarity = 4
		case "QUALITY_ORANGE", "QUALITY_ORANGE_SP":
			char.Rarity = 5
		default:
			log.Fatal("Unknown rarity type for character ", strings.Replace(avatar.IconName, "UI_AvatarIcon_", "", 1), ": ", avatar.QualityType)
		}

		switch avatar.BodyType {
		case "BODY_LOLI":
			char.Body = profile.BodyLoli
		case "BODY_GIRL":
			char.Body = profile.BodyGirl
		case "BODY_LADY":
			char.Body = profile.BodyLady
		case "BODY_BOY":
			char.Body = profile.BodyBoy
		case "BODY_MALE":
			char.Body = profile.BodyMale
		default:
			log.Fatal("Unknown BodyType")
		}

		switch avatar.WeaponType {
		case "WEAPON_CLAYMORE":
			char.WeaponType = weapon.WeaponClassClaymore
		case "WEAPON_BOW":
			char.WeaponType = weapon.WeaponClassBow
		case "WEAPON_SWORD_ONE_HAND":
			char.WeaponType = weapon.WeaponClassSword
		case "WEAPON_CATALYST":
			char.WeaponType = weapon.WeaponClassCatalyst
		case "WEAPON_POLE":
			char.WeaponType = weapon.WeaponClassSpear
		default:
			log.Fatal("Unknown WeaponType")
		}

		switch locationMap[avatar.ID] {
		case "ASSOC_TYPE_INAZUMA":
			char.Region = profile.ZoneInazuma
		case "ASSOC_TYPE_LIYUE":
			char.Region = profile.ZoneInazuma
		case "ASSOC_TYPE_MONDSTADT":
			char.Region = profile.ZoneInazuma
		case "ASSOC_TYPE_SUMERU":
			char.Region = profile.ZoneSumeru
		default:
			char.Region = profile.ZoneUnknown
		}

		switch specializedStatMap[avatar.AvatarPromoteID] {
		case "FIGHT_PROP_CRITICAL_HURT":
			char.Specialized = attributes.CD
		case "FIGHT_PROP_HEAL_ADD":
			char.Specialized = attributes.Heal
		case "FIGHT_PROP_ATTACK_PERCENT":
			char.Specialized = attributes.ATKP
		case "FIGHT_PROP_ELEMENT_MASTERY":
			char.Specialized = attributes.EM
		case "FIGHT_PROP_HP_PERCENT":
			char.Specialized = attributes.HPP
		case "FIGHT_PROP_CHARGE_EFFICIENCY":
			char.Specialized = attributes.ER
		case "FIGHT_PROP_CRITICAL":
			char.Specialized = attributes.CR
		case "FIGHT_PROP_PHYSICAL_ADD_HURT":
			char.Specialized = attributes.PhyP
		case "FIGHT_PROP_ELEC_ADD_HURT":
			char.Specialized = attributes.ElectroP
		case "FIGHT_PROP_ROCK_ADD_HURT":
			char.Specialized = attributes.GeoP
		case "FIGHT_PROP_FIRE_ADD_HURT":
			char.Specialized = attributes.PyroP
		case "FIGHT_PROP_WATER_ADD_HURT":
			char.Specialized = attributes.HydroP
		case "FIGHT_PROP_DEFENSE_PERCENT":
			char.Specialized = attributes.DEFP
		case "FIGHT_PROP_ICE_ADD_HURT":
			char.Specialized = attributes.CryoP
		case "FIGHT_PROP_WIND_ADD_HURT":
			char.Specialized = attributes.AnemoP
		case "FIGHT_PROP_GRASS_ADD_HURT":
			char.Specialized = attributes.DendroP
		default:
			log.Fatal("Unknown Specialized Stat")
		}
		// characterArray[index].Element = attributes.Element(avatar.AvatarElementalType)
		if strings.Contains(avatar.PropGrowCurves[0].GrowCurve, "S5") {
			char.HPCurve = curves.GROW_CURVE_HP_S5
			char.AtkCurve = curves.GROW_CURVE_ATTACK_S5
			char.DefCurve = curves.GROW_CURVE_HP_S5
		} else {
			char.HPCurve = curves.GROW_CURVE_HP_S4
			char.AtkCurve = curves.GROW_CURVE_ATTACK_S4
			char.DefCurve = curves.GROW_CURVE_HP_S4
		}
		char.BaseHP = avatar.HpBase
		char.BaseAtk = avatar.AttackBase
		char.BaseDef = avatar.DefenseBase
		char.PromotionBonus = promoDataMap[avatar.AvatarPromoteID]

		characterArray = append(characterArray, char)
		//print out the character name
		fmt.Println(strings.Replace(avatar.IconName, "UI_AvatarIcon_", "", 1))
		fmt.Printf("%+v\n", char)
		// fmt.Printf("%+v\n", avatar.PropGrowCurves)

	}

	//print characterArray

}

func fetchJsonFromUrl(path string) (string, error) {

	resp, err := http.Get(path)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("%v: %v", resp.Status, path)
	}

	out, err := io.ReadAll(resp.Body)
	return string(out), err
}

var tmpl = `package curves

import (
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/keys"
	"github.com/genshinsim/gcsim/pkg/core/player/character/profile"
	"github.com/genshinsim/gcsim/pkg/core/player/weapon"
)


var CharBaseMap = map[keys.Char]CharBase{
	{{- range $key, $value := . }}
	{{- if $value.Key }}
	keys.{{$value.Key}}: {
		Rarity: {{$value.Rarity}},
		Body: profile.Body {{- $value.Body}},
		Element: attributes. {{- $value.Element}},
		Region: profile.Zone {{- $value.Region}},
		WeaponType: weapon.WeaponClass {{- $value.WeaponType}},
		HPCurve: {{$value.Curve.HP}},
		AtkCurve: {{$value.Curve.Atk}},
		DefCurve: {{$value.Curve.Def}},
		BaseHP: {{$value.Base.HP}},
		BaseAtk: {{$value.Base.Atk}},
		BaseDef: {{$value.Base.Def}},
		Specialized: {{$value.Specialized}},
		PromotionBonus: []PromoData{
			{{- range $e := $value.PromotionData}}
			{
				MaxLevel: {{$e.Max}},
				HP:       {{$e.HP}},
				Atk:      {{$e.Atk}},
				Def:      {{$e.Def}},
				Special:  {{$e.Specialized}},
			},
			{{- end }}
		},
	},
	{{- end }}
	{{- end }}
}

`

var SpecKeyToStat = map[string]string{
	"FIGHT_PROP_CRITICAL_HURT":     "attributes.CD",
	"FIGHT_PROP_HEAL_ADD":          "attributes.Heal",
	"FIGHT_PROP_ATTACK_PERCENT":    "attributes.ATKP",
	"FIGHT_PROP_ELEMENT_MASTERY":   "attributes.EM",
	"FIGHT_PROP_HP_PERCENT":        "attributes.HPP",
	"FIGHT_PROP_CHARGE_EFFICIENCY": "attributes.ER",
	"FIGHT_PROP_CRITICAL":          "attributes.CR",
	"FIGHT_PROP_PHYSICAL_ADD_HURT": "attributes.PhyP",
	"FIGHT_PROP_ELEC_ADD_HURT":     "attributes.ElectroP",
	"FIGHT_PROP_ROCK_ADD_HURT":     "attributes.GeoP",
	"FIGHT_PROP_FIRE_ADD_HURT":     "attributes.PyroP",
	"FIGHT_PROP_WATER_ADD_HURT":    "attributes.HydroP",
	"FIGHT_PROP_DEFENSE_PERCENT":   "attributes.DEFP",
	"FIGHT_PROP_ICE_ADD_HURT":      "attributes.CryoP",
	"FIGHT_PROP_WIND_ADD_HURT":     "attributes.AnemoP",
	"FIGHT_PROP_GRASS_ADD_HURT":    "attributes.DendroP",
}

func promoteLevelToMaxLevel(level int) int {
	switch level {
	case 0:
		return 20
	case 1:
		return 40
	case 2:
		return 50
	case 3:
		return 60
	case 4:
		return 70
	case 5:
		return 80
	case 6:
		return 90
	default:
		return 0
	}
}

func addPropArraytoPromoData(addProps []AddProp) curves.PromoData {
	var out curves.PromoData
	for _, prop := range addProps {

		switch prop.PropType {
		case "FIGHT_PROP_BASE_HP":
			out.HP = prop.Value
		case "FIGHT_PROP_BASE_ATTACK":
			out.Atk = prop.Value
		case "FIGHT_PROP_BASE_DEFENSE":
			out.Def = prop.Value
		default:
			out.Special = prop.Value
		}
	}
	return out
}
