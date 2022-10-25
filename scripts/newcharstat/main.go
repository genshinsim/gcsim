package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
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
	Critical        float64 `json:"critical"`
	CriticalHurt    float64 `json:"criticalHurt"`
	PropGrowCurves  []struct {
		Type      string `json:"type"`
		GrowCurve string `json:"growCurve"`
	} `json:"propGrowCurves"`
	ID              int    `json:"id"`
	NameTextMapHash int    `json:"nameTextMapHash"`
	UseType         string `json:"useType,omitempty"`
}
type AvatarPromotes []struct {
	AvatarPromoteID int `json:"avatarPromoteId"`
	AddProps        []struct {
		PropType string  `json:"propType"`
		Value    float64 `json:"value"`
	} `json:"addProps"`
	PromoteLevel int `json:"promoteLevel,omitempty"`
	// PromoteAudio    string `json:"promoteAudio"`
	// CostItems       []struct {
	// } `json:"costItems"`
	// UnlockMaxLevel int `json:"unlockMaxLevel"`
	// ScoinCost           int `json:"scoinCost,omitempty"`
	// RequiredPlayerLevel int `json:"requiredPlayerLevel,omitempty"`
}

type FetterInfo []struct {
	AvatarAssocType string `json:"avatarAssocType"`
	AvatarId        int    `json:"avatarId"`
}

type CharStatCurve int

type CharBase struct {
	Rarity     int                `json:"rarity"`
	Body       profile.BodyType   `json:"-"`
	Element    attributes.Element `json:"element"`
	Region     profile.ZoneType   `json:"-"`
	WeaponType weapon.WeaponClass `json:"weapon_class"`

	HPCurve        CharStatCurve   `json:"-"`
	AtkCurve       CharStatCurve   `json:"-"`
	DefCurve       CharStatCurve   `json:"-"`
	BaseHP         float64         `json:"-"`
	BaseAtk        float64         `json:"-"`
	BaseDef        float64         `json:"-"`
	Specialized    attributes.Stat `json:"-"`
	PromotionBonus []PromoData     `json:"-"`
}

type PromoData struct {
	MaxLevel int
	HP       float64
	Atk      float64
	Def      float64
	Special  float64
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

	characterArray := make([]CharBase, len(avatars))

	for _, avatar := range avatars {
		if avatar.UseType != "AVATAR_FORMAL" {
			continue
		}
		char := CharBase{}

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

		// characterArray[index].Region = profile.ZoneType(avatar.AvatarRegion)
		// characterArray[index].Element = attributes.Element(avatar.AvatarElementalType)
		char.HPCurve = CharStatCurve(avatar.HpBase)
		char.AtkCurve = CharStatCurve(avatar.AttackBase)
		char.DefCurve = CharStatCurve(avatar.DefenseBase)
		char.BaseHP = avatar.HpBase
		char.BaseAtk = avatar.AttackBase
		char.BaseDef = avatar.DefenseBase
		// char.Specialized = attributes.Stat(avatar.Critical)

		characterArray = append(characterArray, char)
		//print out the character name
		fmt.Println(strings.Replace(avatar.IconName, "UI_AvatarIcon_", "", 1))
		fmt.Printf("%+v\n", char)

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
