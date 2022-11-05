package main

import (
	"github.com/genshinsim/gcsim/pkg/core/curves"
)

type WeaponPromoteConfigs []struct {
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

func getWeaponPromoData() map[int][]curves.PromoData {
	weaponPromoteConfigs := getJsonFromFile[WeaponPromoteConfigs]("../ExcelBinOutput/WeaponPromoteExcelConfigData.json")

	// reshape avatarPromotes to map of avatarPromoteId to PromoData array of 7 items(1 for each ascension)
	promoDataMap := make(map[int][]curves.PromoData)
	for _, weaponPromoteConfig := range weaponPromoteConfigs {
		promoData := addPropArraytoPromoData(weaponPromoteConfig.AddProps)
		promoData.MaxLevel = (weaponPromoteConfig.UnlockMaxLevel)
		// fmt.Printf("%+v\n", promoData)
		// fmt.Printf("%+v\n", weaponPromoteConfig.AddProps)

		promoDataMap[weaponPromoteConfig.WeaponPromoteID] = append(promoDataMap[weaponPromoteConfig.WeaponPromoteID], promoData)

	}

	return promoDataMap
}

//not entirely correct, but good enough for now
func addPropArraytoPromoData(addProps []AddProp) curves.PromoData {
	var out curves.PromoData
	for _, prop := range addProps {
		switch prop.PropType {
		case "FIGHT_PROP_BASE_HP":
			out.HP = prop.Value
		case "FIGHT_PROP_BASE_ATTACK":
			out.Atk = prop.Value
			return out
		case "FIGHT_PROP_BASE_DEFENSE":
			out.Def = prop.Value
		default:
			out.Special = prop.Value
		}
	}
	return out
}
