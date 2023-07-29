package main

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/curves"
)

type WeaponPromoteConfig []struct {
	WeaponPromoteID int `json:"weaponPromoteId"`
	// CostItems       []struct {
	// } `json:"costItems"`
	AddProps       []AddProp
	UnlockMaxLevel int `json:"unlockMaxLevel"`
	PromoteLevel   int `json:"promoteLevel,omitempty"`
	// RequiredPlayerLevel int `json:"requiredPlayerLevel,omitempty"`
	// CoinCost            int `json:"coinCost,omitempty"`
}
type AddProp struct {
	PropType string  `json:"propType"`
	Value    float64 `json:"value"`
}

func getCharSpecializedStatandPromoData() map[int][]curves.PromoData {

	weaponPromoteJson, err := fetchJsonFromUrl("https://raw.githubusercontent.com/Dimbreath/GenshinData/master/ExcelBinOutput/WeaponPromoteExcelConfigData.json")
	if err != nil {
		log.Fatal(err)
	}

	var weaponPromotes WeaponPromoteConfig
	if err := json.Unmarshal([]byte(weaponPromoteJson), &weaponPromotes); err != nil {
		log.Fatal(err)
	}

	// reshape avatarPromotes to map of avatarPromoteId to PromoData array of 7 items(1 for each ascension)
	promoDataMap := make(map[int][]curves.PromoData)
	for _, weaponPromo := range weaponPromotes {
		partialPromoData := addPropArraytoPromoData(weaponPromo.AddProps)
		partialPromoData.MaxLevel = (weaponPromo.UnlockMaxLevel)
		promoDataMap[weaponPromo.WeaponPromoteID] = append(promoDataMap[weaponPromo.WeaponPromoteID], partialPromoData)
	}

	return promoDataMap
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

func determineStat(specializedStat string) (attributes.Stat, error) {
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
		return attributes.EndStatType, nil
	default:
		return attributes.EndStatType, errors.New("unknown stat type")

	}
}
