package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/genshinsim/gcsim/pkg/core/curves"
	"github.com/genshinsim/gcsim/pkg/core/info"
)

type BaseStatCurves struct {
	HpCurve  string
	AtkCurve string
	DefCurve string
}

func main() {
	var err error
	avatars, textMapId := getAvatarArray()
	locationMap := getCharLocationMap()
	specializedStatMap, promoDataMap := getCharSpecializedStatandPromoData()
	elementMap := getAvatarElementMap()

	//currently only gets english
	_ = generateAvatarNameMap(textMapId)

	characterArray := make([]curves.CharBase, len(avatars))

	for _, avatar := range avatars {
		char := curves.CharBase{}
		charBaseStatCurves := extractCharStatCurves(avatar.PropGrowCurves)
		charName := determineCharName(avatar.IconName)

		char.BaseHP = avatar.HpBase
		char.BaseAtk = avatar.AttackBase
		char.BaseDef = avatar.DefenseBase
		char.PromotionBonus = promoDataMap[avatar.AvatarPromoteID]
		char.Element = convertElement(elementMap[avatar.SkillDepotID])
		char.Specialized, err = determineStat(specializedStatMap[avatar.AvatarPromoteID])
		if err != nil {
			log.Fatal("Unknown specialized stat for character ", charName, ": ", specializedStatMap[avatar.ID])
		}
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
			log.Fatal("Unknown weapon type for character ", charName, ": ", avatar.WeaponType)
		}
		char.HPCurve, err = determineCharStatCurves(charBaseStatCurves.HpCurve)
		if err != nil {
			log.Fatal("Unknown stat curve for character ", charName, ": ", charBaseStatCurves.HpCurve)
		}
		char.AtkCurve, err = determineCharStatCurves(charBaseStatCurves.AtkCurve)
		if err != nil {
			log.Fatal("Unknown stat curve for character ", charName, ": ", charBaseStatCurves.AtkCurve)
		}
		char.DefCurve, err = determineCharStatCurves(charBaseStatCurves.DefCurve)
		if err != nil {
			log.Fatal("Unknown stat curve for character ", charName, ": ", charBaseStatCurves.DefCurve)
		}

		characterArray = append(characterArray, char)
		// fmt.Println(charName)
		// fmt.Printf("%+v\n", char)
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

func determineCharName(iconName string) string {
	return strings.Replace(iconName, "UI_AvatarIcon_", "", 1)
}

func determineCharBody(bodyType string) (info.BodyType, error) {
	switch bodyType {
	case "BODY_LOLI":
		return info.BodyLoli, nil
	case "BODY_GIRL":
		return info.BodyGirl, nil
	case "BODY_LADY":
		return info.BodyLady, nil
	case "BODY_BOY":
		return info.BodyBoy, nil
	case "BODY_MALE":
		return info.BodyMale, nil
	default:
		return info.BodyBoy, errors.New("unknown bodytype")
	}
}

func determineCharWeaponType(weaponType string) (info.WeaponClass, error) {
	switch weaponType {
	case "WEAPON_CLAYMORE":
		return info.WeaponClassClaymore, nil
	case "WEAPON_BOW":
		return info.WeaponClassBow, nil
	case "WEAPON_SWORD_ONE_HAND":
		return info.WeaponClassSword, nil
	case "WEAPON_CATALYST":
		return info.WeaponClassCatalyst, nil
	case "WEAPON_POLE":
		return info.WeaponClassSpear, nil
	default:
		return info.WeaponClassSword, errors.New("unknown weapontype")
	}
}

func extractCharStatCurves(propGrowCurves PropGrowCurves) BaseStatCurves {
	//reshape avatar.PropGrowCurves to be a map of Type to GrowCurve
	growCurveMap := make(map[string]string)
	for _, growCurve := range propGrowCurves {
		growCurveMap[growCurve.Type] = growCurve.GrowCurve
	}

	return BaseStatCurves{
		HpCurve:  growCurveMap["FIGHT_PROP_BASE_HP"],
		AtkCurve: growCurveMap["FIGHT_PROP_BASE_ATTACK"],
		DefCurve: growCurveMap["FIGHT_PROP_BASE_DEFENSE"],
	}
}

func determineCharStatCurves(statCurve string) (curves.CharStatCurve, error) {
	//print statCurve
	switch statCurve {
	case "GROW_CURVE_ATTACK_S5":
		return curves.GROW_CURVE_ATTACK_S5, nil
	case "GROW_CURVE_ATTACK_S4":
		return curves.GROW_CURVE_ATTACK_S4, nil
	case "GROW_CURVE_HP_S5":
		return curves.GROW_CURVE_HP_S5, nil
	case "GROW_CURVE_HP_S4":
		return curves.GROW_CURVE_HP_S4, nil

	default:
		return curves.GROW_CURVE_HP_S5, errors.New("unknown stat curve")
	}
}
