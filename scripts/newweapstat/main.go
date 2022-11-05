package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/curves"
)

type WeaponStat struct {
	SpecializedStat  attributes.Stat
	AtkCurve         curves.WeaponStatCurve
	SpecializedCurve curves.WeaponStatCurve
	BaseAtk          float64
	BaseSpecialized  float64
}

func main() {
	// var weapons []curves.WeaponBase
	weaponConfigArray := getWeapon()
	weaponPromoDataMap := getWeaponPromoData()

	for _, weaponConfig := range weaponConfigArray {
		var weapon curves.WeaponBase
		weaponStats := convertWeaponPropsToWeaponStats(weaponConfig.WeaponProps)

		weapon.AtkCurve = weaponStats.AtkCurve
		weapon.BaseAtk = weaponStats.BaseAtk
		weapon.Specialized = weaponStats.SpecializedStat
		weapon.SpecializedCurve = weaponStats.SpecializedCurve
		weapon.BaseSpecialized = weaponStats.BaseSpecialized
		weapon.PromotionBonus = weaponPromoDataMap[weaponConfig.WeaponPromoteID]
		// fmt.Println(weaponConfig.ID)
		// fmt.Printf("%+v\n", weapon)
		// weapons = append(weapons, weapon)
	}
	// weapons[]
}

func getJsonFromFile[V WeaponConfigs | WeaponPromoteConfigs](path string) V {
	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var configData V
	json.Unmarshal(byteValue, &configData)

	return configData
}

func convertWeaponPropsToWeaponStats(weaponProps []WeaponProp) WeaponStat {
	var weaponStat WeaponStat
	for index, weaponProp := range weaponProps {
		if index >= 2 {
			log.Fatal("Too many weapon props")
		}
		if weaponProp.PropType == "FIGHT_PROP_BASE_ATTACK" {
			weaponStat.BaseAtk = weaponProp.InitValue
			attackStatCurve, err := determineWeaponStatCurve(weaponProp.Type)
			if err != nil {
				log.Fatal(err)
			}
			weaponStat.AtkCurve = attackStatCurve
		} else {
			specializedStat, err := determineWeaponStat(weaponProp.PropType)
			if err != nil {
				log.Fatal(err)
			}
			specializedStatCurve, err := determineWeaponStatCurve(weaponProp.Type)
			if err != nil {
				log.Fatal(err)
			}
			weaponStat.SpecializedStat = specializedStat
			weaponStat.BaseSpecialized = weaponProp.InitValue
			weaponStat.SpecializedCurve = specializedStatCurve
		}

	}
	return weaponStat
}

func determineWeaponStat(specializedStat string) (attributes.Stat, error) {
	switch specializedStat {
	case "FIGHT_PROP_CRITICAL_HURT":
		return attributes.CD, nil
	case "FIGHT_PROP_HEAL_ADD":
		return attributes.Heal, nil
	case "FIGHT_PROP_ATTACK_PERCENT":
		return attributes.ATKP, nil
	case "FIGHT_PROP_ELEMENT_MASTERY":
		return attributes.EM, nil
	case "FIGHT_PROP_HP_PERCENT":
		return attributes.HPP, nil
	case "FIGHT_PROP_CHARGE_EFFICIENCY":
		return attributes.ER, nil
	case "FIGHT_PROP_CRITICAL":
		return attributes.CR, nil
	case "FIGHT_PROP_PHYSICAL_ADD_HURT":
		return attributes.PhyP, nil
	case "FIGHT_PROP_ELEC_ADD_HURT":
		return attributes.ElectroP, nil
	case "FIGHT_PROP_ROCK_ADD_HURT":
		return attributes.GeoP, nil
	case "FIGHT_PROP_FIRE_ADD_HURT":
		return attributes.PyroP, nil
	case "FIGHT_PROP_WATER_ADD_HURT":
		return attributes.HydroP, nil
	case "FIGHT_PROP_DEFENSE_PERCENT":
		return attributes.DEFP, nil
	case "FIGHT_PROP_ICE_ADD_HURT":
		return attributes.CryoP, nil
	case "FIGHT_PROP_WIND_ADD_HURT":
		return attributes.AnemoP, nil
	case "FIGHT_PROP_GRASS_ADD_HURT":
		return attributes.DendroP, nil
	case "":
		return attributes.NoStat, nil
	default:
		return attributes.NoStat, errors.New("unknown stat type")

	}
}

func determineWeaponStatCurve(input string) (curves.WeaponStatCurve, error) {
	switch input {
	case "GROW_CURVE_CRITICAL_101":
		return curves.GROW_CURVE_CRITICAL_101, nil
	case "GROW_CURVE_CRITICAL_201":
		return curves.GROW_CURVE_CRITICAL_201, nil
	case "GROW_CURVE_CRITICAL_301":
		return curves.GROW_CURVE_CRITICAL_301, nil
	case "GROW_CURVE_ATTACK_101":
		return curves.GROW_CURVE_ATTACK_101, nil
	case "GROW_CURVE_ATTACK_102":
		return curves.GROW_CURVE_ATTACK_102, nil
	case "GROW_CURVE_ATTACK_103":
		return curves.GROW_CURVE_ATTACK_103, nil
	case "GROW_CURVE_ATTACK_104":
		return curves.GROW_CURVE_ATTACK_104, nil
	case "GROW_CURVE_ATTACK_105":
		return curves.GROW_CURVE_ATTACK_105, nil
	case "GROW_CURVE_ATTACK_201":
		return curves.GROW_CURVE_ATTACK_201, nil
	case "GROW_CURVE_ATTACK_202":
		return curves.GROW_CURVE_ATTACK_202, nil

	case "GROW_CURVE_ATTACK_203":
		return curves.GROW_CURVE_ATTACK_203, nil
	case "GROW_CURVE_ATTACK_204":
		return curves.GROW_CURVE_ATTACK_204, nil
	case "GROW_CURVE_ATTACK_205":
		return curves.GROW_CURVE_ATTACK_205, nil
	case "GROW_CURVE_ATTACK_301":
		return curves.GROW_CURVE_ATTACK_301, nil
	case "GROW_CURVE_ATTACK_302":
		return curves.GROW_CURVE_ATTACK_302, nil
	case "GROW_CURVE_ATTACK_303":
		return curves.GROW_CURVE_ATTACK_303, nil
	case "GROW_CURVE_ATTACK_304":
		return curves.GROW_CURVE_ATTACK_304, nil
	case "GROW_CURVE_ATTACK_305":
		return curves.GROW_CURVE_ATTACK_305, nil
	default:
		return curves.GROW_CURVE_ATTACK_101, errors.New("unknown curve type")
	}
}
