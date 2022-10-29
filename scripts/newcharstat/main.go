package main

import (
	"errors"
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

func main() {
	var err error
	avatars := getAvatarArray()
	locationMap := getCharLocationMap()
	specializedStatMap, promoDataMap := getCharSpecializedStatandPromoData()

	characterArray := make([]curves.CharBase, len(avatars))

	for _, avatar := range avatars {
		char := curves.CharBase{}
		charName := determineCharName(avatar.IconName)
		char.Rarity, err = determineCharRarity(avatar.QualityType)
		if err != nil {
			log.Fatal("Unknown rarity type for character ", charName, ": ", avatar.QualityType)
		}

		char.Body, err = determineCharBody(avatar.BodyType)
		if err != nil {
			log.Fatal("Unknown body type for character ", charName, ": ", avatar.BodyType)
		}

		char.WeaponType, err = determineCharWeaponType(avatar.WeaponType)
		if err != nil {
			log.Fatal("Unknown weapon type for character ", charName, ": ", avatar.WeaponType)
		}

		char.Region, err = determineCharRegion(locationMap[avatar.ID])
		if err != nil {
			log.Fatal("Unknown region for character ", charName, ": ", locationMap[avatar.ID])
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
		// fmt.Println(strings.Replace(avatar.IconName, "UI_AvatarIcon_", "", 1))
		// fmt.Printf("%+v\n", char)
		// fmt.Printf("%+v\n", avatar.PropGrowCurves)
	}
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

func determineCharRarity(qualityType string) (int, error) {
	switch qualityType {
	case "QUALITY_PURPLE":
		return 4, nil
	case "QUALITY_ORANGE", "QUALITY_ORANGE_SP":
		return 5, nil
	default:
		return 0, errors.New("unknown quality type")
	}
}

func determineCharRegion(location string) (profile.ZoneType, error) {
	switch location {
	case "ASSOC_TYPE_INAZUMA":
		return profile.ZoneInazuma, nil
	case "ASSOC_TYPE_LIYUE":
		return profile.ZoneInazuma, nil
	case "ASSOC_TYPE_MONDSTADT":
		return profile.ZoneInazuma, nil
	case "ASSOC_TYPE_SUMERU":
		return profile.ZoneSumeru, nil
	case "ASSOC_TYPE_MAINACTOR", "ASSOC_TYPE_RANGER":
		return profile.ZoneUnknown, nil
	case "ASSOC_TYPE_FATUI":
		return profile.ZoneSnezhnaya, nil

	default:
		return profile.ZoneUnknown, errors.New("unknown location")
	}

}

func determineCharName(iconName string) string {
	return strings.Replace(iconName, "UI_AvatarIcon_", "", 1)
}

func determineCharBody(bodyType string) (profile.BodyType, error) {
	switch bodyType {
	case "BODY_LOLI":
		return profile.BodyLoli, nil
	case "BODY_GIRL":
		return profile.BodyGirl, nil
	case "BODY_LADY":
		return profile.BodyLady, nil
	case "BODY_BOY":
		return profile.BodyBoy, nil
	case "BODY_MALE":
		return profile.BodyMale, nil
	default:
		return profile.BodyBoy, errors.New("unknown bodytype")
	}
}

func determineCharWeaponType(weaponType string) (weapon.WeaponClass, error) {
	switch weaponType {
	case "WEAPON_CLAYMORE":
		return weapon.WeaponClassClaymore, nil
	case "WEAPON_BOW":
		return weapon.WeaponClassBow, nil
	case "WEAPON_SWORD_ONE_HAND":
		return weapon.WeaponClassSword, nil
	case "WEAPON_CATALYST":
		return weapon.WeaponClassCatalyst, nil
	case "WEAPON_POLE":
		return weapon.WeaponClassSpear, nil
	default:
		return weapon.WeaponClassSword, errors.New("unknown weapontype")
	}
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
